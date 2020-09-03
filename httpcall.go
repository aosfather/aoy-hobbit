package main

import (
	"io"
	"net/http"
	"strings"
	"time"
)

var _defaultClient = &http.Client{Timeout: 18 * time.Second}

//http访问

func DoGet(url string, writer io.Writer, headers map[string]string) error {

	return DoMethod("GET", url, "", writer, headers)
}

func DoPost(url string, content string, writer io.Writer, headers map[string]string) error {
	return DoMethod("POST", url, content, writer, headers)
}

func DoMethod(method string, url string, content string, writer io.Writer, headers map[string]string) error {
	for times := 0; times < 3; times++ {
		err := doMethod(method, url, content, writer, headers)
		if err == nil {
			break
		}
		return err
	}
	return nil
}

//当遇到网络错误的时候，会尝试三次，每次间隔1秒
func doMethod(method string, url string, content string, writer io.Writer, headers map[string]string) error {
	//请求客户端
	c := _defaultClient
	var body io.Reader
	if content == "" {
		body = nil
	} else {
		body = strings.NewReader(string(content))
	}

	//创建请求
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		errs(err.Error())
		return err
	}

	//设置header
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	//访问请求
	resp, err := c.Do(req)
	if err != nil {
		errs(err.Error())
		return err
	}

	//读取返回body内容
	defer resp.Body.Close()
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		errs(err.Error())
	}
	return nil
}
