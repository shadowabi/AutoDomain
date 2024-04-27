package net2

import (
	"bytes"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Error"
	"io"
	"net/http"
)

func HandleResponse(resp *http.Response) (bodyString string) {
	bodyBuf := new(bytes.Buffer)
	_, err := io.Copy(bodyBuf, resp.Body)
	Error.HandleError(err)
	bodyString = bodyBuf.String()
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