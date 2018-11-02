# jsonValidator
Package qeshion2104/jsonValidator check struct members required and existed from application/json Type

# Example

Here's a quick example: we parse POST body value then 
check all required are exist:

```go
type Person struct {
	Id int `json:"id, required"`
	Name string `json:"name, required"` 
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
    body, err  := ioutil.ReadAll(r.Body)

	defer r.Body.Close()
	if err !=  nil {
		return err
	}
	person := &Person{}
	// body is []byte as jsonData
	err := jsonValidator.GetValidateJsonData(person, body)
	if err !=  nil {
		//Handle error
		return err
	}

    // Do something with person.Name or person.Id
}
```
To define a member as required, use a struct tag "json" and contain "required". If not use "required", member will be ignore when doing check :
```go
type  exceptJson002  struct {
	Ex002 int  `json:"ex002, required, omitempty"`
	Ex003 float64  `json:"ex003, required, omitempty"`
}

type  exceptJson001  struct {
	Ex001 int  `json:"ex001, required, omitempty"`
	Ej *exceptJson002 `json:"ej, required, omitempty"`
}

type  exceptJson_ArrayRequired  struct {
	Arr []*exceptJson002 `json:"shouldExist, required, omitempty"`
}

type  exceptJson  struct {
	Ex1 float64  `json:"ex1, required, omitempty"`
	Ex2 string  `json:"ex2, omitempty"`
	Strs []string  `json:"strs, required, omitempty"`
	Obj *exceptJson002 `json:"obj, required, omitempty"`
	ObjArr []*exceptJson002 `json:"objArr, required, omitempty"`
	ObjArr2 [][]*exceptJson002 `json:"objArr2, required, omitempty"`
}
```

# License

MIT licensed. See the LICENSE file for detail.