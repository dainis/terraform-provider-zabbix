package reflector_test

import (
	. "."
	"reflect"
	"testing"
)

func BenchmarkStructToMap(b *testing.B) {
	s := T{42, 8, 0xbadcafe, 3.14, "str", nil, 13}
	m := make(map[string]interface{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StructToMap(&s, m, "json")
	}
	b.StopTimer()
	expected := map[string]interface{}{
		"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str", "Pstring": nil,
	}
	if !reflect.DeepEqual(expected, m) {
		b.Fatalf("%#v\n%#v", expected, m)
	}
}

func BenchmarkStructValueToMap(b *testing.B) {
	s := T{42, 8, 0xbadcafe, 3.14, "str", nil, 13}
	m := make(map[string]interface{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StructValueToMap(s, m, "json")
	}
	b.StopTimer()
	expected := map[string]interface{}{
		"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str", "Pstring": nil,
	}
	if !reflect.DeepEqual(expected, m) {
		b.Fatalf("%#v\n%#v", expected, m)
	}
}

func BenchmarkStructsToMaps(b *testing.B) {
	s := []T{{42, 8, 0xbadcafe, 3.14, "str", nil, 13}}
	var m []map[string]interface{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StructsToMaps(s, &m, "json")
	}
	b.StopTimer()
	expected := []map[string]interface{}{{
		"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str", "Pstring": nil,
	}}
	if !reflect.DeepEqual(expected, m) {
		b.Fatalf("%#v\n%#v", expected, m)
	}
}

func BenchmarkMapToStruct(b *testing.B) {
	var s T
	m := map[string]interface{}{
		"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe),
		"f32": float32(3.14), "String": "str", "foo": 13,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapToStruct(m, &s, NoConvert, "json")
	}
	b.StopTimer()
	expected := T{42, 8, 0xbadcafe, 3.14, "str", nil, 0}
	if !reflect.DeepEqual(expected, s) {
		b.Fatalf("Expected %#v, got %#v", expected, s)
	}
}

func BenchmarkMapsToStructs(b *testing.B) {
	var s []T
	maps := []map[string]interface{}{
		{"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe)},
		{"f32": float32(3.14), "String": "str", "foo": 13},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapsToStructs(maps, &s, NoConvert, "json")
	}
	b.StopTimer()
	expected := []T{{42, 8, 0xbadcafe, 0, "", nil, 0}, {0, 0, 0, 3.14, "str", nil, 0}}
	if !reflect.DeepEqual(expected, s) {
		b.Fatalf("Expected %#v, got %#v", expected, s)
	}
}

func BenchmarkMapsToStructs2(b *testing.B) {
	var s []T
	maps := []interface{}{
		map[string]interface{}{"Int": 42, "Uint8": uint8(8), "Uintptr": uintptr(0xbadcafe)},
		map[string]interface{}{"f32": float32(3.14), "String": "str", "foo": 13},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapsToStructs2(maps, &s, NoConvert, "json")
	}
	b.StopTimer()
	expected := []T{{42, 8, 0xbadcafe, 0, "", nil, 0}, {0, 0, 0, 3.14, "str", nil, 0}}
	if !reflect.DeepEqual(expected, s) {
		b.Fatalf("Expected %#v, got %#v", expected, s)
	}
}
