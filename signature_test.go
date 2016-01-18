package AliyunOpenSearch

import (
	"testing"
)

func TestMain(t *testing.T) {
	// var p map[string]string = make(map[string]string)
	// 公共参数：Version=v2&AccessKeyId=testid&SignatureMethod=HMAC－SHA1&SignatureVersion=1.0&SignatureNonce=14053016951271226&Timestamp=2014-07-14T01:34:55Z
	// 请求参数为：query=config=format:json,start:0,hit:20&&query=default:'的'&index_name=ut_3885312&format=json&fetch_fields=title;gmt_modified
	var params map[string]string = make(map[string]string)

	params["Version"] = "v2"
	params["AccessKeyId"] = "testid"
	params["SignatureMethod"] = "HMAC-SHA1"
	params["SignatureVersion"] = "1.0"

	params["Timestamp"] = "2014-07-14T01:34:55Z"
	params["SignatureNonce"] = "14053016951271226"
	params["fetch_fields"] = "title;gmt_modified"
	params["format"] = "json"
	params["index_name"] = "ut_3885312"
	params["query"] = "config=format:json,start:0,hit:20&&query=default:'的'"
	params["Signature"] = signature(params, "GET", "testsecret")

	if params["Signature"] != "/GWWQkztlp/9Qg7rry2DuCSfKUQ=" {
		t.Error("Signature must: /GWWQkztlp/9Qg7rry2DuCSfKUQ=")
	}
}
