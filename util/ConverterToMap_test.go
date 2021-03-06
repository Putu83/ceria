package util

import (
	"reflect"
	"testing"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/stretchr/testify/assert"
)

func TestConvertToMap(t *testing.T) {

	util := NewUtilConvertToMap()

	columns := []string{
		"name",
		"age",
	}

	values := []interface{}{
		"Maulana",
		50,
	}

	newMap := map[string]string{
		"name": "Maulana",
		"age":  "50",
	}

	assert.Equal(t, newMap, util.ConvertToDynamicMap(columns, values))

	// key test
	type Keys struct {
		name string
		age  string
	}

	keyTest := Keys{}

	assert.Equal(t, columns, util.ConvertInterfaceToKeyStr(keyTest))

	// test from a value array model

	listKeys := []Keys{
		{name: "triadi", age: "10"},
		{name: "udin", age: "15"},
	}

	tesRes := util.ConvertMultiStructToMap(listKeys)
	assert.Len(t, tesRes, 2)

	// PTR test
	tesRes = util.ConvertMultiStructToMap(&listKeys)
	assert.Len(t, tesRes, 2)

	// test a single struct
	assert.Len(t, util.ConvertStructToSingeMap(&Keys{}), 2)

	// test a struct nested
	type StructNestedTest struct {
		gorm.Model
		Name string
	}

	testStructNested := &StructNestedTest{
		Name: "perjaka witan",
	}

	assert.Len(t, util.ConvertStructToSingeMap(testStructNested), 5)

	var uiMode uint
	uiMode = 64

	type testMultiField struct {
		Name    string
		Age     int
		YourAge uint
		Times   time.Time
		Strb    []byte
	}
	strEx := "abc"
	var strb []byte
	copy(strb[:], strEx)

	// test a single with value
	assert.Len(t, util.ConvertStructToSingeMap(&testMultiField{
		Name:    "triadi",
		Age:     40,
		YourAge: uiMode,
		Times:   time.Now(),
		Strb:    strb,
	}), 5)

	// ignore the interface not a struct

	var nullParam interface{}
	assert.Len(t, util.ConvertStructToSingeMap(nullParam), 0)

	// check struct with array and non array
	type KeysParent struct {
		Non Keys
		Arr []Keys
		Str []string
		Mp  map[string]string
	}

	testStruct := KeysParent{
		Non: Keys{
			name: "triadi",
			age:  "40",
		},
		Arr: []Keys{
			Keys{name: "udin", age: "40"},
			Keys{name: "paijo", age: "70"},
		},
		Str: []string{
			"udin",
			"paijo",
		},
		Mp: map[string]string{
			"title": "hello",
		},
	}

	resMulti := util.ConvertStructToSingeMap(&testStruct)

	mp := reflect.ValueOf(resMulti["mp"])
	knd := mp.Kind()
	if knd == reflect.Map {
		for _, e := range mp.MapKeys() {
			v := mp.MapIndex(e)
			assert.Equal(t, "hello", util.ConvertDataToString(v))
		}
	}

	assert.Len(t, resMulti, 4)

	// convert datetime to string test
	var interfaceToTest interface{}
	now := time.Now()
	interfaceToTest = &now

	assert.NotEmpty(t, util.ConvertDataToString(interfaceToTest))

	// reuseable test
	util = NewUtilConvertToMap()

	// check iterature in reflect
	var param reflect.Value

	assert.Empty(t, util.RefValueToInterface(param))

	// check negative case only struct
	a := 5
	param = reflect.ValueOf(a)
	assert.Empty(t, util.RefValueToInterface(param))

	// check positif case
	param = reflect.ValueOf(&Keys{})
	assert.NotEmpty(t, util.RefValueToInterface(param))

	// Test interfacemap to map
	newMapParams := map[string]interface{}{
		"data":    "student",
		"length":  50,
		"Times":   time.Now(),
		"Strb":    strb,
		"YourAge": uiMode,
	}

	newUtil := NewUtilConvertToMap()
	// Tes error when failed map
	resM, errM := newUtil.ConvertInterfaceMaptoMap("is_string")
	assert.Error(t, errM)
	assert.Empty(t, resM)

	// Tes with real map
	resValidM, errValidM := util.ConvertInterfaceMaptoMap(newMapParams)

	assert.NoError(t, errValidM)
	assert.NotEmpty(t, resValidM)

	// test with non slice type
	var testFailedConvert string
	assert.Nil(t, util.RebuildToSlice(&testFailedConvert))

	// test convert struct to slice struct
	makeItSlice := util.RebuildToSlice(&Keys{"udin", "50"})

	assert.NotNil(t, makeItSlice)
	// assert.True(t, )

	// Test with ignore elastic
	type KeysExampleIgnore struct {
		Non Keys `ceria:"ignoreStructField"`
	}

	keysIgnoreElastic := KeysExampleIgnore{
		Non: Keys{
			name: "triadi",
			age:  "40",
		},
	}
	resultSetStr := util.SetFieldNullByTag(&keysIgnoreElastic)
	valStrAfterSet := reflect.ValueOf(resultSetStr).Elem()
	fieldToTest := valStrAfterSet.FieldByName("non")
	assert.Empty(t, fieldToTest)

	// Non pointer
	keysIgnoreElasticNon := KeysExampleIgnore{
		Non: Keys{
			name: "triadi",
			age:  "40",
		},
	}
	resultSetStrNon := util.SetFieldNullByTag(keysIgnoreElasticNon)
	valStrAfterSetNon := reflect.ValueOf(resultSetStrNon)
	fieldToTestNon := valStrAfterSetNon.FieldByName("non")
	assert.Empty(t, fieldToTestNon)

	// tested not struct and not ptr whichis is string or other
	stringTested := "abcd"
	assert.Empty(t, util.SetFieldNullByTag(stringTested))
}
