package main

import (
	"fmt"
	"time"
)

type result struct {
	sumValue      int
	multiplyValue int
}

func main() {
	resultChan := make(chan result, 1)
	go sumAndMultiply(2, 3, resultChan)

	res := <-resultChan
	fmt.Printf("Sum value: %d\n", res.sumValue)
	fmt.Printf("Multiply Value %d\n", res.multiplyValue)
	close(resultChan)
}

func sumAndMultiply(a, b int, resultChan chan result) {
	sumValue := a + b
	multiplyValue := a * b
	res := result{sumValue: sumValue, multiplyValue: multiplyValue}
	time.Sleep(time.Second * 2)
	resultChan <- res
	return
}
