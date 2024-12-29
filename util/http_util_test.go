package util

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestGet(t *testing.T) {
	//url := `http://t.open.api.heguang.club/product-service/app/asset/detail?propertyId=318898316014391315383855&test=1212`
	url := `http://t.open.api.heguang.club/product-service/app/asset/detail`

	header := map[string]string{}
	header["appId"] = "100005-XIAORU"
	header["random"] = "21212"
	header["timestamp"] = "21212"
	header["sign"] = "c9fAOIEJ/QpTadX2rrRnx/Ao9JGtXMg/jFUpJqvjliv9Yu8rQm9w2og6gBmknQzfVN9xfdw86uoHbnGESgVQpVj93jt1BTjwHNnSqlXBvDjOyg4uSe2Qbrw2/uRGnYqrx07Xu9dpG+GsrwOb3n2u2EHwEdzUXCsNEYqtxesC+MHmVhPpmfQy0nRxoAH1uc0XXM/b7vJ8vlWCocb3QFLPbW77yM4AdZg8XiJwzPZ25RUPJ4VNLPiad771FJVcuQ0WNewBN+yb3smQVUAX0nR+NzhEoAztm+hUvGjc6sq2A/ybXJmJSdNaBCBYyW/1ioYCAfep5h1frKMbAwbvHGS/CA=="
	header["token"] = `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySWQiOjE3Njk1OTU0ODI2NTAyMzIwMTUsInprVXNlcklkIjoiMjkyMTQxOTM5MzQyNTc3NjY0NDgyOTIxIiwiYXBwSWQiOiJTVEFSU1BBQ0UtWElBT1JVIiwidXNlcm5hbWUiOiLmooHml63po54tMzMwNiIsInVzZXJQaG9uZSI6IjE4MTAwOTQzMzA2IiwiaWF0IjoxNzMwNzE3MzA1fQ.fKJEBPA_fuwa0BAmfy5gaPCr6_cZft5Zqoc6vO1RG_0`

	param := map[string]string{}
	param["propertyId"] = "318898316014391315383855"
	param["test"] = "1212"
	code, message, headers, response, e := HTTPGet(url, header, param)
	if e != nil {
		t.Error(e)
	}

	fmt.Println(code)
	fmt.Println(response)
	fmt.Println(message)
	fmt.Println(headers)

	m := map[string]any{}
	fmt.Println(string(response))
	if ee2 := json.Unmarshal([]byte(response), &m); ee2 != nil {
		log.Fatalf("Parse response failed, reason: %v \n", ee2)
	}

	fmt.Println(m)
}
