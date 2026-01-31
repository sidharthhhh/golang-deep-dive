package main

import "fmt"

// Custom Type
type UserID int

func main() {
	fmt.Println("--- Type Conversions ---")

	var a int = 42
	var b float64 = 42.5

	// c := a + b // ERROR: invalid operation: a + b (mismatched types int and float64)

	// Explicit conversion
	c := float64(a) + b
	fmt.Printf("Sum: %f\n", c)

	// Data loss warning
	d := int(b)
	fmt.Printf("Float to Int (Truncated): %d\n", d) // 42.5 -> 42

	fmt.Println("\n--- Custom Types ---")
	var myID UserID = 100
	var plainInt int = 100

	// if myID == plainInt { } // ERROR: Mismatched types

	fmt.Printf("UserID: %v (Type: %T)\n", myID, myID)

	// Conversion works because underlying types are compatible
	if int(myID) == plainInt {
		fmt.Println("Values match after conversion")
	}

	fmt.Println("\n--- String <-> Bytes ---")
	str := "Hello"
	bytes := []byte(str) // Conversion, memory copy

	fmt.Printf("String: %s\n", str)
	fmt.Printf("Bytes: %v\n", bytes)

	bytes[0] = 'M' // Change 'H' to 'M'
	str2 := string(bytes)
	fmt.Println("New String:", str2)
	fmt.Println("Old String (Immutable):", str)
}
