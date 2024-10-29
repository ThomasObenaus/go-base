package stop

import (
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_stores_shutdown_state(t *testing.T) {
	// GIVEN
	synchronizedList := OrderedStoppableList{}

	// WHEN
	synchronizedList.StopAllInOrder(zerolog.Logger{})

	// THEN
	value := synchronizedList.isShuttingDown
	assert.Equal(t, true, value)
}

func Test_returns_error_on_add_if_service_is_shutting_down(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	synchronizedList := OrderedStoppableList{}
	item1 := NewMockStoppable(mockCtrl)

	// WHEN
	synchronizedList.isShuttingDown = false
	err := synchronizedList.AddToFront(item1)

	// THEN
	assert.NoError(t, err)

	// WHEN
	synchronizedList.isShuttingDown = false
	err = synchronizedList.AddToBack(item1)

	// THEN
	assert.NoError(t, err)

	// WHEN
	synchronizedList.isShuttingDown = true
	err = synchronizedList.AddToFront(item1)

	// THEN
	assert.Error(t, err)

	// WHEN
	synchronizedList.isShuttingDown = true
	err = synchronizedList.AddToBack(item1)

	// THEN
	assert.Error(t, err)
}
