// 阿里云开放搜索 golang sdk
// Doc:https://help.aliyun.com/document_detail/opensearch/api-reference/api-interface/search-related.html?spm=5176.docopensearch/api-reference/terminology.2.6.8gxPsl

package AliyunOpenSearch

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"
	"time"

	// "github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
)

type OpenSearchClient struct {
	baseUrl     string
	acId        string
	secret      string
	version     string
	signMethod  string
	signVersion string
}

// baseUrl: like: "http://opensearch.aliyun.com"
func NewOpenSearchClient(baseUrl, acId, secret string) *OpenSearchClient {
	cli := new(OpenSearchClient)
	cli.version = "v2"
	cli.signMethod = "HMAC-SHA1"
	cli.signVersion = "1.0"

	cli.baseUrl = baseUrl
	cli.acId = acId
	cli.secret = secret
	return cli
}

func (self *OpenSearchClient) Search(index_name, key_word string, page int, per_page int, sort string) (json string) {
	key_word = strings.TrimSpace(key_word)
	if key_word == "" {
		log.Printf("search keyword  can not be empty string.")
		return ""
	}

	mehtod, path := "GET", "/search"
	var p map[string]string = make(map[string]string)
	//query子句，必须有
	p["query"] = fmt.Sprintf("query=%s", key_word)
	//config子句，主要分页用。format:xml||json||fulljson
	p["query"] = fmt.Sprintf("%s&&config=start:%d,hit:%d,format:%s", p["query"], (page-1)*per_page, per_page, "json")
	if sort != "" {
		arrSort := strings.Split(sort, ",")
		var sliceSort []string = arrSort[:]
		for i := 0; i < len(arrSort); i++ {
			if strings.HasPrefix(arrSort[i], "-") {
				sliceSort[i] = arrSort[i]
			} else {
				sliceSort[i] = fmt.Sprintf("+%s", arrSort[i])
			}
		}
		//排序子句格式为：+field1;-field2
		p["query"] = fmt.Sprintf("%s&&sort=%s", p["query"], strings.Join(sliceSort, ";"))
	}
	//filter子句还没实现，现在没用到。doc：https://help.aliyun.com/document_detail/opensearch/api-reference/query-clause/filter-clause.html?spm=5176.docopensearch/api-reference/api-interface/search-related.2.4.V4LQ6c
	p["index_name"] = index_name

	return self.call(mehtod, path, p)
}

func percentEncode(s string) string {
	// return quote(str(string)).replace('+', '%20').replace('*', '%%2A').replace('%%7E', '~')
	s = url.QueryEscape(s)
	s = strings.Replace(s, "+", "%20", -1)
	s = strings.Replace(s, "*", "%2A", -1)
	s = strings.Replace(s, "%7E", "~", -1)
	return s
}
func base64sha1(message string, key string) string {

	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(message))
	hash := h.Sum(nil)

	base64signature := base64.StdEncoding.EncodeToString(hash)

	return base64signature
}

func signature(params map[string]string, method, secret string) (sign string) {
	//Doc: https://help.aliyun.com/document_detail/opensearch/api-reference/call-method/signature.html?spm=5176.docopensearch/api-reference/call-method/common-params.2.1.N0qFtk
	var slice sort.StringSlice
	for k, _ := range params {
		slice = append(slice, k)
	}
	slice.Sort()

	var a []string = make([]string, 0)
	for _, k := range slice {
		a = append(a, percentEncode(k)+"="+percentEncode(params[k]))
	}

	query := strings.Join(a, "&")
	base := method + "&%2F&" + percentEncode(query)
	return base64sha1(base, secret+"&")
}

func (self *OpenSearchClient) call(method string, path string, params map[string]string) string {
	if method != "GET" {
		log.Printf("aliyun OpenSearchClient just accessed method of GET or POST. now=[%s]", method)
		return ""
	}

	params["Version"] = self.version
	params["AccessKeyId"] = self.acId
	params["SignatureMethod"] = self.signMethod
	params["SignatureVersion"] = self.signVersion

	now := time.Now()
	params["Timestamp"] = now.UTC().Format("2006-01-02T15:04:05Z")
	params["SignatureNonce"] = fmt.Sprintf("%d", now.UnixNano())
	params["Signature"] = signature(params, method, self.secret)

	connTimeout := 2 * time.Second
	readTimeout := 10 * time.Second

	var s string
	var err error
	if method == "GET" {
		val := url.Values{}
		for k, v := range params {
			val.Add(k, v)
		}
		log.Printf("Called URL: ", self.baseUrl+path+"?"+val.Encode())
		s, err = httplib.Get(self.baseUrl+path+"?"+val.Encode()).Debug(true).SetTimeout(connTimeout, readTimeout).String()
	} else {
		req := httplib.Post(self.baseUrl+path).SetTimeout(connTimeout, readTimeout)
		for k, v := range params {
			req.Param(k, v)
		}

		s, err = req.String()
	}
	if err != nil {
		log.Printf("OpenSearchClient [%s] [%s] [%s] err=[%s]",
			method, self.baseUrl, path, err)
	}
	return s
}
