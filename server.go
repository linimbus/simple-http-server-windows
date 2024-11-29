package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"html/template"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/astaxie/beego/logs"
)

const (
	zipKey         = "zip"
	zipValue       = "true"
	zipContentType = "application/zip"

	osPathSeparator = string(filepath.Separator)
)

const directoryListingTemplateText = `
<html>
<head>
	<title>{{ .Title }}</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<style>body{font-family: sans-serif;}td{padding:.5em;}a{display:block;}tbody tr:nth-child(odd){background:#eee;}.number{text-align:right}.text{text-align:left;word-break:break-all;}canvas,table{width:100%;max-width:100%;}</style>
</head>
<body>
<h1>{{ .Title }}</h1>
{{ if or .Files .AllowUpload }}
<table>
	<thead>
		<th></th>
		<th colspan=2 class=number>Size (bytes)</th>
	</thead>
	<tbody>
	{{- if .Files }}
	<tr><td colspan=3><a href="{{ .ZipURL }}">.zip of all files</a></td></tr>
	{{- end }}
	{{- range .Files }}
	<tr>
		{{ if (not .IsDir) }}
		<td class=text><a href="{{ .URL.String }}">{{ .Name }}</td>
		<td class=number>{{.Size.String }}</td>
		<td class=number>({{ .Size | printf "%d" }})</td>
		{{ else }}
		<td colspan=3 class=text><a href="{{ .URL.String }}">{{ .Name }}</td>
		{{ end }}
	</tr>
	{{- end }}
	{{- if .AllowUpload }}
	<tr><td colspan=3><form method="post" enctype="multipart/form-data"><input required name="file" type="file"/><input value="Upload" type="submit"/></form></td></tr>
	{{- end }}
	</tbody>
</table>
{{ end }}
</body>
</html>
`

type fileSizeBytes int64

func (f fileSizeBytes) String() string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)
	divBy := func(x int64) int {
		return int(math.Round(float64(f) / float64(x)))
	}
	switch {
	case f < KB:
		return fmt.Sprintf("%d", f)
	case f < MB:
		return fmt.Sprintf("%dK", divBy(KB))
	case f < GB:
		return fmt.Sprintf("%dM", divBy(MB))
	case f >= GB:
		fallthrough
	default:
		return fmt.Sprintf("%dG", divBy(GB))
	}
}

type directoryListingFileData struct {
	Name  string
	Size  fileSizeBytes
	IsDir bool
	URL   *url.URL
}

type directoryListingData struct {
	Title       string
	ZipURL      *url.URL
	Files       []directoryListingFileData
	AllowUpload bool
}

var (
	directoryListingTemplate = template.Must(template.New("").Parse(directoryListingTemplateText))
)

type fileHandler struct {
	route       string
	path        string
	allowUpload bool
	allowDelete bool

	timeout int
	address string
	server  *http.Server

	flowbytes int64
	requests  int64
	sessions  int64

	sync.WaitGroup
}

func (f *fileHandler) serveStatus(w http.ResponseWriter, r *http.Request, status int) error {
	w.WriteHeader(status)
	_, err := w.Write([]byte(http.StatusText(status)))
	if err != nil {
		return err
	}
	return nil
}

