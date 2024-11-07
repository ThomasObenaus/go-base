package stop

import (
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_stores_shutdown_state(t *testing.T) {
	// GIVEN
	synchronizedList := Registry{}

	// WHEN
	synchronizedList.StopAllInOrder(zerolog.Nop())

	// THEN
	value := synchronizedList.shutdownInProgressOrComplete
	assert.Equal(t, true, value)
}

func Test_returns_error_on_add_if_service_is_shutting_down(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	synchronizedList := Registry{}
	item1 := NewMockStoppable(mockCtrl)

	// WHEN
	synchronizedList.shutdownInProgressOrComplete = false
	err := synchronizedList.AddToFront(item1)

	// THEN
	assert.NoError(t, err)

	// WHEN
	synchronizedList.shutdownInProgressOrComplete = false
	err = synchronizedList.AddToBack(item1)

	// THEN
	assert.NoError(t, err)

	// WHEN
	synchronizedList.shutdownInProgressOrComplete = true
	err = synchronizedList.AddToFront(item1)

	// THEN
	assert.Error(t, err)

	// WHEN
	synchronizedList.shutdownInProgressOrComplete = true
	err = synchronizedList.AddToBack(item1)

	// THEN
	assert.Error(t, err)
}
