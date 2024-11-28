package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"
)

const (
	rootRoute = "/"
)

func HttpServer(addr string, routes Routes, cert, key string) error {
	mux := http.DefaultServeMux
	handlers := make(map[string]http.Handler)
	paths := make(map[string]string)

	if len(routes.Values) == 0 {
		_ = routes.Set(".")
	}

	for _, route := range routes.Values {
		handlers[route.Route] = &fileHandler{
			route:       route.Route,
			path:        route.Path,
			allowUpload: true,
			allowDelete: true,
		}
		paths[route.Route] = route.Path
	}

	for route, path := range paths {
		mux.Handle(route, handlers[route])
		logs.Info("serving local path %q on %q", path, route)
	}

	_, rootRouteTaken := handlers[rootRoute]
	if !rootRouteTaken {
		route := routes.Values[0].Route
		mux.Handle(rootRoute, http.RedirectHandler(route, http.StatusTemporaryRedirect))
		logs.Info("redirecting to %q from %q", route, rootRoute)
	}

	binaryPath, _ := os.Executable()
	if binaryPath == "" {
		binaryPath = "server"
	}

	if cert != "" && key != "" {
		logs.Info("%s (HTTPS) listening on %q", filepath.Base(binaryPath), addr)
		return http.ListenAndServeTLS(addr, cert, key, mux)
	}

	logs.Info("%s listening on %q", filepath.Base(binaryPath), addr)
	return http.ListenAndServe(addr, mux)
}
