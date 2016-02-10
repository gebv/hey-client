package api

import (
	"golang.org/x/oauth2"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/golang/glog"
	// "net/url"
	// "io/ioutil"
	"errors"
)

var heyOauth2Config *oauth2.Config
var heyClient *HeyClient
var clientId = "b4c8dd5b-852c-460a-9b4a-26109f9162a2"

func InitOauthCallback(r *mux.Router) {
	sr := r.PathPrefix("/oauth2").Subrouter()

	sr.HandleFunc("/callback", CallbackHandler)
	sr.HandleFunc("/me", MeHandler)
}

func init() {
	
	heyOauth2Config = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: "demo",
		Scopes:       []string{},
		RedirectURL:  "http://192.168.1.36:65002/api/v1/oauth2/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://192.168.1.36:65001/api/v1/oauth2/authorize",
			TokenURL: "http://192.168.1.36:65001/api/v1/oauth2/token",
		},
	}

	heyClient = NewHeyClient(heyOauth2Config)
}

func AuthorizedAtHeyService() error {
	authCode := oauth2.SetAuthURLParam("client_key", "demo")
	glog.Infof("oauth login url='%v'", heyOauth2Config.AuthCodeURL("csrf", authCode))

	resp, err := http.Get(heyOauth2Config.AuthCodeURL("csrf", authCode))

	if err != nil {
		glog.Infof("error connect client_id=%v, err=%v", clientId, err)
		return err
	}

	if resp.StatusCode == http.StatusOK {
		glog.Infof("connect success client_id=%v", clientId)
		return nil
	}

	return errors.New("authorization error. unknown reason")
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	data, err := heyClient.Me()

	if err != nil {
		glog.Errorf("hey: action='/api/v1/me', err='%s'", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(data.([]byte))
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	t, err := heyOauth2Config.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))

	if err != nil {
		glog.Warningf("oauth error: err='%s'", err.Error())
		return
	}

	heyClient.Client = heyOauth2Config.Client(oauth2.NoContext, t)
	w.WriteHeader(http.StatusOK)
}