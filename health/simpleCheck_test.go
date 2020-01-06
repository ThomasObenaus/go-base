package health

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewSimpleCheck(t *testing.T) {

	// GIVEN
	err1 := fmt.Errorf("The Check Failed")
	checkFun := func() error {
		return err1
	}
	checkName := "check1"

	// WHEN
	check, err := NewSimpleCheck(checkName, checkFun)

	// THEN
	assert.NoError(t, err)
	assert.NotNil(t, check)
	assert.Equal(t, err1, check.IsHealthy())
	assert.Equal(t, checkName, check.String())
}

func Test_NewSimpleCheckShouldFail(t *testing.T) {

	// WHEN
	check, err := NewSimpleCheck("name", nil)

	// THEN
	assert.Error(t, err)
	assert.Nil(t, check)

	// WHEN
	check, err = NewSimpleCheck("", nil)

	// THEN
	assert.Error(t, err)
	assert.Nil(t, check)
}
