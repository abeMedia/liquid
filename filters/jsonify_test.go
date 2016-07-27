package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestJsonifyMap(t *testing.T) {
	spec := gspec.New(t)
	filter := JsonifyFactory(nil)
	spec.Expect(filter(map[string]interface{}{"foo":"bar"}, nil).(string)).ToEqual(`{"foo":"bar"}`)
}
