package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/magicnana999/im/pkg/jsonext"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	httpclient *http.Client
)

func init() {
	transport := &http.Transport{
		// 最大空闲连接数
		MaxIdleConns: 100,
		// 每个host最大空闲连接数
		MaxIdleConnsPerHost: 10,
		// 每个host最大连接数
		MaxConnsPerHost: 50,
		// 空闲连接存活时间
		IdleConnTimeout: 90 * time.Second,
		// 限制连接建立时间
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second, // 连接超时时间
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	httpclient = &http.Client{
		Transport: transport,
		Timeout:   time.Second * 30, // 请求超时时间
	}
}

func Get(urlStr string, headers map[string]string, params map[string]any) ([]byte, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	query := u.Query()
	for k, v := range params {
		query.Set(k, string(jsonext.MarshalNoErr(v)))
	}
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return execute(req)
}

func execute(request *http.Request) ([]byte, error) {

	resp, err := httpclient.Do(request)

	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return data, nil
	} else {
		return nil, errors.New("http status: " + resp.Status)
	}
}

func PostJson(
	urlStr string,
	headers map[string]string,
	params map[string]any) ([]byte, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	buf = bytes.NewBuffer(b)

	// 创建请求
	req, err := http.NewRequest(http.MethodPost, u.String(), buf)

	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	req.Header.Set("Content-Type", "application/json")

	return execute(req)
}

func PostForm(
	urlStr string,
	headers map[string]string,
	params map[string]any) ([]byte, error) {

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	payload := make(url.Values)
	for k, v := range params {
		payload.Set(k, string(jsonext.MarshalNoErr(v)))
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(payload.Encode()))

	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return execute(req)
}
