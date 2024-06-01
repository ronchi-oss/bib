package conf

import (
	"bytes"
	"gopkg.in/yaml.v3"
)

func YAMLEncode(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	e := yaml.NewEncoder(&b)
	defer e.Close()
	e.SetIndent(2)
	if err := e.Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
