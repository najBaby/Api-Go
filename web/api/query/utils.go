package query

import (
	// "encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// Operation Query Verify Fields
func ParseFields(query *Query, etype reflect.Type) {
	for _, name := range query.Fields {
		verifyNames(name, etype, &query.Errors)
	}
	// b, err := json.Marshal(assignFields(query, etype))
	// if err != nil {
	// 	query.Errors = append(query.Errors, err.Error())
	// }
	// entity := reflect.New(etype).Interface()
	// err = json.Unmarshal(b, entity)

	// if err != nil {
	// 	query.Errors = append(query.Errors, err.Error())
	// }
	// query.Fields = nil
	// query.Fields = entity
}

func verifyNames(s string, etype reflect.Type, errors *[]string) {
	names := strings.Split(s, "__")
	if f, ok := etype.FieldByName(names[0]); ok {
		if others := strings.Join(names[1:], "__"); len(others) > 0 {
			verifyNames(others, f.Type.Elem(), errors)
		}
	} else {
		*errors = append(*errors, fmt.Sprintf("unknow fields's name '%v' in 'fields'", s))
	}
}

// func assignFields(query *Query, etype reflect.Type) map[string]interface{} {
// 	mp := make(map[string]interface{})
// 	if query.Fields != nil {
// 		for _, v := range query.Fields {
// 			mp[v] = assignFieldsOfStruct(v, etype, []string{"0"}, &query.Errors).Interface()
// 		}
// 	}
// 	return corrige(mp)
// }

// Operation Query Assign Filter
// func ParseFilter(query *Query, etype reflect.Type) {
// 	b, err := json.Marshal(assignFilter(query, etype))
// 	if err != nil {
// 		query.Errors = append(query.Errors, err.Error())
// 	}
// 	entity := reflect.New(etype).Interface()
// 	err = json.Unmarshal(b, entity)
// 	if err != nil {
// 		query.Errors = append(query.Errors, err.Error())
// 	}
// 	query.Filter = nil
// 	query.Filter = entity
// }

func ParseFilter(query *Query, etype reflect.Type) {
	mp := make(map[string]interface{})
	for k, v := range query.Filter.(url.Values) {
		mp[k] = assignFieldsOfStruct(k, etype, v, &query.Errors).Interface()
	}
	query.Filter = mp
}

func corrige(data map[string]interface{}) map[string]interface{} {
	mp := make(map[string]interface{})
	for k, v := range data {
		s := strings.Split(k, "__")
		if first, others := s[0], strings.Join(s[1:], "__"); len(others) == 0 {
			mp[first] = v
		} else {
			m := make(map[string]interface{})
			m[others] = v
			mp[first] = corrige(m)
		}
	}
	return mp
}

func assignByType(s string, etype reflect.Type, values []string, errors *[]string) reflect.Value {
	vz := reflect.Zero(etype)
	switch etype.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Println(values)
		vz = assignInt(etype, values[0], errors)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		vz = assignUint(etype, values[0], errors)
	case reflect.Bool:
		vz = assignBool(etype, values[0], errors)
	case reflect.String:
		vz = reflect.ValueOf(values[0])
	case reflect.TypeOf(primitive.ObjectID{}).Kind():
		id, err := primitive.ObjectIDFromHex(values[0])
		if err != nil {
			panic(err)
		}
		vz = reflect.ValueOf(id)
	case reflect.Ptr:
		vz = assignByType(s, etype.Elem(), values, errors)
	case reflect.Struct:
		vz = assignFieldsOfStruct(s, etype, values, errors)
	case reflect.Slice:
		slice := reflect.MakeSlice(etype, 0, len(values))
		child := etype.Elem()
		for _, value := range values {
			slice = reflect.Append(slice, assignByType(s, child, []string{value}, errors))
		}
		vz = slice
	}
	return vz
}

func assignFieldsOfStruct(s string, etype reflect.Type, values []string, errors *[]string) reflect.Value {
	chaines := strings.Split(s, "__")
	if f, ok := etype.FieldByName(chaines[0]); ok {
		if others := strings.Join(chaines[1:], "__"); len(others) > 0 {
			return assignByType(others, f.Type, values, errors)
		}
		return assignByType("", f.Type, values, errors)
	}
	*errors = append(*errors, fmt.Sprintf("unknow fields's name '%v' in '%v'", s, etype.Name()))
	return reflect.Zero(etype)
}

func assignInt(etype reflect.Type, value string, errors *[]string) reflect.Value {
	intvalue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		*errors = append(*errors, fmt.Sprintf("Invalid syntax: Cannot convert '%v' in Int", value))
		return reflect.Zero(etype)
	}
	evalue := reflect.New(etype).Elem()
	evalue.SetInt(intvalue)
	return evalue
}

func assignUint(etype reflect.Type, value string, errors *[]string) reflect.Value {
	uintvalue, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		*errors = append(*errors, fmt.Sprintf("Invalid syntax: Cannot convert '%v' in Uint", value))
		return reflect.Zero(etype)
	}
	evalue := reflect.New(etype).Elem()
	evalue.SetUint(uintvalue)
	return evalue
}

func assignBool(etype reflect.Type, value string, errors *[]string) reflect.Value {
	boolvalue, err := strconv.ParseBool(value)
	if err != nil {
		*errors = append(*errors, fmt.Sprintf("Invalid syntax: Cannot convert '%v' in boolean", value))
		return reflect.Zero(etype)
	}
	evalue := reflect.New(etype).Elem()
	evalue.SetBool(boolvalue)
	return evalue
}
