package main

var server_status bool

func ServerEnable() bool {
	if server_status {
		server_status = false
	} else {
		server_status = true
	}
	return server_status
}
