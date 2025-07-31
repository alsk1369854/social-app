package pkg

import (
	"bytes"
	"encoding/json"
	"sync"
)

type HTTPUtils struct{}

var httpUtilsOnce sync.Once
var httpUtils *HTTPUtils

func NewHTTPUtils() *HTTPUtils {
	httpUtilsOnce.Do(func() {
		httpUtils = &HTTPUtils{}
	})
	return httpUtils
}

func (u *HTTPUtils) ToJSONBuffer(data interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return nil, err
	}
	return buf, nil
}
