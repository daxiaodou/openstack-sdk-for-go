package http

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	POST   = "POST"
	GET    = "GET"
	DELETE = "DELETE"
)

type RequestParam struct {
	Url          string
	AccessToken  string
	Content_Type string
	HttpAction   string
	Params       string
	Attr         map[string]interface{}
}

func (r *RequestParam) DoGetRequest() (result string, err error) {
	if len(r.Url) <= 0 {
		fmt.Errorf("Url is null")
		err = errors.New("Url is null")
		return "", err
	}

	req := r.getRequest()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http get error! ", *req, err.Error())
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("http ioutil error! ", req, err.Error())
		return "", err
	}
	result = string(body)
	return result, err
}

func (r *RequestParam) DoPost() (result *http.Response, err error) {
	if len(r.Url) <= 0 {
		fmt.Println("Url is null! ", r.Url)
		err = errors.New("Url is null! ")
		return new(http.Response), err
	}

	req := r.getRequest()
	client := &http.Client{}
	result, err = client.Do(req)
	return result, err
}

func (r *RequestParam) DoPostRequest() (result string, err error) {
	if len(r.Url) <= 0 {
		fmt.Println("Url is null! ", r.Url)
		err = errors.New("Url is null! ")
		return "", err
	}

	req := r.getRequest()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http post error! ", *req, err.Error())
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("http ioutil error! ", req, err.Error())
		return "", err
	}
	result = string(body)
	return result, err
}

func (r *RequestParam) DoDeleteRequest() (result string, err error) {
	if len(r.Url) <= 0 {
		fmt.Println("Url is null! ", r.Url)
		err = errors.New("Url is null! ")
		return "", err
	}

	req := r.getRequest()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http del error! ", *req, err.Error())
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("http ioutil error! ", req, err.Error())
		return "", err
	}
	result = string(body)
	return result, err
}

func (r *RequestParam) getRequest() (req *(http.Request)) {
	reqBody := strings.NewReader(string(r.Params))
	req, err := http.NewRequest(r.HttpAction, r.Url, reqBody)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", r.Content_Type)
	req.Header.Set("X-Auth-Token", r.AccessToken)
	return req
}

func (r *RequestParam) DoPostIdentity() (result *http.Response, err error) {
	if len(r.Url) <= 0 {
		fmt.Println("Url is null! ", r.Url)
		err = errors.New("Url is null! ")
		return new(http.Response), err
	}

	req := r.getRequestIdentity()
	client := &http.Client{}
	result, err = client.Do(req)
	return result, err
}

func (r *RequestParam) getRequestIdentity() (req *(http.Request)) {
	reqBody := strings.NewReader(string(r.Params))
	req, err := http.NewRequest(r.HttpAction, r.Url, reqBody)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", r.Content_Type)
	return req
}
