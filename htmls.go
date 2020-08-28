package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
)

//节点属性 to array
func NodeChildAttrToArray(node *goquery.Selection, attr string) []string {
	var urls []string
	var href string
	var exists bool
	node = node.Children()
	node.Each(func(index int, s *goquery.Selection) {
		href, exists = s.Attr(attr)
		if exists {
			urls = append(urls, href)
		}
	})
	return urls
}

func NodesChildrenAttrToArray(node *goquery.Selection, attr string) []string {
	var urls []string
	var href string
	var exists bool
	node.Each(func(index int, s *goquery.Selection) {
		href, exists = s.Children().Attr(attr)
		if exists {
			urls = append(urls, href)
		}
	})
	return urls
}

//节点 获取文本内容
func NodeText(node *goquery.Selection) string {
	if node != nil {
		return node.Text()
	}
	return ""
}

//根据path查找节点，类似："div.bookname h1" "div#content"
func DocFind(doc *goquery.Document, path string) *goquery.Selection {
	if doc != nil {
		return doc.Find(path)
	}
	return nil
}

//处理Html
func GetDocument(url string) *goquery.Document {
	var root *html.Node
	//判断url是否有http请求头
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		root, _ = GetUrl(url)
	} else { //作为文件处理
		file, _ := os.Open(url)
		root, _ = html.Parse(file)
	}

	return goquery.NewDocumentFromNode(root)
}

//获取url的内容
func GetUrl(url string) (*html.Node, error) {
	res, e := http.Get(url)
	if e != nil {
		return nil, e
	}

	defer res.Body.Close()
	if res.Request == nil {
		return nil, fmt.Errorf("Response.Request is nil")
	}

	// Parse the HTML into nodes
	root, e := html.Parse(res.Body)
	if e != nil {
		return nil, e
	}
	return root, nil
}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}
