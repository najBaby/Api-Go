package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"strconv"
)

func parse(s string, etype reflect.Type) reflect.Type {
	chaines := strings.Split(s, "_")
	if f, ok := etype.FieldByName(chaines[0]); ok {
		if f.Type.Kind() == reflect.Ptr && f.Type.Elem().Kind() == reflect.Struct {
			if others := strings.Join(chaines[1:], "_"); len(others) > 0 {
				return parse(others, f.Type.Elem())
			}
		} else {
			return f.Type
		}
	}
	return nil
}

type U struct {
	Age      int
	Name     string
	LastName string
	U        *U
}

func main() {
	m := make(map[string]interface{})
	m["Age"] = 10
	m["LastName"] = "Coulibaly"
	n := make(map[string]interface{})
	n["Name"] = "baby"
	n["Age"] = 20
	m["U"] = n
	u := U{}
	json.Unmarshal([]byte(string(m)), &u)
	fmt.Println(u.U)
}

func Convert(etype reflect.Type, m map[string]string) reflect.Value {
	mapValue := reflect.New(reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf(interface{})))
	for k, v := range m {
		mapV := reflect.New(reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf(interface{})))
		fields := strings.Split(k, ",")
	}
	return mapValue
}

func assign(s string, mp map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	m[s] = mp
	return m
}

func Assign(etype reflect.Type, i interface{}) reflect.Value {
	evalue := reflect.Zero(etype)
	itype := reflect.TypeOf(i)
	switch {
	case itype.Kind() == etype.Kind():
		evalue.Set(reflect.ValueOf(i))
	case itype.Kind() == reflect.String:
		evalue.Set(AssignByType(etype, i.(string)))
	case itype.Kind() == reflect.Map:
	}
	return evalue
}

func AssignValues(etype reflect.Type, m interface{}) reflect.Value {
	mp := m.(map[string]interface{})
	evalue := reflect.New(etype)
	for k, v := range mp {
		if f, ok := etype.FieldByName(k); ok {
			if f.Type.Kind() == reflect.Ptr && f.Type.Elem().Kind() == reflect.Struct {
				evalue.Elem().FieldByName(k).Set(AssignValues(f.Type.Elem(), v))
			} else {
				evalue.Elem().FieldByName(k).Set(Assign(f.Type, v))
			}
		}
	}
	return evalue
}

func AssignByType(etype reflect.Type, value string) reflect.Value {
	vz := reflect.Zero(etype)
	switch etype.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		vz = AssignInt(etype, value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		vz = AssignUint(etype, value)
	case reflect.String:
		vz = reflect.ValueOf(value)
	}
	return vz
}

func AssignInt(etype reflect.Type, value string) reflect.Value {
	intvalue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return reflect.Zero(etype)
	}
	evalue := reflect.New(etype).Elem()
	evalue.SetInt(intvalue)
	return evalue
}

func AssignUint(etype reflect.Type, value string) reflect.Value {
	uintvalue, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return reflect.Zero(etype)
	}
	evalue := reflect.New(etype).Elem()
	evalue.SetUint(uintvalue)
	return evalue
}

func AssignBool(etype reflect.Type, value string) reflect.Value {
	return reflect.Zero(etype)
}

//
// type Request struct {
// 	Entity interface{}
// 	*RequestOptions
// }
//
// type RequestOptions struct {
// 	Operation string                 `json:"operation"`
// 	Model     string                 `json:"model"`
// 	Fields    []string               `json:"fields"`
// 	Filter    map[string]interface{} `json:"filter"`
// 	Errors    []string               `json:"-"`
// 	Method    string                 `json:"-"`
// }
//
// func NewRequest(httprequest *http.Request) *Request {
// 	reqOpts := RequestOptions{
// 		Method: httprequest.Method,
// 	}
// 	switch query, mutation := httprequest.FormValue("query"), httprequest.FormValue("mutation"); {
// 	case query != "":
// 		reqOpts.Model = query
// 		reqOpts.Operation = "query"
// 		httprequest.Form.Del("query")
// 	case mutation != "":
// 		reqOpts.Model = mutation
// 		reqOpts.Operation = "mutation"
// 		httprequest.Form.Del("mutation")
// 	default:
// 		reqOpts.Errors = []string{"invalid opeartion"}
// 		goto direct
// 	}
// 	reqOpts.Fields = strings.Split(httprequest.FormValue("fields"), ",")
// 	httprequest.Form.Del("fields")
// 	reqOpts.Filter = make(map[string]interface{})
// 	for k := range httprequest.Form {
// 		reqOpts.Filter[k] = httprequest.FormValue(k)
// 	}
// direct:
// 	return &Request{
// 		RequestOptions: &reqOpts,
// 	}
// }
//
// func (req *Request) ParseEntities(entities []interface{}) {
// 	for _, entity := range entities {
// 		entityType := reflect.TypeOf(entity).Elem()
// 		if entityType.Name() != req.Model {
// 			continue
// 		}
//
// 		for _, fieldName := range req.Fields {
// 			if _, ok := entityType.FieldByName(fieldName); !ok {
// 				req.Errors = append(req.Errors, fmt.Sprintf("unknow field's name '%v'", fieldName))
// 			}
// 		}
// 		entityValue := reflect.ValueOf(entity).Elem()
// 		for filterName, filterValue := range req.Filter {
// 			if field, ok := entityType.FieldByName(filterName); !ok {
// 				req.Errors = append(req.Errors, fmt.Sprintf("unknow field's name '%v'", filterName))
// 			} else {
// 				fieldValue := entityValue.FieldByName(filterName)
// 				switch field.Type.Kind() {
// 				case reflect.Ptr:
// 					fieldValue.Set(reflect.New(field.Type.Elem()))
// 					if fieldID := fieldValue.Elem().FieldByName("Id"); fieldID.IsValid() {
// 						if id, b := strconv.Atoi(filterValue.(string)); b == nil {
// 							fieldID.SetInt(int64(id))
// 						}
// 					}
// 				case reflect.String:
// 					entityValue.FieldByName(filterName).SetString(filterValue.(string))
// 				}
// 			}
// 		}
//
// 		req.Entity = entityValue.Interface()
// 		return
// 	}
// 	req.Errors = append(req.Errors, fmt.Sprintf("unknow entity's name '%v'", req.Model))
// 	return
// }
