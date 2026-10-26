// Package model provides native JSON preservation shared by the service models.
package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"
)

type Extra struct {
	ExtraFields map[string]json.RawMessage `json:"-"`
}

func knownFields(value interface{}) map[string]bool {
	t := reflect.TypeOf(value)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	known := map[string]bool{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := strings.Split(f.Tag.Get("json"), ",")[0]
		if name != "" && name != "-" {
			known[name] = true
		}
	}
	return known
}

func Marshal(value interface{}, extra map[string]json.RawMessage) ([]byte, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var fields map[string]json.RawMessage
	if err = json.Unmarshal(data, &fields); err != nil {
		return nil, err
	}
	known := knownFields(value)
	// Preserve explicit empty native collections while keeping nil fields omitted.
	v := reflect.ValueOf(value)
	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		item := v.Field(i)
		if !known[name] || (item.Kind() != reflect.Slice && item.Kind() != reflect.Map) || item.IsNil() || item.Len() != 0 {
			continue
		}
		if item.Kind() == reflect.Slice {
			fields[name] = json.RawMessage(`[]`)
		} else {
			fields[name] = json.RawMessage(`{}`)
		}
	}
	for key, value := range extra {
		if !known[key] {
			fields[key] = value
		}
	}
	return json.Marshal(fields)
}

func Unmarshal(data []byte, value interface{}, extra *map[string]json.RawMessage) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		return fmt.Errorf("expected model object")
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return fmt.Errorf("invalid model JSON")
	}
	if fields == nil {
		return fmt.Errorf("expected model object")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(value); err != nil {
		return fmt.Errorf("model field type does not match native contract")
	}
	known := knownFields(value)
	if *extra == nil {
		*extra = map[string]json.RawMessage{}
	}
	for key, item := range fields {
		if !known[key] {
			(*extra)[key] = append(json.RawMessage(nil), item...)
		}
	}
	for key := range known {
		delete(*extra, key)
	}
	return nil
}

func ToMap(value interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("model cannot be encoded as finite JSON")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var result map[string]interface{}
	err = decoder.Decode(&result)
	return result, err
}
func FromMap(value map[string]interface{}, target interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("map cannot be encoded as finite JSON")
	}
	return json.Unmarshal(data, target)
}
func Required(value *string, field string) error {
	if value == nil || strings.TrimSpace(*value) == "" {
		return fmt.Errorf("%s is required", field)
	}
	return nil
}
func ValidateChildren(value interface{}) error {
	v := reflect.ValueOf(value)
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return fmt.Errorf("model is nil")
		}
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.CanInterface() {
			continue
		}
		if f.Kind() == reflect.Pointer && !f.IsNil() {
			if child, ok := f.Interface().(interface{ Validate() error }); ok {
				if err := child.Validate(); err != nil {
					return err
				}
			}
		}
		if f.Kind() == reflect.Slice {
			for j := 0; j < f.Len(); j++ {
				item := f.Index(j)
				if item.Kind() == reflect.Pointer && item.IsNil() {
					return fmt.Errorf("model list contains nil")
				}
				if child, ok := item.Interface().(interface{ Validate() error }); ok {
					if err := child.Validate(); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
func Positive(value *float64, field string) error {
	if value == nil || *value <= 0 || math.IsInf(*value, 0) || math.IsNaN(*value) {
		return fmt.Errorf("%s must be finite and positive", field)
	}
	return nil
}
func Page(page, size *int32, all bool, max int32) error {
	if page != nil && *page < 1 {
		return fmt.Errorf("page must be positive")
	}
	if size != nil && ((!all || *size != -1) && (*size < 1 || (max > 0 && *size > max))) {
		return fmt.Errorf("invalid page_size")
	}
	return nil
}
func Total(value interface{}) (int64, bool, error) {
	if value == nil {
		return 0, false, nil
	}
	var text string
	switch v := value.(type) {
	case string:
		text = v
	case json.Number:
		text = string(v)
	case int:
		text = fmt.Sprint(v)
	case int64:
		text = fmt.Sprint(v)
	case float64:
		if math.Trunc(v) != v {
			return 0, false, fmt.Errorf("invalid total")
		}
		text = fmt.Sprintf("%.0f", v)
	default:
		return 0, false, fmt.Errorf("invalid total")
	}
	var n json.Number = json.Number(text)
	i, err := n.Int64()
	if err != nil || i < 0 {
		return 0, false, fmt.Errorf("invalid total")
	}
	return i, true, nil
}
