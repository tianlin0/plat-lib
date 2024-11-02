package startupconfig

// Trace tracing配置抽象
type Trace interface {
	// TService service
	TService() string
	// TTenantID TenantID
	TTenantID() string
	// TAddress Address
	TAddress() string
	// THTTPPort HTTPPort
	THTTPPort() string
	// TGRPCPort GRPCPort
	TGRPCPort() string
	// TSampleRatio SampleRatio
	TSampleRatio() float64
}

// TracingConfig 配置
type TracingConfig struct {
	ServiceName string  `json:"service" yaml:"service"`
	TenantID    string  `json:"tenantId" yaml:"tenantId"`
	Address     string  `json:"address" yaml:"address"`
	HTTPPort    string  `json:"httpPort" yaml:"httpPort"`
	GRPCPort    string  `json:"grpcPort" yaml:"grpcPort"`
	SampleRatio float64 `json:"sampleRatio" yaml:"sampleRatio"`
}

// TService service
func (c *TracingConfig) TService() string {
	if c == nil {
		return ""
	}
	return c.ServiceName
}

// TTenantID TenantID
func (c *TracingConfig) TTenantID() string {
	if c == nil {
		return ""
	}
	return c.TenantID
}

// TAddress Address
func (c *TracingConfig) TAddress() string {
	if c == nil {
		return ""
	}
	return c.Address
}

// THTTPPort HTTPPort
func (c *TracingConfig) THTTPPort() string {
	if c == nil {
		return ""
	}
	return c.HTTPPort
}

// TGRPCPort GRPCPort
func (c *TracingConfig) TGRPCPort() string {
	if c == nil {
		return ""
	}
	return c.GRPCPort
}

// TSampleRatio SampleRatio
func (c *TracingConfig) TSampleRatio() float64 {
	if c == nil {
		return 0
	}
	return c.SampleRatio
}
