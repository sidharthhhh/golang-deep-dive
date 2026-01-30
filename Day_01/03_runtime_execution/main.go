package main

// 03_runtime_execution.go
// -----------------------------------------------------------------------------
// This file demonstrates the Go RUNTIME.
//
// MENTAL MODEL:
// A Go program is NOT just your 'main' function.
// It is wrapped inside a RUNTIME that manages:
// 1. The HEAP (Memory allocation).
// 2. The STACK (Function calls).
// 3. The SCHEDULER (running Goroutines).
// 4. THE GARBAGE COLLECTOR (cleaning memory).
//
// execution flow:
// [OS Process Start] -> [Runtime Init] -> [main.main()] -> [Runtime Shutdown]
// -----------------------------------------------------------------------------

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	// The runtime started BEFORE this line!

	fmt.Println("=== Runtime Characteristics ===")

	// 1. CPU & Threads
	// Go defaults to using all available CPU cores.
	// This is controlled by GOMAXPROCS.
	fmt.Printf("CPU Cores Available (GOMAXPROCS): %d\n", runtime.GOMAXPROCS(0))

	// 2. Goroutines
	// Even in this simple program, runtime background tasks exist.
	// At least 1 (main) + GC background workers.
	fmt.Printf("Active Goroutines (system threads + user): %d\n", runtime.NumGoroutine())

	// 3. Memory Stats (Stack & Heap)
	// We can inspect the memory allocator.
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Memory Allocated (Heap): %d bytes\n", m.Alloc)
	fmt.Printf("System Memory Obtained: %d bytes\n", m.Sys)

	// Simulate work to see runtime change
	go func() {
		// New goroutine! managed by scheduler.
		time.Sleep(100 * time.Millisecond)
	}()

	time.Sleep(10 * time.Millisecond) // Give scheduler time to start the goroutine
	fmt.Printf("Active Goroutines (after spawn): %d\n", runtime.NumGoroutine())
}
