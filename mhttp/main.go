package mhttp

import "net/http"

func main() {
	srv := &http.Server{Addr: ":8080"}
	srv.ListenAndServe()

	srv.Shutdown(nil) // Shutdown gracefully shuts down the server without interrupting any
	//srv.RegisterOnShutdown()
}
