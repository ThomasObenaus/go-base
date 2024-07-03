package list

import (
	"github.com/ThomasObenaus/go-base/shutdown/v2/stop"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func Test_can_add_items_to_front(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	synchronizedList := SynchronizedList{}
	item1 := NewMockStoppable(mockCtrl)
	item2 := NewMockStoppable(mockCtrl)
	item3 := NewMockStoppable(mockCtrl)

	synchronizedList.AddToFront(item1)
	synchronizedList.AddToFront(item2)
	synchronizedList.AddToFront(item3)

	items := synchronizedList.GetItems()

	assert.Equal(t, items, []stop.Stoppable{item3, item2, item1})
}

func Test_can_add_items_to_back(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	synchronizedList := SynchronizedList{}
	item1 := NewMockStoppable(mockCtrl)
	item2 := NewMockStoppable(mockCtrl)
	item3 := NewMockStoppable(mockCtrl)

	synchronizedList.AddToBack(item1)
	synchronizedList.AddToBack(item2)
	synchronizedList.AddToBack(item3)

	items := synchronizedList.GetItems()

	assert.Equal(t, items, []stop.Stoppable{item1, item2, item3})
}

func Test_does_allow_concurrent_add_to_front(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	synchronizedList := SynchronizedList{}
	waitGroup := sync.WaitGroup{}

	for i := 0; i < 10000; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			synchronizedList.AddToFront(NewMockStoppable(mockCtrl))
		}()
	}

	waitGroup.Wait()

	items := synchronizedList.GetItems()
	assert.Equal(t, 10000, len(items))
}

func Test_does_allow_concurrent_add_to_back(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	synchronizedList := SynchronizedList{}
	waitGroup := sync.WaitGroup{}

	for i := 0; i < 10000; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			synchronizedList.AddToBack(NewMockStoppable(mockCtrl))
		}()
	}

	waitGroup.Wait()

	items := synchronizedList.GetItems()
	assert.Equal(t, 10000, len(items))
}

func Test_does_allow_concurrent_get_items(t *testing.T) {
	// TODO: this trips the data race detector of go but how to test this explicitly ?
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStoppable := NewMockStoppable(mockCtrl)

	synchronizedList := SynchronizedList{}
	waitGroup := sync.WaitGroup{}

	for i := 0; i < 10000; i++ {
		waitGroup.Add(2)
		go func() {
			defer waitGroup.Done()
			synchronizedList.AddToBack(mockStoppable)
		}()

		go func() {
			defer waitGroup.Done()
			items := synchronizedList.GetItems()
			if len(items) > 0 {
				assert.Equal(t, mockStoppable, items[0])
			}
		}()
	}

	waitGroup.Wait()
}
