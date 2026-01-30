package main

// 02_compilation_flow.go
// -----------------------------------------------------------------------------
// This file explains and SIMULATES the journey from SOURCE CODE to BINARY.
//
// MENTAL MODEL:
// Go uses a "static linking" model.
// When you run `go build`, the following happens:
// 1. SCANNING & PARSING:
//    - The compiler reads .go files.
//    - Checks for syntax errors (e.g., missing brackets).
//
// 2. TYPE CHECKING:
//    - Ensures types match (e.g., cannot add int + string).
//    - AST (Abstract Syntax Tree) generation.
//
// 3. INTERMEDIATE REPRESENTATION (SSA):
//    - Converts code to a machine-independent form (Static Single Assignment).
//    - This is where Go optimizations happen (dead code elimination).
//
// 4. MACHINE CODE GENERATION:
//    - Converts SSA to assembly for the target CPU (x86_64 or ARM64).
//    - This happens in the `pkg/` directory for dependencies.
//
// 5. LINKING:
//    - The `link` tool combines your code + standard library (fmt, runtime).
//    - Result: A single, standalone binary executable.
//    - No external DLLs needed!
// -----------------------------------------------------------------------------

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("=== Compiling Simulation ===")

	// Simulating the environment context
	fmt.Printf("1. Compiler detects Architecture: %s\n", runtime.GOARCH)
	fmt.Printf("2. Compiler detects OS: %s\n", runtime.GOOS)

	fmt.Println("3. Linking 'fmt' package... [Standard Lib]")
	fmt.Println("4. Linking 'runtime' package... [Scheduler, GC]")

	fmt.Println("5. Generating executable binary...")
	fmt.Println("   Success! The binary includes everything needed to run.")
}
