package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testParseConfigTag struct {
	configTagStr string
	typeOfEntry  reflect.Type
	nameOfParent string
}

func Test_parseConfigTag_NoParent(t *testing.T) {
	// GIVEN
	simpleString := testParseConfigTag{
		configTagStr: "{'name':'field-string','desc':'string field','default':'default'}",
		typeOfEntry:  reflect.TypeOf(""),
	}

	// WHEN + THEN
	tagStr, err := parseConfigTag(simpleString.configTagStr, simpleString.typeOfEntry, simpleString.nameOfParent)
	assert.NoError(t, err)
	assert.Equal(t, "field-string", tagStr.Name)
	assert.Equal(t, "string field", tagStr.Description)
	assert.Equal(t, "default", tagStr.Def)
	assert.Equal(t, simpleString.typeOfEntry, reflect.TypeOf(tagStr.Def))
	assert.False(t, tagStr.IsRequired())

	// GIVEN
	simpleInt := testParseConfigTag{
		configTagStr: "{'name':'field-int','desc':'int field','default':1111}",
		typeOfEntry:  reflect.TypeOf(int(0)),
	}

	// WHEN + THEN
	tagInt, err := parseConfigTag(simpleInt.configTagStr, simpleInt.typeOfEntry, simpleInt.nameOfParent)
	assert.NoError(t, err)
	assert.Equal(t, "field-int", tagInt.Name)
	assert.Equal(t, "int field", tagInt.Description)
	assert.Equal(t, 1111, tagInt.Def)
	assert.Equal(t, simpleInt.typeOfEntry, reflect.TypeOf(tagInt.Def))
	assert.False(t, tagInt.IsRequired())

	// GIVEN
	simpleFloat := testParseConfigTag{
		configTagStr: "{'name':'field-float','desc':'float field','default':22.22}",
		typeOfEntry:  reflect.TypeOf(float64(0)),
	}

	// WHEN + THEN
	tagFloat, err := parseConfigTag(simpleFloat.configTagStr, simpleFloat.typeOfEntry, simpleFloat.nameOfParent)
	assert.NoError(t, err)
	assert.Equal(t, "field-float", tagFloat.Name)
	assert.Equal(t, "float field", tagFloat.Description)
	assert.Equal(t, 22.22, tagFloat.Def)
	assert.Equal(t, simpleFloat.typeOfEntry, reflect.TypeOf(tagFloat.Def))
	assert.False(t, tagFloat.IsRequired())

	// GIVEN
	simpleBool := testParseConfigTag{
		configTagStr: "{'name':'field-bool','desc':'bool field','default':true}",
		typeOfEntry:  reflect.TypeOf(bool(true)),
	}

	// WHEN + THEN
	tagBool, err := parseConfigTag(simpleBool.configTagStr, simpleBool.typeOfEntry, simpleBool.nameOfParent)
	assert.NoError(t, err)
	assert.Equal(t, "field-bool", tagBool.Name)
	assert.Equal(t, "bool field", tagBool.Description)
	assert.Equal(t, true, tagBool.Def)
	assert.Equal(t, simpleBool.typeOfEntry, reflect.TypeOf(tagBool.Def))
	assert.False(t, tagBool.IsRequired())
}

func Test_extractConfigTags_Primitives(t *testing.T) {

	// GIVEN
	type primitives struct {
		ShouldBeSkipped string
		SomeFieldStr    string  `cfg:"{'name':'field-str','desc':'a string field','default':'default value'}"`
		SomeFieldInt    int     `cfg:"{'name':'field-int','desc':'a int field','default':11}"`
		SomeFieldFloat  float64 `cfg:"{'name':'field-float','desc':'a float field','default':22.22}"`
		SomeFieldBool   bool    `cfg:"{'name':'field-bool','desc':'a bool field','default':true}"`
	}

	sType := reflect.TypeOf(primitives{})

	// WHEN
	entries, err := extractConfigTags(sType, "", configTag{})

	// THEN
	assert.NoError(t, err)
	assert.Len(t, entries, 4)
	assert.Equal(t, "field-str", entries[0].Name)
	assert.Equal(t, "a string field", entries[0].Description)
	assert.Equal(t, "default value", entries[0].Def)
	assert.Equal(t, reflect.TypeOf(""), reflect.TypeOf(entries[0].Def))
	assert.False(t, entries[0].IsRequired())

	assert.Equal(t, "field-int", entries[1].Name)
	assert.Equal(t, "a int field", entries[1].Description)
	assert.Equal(t, reflect.TypeOf(int(0)), reflect.TypeOf(entries[1].Def))
	assert.Equal(t, 11, entries[1].Def)
	assert.False(t, entries[1].IsRequired())

	assert.Equal(t, "field-float", entries[2].Name)
	assert.Equal(t, "a float field", entries[2].Description)
	assert.Equal(t, reflect.TypeOf(float64(0)), reflect.TypeOf(entries[2].Def))
	assert.Equal(t, 22.22, entries[2].Def)
	assert.False(t, entries[2].IsRequired())

	assert.Equal(t, "field-bool", entries[3].Name)
	assert.Equal(t, "a bool field", entries[3].Description)
	assert.Equal(t, reflect.TypeOf(bool(true)), reflect.TypeOf(entries[3].Def))
	assert.Equal(t, true, entries[3].Def)
	assert.False(t, entries[3].IsRequired())
}

func Test_extractConfigTags_Required(t *testing.T) {

	// GIVEN
	type primitives struct {
		SomeFielOptional  string `cfg:"{'name':'field-str','desc':'a string field','default':'default value'}"`
		SomeFieldRequired string `cfg:"{'name':'field-str','desc':'a string field'}"`
	}

	sType := reflect.TypeOf(primitives{})

	// WHEN
	entries, err := extractConfigTags(sType, "", configTag{})

	// THEN
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.False(t, entries[0].IsRequired())
	assert.True(t, entries[1].IsRequired())
}
