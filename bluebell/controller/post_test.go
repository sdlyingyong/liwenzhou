package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostHandler(t *testing.T) {
	//control的基本就是参数检查需要测试
	//自己开一个服务器
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/api/v1/post"
	r.POST(url, CreatePostHandler)

	//测试发送的参数
	body := `{
		"community_id":1,
		"title":"test",
		"content":"just a test"
	}`
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	//判断是不是返回了预期的响应
	assert.Contains(t, w.Body.String(), "需要登陆")

	//判断读取的json是否是预期的结果
	res := new(ResponseData)
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
		t.Fatalf("json.Unmarshal(w.Body.Bytes(),res) failed, err :%v \n", err)
		return
	}
	assert.Equal(t, res.Code, CodeNeedAuth)
}
