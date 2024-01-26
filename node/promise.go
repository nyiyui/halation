package node

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Promise struct {
	FieldName      string
	SupplyNodeName NodeName
}

// setValue assign a value on the target via a field name.
// The field name may contain identifiers (note: only ASCII letters and numbers for simplicity), strings, and integers (e.g. no parentheses).
// Examples: FieldName, SupplyNodeName.Package, NodeTypes["abc"].Function
func setValue(target interface{}, fieldName string, value interface{}) error {
	v := reflect.ValueOf(target)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if !v.CanAddr() {
		return errors.New("target not addressable, must be a pointer")
	}
	origFieldName := fieldName
	_ = origFieldName
	start := 0
	for {
		if v.Kind() == reflect.Interface {
			v = v.Elem()
		}
		first, end, t, err := getFirst(fieldName)
		if err != nil {
			return err
		}
		lastFirst := end == len(fieldName)
		switch t {
		case tokenTypeIdent:
			first := first.(string)
			if v.Kind() != reflect.Struct {
				return fmt.Errorf("%s: expected struct, got %s", origFieldName[start:end], v.Kind())
			}
			if _, ok := v.Type().FieldByName(first); !ok {
				return fmt.Errorf("field %s not found in %T %v", first, v.Interface(), v.Interface())
			}
			v = v.FieldByName(first)
			if lastFirst {
				v.Set(reflect.ValueOf(value))
			}
		case tokenTypeInt | tokenTypeKey:
			first := first.(int)
			if k := v.Kind(); !(k == reflect.Array || k == reflect.Slice || k == reflect.String) {
				return fmt.Errorf("%s: expected array or slice or string, got %s", origFieldName[start:end], k)
			}
			if first >= v.Len() {
				return fmt.Errorf("%s: index %d out of range", origFieldName[start:end], first)
			}
			v = v.Index(first)
			if lastFirst {
				v.Set(reflect.ValueOf(value))
			}
		case tokenTypeString | tokenTypeKey:
			first := first.(string)
			if k := v.Kind(); k != reflect.Map {
				return fmt.Errorf("%s: expected map, got %s", origFieldName[start:end], k)
			}
			if !lastFirst {
				v = v.MapIndex(reflect.ValueOf(first))
			} else {
				v.SetMapIndex(reflect.ValueOf(first), reflect.ValueOf(value))
			}
		default:
			panic("not implemented yet")
		}
		start = end
		if lastFirst {
			break
		} else {
			fieldName = fieldName[end:]
		}
	}
	return nil
}

type tokenType int

const (
	tokenTypeIdent  = 1
	tokenTypeString = 2
	tokenTypeInt    = 3
	tokenTypeKey    = 4
)

// end is the index of the first character not included to generate first.
func getFirst(fieldName string) (first interface{}, end int, t tokenType, err error) {
	if fieldName[0] == '[' {
		k := strings.Index(fieldName, "]")
		if k == -1 {
			return nil, 0, 0, errors.New("closing bracket not found")
		}
		inner := fieldName[1:k]
		switch {
		case inner[0] == '"':
			val, err := strconv.Unquote(inner)
			return val, k + 1, tokenTypeString | tokenTypeKey, err
		case '0' <= inner[0] && inner[0] <= '9':
			val, err := strconv.ParseInt(inner, 0, 0)
			return int(val), k + 1, tokenTypeInt | tokenTypeKey, err
		default:
			return nil, 0, 0, errors.New("invalid key")
		}
	}
	i := strings.Index(fieldName, ".")
	if i != -1 {
		return fieldName[:i], i, tokenTypeIdent, nil
	}
	j := strings.Index(fieldName, "[")
	if j != -1 {
		return fieldName[:j], j, tokenTypeIdent, nil
	}
	return fieldName, len(fieldName), tokenTypeIdent, nil
}
