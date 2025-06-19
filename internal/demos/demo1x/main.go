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
	const topic = "emqx-demo1x-topic"

	{
		config := &mqttgo.Config{
			BrokerServer: "ws://127.0.0.1:8083/mqtt",
			Username:     "username",
			Password:     "password",
			OrderMatters: false,
		}
		onConnect := func(c mqttgo.Client, retryTimes uint64) (mqttgo.RetryType, error) {
			token := c.Subscribe(topic, 1, func(client mqttgo.Client, message mqttgo.Message) {
				zaplog.SUG.Debugln("subscribe-msg:", neatjsons.SxB(message.Payload()))
			})
			tokenState, err := mqttgo.WaitToken(token)
			if err != nil {
				return mqttgo.RetryTypeRetries, erero.Wro(err)
			}
			must.Same(tokenState, mqttgo.TokenStateSuccess)
			return mqttgo.RetryTypeSuccess, nil
		}
		client1 := rese.V1(mqttgo.NewClient(config, utils.NewUUID(), onConnect))
		defer client1.Disconnect(500)
	}

	time.Sleep(time.Millisecond * 500)

	{
		emqxHttpClient := emqxgo.NewEmqxHttpClient(&emqxgo.Config{
			BaseUrl:     "http://127.0.0.1:18083",
			ApiUsername: must.Nice(os.Getenv("EMQX_API_USERNAME")),
			ApiPassword: must.Nice(os.Getenv("EMQX_API_PASSWORD")),
		})

		type Message struct {
			Uuid string `json:"uuid"`
		}
		payload := neatjsonm.S(&Message{Uuid: utils.NewUUID()})

		var msgBatch []*emqxgo.BulkMessage
		for i := 0; i < 3; i++ {
			msgBatch = append(msgBatch, &emqxgo.BulkMessage{
				PayloadEncoding: "plain",
				Topic:           topic,
				Qos:             1,
				Payload:         payload,
				Properties:      map[string]string{},
				Retain:          false,
			})
		}
		zaplog.SUG.Debugln("message:", neatjsons.S(msgBatch))
		results := rese.A1(emqxHttpClient.PublishBulk(context.Background(), msgBatch))
		zaplog.SUG.Debugln("results:", neatjsons.S(results))
	}

	time.Sleep(time.Millisecond * 500)
}
