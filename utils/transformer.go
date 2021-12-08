package utils

import (
	"github.com/huypher/kit/container"
	"reflect"
	"strconv"
)

const (
	structToMapTag      = "keyname"
	ommitEmptyTagOption = "omitempty"
)

func IntsToStrings(arr []int) []string {
	res := make([]string, len(arr))
	for i, n := range arr {
		res[i] = strconv.Itoa(n)
	}

	return res
}

// FlattenStructToMap convert struct to map with key name is defined int tag named `keyname`.
// Parameter must be a struct pointer.
func FlattenStructToMap(obj interface{}) map[string]interface{} {
	mapResult := map[string]interface{}{}

	valueOf := reflect.ValueOf(obj)
	if valueOf.Type().Kind() != reflect.Interface &&
		valueOf.Type().Kind() != reflect.Ptr {
		return map[string]interface{}{}
	}

	structValue := valueOf.Elem()
	numsOfField := structValue.NumField()
	for i := 0; i < numsOfField; i++ {
		field := structValue.Type().Field(i)
		name := field.Name
		structFieldValue := structValue.FieldByName(name)
		if !structFieldValue.CanSet() {
			continue
		}

		var keyName string
		var tagOptions tagOptions
		tag, ok := field.Tag.Lookup(structToMapTag)
		if ok {
			keyName, tagOptions = parseTag(tag)
			if tagOptions.contains(ommitEmptyTagOption) && structFieldValue.IsZero() {
				continue
			}
		}

		if keyName == "" {
			keyName = name
		}

		mapResult[keyName] = structFieldValue.Interface()
	}

	return mapResult
}

// FlattenStructToContainerMap convert struct to container map with key name is defined int tag named `keyname`.
// Parameter must be a struct pointer.
func FlattenStructToContainerMap(obj interface{}) container.Map {
	containerMap := container.Map{}
	toMap := FlattenStructToMap(obj)
	for k, v := range toMap {
		containerMap[k] = v
	}

	return containerMap
}
