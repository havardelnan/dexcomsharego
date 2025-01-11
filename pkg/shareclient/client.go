package shareclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

var (
	baseUrl = "https://shareous1.dexcom.com/ShareWebServices/Services/"
)

type ShareClient struct {
	client *http.Client
}

func NewShareClient() *ShareClient {
	return &ShareClient{
		client: &http.Client{},
	}
}

func (c *ShareClient) PostJSON(path string, body interface{}, response interface{}) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", baseUrl+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	setCommonHeaders(req)
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = handleResponse(res, response)
	if err != nil {
		return err
	}

	return nil

}

func handleResponse(res *http.Response, out any) error {

	if res.StatusCode > 399 || res.StatusCode < 200 {
		return fmt.Errorf("http error: %s from %s", res.Status, res.Request.URL)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.Header.Get("Content-Type") == "text/plain" {
		v := reflect.ValueOf(out)
		if v.Kind() != reflect.Ptr || v.IsNil() {
			return fmt.Errorf("out must be a pointer and not nil")
		}
		v.Elem().Set(reflect.ValueOf(string(body)))
		return nil
	}

	err = json.Unmarshal(body, out)
	if err != nil {
		return err
	}

	return nil
}

func setCommonHeaders(req *http.Request) {
	req.Header.Set("Accept-Encoding", "application/json")
	req.Header.Set("Content-Type", "application/json")
}
