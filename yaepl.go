package yaepl

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Tag holds the parsed components of the tag
type Tag struct {
	Key      string
	Required bool
}

// parseTag parses the tag value of the format `mylib:"key:MY_KEY;required"`
func parseTag(tag string) Tag {
	var result Tag
	parts := strings.SplitSeq(tag, ";")

	for part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "key:") {
			key := strings.TrimPrefix(part, "key:")
			result.Key = key
		} else if part == "required" {
			result.Required = true
		}
	}
	return result
}

func Unmarshal(destination any) error {
	v := reflect.ValueOf(destination)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		tagVal := field.Tag.Get("yaepl")
		if tagVal == "" {
			continue
		}
		tag := parseTag(tagVal)
		if tag.Key == "" {
			return fmt.Errorf("Key must not be empty")
		}
		field_env_value := os.Getenv(tag.Key)
		if tag.Required && field_env_value == "" {
			return fmt.Errorf("Required environment variable %s is not set", tag.Key)
		}
		field_value := v.Field(i)
		switch field.Type.Kind() {
		case reflect.Bool:
			val, err := strconv.ParseBool(field_env_value)
			if err != nil {
				return err
			}
			field_value.SetBool(val)
		case reflect.Uint:
			val, err := strconv.ParseUint(field_env_value, 10, 64)
			if err != nil {
				return err
			}
			field_value.SetUint(val)
		case reflect.Int:
			val, err := strconv.ParseInt(field_env_value, 10, 64)
			if err != nil {
				return err
			}
			field_value.SetInt(val)
		case reflect.Float32:
			val, err := strconv.ParseFloat(field_env_value, 32)
			if err != nil {
				return err
			}
			field_value.SetFloat(val)
		case reflect.Float64:
			val, err := strconv.ParseFloat(field_env_value, 64)
			if err != nil {
				return err
			}
			field_value.SetFloat(val)
		case reflect.String:
			field_value.SetString(field_env_value)
		default:
			return fmt.Errorf("Unsupported type %s", field.Type.Kind().String())
		}
	}
	return nil
}
