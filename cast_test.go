package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_createAndFillStruct_Fail(t *testing.T) {
	// GIVEN
	type my struct {
		Field1 string `cfg:"{'name':'field_1'}"`
	}

	// WHEN
	_, errMissingRequired := createAndFillStruct(reflect.TypeOf(my{}), map[string]interface{}{})

	// THEN
	assert.Error(t, errMissingRequired)

	// GIVEN
	type myInvalidConfig struct {
		Field1 string `cfg:"{}"`
	}

	// WHEN
	_, errInvalidConfig := createAndFillStruct(reflect.TypeOf(myInvalidConfig{}), map[string]interface{}{})

	// THEN
	assert.Error(t, errInvalidConfig)

	// GIVEN
	type myUnexportedField struct {
		field1 string `cfg:"{'name':'field_1','default':'value'}"`
	}

	// WHEN
	_, errUnexportedField := createAndFillStruct(reflect.TypeOf(myUnexportedField{}), map[string]interface{}{})

	// THEN
	assert.Error(t, errUnexportedField)

	// WHEN
	_, errNoStruct := createAndFillStruct(reflect.TypeOf(int(0)), map[string]interface{}{})

	// THEN
	assert.Error(t, errNoStruct)
}

func Test_createAndFillStruct(t *testing.T) {
	// GIVEN
	type nested struct {
		FieldA          float64 `cfg:"{'name':'field_a'}"`
		ShouldBeIgnored bool
	}

	type my struct {
		Field1          string `cfg:"{'name':'field_1'}"`
		Field2          nested `cfg:"{'name':'field_2'}"`
		ShouldBeIgnored bool
	}

	expected := my{Field1: "field-1", Field2: nested{FieldA: 22.22}}

	data := map[string]interface{}{
		"field_1": expected.Field1,
		"field_2": map[string]interface{}{
			"field_a": expected.Field2.FieldA,
		},
	}

	// WHEN
	structVal, err := createAndFillStruct(reflect.TypeOf(my{}), data)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, expected, structVal.Interface())
}

func Test_castToSlice_Fail(t *testing.T) {
	// GIVEN
	vIntSlice := []int{11, 22}
	vIntSliceIf := []interface{}{vIntSlice[0], vIntSlice[1]}

	// WHEN
	casted1, err1 := castToSlice(vIntSliceIf, reflect.TypeOf(int(0)))

	// THEN
	assert.Error(t, err1)
	assert.Nil(t, casted1)

	// GIVEN
	vInt := int(10)

	// WHEN
	casted2, err2 := castToSlice(vInt, reflect.TypeOf([]int{}))

	// THEN
	assert.Error(t, err2)
	assert.Nil(t, casted2)
}

func Test_castToSlice(t *testing.T) {
	// GIVEN
	vIntSlice := []int{11, 22}
	vIntSliceIf := []interface{}{vIntSlice[0], vIntSlice[1]}

	// WHEN
	casted1, err1 := castToSlice(vIntSliceIf, reflect.TypeOf([]int{}))

	// THEN
	assert.NoError(t, err1)
	assert.Equal(t, vIntSlice, casted1)

	// GIVEN
	type my struct {
		Field1 string `cfg:"{'name':'field_1'}"`
	}
	expected := []my{
		{Field1: "field-1"},
		{Field1: "field-2"},
	}

	vStructIf := []interface{}{
		map[string]interface{}{"field_1": expected[0].Field1},
		map[string]interface{}{"field_1": expected[1].Field1},
	}

	// WHEN
	casted2, err2 := castToSlice(vStructIf, reflect.TypeOf([]my{}))

	// THEN
	assert.NoError(t, err2)
	assert.Equal(t, expected, casted2)
}

func Test_castToStruct_Fail(t *testing.T) {
	// GIVEN
	type my struct {
		Field1 string `cfg:"{'name':'field_1'}"`
	}
	vStructMissingRequired := map[string]interface{}{}
	vNoStruct := int(11)

	// WHEN
	casted1, err1 := castToStruct(vStructMissingRequired, reflect.TypeOf(my{}))
	casted2, err2 := castToStruct(vStructMissingRequired, reflect.TypeOf(int(0)))
	casted3, err3 := castToStruct(vNoStruct, reflect.TypeOf(my{}))

	// THEN
	assert.Error(t, err1)
	assert.Nil(t, casted1)
	assert.Error(t, err2)
	assert.Nil(t, casted2)
	assert.Error(t, err3)
	assert.Nil(t, casted3)
}

