package template

// Functions available to gondola templates

import (
	"encoding/json"
	"fmt"
	"gondola/assets"
	"html/template"
	"reflect"
	"strconv"
	"strings"
)

func eq(args ...interface{}) bool {
	if len(args) == 0 {
		return false
	}
	x := args[0]
	switch x := x.(type) {
	case string, int, int64, byte, float32, float64:
		for _, y := range args[1:] {
			if x == y {
				return true
			}
		}
		return false
	}

	for _, y := range args[1:] {
		if reflect.DeepEqual(x, y) {
			return true
		}
	}
	return false
}

func neq(args ...interface{}) bool {
	return !eq(args...)
}

func _json(arg interface{}) string {
	if arg == nil {
		return ""
	}
	b, err := json.Marshal(arg)
	if err == nil {
		return string(b)
	}
	return ""
}

func nz(x interface{}) bool {
	switch x := x.(type) {
	case int, uint, int64, uint64, byte, float32, float64:
		return x != 0
	case string:
		return len(x) > 0
	}
	return false
}

func lower(x string) string {
	return strings.ToLower(x)
}

func join(x []string, sep string) string {
	s := ""
	for _, v := range x {
		s += fmt.Sprintf("%v%s", v, sep)
	}
	if len(s) > 0 {
		return s[:len(s)-len(sep)]
	}
	return ""
}

func _map(args ...interface{}) (map[string]interface{}, error) {
	var key string
	m := make(map[string]interface{})
	for ii, v := range args {
		if ii%2 == 0 {
			if s, ok := v.(string); ok {
				key = s
			} else {
				return nil, fmt.Errorf("Invalid argument to map at index %d, %t instead of string", ii, v)
			}
		} else {
			m[key] = v
		}
	}
	return m, nil
}

func mult(args ...interface{}) (float64, error) {
	val := 1.0
	for ii, v := range args {
		value := reflect.ValueOf(v)
		switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val *= float64(value.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val *= float64(value.Uint())
		case reflect.Float32, reflect.Float64:
			val *= value.Float()
		case reflect.String:
			v, err := strconv.ParseFloat(value.String(), 64)
			if err != nil {
				return 0, fmt.Errorf("Error parsing string passed to mult at index %d: %s", ii, err)
			}
			val *= v
		default:
			return 0, fmt.Errorf("Invalid argument of type %T passed to mult at index %d", v, ii)
		}
	}
	return val, nil

}

func concat(args ...interface{}) string {
	var s []string
	for _, v := range args {
		s = append(s, fmt.Sprintf("%v", v))
	}
	return strings.Join(s, "")
}

var templateFuncs template.FuncMap = template.FuncMap{
	"eq":     eq,
	"neq":    neq,
	"json":   _json,
	"nz":     nz,
	"lower":  lower,
	"join":   join,
	"map":    _map,
	"mult":   mult,
	"render": assets.Render,
	"concat": concat,
}