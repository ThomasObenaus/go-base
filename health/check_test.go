package health

import (
	"testing"

	mock_health "github.com/ThomasObenaus/go-base/test/mocks/health"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_NewRegistry(t *testing.T) {

	// GIVEN

	// WHEN
	registry := NewRegistry()

	// THEN
	assert.NotNil(t, registry.healthChecks)
}

func Test_ShouldRegister(t *testing.T) {

	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	check1 := mock_health.NewMockCheck(mockCtrl)
	registry := NewRegistry()

	// WHEN
	check1.EXPECT().Name().Return("check1")
	err := registry.Register(check1)

	// THEN
	assert.NoError(t, err)
	assert.Len(t, registry.healthChecks, 1)
}

func Test_ShouldNotRegister(t *testing.T) {

	// GIVEN
	registry := NewRegistry()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	check1 := mock_health.NewMockCheck(mockCtrl)

	// WHEN
	check1.EXPECT().Name().Return("")

	err := registry.Register(check1)

	// THEN
	assert.Error(t, err)
	assert.Len(t, registry.healthChecks, 0)

	// WHEN
	err = registry.Register(nil)

	// THEN
	assert.Error(t, err)
	assert.Len(t, registry.healthChecks, 0)
}