func Test_castToStruct(t *testing.T) {
	// GIVEN
	type nested struct {
		FieldA float64 `cfg:"{'name':'field_a'}"`
	}
	type my struct {
		Field1 string   `cfg:"{'name':'field_1'}"`
		Field2 int      `cfg:"{'name':'field_2'}"`
		Field3 nested   `cfg:"{'name':'field_3'}"`
		Field4 []int    `cfg:"{'name':'field_4'}"`
		Field5 []nested `cfg:"{'name':'field_5'}"`
	}

	expected := my{
		Field1: "a field",
		Field2: 11,
		Field3: nested{
			FieldA: 22.22,
		},
		Field4: []int{11, 22},
		Field5: []nested{{FieldA: 22.22}},
	}

	vStruct := map[string]interface{}{
		"field_1": expected.Field1,
		"field_2": expected.Field2,
		"field_3": map[string]interface{}{
			"field_a": expected.Field3.FieldA,
		},
		"field_4": []interface{}{
			expected.Field4[0],
			expected.Field4[1],
		},
		"field_5": []interface{}{
			map[string]interface{}{"field_a": expected.Field5[0].FieldA},
		},
	}

	// WHEN
	casted1, err1 := castToStruct(vStruct, reflect.TypeOf(my{}))

	// THEN
	assert.NoError(t, err1)
	assert.Equal(t, reflect.TypeOf(my{}), reflect.TypeOf(casted1))
	assert.Equal(t, expected, casted1)
}

func Test_castToPrimitive(t *testing.T) {
	// GIVEN
	vInt := int(11)
	vIntSlice := []int{11, 22}
	vString := "something"

	// WHEN
	castedInt, errInt := castToPrimitive(vInt, reflect.TypeOf(int(0)))
	castedIntSlice, errIntSlice := castToPrimitive(vIntSlice, reflect.TypeOf([]int{}))
	castedWrongType, errWrongType := castToPrimitive(vString, reflect.TypeOf(int(0)))

	// THEN
	assert.NoError(t, errInt)
	assert.Equal(t, 11, castedInt)
	assert.NoError(t, errIntSlice)
	assert.Equal(t, []int{11, 22}, castedIntSlice)
	assert.Error(t, errWrongType)
	assert.Nil(t, castedWrongType)
}

func Test_isFieldExported(t *testing.T) {
	// GIVEN
	type my struct {
		ExportedField   string
		unExportedField string
	}
	reflectedType := reflect.TypeOf(my{})
	exportedField, _ := reflectedType.FieldByName("ExportedField")
	unExportedField, _ := reflectedType.FieldByName("unExportedField")

	// WHEN + THEN
	assert.True(t, isFieldExported(exportedField))
	assert.False(t, isFieldExported(unExportedField))
}

func Test_fullFieldName(t *testing.T) {
	assert.Equal(t, "root", fullFieldName("", "root"))
	assert.Equal(t, "root.child", fullFieldName("root", "child"))
	assert.Equal(t, "root.children.child", fullFieldName("root.children", "child"))
}

func Test_parseStringContainingSliceOfMaps(t *testing.T) {
	v1 := `[{"name":"name1","key":"key1","count":1},{"name":"name2","key":"key2","count":2}]`
	v2 := `[]`
	v3 := `invalid`
	v4 := `{}`
	v5 := `[1 2 3 4]`

	// WHEN
	r1, err1 := parseStringContainingSliceOfMaps(v1)
	r2, err2 := parseStringContainingSliceOfMaps(v2)
	r3, err3 := parseStringContainingSliceOfMaps(v3)
	r4, err4 := parseStringContainingSliceOfMaps(v4)
	r5, err5 := parseStringContainingSliceOfMaps(v5)

	// THEN
	require.NoError(t, err1)
	require.Len(t, r1, 2)
	assert.Equal(t, "name1", r1[0]["name"])
	assert.Equal(t, "key1", r1[0]["key"])
	assert.Equal(t, float64(1), r1[0]["count"])
	assert.Equal(t, "name2", r1[1]["name"])
	assert.Equal(t, "key2", r1[1]["key"])
	assert.Equal(t, float64(2), r1[1]["count"])
	require.NoError(t, err2)
	assert.Empty(t, r2)
	assert.Error(t, err3)
	assert.Empty(t, r3)
	assert.Error(t, err4)
	assert.Empty(t, r4)
	assert.Error(t, err5)
	assert.Empty(t, r5)
}
