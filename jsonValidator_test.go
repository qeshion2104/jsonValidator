package jsonValidator

import (
	// "encoding/json"
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
)


type exceptJson002 struct {
	Ex002 int			`json:"ex002, required, omitempty"`
	Ex003 float64		`json:"ex003, required, omitempty"`
}

type exceptJson001 struct {
	Ex001 int			`json:"ex001, required, omitempty"`
	Ej *exceptJson002	`json:"ej, required, omitempty"`
}

type exceptJson_required struct {
	ShouldExist int	`json:"shouldExist, required, omitempty"`
}
type exceptJson_ArrayRequired struct {
	Arr []*exceptJson002 	`json:"shouldExist, required, omitempty"`
}

type exceptJson struct {
	Ex1 float64					`json:"ex1, required, omitempty"`
	Ex2 string					`json:"ex2, omitempty"`
	Strs []string				`json:"strs, required, omitempty"`
	Obj *exceptJson002			`json:"obj, required, omitempty"`
	ObjArr []*exceptJson002 	`json:"objArr, required, omitempty"`
	ObjArr2 [][]*exceptJson002 	`json:"objArr2, required, omitempty"`
	
}

func TestJsonValidator(t *testing.T) {
	testJson := []byte(`{
		"ex1":1,
		"ex2":"ex2",
		"strs":["a","b"],
		"obj":{"ex002":5, "ex003": 7}, 
		"objArr":[{"ex002":50, "ex003": 70}, {"ex002":52, "ex003": 61}],
		"objArr2":[[{"ex002":50, "ex003": 70}, {"ex002":52, "ex003": 61}]]}`)

	result := &exceptJson{}
	err := GetValidateJsonData(result, testJson)
	assert.Nil(t, err)

	assert.Equal(t, float64(1), result.Ex1, "result.Ex1 should be 1")
	assert.Equal(t, "ex2", result.Ex2, "result.Ex2 should be ex2")
	
	assert.Equal(t, "a", result.Strs[0])
	assert.Equal(t, "b", result.Strs[1])

	obj := reflect.ValueOf(result.Obj)
	objType := obj.Type()
	if objType != reflect.TypeOf(&exceptJson002{}) {
		t.Error("obj should be exceptJson002 Type")
	}

	assert.Equal(t, 50, result.ObjArr[0].Ex002, "result.Ex1 should be 1")
	assert.Equal(t, float64(70), result.ObjArr[0].Ex003)
	assert.Equal(t, 52, result.ObjArr[1].Ex002)
	assert.Equal(t, float64(61), result.ObjArr[1].Ex003)
	
	obj = reflect.ValueOf(result.ObjArr)
	objType = obj.Type()
	if objType != reflect.TypeOf(make([]*exceptJson002, 0)) {
		t.Error("ObjArr should be []*exceptJson002 Type")
	}

	obj = reflect.ValueOf(result.ObjArr2)
	objType = obj.Type()
	if objType != reflect.TypeOf(make([][]*exceptJson002, 0)) {
		t.Error("ObjArr2 should be [][]*exceptJson002 Type")
	}

	obj = reflect.ValueOf(result)
	objType = obj.Type()
	if objType != reflect.TypeOf(&exceptJson{}) {
		t.Error("result should be *exceptJson Type")
	}
}

func TestJsonValidator_Required(t *testing.T) {
	except := &exceptJson_required{}
	jsonData := []byte(
		`{
			"ShouldExist":5
		}`)
	err := GetValidateJsonData(except, jsonData)
	assert.Nil(t, err)
	assert.Equal(t, 5, except.ShouldExist)
}

func TestJsonValidator_LoseRequired(t *testing.T) {
	except := &exceptJson_required{}
	jsonData := []byte(
		`{
			"Should":5
		}`)
	err := GetValidateJsonData(except, jsonData)
	assert.NotNil(t, err)
}

func TestJsonValidator_ArrayRequired(t *testing.T) {
	except := &exceptJson_ArrayRequired{}
	jsonData := []byte(
		`{
			"arr":[{"ex002":50, "ex008": 70}, {"ex002":52, "ex003": 61}]
		}`)
	err := GetValidateJsonData(except, jsonData)
	assert.NotNil(t, err)
}

