package config

type appStruct struct {
	appId   string
	appName string
}

var currentAppInstance = &appStruct{}

// SetApp 设置AppId
func SetApp(appId string, appName string) {
	if appId != "" {
		currentAppInstance.appId = appId
	}
	if appName != "" {
		currentAppInstance.appName = appName
	}
}

// GetAppId 获取AppId
func GetAppId() string {
	return currentAppInstance.appId
}

// GetAppName 获取AppName
func GetAppName() string {
	return currentAppInstance.appName
}
