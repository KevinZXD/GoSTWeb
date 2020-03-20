package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)
	defer  close(ch)
	go sendData(ch)
	go getData(ch)
	//out := make(chan int)  //产生死锁
	//out <- 2
	//go f2(out)
	//go f1(out)
	time.Sleep(1e9)
	//print(<-ch)

}

func sendData(ch chan string) {
	ch <- "Washington"
	ch <- "Tripoli"
	ch <- "London"
	ch <- "Beijing"
	ch <- "Tokyo"

}
func f2(in chan int) {
	in<-1
}
func f1(in chan int) {
	fmt.Println(<-in)
}

func getData(ch chan string) {

	//time.Sleep(2e9)
	for {
		input,ok := <-ch
		if  ok{
			fmt.Printf("%s ", input)
		}

	}
	//v1,v2,v3:=<-ch,<-ch,<-ch
	//print(v1)
	//print(v2)
	//print(v3)
}