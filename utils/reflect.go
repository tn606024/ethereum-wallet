package utils

import (
	"fmt"
	"reflect"
	"strconv"
)

func ExtractValue(input interface{}) (output string) {
	v := direct(reflect.ValueOf(input))
	if v.Kind() == reflect.Slice {
		length := v.Len()
		for i := 0; i < length; i++ {
			if i == 0 {
				output = valueToStr(v.Index(i))
			} else {
				output += "," + valueToStr(v.Index(i))
			}
		}
	} else {
		output = valueToStr(v)
	}

	return
}

func valueToStr(v reflect.Value) (str string) {
	switch v.Kind() {
	case reflect.String:
		str = v.String()
	case reflect.Int:
		str = strconv.FormatInt(v.Int(), 10)
	default:
		panic(fmt.Sprintf("valueToStr: %v is of unexpected kind %q", v, v.Kind()))
	}
	return
}


func direct(v reflect.Value) reflect.Value {
	for ; v.Kind() == reflect.Ptr; v = v.Elem() {
		// relax
	}
	return v
}

