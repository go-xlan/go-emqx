package emqxgo

// ConnectedEvent 当客户端上线时发出的消息
// https://www.emqx.io/docs/en/v5/mqtt/mqtt-system-topics.html#client-online-and-offline-events
// 这个结构体是从网页里的json直接转过来的
type ConnectedEvent struct {
	Username       string `json:"username"`
	Ts             string `json:"ts"`       //没有
	SockPort       string `json:"sockport"` //没有
	ProtoVer       int    `json:"proto_ver"`
	ProtoName      string `json:"proto_name"`
	Keepalive      int    `json:"keepalive"`
	Ipaddress      string `json:"ipaddress"` //没有
	ExpiryInterval int    `json:"expiry_interval"`
	ConnectedAt    int64  `json:"connected_at"`
	ConnAck        string `json:"connack"` //没有
	ClientID       string `json:"clientid"`
	CleanStart     bool   `json:"clean_start"`
}

// DisconnectedEvent 当客户端离线时发出的消息
// https://www.emqx.io/docs/en/v5/mqtt/mqtt-system-topics.html#client-online-and-offline-events
// 这个结构体是从网页里的json直接转过来的
type DisconnectedEvent struct {
	Username       string `json:"username"`
	Ts             string `json:"ts"`       //没有
	SockPort       string `json:"sockport"` //没有
	Reason         string `json:"reason"`
	ProtoVer       int    `json:"proto_ver"`
	ProtoName      string `json:"proto_name"`
	Ipaddress      string `json:"ipaddress"` //没有
	DisconnectedAt int64  `json:"disconnected_at"`
	ClientID       string `json:"clientid"`
}