func TestJsonValidator_Fail(t *testing.T) {
	tests := []struct {
		Name   string
		ShouldError bool
		Field []byte
    }{
		{
			Name: "lose ex1 required",
			ShouldError: true,
			Field: []byte(
				`{
				"ex2":"ex2",
				"strs":["a","b"],
				"obj":{"ex002":5, "ex003": 7}, 
				"objArr":[{"ex002":50, "ex003": 70}, {"ex002":52, "ex003": 61}],
				"objArr2":[[{"ex002":50, "ex003": 70}, {"ex002":52, "ex003": 61}]]
				}`),
		},{
			Name: "lose ex2 not required",
			ShouldError: false,
			Field: []byte(
				`{
					"ex1":1,
					"strs":["a","b"],
					"obj":{"ex002":5, "ex003": 7}, 
					"objArr":[{"ex002":50, "ex003": 70}, {"ex002":52, "ex003": 61}],
					"objArr2":[[{"ex002":50, "ex003": 70}, {"ex002":52, "ex003": 61}]]
				}`),
		},{
			Name: "lose strs required",
			ShouldError: true,
			Field: []byte(
				`{
					"ex1":1,
					"ex2":"ex2",
					"obj":{"ex002":5, "ex003": 7}, 
					"objArr":[{"ex002":50, "ex003": 70}, {"ex002":52, "ex003": 61}],
					"objArr2":[[{"ex002":50, "ex003": 70}, {"ex002":52, "ex003": 61}]]
				}`),
		},{
			Name: "lose obj required",
			ShouldError: true,
			Field: []byte(
				`{
					"ex1":1,
					"ex2":"ex2",
					"strs":["a","b"],
					"objArr":[{"ex002":50, "ex003": 70}, {"ex002":52, "ex003": 61}],
					"objArr2":[[{"ex002":50, "ex003": 70}, {"ex002":52, "ex003": 61}]]
				}`),
		},{
			Name: "lose objArr required",
			ShouldError: true,
			Field: []byte(
				`{
					"ex1":1,
					"ex2":"ex2",
					"strs":["a","b"],
					"obj":{"ex002":5, "ex003": 7}, 
					"objArr2":[[{"ex002":50, "ex003": 70}, {"ex002":52, "ex003": 61}]]
				}`),
		},{
			Name: "lose objArr2 required",
			ShouldError: true,
			Field: []byte(
				`{
					"ex1":1,
					"ex2":"ex2",
					"strs":["a","b"],
					"obj":{"ex002":5, "ex003": 7}, 
					"objArr":[{"ex002":50, "ex003": 70}, {"ex002":52, "ex003": 61}]
				}`),
		},{
			Name: "obj with wrong schema ex008 lost ex003",
			ShouldError: true,
			Field: []byte(
				`{
					"ex1":1,
					"ex2":"ex2",
					"strs":["a","b"],
					"obj":{"ex002":5, "ex008": 7}, 
					"objArr":[{"ex002":50, "ex003": 70}, {"ex002":52, "ex003": 61}],
					"objArr2":[[{"ex002":50, "ex003": 70}, {"ex002":52, "ex003": 61}]]
				}`),
		},{
			Name: "objArr with wrong schema ex008 lost ex003 in objArr[0]",
			ShouldError: true,
			Field: []byte(
				`{
					"ex1":1,
					"ex2":"ex2",
					"strs":["a","b"],
					"obj":{"ex002":5, "ex003": 7}, 
					"objArr":[{"ex002":50, "ex008": 70}, {"ex002":52, "ex003": 61}],
					"objArr2":[[{"ex002":50, "ex003": 70}, {"ex002":52, "ex003": 61}]]
				}`),
		},{
			Name: "objArr with wrong schema ex008 lost ex002 in objArr[1]",
			ShouldError: true,
			Field: []byte(
				`{
					"ex1":1,
					"ex2":"ex2",
					"strs":["a","b"],
					"obj":{"ex002":5, "ex003": 7}, 
					"objArr":[{"ex002":50, "ex003": 70}, {"ex008":52, "ex003": 61}],
					"objArr2":[[{"ex002":50, "ex003": 70}, {"ex002":52, "ex003": 61}]]
				}`),
		},
	}
     

	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func (t *testing.T) {
			t.Parallel()		//Parallel run
			result := &exceptJson{}
			err := GetValidateJsonData(result, tt.Field)
			if tt.ShouldError {
				assert.NotNil(t, err)
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func BenchmarkJsonValidator(b *testing.B) {
    for i := 0; i < b.N; i++ {
        testJson := []byte(`{
			"ex1":1,
			"ex2":"ex2",
			"strs":["a","b"],
			"obj":{"ex002":5, "ex003": 7}, 
			"objArr":[{"ex002":50, "ex003": 70}, {"ex002":52, "ex003": 61}],
			"objArr2":[[{"ex002":50, "ex003": 70}, {"ex002":52, "ex003": 61}]]}`)
	
		result := &exceptJson{}
		err := GetValidateJsonData(result, testJson)
		if err != nil {
			panic(err)
		}
    }
}


