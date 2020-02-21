package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 开启跨域资源共享
func Cors(f func(*gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		w, r := c.Writer, c.Request
		// 允许访问所有域, 可以换成具体url, 注意仅具体url才能带cookie信息
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// header的类型
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		// 设置为true, 允许ajax异步请求带cookie信息
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		// 允许请求方法
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		// 返回数据格式是json
		w.Header().Set("content-type", "application/json;charset=UTF-8")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		f(c)
	}
}
