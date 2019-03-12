package main

import (
	"fmt"
	"runtime"
)

const attempts = 10000

func main() {
	runtime.GOMAXPROCS(2)

	obj := &demoRecord{}
	testScenario("two int", obj.setTwoInt, obj.readTwoInt, obj.reset)
	testScenario("int arr", obj.setIntArr, obj.readIntArr, obj.reset)
}

func testScenario(label string, set func(), get func() (bool, bool), reset func()) {
	foundReorder := false
	var cnt int

	for cnt = 0; !foundReorder && cnt < attempts; cnt++ {
		reset()

		doneCh := make(chan struct{})
		go func() {
			defer close(doneCh)
			set()
		}()

		getOkCh := make(chan bool)
		go func() {
			defer close(getOkCh)
			done, ok := get()
			for !done && !ok {
				done, ok = get()
			}
			getOkCh <- ok
		}()

		foundReorder = <-getOkCh
		<-doneCh
	}

	if foundReorder {
		fmt.Printf("%d - FOUND %s\n", cnt, label)
	} else {
		fmt.Printf("%d - NOT FOUND %s\n", cnt, label)
	}
}

type demoRecord struct {
	val1   int
	val2   int
	intArr []int
}

func (t *demoRecord) reset() {
	t.val1 = 0
	t.val2 = 0
	t.intArr = make([]int, 2)
}

func (t *demoRecord) setTwoInt() {
	t.val1 = 1
	t.val2 = 2
}

func (t *demoRecord) readTwoInt() (bool, bool) {
	v1 := t.val1
	v2 := t.val2

	if v1 == 1 && v2 == 2 {
		return true, false
	}
	if v1 != 1 && v2 == 2 {
		fmt.Printf("v1: %d v2: %d\n", v1, v2)
		return false, true
	}
	return false, false
}

func (t *demoRecord) setIntArr() {
	t.intArr[0] = 1
	t.intArr[1] = 2
}

func (t *demoRecord) readIntArr() (bool, bool) {
	v1 := t.intArr[0]
	v2 := t.intArr[1]

	if v1 == 1 && v2 == 2 {
		return true, false
	}
	if v1 != 1 && v2 == 2 {
		fmt.Printf("v1: %d v2: %d\n", v1, v2)
		return false, true
	}
	return false, false
}
