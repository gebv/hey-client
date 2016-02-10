package api

import (
	"net/http"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

func InitUsers(r *mux.Router) {
	sr := r.PathPrefix("/users").Subrouter()

	sr.HandleFunc("/{ext_id}", GetUserHandler)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {

	data, err := heyClient.GetUser(mux.Vars(r)["ext_id"])

	if err != nil {
		glog.Errorf("hey: action='GetUser', err='%s'", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(data)
}
