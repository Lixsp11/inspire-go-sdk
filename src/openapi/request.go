package openapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/Lixsp11/inspire-go-sdk/src/cli"
)

type requestModel interface{ Validate() error }
type emptyRequest struct{}

func (*emptyRequest) Validate() error { return nil }

func modelFields(t reflect.Type) map[string]reflect.StructField {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	fields := map[string]reflect.StructField{}
	if t.Kind() != reflect.Struct {
		return fields
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := strings.Split(f.Tag.Get("json"), ",")[0]
		if f.PkgPath == "" && tag != "-" && tag != "" {
			fields[f.Name] = f
			fields[tag] = f
		}
	}
	return fields
}

// Decode a single JSON value while retaining numbers and rejecting duplicate keys.
func decodeJSON(data []byte) (interface{}, error) {
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	var read func(int) (interface{}, error)
	read = func(depth int) (interface{}, error) {
		if depth > 64 {
			return nil, fmt.Errorf("JSON nesting is too deep")
		}
		token, err := d.Token()
		if err != nil {
			return nil, err
		}
		switch token {
		case json.Delim('{'):
			out := map[string]interface{}{}
			for d.More() {
				key, err := d.Token()
				if err != nil {
					return nil, err
				}
				name, ok := key.(string)
				if !ok {
					return nil, fmt.Errorf("invalid JSON object")
				}
				if _, exists := out[name]; exists {
					return nil, fmt.Errorf("duplicate JSON member")
				}
				value, err := read(depth + 1)
				if err != nil {
					return nil, err
				}
				out[name] = value
			}
			end, err := d.Token()
			if err != nil || end != json.Delim('}') {
				return nil, fmt.Errorf("invalid JSON object")
			}
			return out, nil
		case json.Delim('['):
			out := []interface{}{}
			for d.More() {
				value, err := read(depth + 1)
				if err != nil {
					return nil, err
				}
				out = append(out, value)
			}
			end, err := d.Token()
			if err != nil || end != json.Delim(']') {
				return nil, fmt.Errorf("invalid JSON array")
			}
			return out, nil
		default:
			if _, ok := token.(json.Delim); ok {
				return nil, fmt.Errorf("unexpected JSON delimiter")
			}
			return token, nil
		}
	}
	value, err := read(0)
	if err != nil {
		return nil, err
	}
	if _, err = d.Token(); err != io.EOF {
		return nil, fmt.Errorf("expected one JSON value")
	}
	return value, nil
}

func normalizeModel(value interface{}, t reflect.Type) (interface{}, error) {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if value == nil {
		return nil, nil
	}
	switch t.Kind() {
	case reflect.Struct:
		object, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected a JSON object")
		}
		fields := modelFields(t)
		out := map[string]interface{}{}
		for name, item := range object {
			field, ok := fields[name]
			if !ok {
				return nil, fmt.Errorf("unknown request member")
			}
			tag := strings.Split(field.Tag.Get("json"), ",")[0]
			if _, duplicate := out[tag]; duplicate {
				return nil, fmt.Errorf("duplicate request member")
			}
			converted, err := normalizeModel(item, field.Type)
			if err != nil {
				return nil, err
			}
			out[tag] = converted
		}
		return out, nil
	case reflect.Slice:
		array, ok := value.([]interface{})
		if !ok {
			return nil, fmt.Errorf("expected a JSON array")
		}
		for i, item := range array {
			converted, err := normalizeModel(item, t.Elem())
			if err != nil {
				return nil, err
			}
			array[i] = converted
		}
		return array, nil
	case reflect.Map:
		object, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected a JSON object")
		}
		for name, item := range object {
			converted, err := normalizeModel(item, t.Elem())
			if err != nil {
				return nil, err
			}
			object[name] = converted
		}
		return object, nil
	default:
		return value, nil
	}
}

func assignParameter(value reflect.Value, path []string, raw string) error {
	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}
		value = value.Elem()
	}
	if len(path) == 0 {
		var data []byte
		if value.Kind() == reflect.String {
			data, _ = json.Marshal(raw)
		} else {
			parsed, err := decodeJSON([]byte(raw))
			if err != nil {
				return err
			}
			parsed, err = normalizeModel(parsed, value.Type())
			if err != nil {
				return err
			}
			data, err = json.Marshal(parsed)
			if err != nil {
				return err
			}
		}
		return json.Unmarshal(data, value.Addr().Interface())
	}
	switch value.Kind() {
	case reflect.Struct:
		field, ok := modelFields(value.Type())[path[0]]
		if !ok || field.Name != path[0] {
			return fmt.Errorf("unknown request parameter")
		}
		return assignParameter(value.FieldByIndex(field.Index), path[1:], raw)
	case reflect.Slice:
		index, err := strconv.Atoi(path[0])
		if err != nil || index < 1 || index > 1000 || strconv.Itoa(index) != path[0] {
			return fmt.Errorf("array indexes must be between 1 and 1000")
		}
		for value.Len() < index {
			value.Set(reflect.Append(value, reflect.Zero(value.Type().Elem())))
		}
		return assignParameter(value.Index(index-1), path[1:], raw)
	default:
		return fmt.Errorf("parameter has no child fields")
	}
}

func loadRequest(ctx *cli.Context, request requestModel) error {
	body, hasBody := ctx.Flags().GetValue("body")
	file, hasFile := ctx.Flags().GetValue("body-file")
	if hasBody && hasFile {
		return cli.UsageError("InvalidParameter", "Choose --body or --body-file.")
	}
	parameters := map[string]string{}
	for _, f := range append(ctx.Flags().Flags(), ctx.UnknownFlags().Flags()...) {
		if f.Category == "parameter" && f.IsAssigned() {
			value, _ := f.GetValue()
			parameters[f.Name] = value
		}
	}
	if (hasBody || hasFile) && len(parameters) > 0 {
		return cli.UsageError("InvalidParameter", "Use a JSON body or API parameter flags in one call.")
	}
	if hasBody || hasFile {
		data := []byte(body)
		if hasFile {
			var err error
			data, err = readBodyFile(ctx, file)
			if err != nil {
				return err
			}
		}
		if len(data) > 2<<20 {
			return cli.UsageError("InvalidParameter", "Request body exceeds 2 MiB.")
		}
		parsed, err := decodeJSON(data)
		if err != nil {
			return cli.UsageError("InvalidParameter", "Body must be one valid JSON object without duplicate members.")
		}
		if _, ok := parsed.(map[string]interface{}); !ok {
			return cli.UsageError("InvalidParameter", "Body must be a JSON object.")
		}
		parsed, err = normalizeModel(parsed, reflect.TypeOf(request))
		if err != nil {
			return cli.UsageError("InvalidParameter", "Body fields do not match this API; inspect its help.")
		}
		data, err = json.Marshal(parsed)
		if err != nil {
			return cli.UsageError("InvalidParameter", "Body cannot be encoded.")
		}
		if err = json.Unmarshal(data, request); err != nil {
			return cli.UsageError("InvalidParameter", "Body value types do not match this API.")
		}
	} else {
		for name, raw := range parameters {
			for other := range parameters {
				if strings.HasPrefix(other, name+".") {
					return cli.UsageError("InvalidParameter", "Overlapping API parameter paths are not allowed.")
				}
			}
			if err := assignParameter(reflect.ValueOf(request), strings.Split(name, "."), raw); err != nil {
				return cli.UsageError("InvalidParameter", "Invalid API parameter --"+name+"; inspect the parameter name and value type.")
			}
		}
	}
	if err := request.Validate(); err != nil {
		return cli.UsageError("InvalidParameter", err.Error())
	}
	return nil
}
