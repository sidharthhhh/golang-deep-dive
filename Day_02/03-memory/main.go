package main

import "fmt"

func main() {
	fmt.Println("--- Memory Addresses ---")

	a := 10
	b := 10

	// %p prints the address in hex
	fmt.Printf("Variable a | Value: %d | Address: %p\n", a, &a)
	fmt.Printf("Variable b | Value: %d | Address: %p\n", b, &b)

	fmt.Println("\nEven though 'a' and 'b' have the same value, they live at different locations.")

	fmt.Println("\n--- Stack Frames ---")
	printAddress(a)
}

// x is a COPY of the value passed in. It lives in printAddress's stack frame.
func printAddress(x int) {
	fmt.Printf("Function param x | Value: %d | Address: %p\n", x, &x)
	fmt.Println("Notice address of 'x' is different from 'a'. It's a copy.")
}
