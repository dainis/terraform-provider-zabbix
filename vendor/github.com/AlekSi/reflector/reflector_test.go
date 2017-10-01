package reflector_test

import (
	. "."
	"testing"
)

type T struct {
	Int     int
	Uint8   uint8
	Uintptr uintptr
	Float32 float32 `json:"f32"`
	String  string
	Pstring *string
	foo     int
}

func TestStructToMapBad1(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "StructToMap: expected pointer to struct as first argument, got struct" {
			t.Error(r)
		}
	}()
	m := make(map[string]interface{})
	StructToMap(T{}, m, "")
	t.Fatal("should panic")
}

func TestStructToMapBad2(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "StructToMap: expected pointer to struct as first argument, got pointer to int" {
			t.Error(r)
		}
	}()
	var i int
	m := make(map[string]interface{})
	StructToMap(&i, m, "")
	t.Fatal("should panic")
}

func TestMapToStructBad1(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "MapToStruct: expected pointer to struct as second argument, got struct" {
			t.Error(r)
		}
	}()
	MapToStruct(map[string]interface{}{}, T{}, NoConvert, "")
	t.Fatal("should panic")
}

func TestMapToStructBad2(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "MapToStruct: expected pointer to struct as second argument, got pointer to int" {
			t.Error(r)
		}
	}()
	var i int
	MapToStruct(map[string]interface{}{}, &i, NoConvert, "")
	t.Fatal("should panic")
}

func TestMapToStructWrongType(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "MapToStruct, field Uint8: interface conversion: interface is int, not uint8" {
			t.Error(r)
		}
	}()
	type T struct {
		Uint8 uint8
	}
	var s T
	m := map[string]interface{}{"Uint8": 8}
	MapToStruct(m, &s, NoConvert, "")
	t.Fatal("should panic")
}

func TestMapsToStructsBad1(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "MapsToStructs: expected pointer to slice of structs as second argument, got slice" {
			t.Error(r)
		}
	}()
	var s []T
	m := map[string]interface{}{
		"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str", "foo": 13,
	}
	MapsToStructs([]map[string]interface{}{m}, s, NoConvert, "json")
	t.Fatal("should panic")
}

func TestMapsToStructsBad2(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "MapsToStructs: expected pointer to slice of structs as second argument, got pointer to int" {
			t.Error(r)
		}
	}()
	var s *int
	m := map[string]interface{}{
		"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str", "foo": 13,
	}
	MapsToStructs([]map[string]interface{}{m}, s, NoConvert, "json")
	t.Fatal("should panic")
}

func TestMapsToStructsBad3(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "MapsToStructs: expected pointer to slice of structs as second argument, got pointer to slice of int" {
			t.Error(r)
		}
	}()
	var s *[]int
	m := map[string]interface{}{
		"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str", "foo": 13,
	}
	MapsToStructs([]map[string]interface{}{m}, s, NoConvert, "json")
	t.Fatal("should panic")
}