func (f *fileHandler) serveZip(w http.ResponseWriter, r *http.Request, osPath string) error {
	w.Header().Set("Content-Type", zipContentType)
	name := filepath.Base(osPath) + ".zip"
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename=%q`, name))
	return FileZip(w, osPath)
}

func (f *fileHandler) serveDir(w http.ResponseWriter, r *http.Request, osPath string) error {
	d, err := os.Open(osPath)
	if err != nil {
		return err
	}
	files, err := d.Readdir(-1)
	if err != nil {
		return err
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return directoryListingTemplate.Execute(w, directoryListingData{
		AllowUpload: f.allowUpload,
		Title: func() string {
			relPath, _ := filepath.Rel(f.path, osPath)
			urlPath := filepath.Join(filepath.Base(f.path), relPath)
			return strings.ReplaceAll(urlPath, osPathSeparator, "/")
		}(),
		ZipURL: func() *url.URL {
			url := *r.URL
			q := url.Query()
			q.Set(zipKey, zipValue)
			url.RawQuery = q.Encode()
			return &url
		}(),
		Files: func() (out []directoryListingFileData) {
			for _, d := range files {
				name := d.Name()
				if d.IsDir() {
					name += "/"
				}
				fileData := directoryListingFileData{
					Name:  name,
					IsDir: d.IsDir(),
					Size:  fileSizeBytes(d.Size()),
					URL: func() *url.URL {
						url := *r.URL
						url.Path = path.Join(url.Path, name)
						if d.IsDir() {
							url.Path += "/"
						}
						return &url
					}(),
				}
				out = append(out, fileData)
			}
			return out
		}(),
	})
}

func (f *fileHandler) serveUploadTo(w http.ResponseWriter, r *http.Request, osPath string) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	in, h, err := r.FormFile("file")
	if err == http.ErrMissingFile {
		w.Header().Set("Location", r.URL.String())
		w.WriteHeader(303)
	}
	if err != nil {
		return err
	}
	outPath := filepath.Join(osPath, filepath.Base(h.Filename))
	out, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	w.Header().Set("Location", r.URL.String())
	w.WriteHeader(303)
	return nil
}

// ServeHTTP is http.Handler.ServeHTTP
func (f *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logs.Info("http server request [%s] %s %s %s", f.path, r.RemoteAddr, r.Method, r.URL.String())

	atomic.AddInt64(&f.requests, 1)

	urlPath := r.URL.Path
	if !strings.HasPrefix(urlPath, "/") {
		urlPath = "/" + urlPath
	}
	urlPath = strings.TrimPrefix(urlPath, f.route)
	urlPath = strings.TrimPrefix(urlPath, "/"+f.route)

	osPath := strings.ReplaceAll(urlPath, "/", osPathSeparator)
	osPath = filepath.Clean(osPath)
	osPath = filepath.Join(f.path, osPath)

	info, err := os.Stat(osPath)
	switch {
	case os.IsNotExist(err):
		_ = f.serveStatus(w, r, http.StatusNotFound)
	case os.IsPermission(err):
		_ = f.serveStatus(w, r, http.StatusForbidden)
	case err != nil:
		_ = f.serveStatus(w, r, http.StatusInternalServerError)
	case !f.allowDelete && r.Method == http.MethodDelete:
		_ = f.serveStatus(w, r, http.StatusForbidden)
	case !f.allowUpload && r.Method == http.MethodPost:
		_ = f.serveStatus(w, r, http.StatusForbidden)
	case r.URL.Query().Get(zipKey) != "":
		err := f.serveZip(w, r, osPath)
		if err != nil {
			_ = f.serveStatus(w, r, http.StatusInternalServerError)
		}
	case f.allowUpload && info.IsDir() && r.Method == http.MethodPost:
		err := f.serveUploadTo(w, r, osPath)
		if err != nil {
			_ = f.serveStatus(w, r, http.StatusInternalServerError)
		}
	case f.allowDelete && !info.IsDir() && r.Method == http.MethodDelete:
		err := os.Remove(osPath)
		if err != nil {
			_ = f.serveStatus(w, r, http.StatusInternalServerError)
		}
	case info.IsDir():
		err := f.serveDir(w, r, osPath)
		if err != nil {
			_ = f.serveStatus(w, r, http.StatusInternalServerError)
		}
	default:
		http.ServeFile(w, r, osPath)
	}
}

func (f *fileHandler) Shutdown() error {
	context, cencel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	err := f.server.Shutdown(context)
	cencel()
	if err != nil {
		logs.Error("http file server ready to shut down fail, %s", err.Error())
	}
	f.Wait()
	return err
}

func CreateHttpServer(addr, folder string,
	upload, delete bool,
	https bool, cert, key string) (*fileHandler, error) {

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logs.Error("http file server listen %s address fail", addr)
		return nil, err
	}
	logs.Info("http file server listening on %s", addr)

	var tlsConfig *tls.Config
	if https {
		tlsConfig, err = CreateTlsConfig(cert, key)
		if err != nil {
			logs.Error("create tls config for http server fail, %s", err.Error())
			return nil, err
		}
		listen = tls.NewListener(listen, tlsConfig)
	}

	fileHandler := &fileHandler{
		route:       "/",
		path:        folder,
		allowUpload: upload,
		allowDelete: delete,
	}

	httpserver := &http.Server{
		Handler:      fileHandler,
		ReadTimeout:  time.Duration(60) * time.Second,
		WriteTimeout: time.Duration(60) * time.Second,
		TLSConfig:    tlsConfig,
	}

	fileHandler.server = httpserver
	fileHandler.Add(1)

	go func() {
		defer fileHandler.Done()
		err = httpserver.Serve(listen)
		if err != nil {
			logs.Error("http server attach listen instance fail, %s", err.Error())
		}
	}()

	return fileHandler, nil
}

// func HttpServer(addr string, folder string, cert, key string) error {

// 	mux := http.DefaultServeMux

// 	fileHandler := &fileHandler{
// 		route:       "/",
// 		path:        folder,
// 		allowUpload: true,
// 		allowDelete: true,
// 	}

// 	mux.Handle("/", fileHandler)

// 	logs.Info("http file server listening on %s", addr)

// 	if cert != "" && key != "" {
// 		return http.ListenAndServeTLS(addr, cert, key, mux)
// 	}

// 	return http.ListenAndServe(addr, mux)
// }