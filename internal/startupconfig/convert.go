package startupconfig

import (
	"fmt"
	"log"

	"github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

// ConvertTo convert part config to target interface by json path
// @receiver api
// @param path json path to convert
// @param to target interface to convert
// @return error
func (api *ConfigAPI) ConvertTo(path string, to interface{}) error {
	if api.configBytes == nil {
		return fmt.Errorf("startup config nil")
	}
	from := gjson.GetBytes(api.configBytes, path).Raw
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal([]byte(from), &to); err != nil {
		return err
	}
	return nil
}

// ConvertFromCustomTo convert part of custom config to target interface by relative(relate of 'custom') json path
// @receiver api
// @param relativePath relative(relate of 'custom') json path to convert
// @param to target interface to convert
// @return error
//func (api *ConfigAPI) ConvertFromCustomTo(relativePath string, to interface{}) error {
//	if api.configBytes == nil {
//		return fmt.Errorf("startup config nil")
//	}
//	from := gjson.GetBytes(api.configBytes, fmt.Sprintf("custom.%s", relativePath)).Raw
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	if err := json.Unmarshal([]byte(from), &to); err != nil {
//		return err
//	}
//	return nil
//}

// ConvertFromCustomNormalTo convert part of custom normal config to target interface by relative(relate of 'custom.normal') json path
// @receiver api
// @param relativePath relative(relate of 'custom.normal') json path to convert
// @param to target interface to convert
// @return error
func (api *ConfigAPI) ConvertFromCustomNormalTo(relativePath string, to interface{}) error {
	if api.configBytes == nil {
		return fmt.Errorf("startup config nil")
	}
	from := gjson.GetBytes(api.configBytes, fmt.Sprintf("custom.normal.%s", relativePath)).Raw
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal([]byte(from), &to); err != nil {
		return err
	}
	return nil
}

// Decrypted get decrypted string from encrypted string in config by json path
// @receiver api
// @param path json path to decrypt
// @return string decrypted string
// @return error
func (api *ConfigAPI) Decrypted(path string) (string, error) {
	if api.configBytes == nil {
		return "", fmt.Errorf("startup config nil")
	}
	encrypted := gjson.GetBytes(api.configBytes, path).String()
	decrypted, err := DecDecrypt(encrypted)
	if err != nil {
		return "", err
	}
	return decrypted, nil
}

// CustomNormalDecrypted get decrypted string from encrypted string in custom normal config by relative(relate of 'custom.normal') json path
// @receiver api
// @param relativePath relative(relate of 'custom.normal') json path to convert
// @return string
// @return error
func (api *ConfigAPI) CustomNormalDecrypted(relativePath string) (string, error) {
	return api.Decrypted(fmt.Sprintf("custom.normal.%s", relativePath))
}

// MustDecrypted get decrypted string from encrypted string in config by json path, if decrypt failed, log.Fatal will be called
// @receiver api
// @param path json path to decrypt
// @return string decrypted string
func (api *ConfigAPI) MustDecrypted(path string) string {
	decrypted, err := api.Decrypted(path)
	if err != nil {
		log.Fatalf("failed to decrypt %s: %s", path, err)
		return ""
	}
	return decrypted
}

// CustomNormalMustDecrypted get decrypted string from encrypted string in custom normal config by relative(relate of 'custom.normal') json path, if decrypt failed, log.Fatal will be called
// @receiver api
// @param relativePath relative(relate of 'custom.normal') json path to convert
// @return string decrypted string
func (api *ConfigAPI) CustomNormalMustDecrypted(relativePath string) string {
	return api.MustDecrypted(fmt.Sprintf("custom.normal.%s", relativePath))
}

// GetValue get part of config by json path
// @receiver api
// @param path json path to get
// @return gjson.Result
func (api *ConfigAPI) GetValue(path string) gjson.Result {
	if api.configBytes == nil {
		log.Fatalln("startup config nil")
	}
	return gjson.GetBytes(api.configBytes, path)
}

// GetValueFromCustomNormal get part of custom normal config by relative(relate of 'custom.normal') json path
// @receiver api
// @param relativePath relative(relate of 'custom.normal') json path to get
// @return gjson.Result
func (api *ConfigAPI) GetValueFromCustomNormal(relativePath string) gjson.Result {
	return api.GetValue(fmt.Sprintf("custom.normal.%s", relativePath))
}
