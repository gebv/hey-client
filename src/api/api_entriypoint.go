package api

import (
	"golang.org/x/oauth2"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/golang/glog"
	"net/url"
	"io/ioutil"
)

var heyOauth2Config *oauth2.Config
var heyClient *http.Client

func InitOauthCallback(r *mux.Router) {
	sr := r.PathPrefix("/oauth2").Subrouter()

	sr.HandleFunc("/callback", CallbackHandler)
	sr.HandleFunc("/me", MeHandler)
}

func init() {
	heyOauth2Config = &oauth2.Config{
		ClientID:     "demo",
		ClientSecret: "demo",
		Scopes:       []string{},
		RedirectURL:  "http://192.168.1.36:65002/api/v1/oauth2/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://192.168.1.36:65001/api/v1/oauth2/authorize",
			TokenURL: "http://192.168.1.36:65001/api/v1/oauth2/token",
		},
	}

	authCode := oauth2.SetAuthURLParam("client_key", "demo")
	glog.Infof("oauth login url='%v'", heyOauth2Config.AuthCodeURL("csrf", authCode))
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	_baseUrl, _ := url.Parse(heyOauth2Config.Endpoint.AuthURL)

	_url := url.URL{}
	_url.Scheme = _baseUrl.Scheme
	_url.Host = _baseUrl.Host
	_url.Path = "/api/v1/oauth2/me"

	q := _url.Query()

	// for key, value := range fields {
	// 	q.Add(key, value)
	// }

	_url.RawQuery = q.Encode()

	resp, err := heyClient.Get(_url.String())

	if err != nil {
		glog.Errorf("hey: action='/api/v1/me', err='%s'", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch resp.StatusCode {
	case http.StatusOK:
		data, _ := ioutil.ReadAll(resp.Body)
		w.Write(data)
	default:
		data, _ := ioutil.ReadAll(resp.Body)

		glog.Errorf("hey: action='/api/v1/me', status_code='%v', body='%s'", resp.StatusCode, data)
		w.WriteHeader(http.StatusBadRequest)
		return	
	}
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	t, err := heyOauth2Config.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))

	if err != nil {
		glog.Warningf("oauth error: err='%s'", err.Error())
		return
	}

	heyClient = heyOauth2Config.Client(oauth2.NoContext, t)
	w.WriteHeader(http.StatusOK)
}