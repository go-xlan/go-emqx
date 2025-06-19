package emqxgo

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	restyv2 "github.com/go-resty/resty/v2"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
)

// EmqxHttpClient emqx_dashboard  (opens new window)是 EMQ X 消息服务器的 Web 管理控制台, 该插件默认开启。
// 当 EMQX 启动成功后，可访问 http://localhost:18083 进行查看，
// 默认用户名/密码: admin/public。
type EmqxHttpClient struct {
	client *restyv2.Client
}

type Config struct {
	BaseUrl     string
	ApiUsername string
	ApiPassword string
}

func NewEmqxHttpClient(config *Config) *EmqxHttpClient {
	client := restyv2.New()
	client.SetBaseURL(must.Nice(config.BaseUrl))
	// 虽然会报出这样的信息
	// WARN RESTY Using Basic Auth in HTTP mode is not secure, use HTTPS
	// 但是忽略它
	client.SetBasicAuth(config.ApiUsername, config.ApiPassword)
	client.SetRetryCount(2)
	client.SetRetryWaitTime(5 * time.Second)
	client.AddRetryCondition(func(response *restyv2.Response, err error) bool {
		return response.StatusCode() >= 500
	})
	return &EmqxHttpClient{
		client: client,
	}
}

type BulkMessage struct {
	PayloadEncoding string            `json:"payload_encoding"`
	Topic           string            `json:"topic"`
	Qos             int               `json:"qos"`
	Payload         string            `json:"payload"`
	Properties      map[string]string `json:"properties"`
	Retain          bool              `json:"retain"`
}

type BulkResult struct {
	ID         string `json:"id"`          //200时有值
	Message    string `json:"message"`     //202时有值
	ReasonCode int    `json:"reason_code"` //202时有值
}

func (m *EmqxHttpClient) PublishBulk(ctx context.Context, msgBatch []*BulkMessage) ([]*BulkResult, error) {
	response, err := m.client.R().SetContext(ctx).
		SetBody(msgBatch).
		SetHeader("Content-type", "application/json").
		Post("/api/v5/publish/bulk")
	if err != nil {
		return nil, erero.Wro(err)
	}
	//根据官方文档 https://www.emqx.io/docs/en/v5.0/admin/api-docs.html#tag/Publish
	//能够返回的正确码是 200 和 202
	switch response.StatusCode() {
	case http.StatusOK, http.StatusAccepted:
		var results []*BulkResult
		if err := json.NewDecoder(bytes.NewBuffer(response.Body())).Decode(&results); err != nil {
			return nil, erero.Wro(err)
		}
		var causes []error
		switch response.StatusCode() {
		case http.StatusAccepted:
			for idx, result := range results {
				if result.ReasonCode != 16 {
					causes = append(causes, erero.Errorf("wrong idx=%v reason_code=%v message=%v", idx, result.ReasonCode, result.Message))
				}
			}
		}
		if len(causes) != 0 {
			return nil, erero.Joins(causes)
		}
		return results, nil
	default:
		return nil, erero.Errorf("status %d != (200|202) status: %s", response.StatusCode(), response.Status())
	}
}
