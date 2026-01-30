package main

import "fmt"

// 01_hello_world.go
// This is the simplest Go program.
//
// WHY:
// - To verify the Go toolchain is installed.
// - To see the structure: package -> import -> function.
//
// MENTAL MODEL:
// - Every runnable Go program must have a package 'main'.
// - The execution entry point is always function 'main'.
// - Imports bring in functionality from the standard library (like 'fmt' for formatting).

func main() {
	fmt.Println("Hello, System!")
}
