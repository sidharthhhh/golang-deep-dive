# Day 1: How Go Programs Are Built and Executed

## A. Concept Overview

### What problem does Go solve?
Go was built at Google to solve **scale**.
- **Scale of Codebase**: Managing millions of lines of code with fast build times.
- **Scale of Computation**: Utilizing multi-core CPUs efficiently (concurrency).
- **Scale of Developers**: Simple language specification that is easy to read and maintain.

It replaces C++ for systems programming where performance is critical but development speed matters. It replaces Python where performance is too slow.

### Why does Go implement it this way?
Go is a **compiled**, **statically typed**, **garbage-collected** language.
- **Compiled**: Source code is translated directly to machine code (binary). No VM (like Java JVM) or Interpreter (like Python).
- **Static Linking**: All dependencies (libraries) are included in the single binary file.
  - *Benefit*: You can copy the binary to a server and run it without installing anything.
  - *Trade-off*: Binaries are larger (e.g., 2MB for Hello World vs 10KB in C).

### Comparison
| Language | Execution Model | Performance | Deployment |
| :--- | :--- | :--- | :--- |
| **Python** | Interpreted (line by line) | Slow | Needs Python installed + `pip install` |
| **Java** | Bytecode on JVM | Fast (after warmup) | Needs JVM installed (JRE) |
| **C** | Compiled to Machine Code | Very Fast | Fast, but complex linking (DLLs/shared libs) |
| **Go** | Compiled to Machine Code | Very Fast | **Single Static Binary** (Everything included) |

---

## B. System-Level Insight

### 1. CPU Perspective
When you run a Go binary, the CPU sees native instructions (assembly).
Go manages **Goroutines** (lightweight threads). The OS sees only a few system threads (e.g., 4 threads for a 4-core CPU), but Go runs thousands of Goroutines on top of them.

### 2. Memory Perspective
- **Stack**: Used for function calls. Go stacks start small (2KB) and grow dynamically. This differs from C/Java where threads have large fixed stacks (e.g., 1MB).
- **Heap**: Used for data that outlives a function. Go has a **Garbage Collector (GC)** that pauses execution briefly to clean up unused memory.

### 3. OS Perspective
Go programs talk directly to the OS kernel via **Syscalls**.
The Go Runtime abstracts these (e.g., `fmt.Println` eventually calls the `write` syscall).

---

## C. Code Simulation

### Files Created:
1. `01_hello_world/main.go`: The minimal structure.
   - Run: `go run Day_01/01_hello_world/main.go`
2. `02_compilation_flow/main.go`: Simulates the build pipeline.
   - Run: `go run Day_01/02_compilation_flow/main.go`
3. `03_runtime_execution/main.go`: Shows the runtime in action.
   - Run: `go run Day_01/03_runtime_execution/main.go`
4. `04_syscall_lifecycle/main.go`: **[NEW]** Visualizes User Space vs Kernel Space.
   - Run: `go run Day_01/04_syscall_lifecycle/main.go`
5. `05_stack_trace_experiment/main.go`: **[NEW]** Visualizes the Call Stack and crashes.
   - Run: `go run Day_01/05_stack_trace_experiment/main.go`

**Action**: Open these files and read the comments. They explain the "WHY" behind the syntax.

---

## D. Execution Explanation

### 1. `go run main.go`
- This is a developer convenience command.
- It compiles your code to a temporary folder in `/tmp`.
- It executes the binary immediately.
- It deletes the binary after execution.
- *Use for*: Development, scripts.

### 2. `go build`
- This is for production.
- It compiles your code to a permanent executable in the current folder.
- *Use for*: Deploying to servers.

### 3. Binary Execution
When `./main` (or `main.exe`) runs:
1. OS loads the binary into memory.
2. **Runtime Initialization**:
   - The Go Runtime starts first.
   - It sets up the Garbage Collector.
   - It asks the OS for threads.
3. **Main Goroutine**:
   - The runtime creates the first goroutine for `main.main()`.
4. **Exit**:
   - When `main()` returns, the program exits (terminating all other goroutines immediately).

### 4. Practical: Build & Run (Example)
Here are the commands to compile and run your code manually:

```powershell
# 1. Build the binary from the root directory
# -o hello.exe tells Go to name the output file "hello.exe"
go build -o hello.exe Day_01/01_hello_world/main.go

# 2. Run the binary
./hello.exe
# Output: Hello, System!
```

**Common Pitfall (Path confusion)**:
If you run `go build -o hello.exe Day_01/...` from the root, the `hello.exe` file is created in the **root** folder.
If you then `cd Day_01` and try to run `./hello.exe`, it will fail because the file is in the folder above.

---

## E. Notes Summary

- **Go is "Batteries Included"**: The standard library is huge (HTTP, JSON, Crypto are built-in).
- **Go is Opinionated**: Formatting (`gofmt`), unused variables are errors, folder structure matters.
- **Go is Simple**: No inheritance, no method overloading, no generics (until recently), no exceptions (uses Errors).

### Mistakes Beginners Make:
- Thinking `main` is the first thing that runs (The Runtime is first).
- Ignoring error returns (Go forces you to handle them).
- Putting `go.mod` in the wrong place (we will cover modules later).
- Expecting Object-Oriented behavior (Go is structured, not OO).

---

## F. Advanced Simulations (Extra)

### 1. The Call Registry (`04_syscall_lifecycle.go`)
Most code you write involves the OS.
-   **User Space**: Logical operations (Math, string manipulation). Safe, fast, CPU-bound.
-   **Kernel Space**: Dangerous operations (Hardware access, Network, Files). Slow, requires permission.
-   **The Switch**: Every time you print or read a file, your program "pauses", the Kernel wakes up, does the work, and wakes your program up again. This switch is expensive!

### 2. The Stack Visualizer (`05_stack_trace_experiment.go`)
This file proves that functions "pile up".
-   When `main` calls `level1`, `main` doesn't vanish. It waits.
-   Running this shows the **Stack Trace**—a snapshot of this frozen memory.
-   **Debugging Tip**: When a Go program crashes, it prints this Trace. Read it from TOP (where it crashed) to BOTTOM (how it got there).

---
**Next Step**: proceed to Day 2 to learn about Variables and Memory Layout.
