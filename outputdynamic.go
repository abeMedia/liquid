package liquid

import (
	"github.com/karlseguin/liquid/filters"
	"github.com/karlseguin/liquid/helpers"
	"reflect"
	"strings"
	"sync"
)

var (
	typeCache = make(map[reflect.Value]map[string]int)
	typeLock  sync.RWMutex
)

type OutputDynamic struct {
	Fields  []string
	Filters []filters.Filter
}

func (o *OutputDynamic) Render(data interface{}) []byte {
	for _, field := range o.Fields {
		if data = resolve(data, field); data == nil {
			return []byte("{{" + strings.Join(o.Fields, ".") + "}}")
		}
	}

	value := data
	if o.Filters != nil {
		for _, filter := range o.Filters {
			value = filter(value)
		}
	}
	return helpers.DataToBytes(value)
}

func resolve(data interface{}, field string) interface{} {
	switch typed := data.(type) {
	case map[string]string:
		return typed[field]
	case map[string]interface{}:
		return typed[field]
	case map[string]int:
		return typed[field]
	case map[string]bool:
		return typed[field]
	case map[string]float64:
		return typed[field]
	case map[string]byte:
		return typed[field]
	case map[string][]byte:
		return typed[field]
	}
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Map:
		return v.MapIndex(reflect.ValueOf(field)).Interface()
	case reflect.Struct:
		return resolveStruct(v, field)
	case reflect.Ptr:
		v = reflect.Indirect(v)
		if v.Kind() != reflect.Struct {
			return nil
		}
		return resolveStruct(v, field)
	default:
		return nil
	}
}

func resolveStruct(value reflect.Value, field string) interface{} {
	typeLock.RLock()
	typeData, exists := typeCache[value]
	typeLock.RUnlock()

	if exists == false {
		typeData = buildTypeData(value)
	}
	if index, exists := typeData[field]; exists {
		return value.Field(index).Interface()
	}
	return nil
}

func buildTypeData(value reflect.Value) map[string]int {
	t := value.Type()
	fieldCount := t.NumField()
	typeData := make(map[string]int, fieldCount)
	for i := 0; i < fieldCount; i++ {
		typeData[strings.ToLower(t.Field(i).Name)] = i
	}

	typeLock.Lock()
	defer typeLock.Unlock()
	if typeData, exists := typeCache[value]; exists {
		return typeData
	}
	typeCache[value] = typeData
	return typeData
}
