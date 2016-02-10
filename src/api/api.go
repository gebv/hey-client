package api

func InitApi() {
	r := Srv.Router.PathPrefix("/api/v1").Subrouter()

	// Entry point
	InitOauthCallback(r)

	InitUsers(r)
}