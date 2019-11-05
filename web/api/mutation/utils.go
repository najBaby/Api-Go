package mutation

import (
	"fmt"
	"reflect"
)

// operation mutation verify fields
func ParseFields(mutation *Mutation, etype reflect.Type) {
	verifyNames(mutation.Fields.(map[string]interface{}), etype, &mutation.Errors)
}

func verifyNames(mp map[string]interface{}, etype reflect.Type, errors *[]string) {
	switch etype.Kind() {
	case reflect.Ptr:
		etype = etype.Elem()
		fallthrough
	case reflect.Struct:
		for k, v := range mp {
			if f, ok := etype.FieldByName(k); ok {
				m := make(map[string]interface{})
				m[k] = v
				verifyNames(m, f.Type, errors)
			} else {
				*errors = append(*errors, fmt.Sprintf("field's name '%v' no found in the entity '%v'", k, etype.Name()))
			}
		}
	case reflect.Slice:
	default:
		for k, v := range mp {
			if etype.Kind() != reflect.TypeOf(v).Kind() {
				*errors = append(*errors, fmt.Sprintf("invalid type '%v' of '%v'", v, k))
			}
		}
	}
}

// operation mutation Assign filter
func ParseFilter(mutation *Mutation, etype reflect.Type) {

}
