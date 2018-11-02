package jsonValidator

// 2.write testCase
// 3.write HowToUse.md
// 4.write ReadMe



import (
	"encoding/json"
	"reflect"
	"strings"
	"errors"
)

const (
	Const_Required = "required"
)


func validator (t reflect.Type, reqData map[string]interface{}) error {
	
	tempReqMap := make(map[string]interface{})
	for k, v := range reqData {
		tempReqMap[strings.ToLower(k)] = v
	}

	for i:=0; i< t.NumField(); i++ {
		if ok := strings.Contains(t.Field(i).Tag.Get("json"), Const_Required); ok {
			// check required data is existed
			if data, ok := tempReqMap[strings.ToLower(t.Field(i).Name)]; ok {
				v_data := reflect.ValueOf(data)
				
				// 3 condition, 
				// pure value ex: float64,			==> json.unMarshel will check type
				// slice as  []interface {}			==> recursive check
				// obj as map [string]interface{}   ==> recursive check

				switch v_data.Kind() {
				case reflect.Slice:
					var sliceChecker func (arrType reflect.Type, reqDataArr reflect.Value) error
					sliceChecker = func (arrType reflect.Type, reqDataArr reflect.Value) error {

						// get type from array, and if get pointer, get Elem() again to get non-pointer struct type
						subStructType := arrType.Elem()
						if subStructType.Kind() == reflect.Ptr {
							subStructType = subStructType.Elem()
						}

						// check every Elemets in array to ensure there are no data lost
						subReqDataArr := reqDataArr.Interface().([]interface{})
						for _, elem := range subReqDataArr {
							elem_data := reflect.ValueOf(elem)
							switch elem_data.Kind() {
							case reflect.Slice:
								err := sliceChecker(subStructType, elem_data)
								if err != nil {
									return err
								}
							case reflect.Map:
								// if our data is struct, check required data for struct
								subReqData := elem_data.Interface().(map[string]interface{})
								err := validator(subStructType, subReqData)
								if err != nil {
									return err
								}
							default:
							}
						}
						return nil
					}
					err := sliceChecker(t.Field(i).Type, v_data)
					if (err != nil) {
						return err
					}
				case reflect.Map:
					// if our data is struct, check required data for struct
					subReqData := v_data.Interface().(map[string]interface{})
					subStructType := t.Field(i).Type.Elem()
					err := validator(subStructType, subReqData)
					if err != nil {
						return err
					}
				default:
				}
			} else {
				return errors.New("required data " + t.Field(i).Name + " is not found in ReqData")
			}
		}
	}



	return nil
}

func GetValidateJsonData (dst interface{}, jsonData []byte) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("json: interface must be a pointer to struct")
	}
	v = v.Elem()
	t := v.Type()

	dataForVaild := make(map[string]interface{})
	err := json.Unmarshal(jsonData, &dataForVaild)
	if err != nil {
		return err
	}

	err = validator(t, dataForVaild)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, dst)
	if err != nil {
		return err
	}
	
	return nil
}