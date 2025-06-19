package emqxgo

import (
	"context"
	"encoding/json"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-xlan/go-emqx/internal/utils"
	"github.com/go-xlan/go-mqtt/mqttgo"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsonm"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"github.com/yyle88/rese/resb"
	"github.com/yyle88/syncmap"
	"github.com/yyle88/zaplog"
)

const caseTopic = "emqx-unit-test-case"

var caseClient *EmqxHttpClient
var caseMapCnt = syncmap.New[string, *int64]()

func TestMain(m *testing.M) {
	//8081 端口已被合并至 18083 端口，
	//API 访问基础路径由 /api/v4 切换到 /api/v5，
	//请通过此端口和路径调用 API；
	//时间相关的数据将使用 RFC3339 (opens new window)格式。
	//在这里创建密钥这样才能访问
	//http://127.0.0.1:18083/#/APIKey

	caseClient = NewEmqxHttpClient(&Config{
		BaseUrl:     "http://127.0.0.1:18083",
		ApiUsername: must.Nice(os.Getenv("EMQX_API_USERNAME")),
		ApiPassword: must.Nice(os.Getenv("EMQX_API_PASSWORD")),
	})

	{
		config := &mqttgo.Config{
			BrokerServer: "ws://127.0.0.1:8083/mqtt",
			Username:     "username",
			Password:     "password",
			OrderMatters: false,
		}
		onConnect := func(c mqttgo.Client, retryTimes uint64) (mqttgo.RetryType, error) {
			token := c.Subscribe(caseTopic, 1, func(client mqttgo.Client, message mqttgo.Message) {
				zaplog.SUG.Debugln("subscribe-msg:", neatjsons.SxB(message.Payload()))

				var msg CaseMessage
				must.Done(json.Unmarshal(message.Payload(), &msg))
				must.Nice(msg.Uuid)

				value, _ := caseMapCnt.LoadOrStore(msg.Uuid, utils.GetValuePointer(int64(0)))
				atomic.AddInt64(value, 1)
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

	m.Run()

	time.Sleep(time.Millisecond * 500)
}

type CaseMessage struct {
	Uuid string `json:"uuid"`
}

func TestEmqxHttpClient_PublishBulk(t *testing.T) {
	msg := &CaseMessage{Uuid: utils.NewUUID()}
	payload := neatjsonm.S(msg)

	var msgBatch []*BulkMessage
	for i := 0; i < 3; i++ {
		msgBatch = append(msgBatch, &BulkMessage{
			PayloadEncoding: "plain",
			Topic:           caseTopic,
			Qos:             1,
			Payload:         payload,
			Properties:      map[string]string{},
			Retain:          false,
		})
	}
	t.Log(neatjsons.S(msgBatch))

	results := rese.A1(caseClient.PublishBulk(context.Background(), msgBatch))
	t.Log(neatjsons.S(results))

	time.Sleep(time.Millisecond * 500)

	count := utils.GetPointerValue(resb.P1(caseMapCnt.LoadAndDelete(msg.Uuid)))
	require.Equal(t, int(count), len(msgBatch))
}
