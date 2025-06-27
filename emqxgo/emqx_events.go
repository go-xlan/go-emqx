package emqxgo

// ConnectedEvent 当客户端上线时发出的消息
// https://docs.emqx.com/en/emqx/latest/data-integration/rule-sql-events-and-fields.html#mqtt-events
// https://docs.emqx.com/en/emqx/latest/data-integration/rule-sql-events-and-fields.html#connection-complete-event-events-client-connected
// 注意不要和这个混淆，这个需要按客户端去订阅
// https://docs.emqx.com/en/emqx/latest/observability/mqtt-system-topics.html#client-online-and-offline-events
type ConnectedEvent struct {
	ClientID       string `json:"clientid"`
	Username       string `json:"username"`
	ProtoName      string `json:"proto_name"`
	ProtoVer       int    `json:"proto_ver"`
	Keepalive      int    `json:"keepalive"`
	CleanStart     bool   `json:"clean_start"`
	ExpiryInterval int    `json:"expiry_interval"`
	ConnectedAt    int64  `json:"connected_at"`
}

// DisconnectedEvent 当客户端离线时发出的消息
// https://docs.emqx.com/en/emqx/latest/data-integration/rule-sql-events-and-fields.html#mqtt-events
// https://docs.emqx.com/en/emqx/latest/data-integration/rule-sql-events-and-fields.html#disconnect-event-events-client-disconnected
type DisconnectedEvent struct {
	Reason         string `json:"reason"`
	ClientID       string `json:"clientid"`
	Username       string `json:"username"`
	DisconnectedAt int64  `json:"disconnected_at"`
}
