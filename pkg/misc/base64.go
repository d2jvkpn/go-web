package misc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

const (
	_FilenameEncoder = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
)

var (
	_base64FilenameEncoding *base64.Encoding
)

func init() {
	_base64FilenameEncoding = base64.NewEncoding(_FilenameEncoder)
}

func Base64EncodeMap(data map[string]interface{}) string {
	bts, _ := json.Marshal(data)
	return base64.StdEncoding.EncodeToString(bts)
}

func Base64DecodeMap(str string) (data map[string]interface{}, err error) {
	var bts []byte
	if len(str) == 0 {
		return nil, fmt.Errorf("empty string")
	}

	if bts, err = base64.StdEncoding.DecodeString(str); err != nil {
		return nil, err
	}

	data = make(map[string]interface{}, 5)
	if err = json.Unmarshal(bts, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// replace +/ with -_
func Base64EncodeFilename(src string) string {
	return _base64FilenameEncoding.EncodeToString([]byte(src))
}

// replace +/ with -_
func Base64DecodeFilename(src string) ([]byte, error) {
	return _base64FilenameEncoding.DecodeString(src)
}
