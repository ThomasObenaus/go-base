package list

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func Test_can_add_items_to_front(t *testing.T) {
	synchronizedList := SynchronizedList[string]{}
	item1 := "itme 1"
	item2 := "itme 2"
	item3 := "itme 3"

	synchronizedList.AddToFront(item1)
	synchronizedList.AddToFront(item2)
	synchronizedList.AddToFront(item3)

	items := synchronizedList.GetItems()

	assert.Equal(t, items, []string{item3, item2, item1})
}

func Test_can_add_items_to_back(t *testing.T) {
	synchronizedList := SynchronizedList[string]{}
	item1 := "itme 1"
	item2 := "itme 2"
	item3 := "itme 3"

	synchronizedList.AddToBack(item1)
	synchronizedList.AddToBack(item2)
	synchronizedList.AddToBack(item3)

	items := synchronizedList.GetItems()

	assert.Equal(t, items, []string{item1, item2, item3})
}

func Test_does_allow_concurrent_add_to_front(t *testing.T) {
	synchronizedList := SynchronizedList[string]{}
	waitGroup := sync.WaitGroup{}

	for i := 0; i < 10000; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			synchronizedList.AddToFront("test-string")
		}()
	}

	waitGroup.Wait()

	items := synchronizedList.GetItems()
	assert.Equal(t, 10000, len(items))
}

func Test_does_allow_concurrent_add_to_back(t *testing.T) {
	synchronizedList := SynchronizedList[string]{}
	waitGroup := sync.WaitGroup{}

	for i := 0; i < 10000; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			synchronizedList.AddToBack("test-string")
		}()
	}

	waitGroup.Wait()

	items := synchronizedList.GetItems()
	assert.Equal(t, 10000, len(items))
}

func Test_does_allow_concurrent_get_items(t *testing.T) {
	// TODO: this trips the data race detector of go but how to test this explicitly ?
	synchronizedList := SynchronizedList[string]{}
	waitGroup := sync.WaitGroup{}

	for i := 0; i < 10000; i++ {
		waitGroup.Add(2)
		go func() {
			defer waitGroup.Done()
			synchronizedList.AddToBack("test-string")
		}()

		go func() {
			defer waitGroup.Done()
			items := synchronizedList.GetItems()
			if len(items) > 0 {
				assert.Equal(t, "test-string", items[0])
			}
		}()
	}

	waitGroup.Wait()
}
