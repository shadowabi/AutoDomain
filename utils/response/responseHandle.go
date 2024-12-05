package net2

import (
	"bytes"
	"io"
	"net/http"
)

func HandleResponse(resp *http.Response) (bodyString string) {
	if resp == nil {
		return ""
	}
	bodyBuf := new(bytes.Buffer)
	_, err := io.Copy(bodyBuf, resp.Body)
	if err != nil {
		return ""
	}
	bodyString = bodyBuf.String()
	//print(bodyString)
	return bodyString
}

func GeneratePageList(pageSize int) (pageList []int) {
	startPage := 2
	pageCount := (pageSize-1)/1000 + 1 // 总共需要请求的页面数量
	pageList = make([]int, pageCount-1)
	for i := 0; i < pageCount-1; i++ {
		pageList[i] = startPage + i
	}
	return pageList
}
