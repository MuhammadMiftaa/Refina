package main

import (
	"fmt"
	"runtime"
)

func main() {
	totalCPU := runtime.NumCPU()
	totalThread := runtime.GOMAXPROCS(-1)
	totalGoRoutine := runtime.NumGoroutine()

	fmt.Println("Total CPU:", totalCPU)
	fmt.Println("Total Threads:", totalThread)
	fmt.Println("Total Goroutines:", totalGoRoutine)
}