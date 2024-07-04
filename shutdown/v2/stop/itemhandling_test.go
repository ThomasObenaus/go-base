package stop

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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

	synchronizedList.AddToFront(item1)
	synchronizedList.AddToFront(item2)
	synchronizedList.AddToFront(item3)

	assert.Equal(t, synchronizedList.items, []Stoppable{item3, item2, item1})
}

func Test_can_add_items_to_back(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	synchronizedList := OrderedStoppableList{}
	item1 := NewMockStoppable(mockCtrl)
	item2 := NewMockStoppable(mockCtrl)
	item3 := NewMockStoppable(mockCtrl)

	synchronizedList.AddToBack(item1)
	synchronizedList.AddToBack(item2)
	synchronizedList.AddToBack(item3)

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
			synchronizedList.AddToFront(NewMockStoppable(mockCtrl))
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
			synchronizedList.AddToBack(NewMockStoppable(mockCtrl))
		}()
	}

	waitGroup.Wait()

	items := synchronizedList.items
	assert.Equal(t, 10000, len(items))
}
