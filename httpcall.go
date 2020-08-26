package main

import (
	"io"
	"net/http"
	"strings"
	"time"
)

//http访问

func DoGet(url string, writer io.Writer, headers map[string]string) error {

	return doMethod("GET", url, "", writer, headers)
}

func DoPost(url string, content string, writer io.Writer, headers map[string]string) error {
	return doMethod("POST", url, content, writer, headers)
}

//当遇到网络错误的时候，会尝试三次，每次间隔1秒
func doMethod(method string, url string, content string, writer io.Writer, headers map[string]string) error {
	//请求客户端
	c := http.DefaultClient
	var body io.Reader
	if content == "" {
		body = nil
	} else {
		body = strings.NewReader(string(content))
	}

	req, err := http.NewRequest(method, url, body)
	//设置header
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	//请求，最多请求三次
	times := 1
DO_POST:
	resp, err := c.Do(req)
	if err != nil {
		errs(err.Error())
		if times < 3 {
			times++
			time.Sleep(time.Second)
			debug("try the ", times, " times!")
			goto DO_POST
		} else {
			return err
		}
	}
	//读取返回body内容
	defer resp.Body.Close()
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		errs(err.Error())
		return err
	}
	return nil
}
