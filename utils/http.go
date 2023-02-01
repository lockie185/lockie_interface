package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ResData struct {
	Status bool        `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func (r *ResData) ToJson() []byte {
	data, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	return data
}

func (r *ResData) ToJsonString() string {
	data := r.ToJson()
	if data == nil {
		return ""
	}
	return string(data)
}

func (r *ResData) ParseJson(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *ResData) ParseJsonString(str string) error {
	return r.ParseJson([]byte(str))
}

type HD map[string]interface{}

func (h *HD) ToJson() string {
	data, err := json.Marshal(h)
	if err != nil {
		return ""
	}
	return string(data)
}

func HttpPost(urlStr string, postData *url.Values) (map[string]interface{}, error) {
	req, err := http.PostForm(urlStr, *postData)
	if err != nil {
		return nil, err
	}

	return getRequestData(req)
}

//HttpPostJson
/**
 * http post JSON 数据,返回一个MAP数据
 */
func HttpPostJson(urlStr string, postData map[string]interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(postData)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(b)
	req, err := http.Post(urlStr, "application/json;charset=utf-8", body)
	if err != nil {
		return nil, err
	}
	return getRequestData(req)
}

//HttpPostJsonString
/*
 * http post JSON 数据,返回一个 string 数据
 */
func HttpPostJsonString(urlStr string, postData map[string]interface{}) (string, error) {
	b, err := json.Marshal(postData)
	if err != nil {
		return "", err
	}
	body := bytes.NewBuffer(b)
	req, err := http.Post(urlStr, "application/json;charset=utf-8", body)
	if err != nil {
		return "", err
	}
	return getRequestString(req)
}

//HttpPostJsonBytes
/**
 * http post JSON 数据,返回一个 []byte 数组
 */
func HttpPostJsonBytes(urlStr string, postData []byte) ([]byte, error) {
	body := bytes.NewBuffer(postData)
	req, err := http.Post(urlStr, "application/json;charset=utf-8", body)
	if err != nil {
		return nil, err
	}
	return getRequestBytes(req)
}

//HttpGet
/*
 * http get 请求
 */
func HttpGet(urlStr string) (string, error) {
	req, err := http.Get(urlStr)
	if err != nil {
		return "", err
	}

	return getRequestString(req)
}

//HttpGetBytes
/*
 * http get 请求
 */
func HttpGetBytes(urlStr string) ([]byte, error) {
	req, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	return getRequestBytes(req)
}

func getRequestData(req *http.Response) (map[string]interface{}, error) {
	r, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body.Close()
	var res map[string]interface{}
	err = json.Unmarshal(r, &res)

	return res, err
}

func getRequestString(req *http.Response) (string, error) {
	r, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	req.Body.Close()
	return string(r), nil
}

func getRequestBytes(req *http.Response) ([]byte, error) {
	defer req.Body.Close()
	r, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	return r, nil
}
