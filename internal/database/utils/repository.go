package utils

import (
	"errors"
	"reflect"
	"strings"
)

type Reflection struct {
	field reflect.StructField
	value reflect.Value
}

func GetReflection[T any](entity *T, tag string, value string) (*Reflection, error) {
	reflection := reflect.TypeOf(entity).Elem()
	for i := 0; i < reflection.NumField(); i++ {
		field := reflection.Field(i)
		bsonTags := field.Tag.Get(tag)
		tags := strings.Split(bsonTags, ",")
		for j := 0; j < len(tags); j++ {
			if tags[j] == value {
				return &Reflection{
					field: field,
					value: reflect.ValueOf(entity).Elem().Field(i),
				}, nil
			}
		}
	}
	return nil, errors.New("cannot fetch the reflection")
}

func (r *Reflection) SetValue(v any) {
	if v == nil {
		r.value.Set(reflect.Zero(r.value.Type()))
		return
	}

	r.value.Set(reflect.ValueOf(v))
}

func (r *Reflection) GetStringValue() string {
	return r.value.String()
}
