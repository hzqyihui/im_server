package util

import (
	"errors"
	"fmt"
	"reflect"
)

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

func ArrayColumn(array interface{}, key string) (result []interface{}, err error) {
	result = []interface{}{}
	t := reflect.TypeOf(array)
	v := reflect.ValueOf(array)
	if t.Kind() != reflect.Slice {
		return nil, errors.New("array type not slice")
	}
	//if v.Len() == 0 {
	//	return nil, errors.New("array len is zero")
	//}

	for i := 0; i < v.Len(); i++ {

		indexv := v.Index(i)
		if v.Index(i).Type().Kind() == reflect.Ptr {
			indexv = v.Index(i).Elem()
		}

		fmt.Println(indexv, indexv.Type().Kind())
		if indexv.Type().Kind() != reflect.Struct {
			return nil, errors.New("element type not struct")
		}
		mapKeyInterface := indexv.FieldByName(key)

		if mapKeyInterface.Kind() == reflect.Invalid {
			return nil, errors.New("key not exist")
		}
		result = append(result, mapKeyInterface)
	}
	return result, err
}

func interfaceToString(v interface{}) (result string, err error) {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Int64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		result = fmt.Sprintf("%v", v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result = fmt.Sprintf("%v", v)
	case reflect.String:
		result = v.(string)
	default:
		err = errors.New("can't transition to string")
	}
	return result, err

}
