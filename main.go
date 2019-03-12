package main

import (
	"errors"
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(2)

	foundReorder := false
	for !foundReorder {
		obj := &twoVal{}

		readOkCh := make(chan bool)

		go func() {
			obj.setup()
		}()

		go func() {
			done, err := obj.read()
			for !done && err == nil {
				done, err = obj.read()
			}
			if err != nil {
				readOkCh <- false
			}
			readOkCh <- true
		}()

		res := <- readOkCh
		fmt.Println(res)

		foundReorder = !res
	}
}

type twoVal struct {
	val1 int
	val2 int
}

func (t *twoVal) setup() {
	t.val1 = 1
	t.val2 = 2
}

func (t *twoVal) read() (bool, error) {
	if t.val1 == 1 && t.val2 == 2 {
		return true, nil
	}
	if t.val1 != 1 && t.val2 == 2 {
		return false, errors.New("")
	}
	return false, nil
}