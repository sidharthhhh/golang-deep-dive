package main

import (
	"fmt"
	"runtime/debug"
)

// 05_stack_trace_experiment.go
// -----------------------------------------------------------------------------
// This file visualizes the CALL STACK.
//
// MENTAL MODEL:
// When a function calls another, it pauses.
// The new function is "stacked" on top.
// When the top function finishes, it is "popped", and the one below resumes.
//
// This "Stack" is a real region of memory.
// -----------------------------------------------------------------------------

func main() {
	fmt.Println("=== Call Stack Visualization ===")
	// Start the chain
	level1()
}

func level1() {
	fmt.Println("-> Entering Level 1")
	level2() // Calls Level 2, so Level 1 simply waits here
	fmt.Println("<- Exiting Level 1")
}

func level2() {
	fmt.Println("  -> Entering Level 2")
	level3() // Calls Level 3, so Level 2 waits here
	fmt.Println("  <- Exiting Level 2")
}

func level3() {
	fmt.Println("    -> Entering Level 3 (The Deepest Level)")

	// PRINT THE ACTUAL STACK
	// This command prints the history of how we got here.
	fmt.Println("\n[ SYSTEM STACK TRACE ]")
	debug.PrintStack()
	fmt.Println("[ END STACK TRACE ]\n")
}
