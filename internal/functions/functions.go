package functions

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"reflect"
	"strings"

	log "github.com/vetyy/kubetools/internal/logging"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Incr(i interface{}, arg int) int64 {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Int:
		return v.Int() + int64(arg)
	}
	return 10
}

func DefaultValue(arg interface{}, value interface{}) interface{} {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		if v.Len() == 0 {
			return arg
		}
	case reflect.Bool:
		if !v.Bool() {
			return arg
		}
	case reflect.Invalid:
		return arg
	default:
		return value
	}

	return value
}

func Length(value interface{}) int {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len()
	case reflect.String:
		return len([]rune(v.String()))
	}

	return 0
}

func Lower(s string) string {
	return strings.ToLower(s)
}

func Upper(s string) string {
	return strings.ToUpper(s)
}

func Base64encode(v string) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}

func Base64decode(v string) string {
	data, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func RandomString(length int) (string, error) {
	if length < 0 {
		return "", fmt.Errorf("expected number bigger than zero")
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b), nil
}

func ErrorFunc(message string) error {
	log.Fatalf("Environment invoked error: %v", message)
	return nil
}
