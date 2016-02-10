package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"golang.org/x/oauth2"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
)

func NewHeyClient(c *oauth2.Config) *HeyClient {
	api := new(HeyClient)
	api.Config = c

	return api
}


type HeyClient struct {
	Client        *http.Client
	Config        *oauth2.Config
}

func (c *HeyClient) Me() (interface{}, error) {
	url := c.buildUrl("/api/v1/oauth2/me", map[string]string{})

	data, err := c.get(url)

	return data, err
}

func (c *HeyClient) CreateUser(ext_id string, info map[string]string) (interface{}, error) {
	url := c.buildUrl("/api/v1/users", map[string]string{})

	params := map[string]interface{}{
		"ext_id":   ext_id,
		"ext_info": info,
	}

	b, _ := json.Marshal(params)

	return c.post(url, b)
}

func (c *HeyClient) GetUser(ext_id string) ([]byte, error) {
	url := c.buildUrl("/api/v1/users/", map[string]string{"ext_id": ext_id})

	data, err := c.get(url)

	glog.Infof("get url='%v', data='%s'", url, data)

	return data, err
}

func (c *HeyClient) buildUrl(path string, fields map[string]string) string {
	_configRedirectUrl, _ := url.Parse(c.Config.Endpoint.AuthURL)

	urlInf := url.URL{}
	urlInf.Scheme = _configRedirectUrl.Scheme
	urlInf.Host = _configRedirectUrl.Host
	urlInf.Path = path

	q := urlInf.Query()

	for key, value := range fields {
		q.Add(key, value)
	}

	urlInf.RawQuery = q.Encode()

	return urlInf.String()
}

func (c *HeyClient) post(url string, body []byte) (data []byte, err error) {
	var resp *http.Response
	resp, err = c.Client.Post(url, "application/octet-stream", bytes.NewBuffer(body))

	if err != nil {
		glog.Errorf("Ответ от запроса, %s", err)
		return []byte{}, err
	}

	defer func() {
		resp.Body.Close()
	}()

	data, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		glog.Errorln("Ошибка обработки данных, %s", err)
		return []byte{}, err
	}

	if resp.StatusCode != http.StatusOK || resp.StatusCode != http.StatusCreated {

		return data, fmt.Errorf("Ошибка создания, см. ответ %d", resp.StatusCode)
	}

	return data, nil
}

func (c *HeyClient) postForm(url string, params map[string]string, datas map[string][]byte) (data []byte, err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// TODO: Что будет если в params и datas будут одинаковые параметры?

	for fieldname, value := range params {
		err = w.WriteField(fieldname, value)
	}

	if err != nil {
		glog.Errorf("Формирование запроса, %s", err)
		return []byte{}, err
	}

	for fieldname, value := range datas {
		fw, errCreateField := w.CreateFormField(fieldname)

		if errCreateField != nil {
			err = errCreateField
			break
		}

		if _, errWriteField := fw.Write(value); errWriteField != nil {
			err = errWriteField
			break
		}
	}

	if err != nil {
		glog.Errorf("Формирование запроса, %s", err)
		return []byte{}, err
	}

	w.Close()

	var resp *http.Response
	resp, err = c.Client.Post(url, w.FormDataContentType(), &b)

	if err != nil {
		glog.Errorf("Ответ от запроса, %s", err)
		return []byte{}, err
	}

	defer func() {
		resp.Body.Close()
	}()

	data, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		glog.Errorln("Ошибка обработки данных, %s", err)
		return []byte{}, err
	}

	if resp.StatusCode != http.StatusOK || resp.StatusCode != http.StatusCreated {

		return data, fmt.Errorf("Ошибка создания, см. ответ %d", resp.StatusCode)
	}

	return data, nil
}

func (c *HeyClient) head(url string) bool {
	resp, err := c.Client.Head(url)

	if err != nil {
		glog.Errorf("Head запрос %s", err)
		return false
	}

	if resp.StatusCode != http.StatusOK {
		return false
	}

	// resp.Header.Get("X-id")

	return true
}

func (c *HeyClient) get(url string) (data []byte, err error) {
	resp, err := c.Client.Get(url)

	if err != nil {
		return []byte{}, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		data, err = ioutil.ReadAll(resp.Body)

		if err != nil {
			return []byte{}, err
		}
	default:
	}

	return
}

func (c *HeyClient) IsConnected() bool {
	return c.Client != nil
}
