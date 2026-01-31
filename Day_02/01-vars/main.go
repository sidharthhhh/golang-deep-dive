package main

import "fmt"

// Package level variables
var version = "1.0.0" // Type inferred

func main() {
	fmt.Println("--- Variable Declaration ---")
	
	// 1. Standard Declaration
	var x int = 10
	var y int      // Zero value
	fmt.Printf("x: %d, y: %d\n", x, y)

	// 2. Short Declaration
	z := 20 // Inferred as int
	fmt.Printf("z: %d, Type: %T\n", z, z)

	// 3. Multi-variable
	var a, b = 1, "hello"
	fmt.Println(a, b)

	fmt.Println("\n--- Zero Values ---")
	var (
		i int
		f float64
		s string
		boolVal bool
	)
	// %q quotes strings so we can see the empty one
	fmt.Printf("int: %d, float: %f, string: %q, bool: %t\n", i, f, s, boolVal)

	fmt.Println("\n--- Variable Shadowing (Gotcha) ---")
	shadow := 100
	fmt.Printf("Outer shadow: %d (Addr: %p)\n", shadow, &shadow)
	
	{
		// New scope
		// := creates a NEW variable, it does NOT update the outer one
		shadow := 50 
		fmt.Printf("Inner shadow: %d (Addr: %p)\n", shadow, &shadow)
	}
	
	fmt.Printf("Outer shadow after block: %d (Unchanged)\n", shadow)
}
