package reflector_test

import (
	. "."
	"fmt"
)

func ExampleStructToMap() {
	type T struct {
		Uint8   uint8
		Float32 float32 `json:"f32"` // tag will be used
		String  string
		Pstring *string // pointer will be followed
		foo     int     // not exported
	}
	pstr := "pstr"
	s := T{8, 3.14, "str", &pstr, 13}
	m := make(map[string]interface{})
	StructToMap(&s, m, "json")
	fmt.Printf("%#v %#v %#v %#v %#v\n", m["Uint8"], m["f32"], m["String"], m["Pstring"], m["foo"])
	// Output:
	// 0x8 3.14 "str" "pstr" <nil>
}

func ExampleStructValueToMap() {
	type T struct {
		Uint8   uint8
		Float32 float32 `json:"f32"` // tag will be used
		String  string
		Pstring *string // nil will be present in map
		foo     int     // not exported
	}
	s := T{8, 3.14, "str", nil, 13}
	m := make(map[string]interface{})
	StructValueToMap(s, m, "json")
	fmt.Printf("%#v %#v %#v %#v %#v\n", m["Uint8"], m["f32"], m["String"], m["Pstring"], m["foo"])
	_, ok := m["Pstring"]
	fmt.Println(ok)
	// Output:
	// 0x8 3.14 "str" <nil> <nil>
	// true
}

func ExampleStructsToMaps() {
	type T struct {
		Uint8   uint8
		Float32 float32 `json:"f32"` // tag will be used
		String  string
		Pstring *string // pointer will be followed
		foo     int     // not exported
	}
	pstr := "pstr"
	s := []T{{8, 3.14, "str", &pstr, 13}}
	var m []map[string]interface{}
	StructsToMaps(s, &m, "json")
	fmt.Printf("%#v %#v %#v %#v %#v\n", m[0]["Uint8"], m[0]["f32"], m[0]["String"], m[0]["Pstring"], m[0]["foo"])
	// Output:
	// 0x8 3.14 "str" "pstr" <nil>
}

func ExampleMapToStruct_noConvert() {
	type T struct {
		Uint8   uint8   // no automatic type conversion
		Float32 float32 `json:"f32"` // tag will be used
		String  string  // not present in map, will not be set
		Pstring *string // value will be deep-copied
		foo     int     // not exported, will not be set
	}
	var s T
	m := map[string]interface{}{"Uint8": uint8(8), "f32": float32(3.14), "Pstring": "pstr", "foo": 13}
	MapToStruct(m, &s, NoConvert, "json")
	m["Pstring"] = "foo"
	fmt.Printf("%#v %#v %#v %#v %#v\n", s.Uint8, s.Float32, s.String, *s.Pstring, s.foo)
	// Output:
	// 0x8 3.14 "" "pstr" 0
}

func ExampleMapToStruct_strconv() {
	type T struct {
		Uint8   uint8   // type conversion via strconv
		Float32 float32 `json:"f32"` // tag will be used
		String  string  // not present in map, will not be set
		Pstring *string // value will be deep-copied
		foo     int     // not exported, will not be set
	}
	var s T
	m := map[string]interface{}{"Uint8": 8, "f32": 3, "Pstring": "pstr", "foo": 13}
	MapToStruct(m, &s, Strconv, "json")
	m["Pstring"] = "foo"
	fmt.Printf("%#v %#v %#v %#v %#v\n", s.Uint8, s.Float32, s.String, *s.Pstring, s.foo)
	// Output:
	// 0x8 3 "" "pstr" 0
}

func ExampleMapsToStructs_noConvert() {
	type T struct {
		Uint8   uint8   // no automatic type conversion
		Float32 float32 `json:"f32"` // tag will be used
		String  string  // not present in first map, will not be set
		Pstring *string // value will be deep-copied
		foo     int     // not exported, will not be set
	}
	var s []T
	maps := []map[string]interface{}{
		{"Uint8": uint8(8), "Pstring": "pstr"},
		{"f32": float32(3.14), "String": "str", "foo": 13},
	}
	MapsToStructs(maps, &s, NoConvert, "json")
	fmt.Printf("%#v %#v %#v %#v %#v\n", s[0].Uint8, s[0].Float32, s[0].String, *s[0].Pstring, s[0].foo)
	fmt.Printf("%#v %#v %#v %#v %#v\n", s[1].Uint8, s[1].Float32, s[1].String, s[1].Pstring, s[1].foo)
	// Output:
	// 0x8 0 "" "pstr" 0
	// 0x0 3.14 "str" (*string)(nil) 0
}

func ExampleMapsToStructs_strconv() {
	type T struct {
		Uint8   uint8   // type conversion via strconv
		Float32 float32 `json:"f32"` // tag will be used
		String  string  // not present in first map, will not be set
		Pstring *string // value will be deep-copied
		foo     int     // not exported, will not be set
	}
	var s []T
	maps := []map[string]interface{}{
		{"Uint8": 8, "f32": 3, "foo": 13, "Pstring": "pstr"},
		{"Uint8": "9", "f32": "4", "String": "43", "foo": "13"},
	}
	MapsToStructs(maps, &s, Strconv, "json")
	fmt.Printf("%#v %#v %#v %#v %#v\n", s[0].Uint8, s[0].Float32, s[0].String, *s[0].Pstring, s[0].foo)
	fmt.Printf("%#v %#v %#v %#v %#v\n", s[1].Uint8, s[1].Float32, s[1].String, s[1].Pstring, s[1].foo)
	// Output:
	// 0x8 3 "" "pstr" 0
	// 0x9 4 "43" (*string)(nil) 0
}
