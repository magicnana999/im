package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func HTTPGet(urlStr string, headers, params map[string]string) (int, string, map[string][]string, string, error) {
	// 创建URL
	u, err := url.Parse(urlStr)
	if err != nil {
		return 0, "", nil, "", err
	}

	// 添加查询参数
	query := u.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	u.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)

	if err != nil {
		return 0, "", nil, "", err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return execute(req)
}

func execute(request *http.Request) (int, string, map[string][]string, string, error) {

	resp, err := http.DefaultClient.Do(request)
	defer resp.Body.Close()

	if err != nil {
		return 0, "", nil, "", err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", nil, "", err
	}

	return resp.StatusCode, resp.Status, resp.Header, string(data), nil
}

func HTTPPostJson(
	urlStr string,
	headers map[string]string,
	params map[string]string) (int, string, map[string][]string, string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return 0, "", nil, "", err
	}

	buf := new(bytes.Buffer)
	b, err := json.Marshal(params)
	if err != nil {
		return 0, "", nil, "", err
	}
	buf = bytes.NewBuffer(b)

	// 创建请求
	req, err := http.NewRequest(http.MethodPost, u.String(), buf)

	if err != nil {
		return 0, "", nil, "", err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	req.Header.Set("Content-Type", "application/json")

	return execute(req)
}

func HTTPPostForm(
	urlStr string,
	headers map[string]string,
	params map[string]string) (int, string, map[string][]string, string, error) {
	// 创建URL
	u, err := url.Parse(urlStr)
	if err != nil {
		return 0, "", nil, "", err
	}

	payload := make(url.Values)
	for k, v := range params {
		payload.Set(k, v)
	}

	// 创建请求
	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(payload.Encode()))

	if err != nil {
		return 0, "", nil, "", err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return execute(req)
}
