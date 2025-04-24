//
// GO DEMO
//
// To run:
//
// 	go run demo.go
//
// To build:
//
// 	go build demo.go

package main

import (
	"fmt"
)

// square calculates the square of a number.
func square(n int) int {
	return n * n
}

func next(n int) (int, int) {
	return n, n + 1
}

// parity returns "even" or "odd" based on the parity of a number.
func parity(n int) string {
	if n%2 == 0 {
		return "even"
	}
	return "odd"
}

// main is the entry point of the program.
func main() {
	fmt.Println("Hello World")

	// Initialize a counter
	var i int = 0

	// Loop from 0 to 19
	for i < 20 {
		result := square(i) // Calculate the square
		p := parity(i)      // Determine parity
		fmt.Printf("The square of %d is %d - it is %s\n", i, result, p)
		i++ // Increment the counter
	}

	// Create an array to store squares of numbers
	var squares [20]int
	for j := 0; j < 20; j++ {
		squares[j] = square(j) // Populate the array with squares
	}

	fmt.Println("Squares array:", squares)

}

func printArray(arr [20]int) {
	fmt.Printf("Array = %v/n", arr)
}

func modifyArray(arr [20]int) {
	arr[0] = 99
}
