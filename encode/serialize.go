package encode

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
)

// Serialize 序列化一个对象
func Serialize(s interface{}) (string, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	err := enc.Encode(s)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// UnSerialize 反序列化一个对象
func UnSerialize(s string, data interface{}) error {
	var b bytes.Buffer
	old, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	_, err = b.Write(old)
	if err != nil {
		return err
	}
	enc := gob.NewDecoder(&b)
	err = enc.Decode(data)
	if err != nil {
		return err
	}
	return nil
}
