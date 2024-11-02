package startupconfig

// Recorder 操作记录配置抽象
type Recorder interface {
	// MQConfigKey 消息队列配置Key
	MQConfigKey() string
	// OperationGroup 操作分组
	OperationGroup(string) *OperationGroup
	// OperationGroups 操作分组列表
	OperationGroups() map[string]*OperationGroup
	// OperationResourceNameMap 操作资源中英文对照表
	OperationResourceNameMap() map[string]string
	// Close 关闭操作记录
	Close() bool
}

var (
	userOperationTopicKey = "userOperation"
)

// UserOperationRecorder  用户操作记录配置
type UserOperationRecorder struct {
	RecordClose     bool                       `json:"close" yaml:"close"`
	MQConfKey       string                     `json:"mqConfigKey" yaml:"mqConfigKey"`
	Groups          map[string]*OperationGroup `json:"groups" yaml:"groups"`
	ResourceNameMap map[string]string          `json:"resourceNameMap" yaml:"resourceNameMap"`
}

// OperationGroup 操作分组
type OperationGroup struct {
	GroupId string            `json:"groupId" yaml:"groupId"`
	Modules map[string]string `json:"modules" yaml:"modules"`
}

// Module 获取分组group 的 指定 module
// @receiver g
// @param key
// @return string
func (g *OperationGroup) Module(key string) string {
	if g == nil || g.Modules == nil {
		return ""
	}
	return g.Modules[key]
}

// Close 关闭操作记录
// @receiver r
// @return bool
func (r *UserOperationRecorder) Close() bool {
	if r != nil {
		return r.RecordClose
	}
	return false
}

// MQConfigKey 消息队列配置Key
// @receiver r
// @return string
func (r *UserOperationRecorder) MQConfigKey() string {
	if r == nil {
		return ""
	}
	return r.MQConfKey
}

// OperationGroup 操作分组列表
// @receiver r
// @return map[string]OperationGroup
func (r *UserOperationRecorder) OperationGroup(key string) *OperationGroup {
	if r == nil || r.Groups == nil {
		return nil
	}
	return r.Groups[key]
}

// OperationGroups 操作分组列表
// @receiver r
// @return map[string]OperationGroup
func (r *UserOperationRecorder) OperationGroups() map[string]*OperationGroup {
	if r == nil {
		return nil
	}
	return r.Groups
}

// OperationResourceNameMap 操作资源中英文对照表
// @receiver r
// @return map[string]string
func (r *UserOperationRecorder) OperationResourceNameMap() map[string]string {
	if r == nil || r.ResourceNameMap == nil {
		return nil
	}
	return r.ResourceNameMap
}
