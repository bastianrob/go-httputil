package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

//HTTPAdapter public access
type HTTPAdapter interface {
	Get(url string) (response *http.Response, err error)
	PostJSON(url string, data interface{}) (response *http.Response, err error)
}

//private implementation of HTTPAdapter
type httpAdapterImpl struct {
	httpClient *http.Client
}

//NewHTTPAdapter new instance of HTTPAdapter
func NewHTTPAdapter(httpClient *http.Client) HTTPAdapter {
	return &httpAdapterImpl{httpClient: httpClient}
}

func (adapter *httpAdapterImpl) Get(url string) (response *http.Response, err error) {
	return adapter.httpClient.Get(url)
}

//PostJSON data to a URL
func (adapter *httpAdapterImpl) PostJSON(url string, data interface{}) (response *http.Response, err error) {
	if data == nil {
		return nil, fmt.Errorf("[HTTPAdapter][URL:%s] Data cannot be nil", url)
	}

	databyte, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("[HTTPAdapter][URL:%s] Unable to marshal data: %+v", url, data)
	}

	var properJSON map[string]interface{}
	err = json.Unmarshal(databyte, &properJSON)
	if err != nil {
		return nil, fmt.Errorf("[HTTPAdapter][URL:%s] Data is not a proper JSON: %+v", url, data)
	}

	response, err = adapter.httpClient.Post(url, "application/json", bytes.NewBuffer(databyte))
	if err != nil {
		return nil, fmt.Errorf("[HTTPAdapter][URL:%s] HTTP Request error: %s", url, err.Error())
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return response, fmt.Errorf("[HTTPAdapter][URL:%s] HTTP Request error code: %d", url, response.StatusCode)
	}
	return
}
