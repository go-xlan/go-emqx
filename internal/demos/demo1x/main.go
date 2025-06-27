package main

import (
	"context"
	"os"
	"time"

	"github.com/go-xlan/go-emqx/emqxgo"
	"github.com/go-xlan/go-emqx/internal/utils"
	"github.com/go-xlan/go-mqtt/mqttgo"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsonm"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
)

func main() {
	const mqttTopic = "emqx-go-demo1x-topic"
	const clientIDX = "emqx-go-demo1x-client-"

	{
		config := &mqttgo.Config{
			BrokerServer: "ws://127.0.0.1:8083/mqtt",
			Username:     "username",
			Password:     "password",
			OrderMatters: false,
		}
		clientID := clientIDX + utils.NewUUID()
		client1 := rese.V1(mqttgo.NewClientWithCallback(config, clientID, mqttgo.NewCallback().
			OnConnect(func(c mqttgo.Client, retryTimes uint64) (mqttgo.CallbackState, error) {
				token := c.Subscribe(mqttTopic, 1, func(client mqttgo.Client, message mqttgo.Message) {
					zaplog.SUG.Debugln("subscribe-msg:", neatjsons.SxB(message.Payload()))
				})
				tokenState, err := mqttgo.WaitToken(token)
				if err != nil {
					return mqttgo.CallbackRetries, erero.Wro(err)
				}
				must.Same(tokenState, mqttgo.TokenStateSuccess)
				return mqttgo.CallbackSuccess, nil
			}),
		))
		defer client1.Disconnect(500)
	}

	time.Sleep(time.Millisecond * 500)

	{
		emqxHttpClient := emqxgo.NewEmqxHttpClient(&emqxgo.Config{
			BaseUrl:     "http://127.0.0.1:18083",
			ApiUsername: must.Nice(os.Getenv("EMQX_API_USERNAME")),
			ApiPassword: must.Nice(os.Getenv("EMQX_API_PASSWORD")),
		})

		type messageType struct {
			Uuid string `json:"uuid"`
		}
		payload := neatjsonm.S(&messageType{Uuid: utils.NewUUID()})

		for i := 0; i < 3; i++ {
			onceMessage := &emqxgo.PublishMessage{
				Payload:         payload,
				PayloadEncoding: "plain",
				Properties:      map[string]any{},
				Qos:             1,
				Retain:          false,
				Topic:           mqttTopic,
			}
			zaplog.SUG.Debugln("message:", neatjsons.S(onceMessage))
			onceResult := rese.P1(emqxHttpClient.Publish(context.Background(), onceMessage))
			zaplog.SUG.Debugln("results:", neatjsons.S(onceResult))
		}
	}

	time.Sleep(time.Millisecond * 500)
}
