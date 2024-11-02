package startupconfig

const (
	// SubscriptionPositionLatest is the latest position which means the start consuming position
	// will be the last message
	SubscriptionPositionLatest int = iota

	// SubscriptionPositionEarliest is the earliest position which means the start consuming position
	// will be the first message
	SubscriptionPositionEarliest
)

// TDMQ 消息队列TDMQ连接参数的抽象
type TDMQ interface {
	// URL 返回创建client所需要的目标pulsar地址URL
	URL() string
	// AuthenticationToken 返回作为消费者时的token
	AuthenticationToken() string
	// SubscriptionInitialPosition 消费的初始位置
	SubscriptionInitialPosition() int
	// SubscriptionName 消费的名称
	SubscriptionName() string
	// Topic 订阅主题
	Topic(name string) string
	// ListenerName 返回创建client所需要的ListenerName字段信息
	// ListenerName() string
}

// TdmqConfig tdmq配置
type TdmqConfig struct {
	BrokerAddress           string            `json:"brokerAddr" yaml:"brokerAddr"`
	JwtTokenDecrypted       Decrypted         `json:"jwtToken" yaml:"jwtToken"`
	InitialPosition         string            `json:"initialPosition" yaml:"initialPosition"`
	ConsumeSubscriptionName string            `json:"subscriptionName" yaml:"subscriptionName"`
	Topics                  map[string]string `json:"topics" yaml:"topics"`
}

// URL 返回创建client所需要的目标pulsar地址URL
func (c *TdmqConfig) URL() string {
	if c == nil {
		return ""
	}
	return c.BrokerAddress
}

// AuthenticationToken 返回作为消费者时的token
func (c *TdmqConfig) AuthenticationToken() string {
	if c == nil {
		return ""
	}
	return c.JwtTokenDecrypted.String()
}

// SubscriptionInitialPosition 消费的初始位置
func (c *TdmqConfig) SubscriptionInitialPosition() int {
	if c == nil {
		return 0
	}
	switch c.InitialPosition {
	case "earliest":
		return SubscriptionPositionEarliest
	case "lasted":
		return SubscriptionPositionLatest
	default:
		return SubscriptionPositionLatest
	}
}

// SubscriptionName 消费的名称
func (c *TdmqConfig) SubscriptionName() string {
	if c == nil {
		return ""
	}
	if c.ConsumeSubscriptionName != "" {
		return c.ConsumeSubscriptionName
	}
	return ""
}

// Topic 订阅主题
func (c *TdmqConfig) Topic(name string) string {
	if c == nil {
		return ""
	}
	if c.Topics != nil {
		return c.Topics[name]
	}
	return ""
}
