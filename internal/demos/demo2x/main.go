package main

import (
	"context"
	"encoding/json"
	"os"
	"sync/atomic"
	"time"

	"github.com/go-xlan/go-emqx/emqxgo"
	"github.com/go-xlan/go-emqx/internal/utils"
	"github.com/go-xlan/go-mqtt/mqttgo"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsonm"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"github.com/yyle88/rese/resb"
	"github.com/yyle88/syncmap"
	"github.com/yyle88/zaplog"
)

const mqttTopic = "emqx-go-demo2x-topic"
const clientIDX = "emqx-go-demo2x-client-"

var msgCountMap = syncmap.New[string, *int64]()

func main() {
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

					var msg messageType
					must.Done(json.Unmarshal(message.Payload(), &msg))
					must.Nice(msg.Uuid)

					value, _ := msgCountMap.LoadOrStore(msg.Uuid, utils.GetValuePointer(int64(0)))
					atomic.AddInt64(value, 1)
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

	publishBulk()
	publishBulk()
	publishBulk()

	time.Sleep(time.Millisecond * 500)
}

type messageType struct {
	Uuid string `json:"uuid"`
}

func publishBulk() {
	emqxHttpClient := emqxgo.NewEmqxHttpClient(&emqxgo.Config{
		BaseUrl:     "http://127.0.0.1:18083",
		ApiUsername: must.Nice(os.Getenv("EMQX_API_USERNAME")),
		ApiPassword: must.Nice(os.Getenv("EMQX_API_PASSWORD")),
	})

	msg := &messageType{Uuid: utils.NewUUID()}
	payload := neatjsonm.S(msg)

	var bulkMessages []*emqxgo.PublishBulkMessage
	for i := 0; i < 3; i++ {
		bulkMessages = append(bulkMessages, &emqxgo.PublishBulkMessage{
			Payload:         payload,
			PayloadEncoding: "plain",
			Properties:      map[string]any{},
			Qos:             1,
			Retain:          false,
			Topic:           mqttTopic,
		})
	}
	zaplog.SUG.Debugln(neatjsons.S(bulkMessages))

	bulkResults := rese.A1(emqxHttpClient.PublishBulk(context.Background(), bulkMessages))
	zaplog.SUG.Debugln(neatjsons.S(bulkResults))

	time.Sleep(time.Millisecond * 500)

	count := utils.GetPointerValue(resb.P1(msgCountMap.LoadAndDelete(msg.Uuid)))
	must.Length(bulkMessages, int(count))
}
