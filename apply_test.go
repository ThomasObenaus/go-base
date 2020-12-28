package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getTargetTypeAndValue(t *testing.T) {
	// GIVEN
	type my struct {
		Field1 string `cfg:"{'name':'field_1'}"`
	}
	m := my{
		Field1: "field1",
	}

	// WHEN
	t1, v1, err1 := getTargetTypeAndValue(&m)

	// THEN
	require.NoError(t, err1)
	assert.Equal(t, reflect.TypeOf(m), t1)
	assert.NotEqual(t, reflect.Ptr, t1.Kind())
	assert.True(t, v1.IsValid())
}

func Test_getTargetTypeAndValue_Fail(t *testing.T) {
	// GIVEN
	type my struct {
		Field1 string `cfg:"{'name':'field_1'}"`
	}
	m := my{
		Field1: "field1",
	}

	// WHEN
	t1, v1, err1 := getTargetTypeAndValue(m)

	// THEN
	require.Error(t, err1)
	assert.Nil(t, t1)
	assert.True(t, v1.IsZero())

	// WHEN
	t2, v2, err2 := getTargetTypeAndValue(nil)

	// THEN
	require.Error(t, err2)
	assert.Nil(t, t2)
	assert.True(t, v2.IsZero())

	// GIVEN
	var nilVal *my

	// WHEN
	t3, v3, err3 := getTargetTypeAndValue(nilVal)

	// THEN
	require.Error(t, err3)
	assert.Nil(t, t3)
	assert.True(t, v3.IsZero())
}
