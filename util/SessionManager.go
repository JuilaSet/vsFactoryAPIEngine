package util

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	Key = []byte("UDqeCwDO!@kKuVj@9Xbd*HOV#SwE!t1q&")
)

type SessionManager struct {
	store *sessions.CookieStore
}

var sessionManager *SessionManager

func init() {
	sessionManager = &SessionManager{store: sessions.NewCookieStore(Key)}
}

// 单例对象
func GetSessionManager() *SessionManager {
	return sessionManager
}

func (c *SessionManager) SaveSession(ctx *gin.Context, sessionName string, k string, v string) {
	w, r := ctx.Writer, ctx.Request
	// 获取一个session对象, session-name是session的名字
	session, err := c.store.Get(r, sessionName)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
	// 在session中存储值
	session.Values[k] = v
	// 保存更改
	if err := session.Save(r, w); err != nil {
		panic(err)
	}
}

func (c *SessionManager) GetSession(ctx *gin.Context, sessionName string, k string) interface{} {
	session, err := c.store.Get(ctx.Request, sessionName)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
	return session.Values[k]
}

func (c *SessionManager) DelSession(ctx *gin.Context, sessionName string) {
	w, r := ctx.Writer, ctx.Request
	session, err := c.store.Get(r, sessionName)
	if err != nil {
		panic(err)
	}
	// 将session的最大存储时间设置为小于零的数即为删除
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		panic(err)
	}
}
