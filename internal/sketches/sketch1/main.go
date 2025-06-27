package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xlan/go-emqx/emqxgo"
	"github.com/yyle88/done"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/reggin/warpginhandle"
	"github.com/yyle88/zaplog"
)

func main() {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), gin.BasicAuth(gin.Accounts{
		"username-abc": "password-123",
	}))
	// 根据这里配置上线（连接完成）的请求
	// https://docs.emqx.com/en/emqx/latest/data-integration/rule-sql-events-and-fields.html#connection-complete-event-events-client-connected
	router.POST("/mac-demo-endpoint/client-online", warpginhandle.PX(handleConnectedEvent, makeResp[resType]))

	// 根据这里配置离线（断开连接）的请求
	// https://docs.emqx.com/en/emqx/latest/data-integration/rule-sql-events-and-fields.html#disconnect-event-events-client-disconnected
	router.POST("/mac-demo-endpoint/client-offline", warpginhandle.PX(handleDisconnectedEvent, makeResp[resType]))

	// 根据这个网页配置客户端的 webhook 鉴权
	// https://docs.emqx.com/en/emqx/latest/access-control/authz/http.html
	router.POST("/client-connect-authentication", warpginhandle.RX(handleAuthentication, wrongResp))
	done.Done(router.Run(":8080"))
}

func handleConnectedEvent(arg *emqxgo.ConnectedEvent) (*resType, error) {
	zaplog.SUG.Debugln("connected-event:", neatjsons.S(arg))
	return &resType{}, nil
}

func handleDisconnectedEvent(arg *emqxgo.DisconnectedEvent) (*resType, error) {
	zaplog.SUG.Debugln("disconnected-event:", neatjsons.S(arg))
	return &resType{}, nil
}

type resType struct{}

type respType struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Data any    `json:"data"`
}

func makeResp[RES any](ctx *gin.Context, res *RES, cause error) *respType {
	if cause != nil {
		return wrongResp(ctx, cause)
	} else {
		return okResp(ctx, res)
	}
}

func okResp[RES any](ctx *gin.Context, res *RES) *respType {
	return &respType{
		Code: 0,
		Desc: "SUCCESS",
		Data: res,
	}
}

func wrongResp(ctx *gin.Context, cause error) *respType {
	return &respType{
		Code: -1,
		Desc: cause.Error(),
		Data: nil,
	}
}

type AuthReq struct {
	ClientID string `json:"clientid"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthRes struct {
	Result string `json:"result"`
}

func handleAuthentication(arg *AuthReq) (*AuthRes, error) {
	zaplog.SUG.Debugln("client-auth:", neatjsons.S(arg))
	return &AuthRes{
		Result: "allow", // 是否允许连接 "result": "allow" | "deny" | "ignore"
	}, nil
}
