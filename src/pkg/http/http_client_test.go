package http

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	url := "https://jsonplaceholder.typicode.com/comments"
	params := map[string]any{"postId": "1"}

	resp, err := Get(url, nil, params)
	if err != nil {
		t.Fatalf("GET 请求失败: %v", err)
	}

	fmt.Println(string(resp))

}

func TestPostForm(t *testing.T) {
	url := "https://reqres.in/api/users"
	params := map[string]any{
		"name": "morpheus",
		"job":  "leader",
	}

	resp, err := PostForm(url, nil, params)
	if err != nil {
		t.Fatalf("POST FORM 请求失败: %v", err)
	}

	fmt.Println(string(resp))

}

func TestPostJson(t *testing.T) {
	url := "https://jsonplaceholder.typicode.com/posts"
	params := map[string]any{
		"title":  "foo",
		"body":   "bar",
		"userId": "1",
	}

	resp, err := PostJson(url, nil, params)
	if err != nil {
		t.Fatalf("POST JSON 请求失败: %v", err)
	}

	fmt.Println(string(resp))

}
