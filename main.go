package main

import "backend/api"

func main() {
	// http.HandleFunc("/", api.Handler)
	// http.ListenAndServe(":3000", nil)
	api.StartServer()
}
