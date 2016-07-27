package filters

import (
	"github.com/karlseguin/liquid/core"
	"encoding/json"
)

// Creates a capitalize filter
func JsonifyFactory(parameters []core.Value) core.Filter {
	return Jsonify
}

// Capitalizes words in the input sentence
func Jsonify(input interface{}, data map[string]interface{}) interface{} {
	res, err := json.Marshal(input)
	if err != nil {
	  return err.Error()
	}
	return string(res)
}