package stop

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func Test_can_add_items_to_front(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	synchronizedList := OrderedStoppableList{}
	item1 := NewMockStoppable(mockCtrl)
	item2 := NewMockStoppable(mockCtrl)
	item3 := NewMockStoppable(mockCtrl)

	err := synchronizedList.AddToFront(item1)
	require.NoError(t, err)
	err = synchronizedList.AddToFront(item2)
	require.NoError(t, err)
	err = synchronizedList.AddToFront(item3)
	require.NoError(t, err)

	assert.Equal(t, synchronizedList.items, []Stoppable{item3, item2, item1})
}

func Test_can_add_items_to_back(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	synchronizedList := OrderedStoppableList{}
	item1 := NewMockStoppable(mockCtrl)
	item2 := NewMockStoppable(mockCtrl)
	item3 := NewMockStoppable(mockCtrl)

	err := synchronizedList.AddToBack(item1)
	require.NoError(t, err)
	err = synchronizedList.AddToBack(item2)
	require.NoError(t, err)
	err = synchronizedList.AddToBack(item3)
	require.NoError(t, err)

	assert.Equal(t, synchronizedList.items, []Stoppable{item1, item2, item3})
}

func Test_does_allow_concurrent_add_to_front(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	synchronizedList := OrderedStoppableList{}
	waitGroup := sync.WaitGroup{}

	for i := 0; i < 10000; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			err := synchronizedList.AddToFront(NewMockStoppable(mockCtrl))
			require.NoError(t, err)
		}()
	}

	waitGroup.Wait()

	items := synchronizedList.items
	assert.Equal(t, 10000, len(items))
}

func Test_does_allow_concurrent_add_to_back(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	synchronizedList := OrderedStoppableList{}
	waitGroup := sync.WaitGroup{}

	for i := 0; i < 10000; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			err := synchronizedList.AddToBack(NewMockStoppable(mockCtrl))
			require.NoError(t, err)
		}()
	}

	waitGroup.Wait()

	items := synchronizedList.items
	assert.Equal(t, 10000, len(items))
}
