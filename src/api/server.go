package api

import (
	"github.com/gorilla/mux"
	"github.com/golang/glog"
	"net/http"
)

type Server struct {
	Router *mux.Router
}

var Srv *Server

func NewServer() {

	glog.Info("Server is initializing...")

	Srv = &Server{}
	Srv.Router = mux.NewRouter()
}

func StartServer() {

	var handler http.Handler = Srv.Router

	go func() {
		http.ListenAndServe("192.168.1.36:65002", handler)
	}()
}

func StopServer() {
	glog.Info("Stopping Server...")

	// manners.Close()

	glog.Info("Server stopped")
}