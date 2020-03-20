package main

import (
	"errors"
	"fmt"
)
func Sqrt(f float64) (float64, error) {
	if f < 0 {
		return 0, errors.New ("math - square root of negative number")
	}
	// implementation of Sqrt
	return 1,nil
}
func badCall() {
	panic("bad end")
}

func test() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("Panicing %s\r\n", e)
		}
	}()
	badCall()
	fmt.Printf("After bad call\r\n") // <-- wordt niet bereikt
}
func main(){
	result , _ :=Sqrt(-1)
	print(result)
	fmt.Printf("Calling test\r\n")
	test()
	fmt.Printf("Test completed\r\n")
}