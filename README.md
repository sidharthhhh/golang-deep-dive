# Golang Deep Dive Course 🚀

Welcome to the **Golang Deep Dive**. This is a structured, 20-day journey to master Go from a system-level perspective. We focus on **WHY** things work, not just syntax.

## 🎯 Learning Philosophy
- **Notes-First**: Understanding mental models before writing code.
- **Simulation**: Every concept is proven with code.
- **System-Level**: Deep dives into Memory, CPU, Scheduler, and OS interaction.
- **Production-Grade**: Best practices for reliability, performance, and maintainability.

---

## 📅 Curriculum Roadmap

### **LEVEL 1 — Foundations**
- [x] **Day 1: How Go Programs Are Built and Executed**
  - [x] Concept: Compilation flow, static linking, runtime initiation.
  - [x] System Insight: OS interface, syscalls.
  - [x] Code: Hello World, Compilation Simulation, Runtime introspection (Refactored to subdirectories).
- [ ] **Day 2: Variables, Types, and Memory Layout**
  - Concept: Zero values, type safety, memory allocation.
  - System Insight: Stack vs Heap allocation, memory alignment.
  - Code: Variable shadowing, type conversions, memory address inspection.
- [ ] **Day 3: Control Flow and Functions (Stack Frames)**
  - Concept: `if`, `for`, `switch`, function calls.
  - System Insight: Stack frames, return addresses, defer execution order.
  - Code: Stack trace simulation, recursion depth.
- [ ] **Day 4: Arrays, Slices, and Maps (Internals)**
  - Concept: Data structures in Go.
  - System Insight: Slice headers (ptr, len, cap), map buckets chains.
  - Code: Slice capacity growth, map collision simulation.

### **LEVEL 2 — Core Go**
- [ ] **Day 5: Structs, Methods, and Composition**
  - Concept: Avoiding inheritance, favoring composition.
  - System Insight: Memory layout of structs, method receivers (value vs pointer).
  - Code: Composition patterns, embedding vs inheritance.
- [ ] **Day 6: Pointers and Memory Management**
  - Concept: References, nil pointers, escape analysis.
  - System Insight: When variables escape to heap, GC pressure.
  - Code: Pointer arithmetic (unsafe), escape analysis demo.
- [ ] **Day 7: Interfaces and Polymorphism**
  - Concept: Duck typing, `interface{}`, type assertions.
  - System Insight: Interface tables (itable), dynamic dispatch cost.
  - Code: Dependency injection, mocking with interfaces.
- [ ] **Day 8: Error Handling and Defer**
  - Concept: Errors as values, `defer`, `panic`, `recover`.
  - System Insight: Stack unwinding involving panic.
  - Code: Custom error types, reliable resource cleanup.
- [ ] **Day 9: Packages and Modules**
  - Concept: `go.mod`, dependency management, visibility rules.
  - System Insight: Linker behavior, dead code stripping.
  - Code: Creating a reusable module, internal packages.

### **LEVEL 3 — Concurrency & Runtime**
- [ ] **Day 10: Goroutines and The Scheduler**
  - Concept: Lightweight threads, `go` keyword.
  - System Insight: M:N Scheduler (M: OS Thread, P: Processor, G: Goroutine).
  - Code: Spawning thousands of goroutines, context switching cost.
- [ ] **Day 11: Channels and Communication**
  - Concept: `chan`, buffered vs unbuffered, `select`.
  - System Insight: Channel blocking mechanisms, lock-free queues basics.
  - Code: Worker pools, pipeline patterns.
- [ ] **Day 12: Sync Package (Mutex, WaitGroup, Atomic)**
  - Concept: Shared memory synchronization.
  - System Insight: CPU caches, memory barriers, race detector.
  - Code: Fixing race conditions, atomic counters.
- [ ] **Day 13: Context and Cancellation**
  - Concept: `context.Context`, timeouts, cancellation propagation.
  - System Insight: Request-scoped data lifecycle.
  - Code: Canceling long-running operations, strict timeouts.

### **LEVEL 4 — Backend & APIs**
- [ ] **Day 14: HTTP Internals**
  - Concept: `net/http` standard library.
  - System Insight: TCP, connection reuse (Keep-Alive), request parsing.
  - Code: Raw TCP server vs HTTP server.
- [ ] **Day 15: Building a REST API**
  - Concept: Routing, Handlers, Middleware.
  - System Insight: Middleware chains as function wrappers.
  - Code: Production-ready CRUD API structure.
- [ ] **Day 16: JSON and Serialization**
  - Concept: `encoding/json`, struct tags.
  - System Insight: Reflection overhead in serialization.
  - Code: Custom JSON marshaling, streaming decoders.
- [ ] **Day 17: Database/SQL and Connection Pools** (Moved from L5)
  - Concept: `database/sql` interface, drivers.
  - System Insight: Connection pooling strategy, prepared statements.
  - Code: SQL injection prevention, transaction management.

### **LEVEL 5 — Persistence & Reliability** (Refined from L5)
- [ ] **Day 18: Testing and Benchmarking** (Moved from L6)
  - Concept: `testing` package, table-driven tests.
  - System Insight: Benchmark memory allocations.
  - Code: Unit tests, fuzzing, performance benchmarks.

### **LEVEL 6 — Production Engineering** (Refined from L6)
- [ ] **Day 19: Profiling and Optimization**
  - Concept: PProf (CPU, Heap, Block profiles).
  - System Insight: Flame graphs, GC pause analysis.
  - Code: Finding and fixing a memory leak.

### **LEVEL 7 — Advanced & Modern Backend**
- [ ] **Day 20: Advanced Patterns (gRPC, Observability)**
  - Concept: gRPC vs REST, structured logging, metrics.
  - System Insight: Protobuf binary efficiency, distributed tracing.
  - Code: Simple gRPC service with OpenTelemetry.

---

## 🛠 Usage
1.  **Clone the Repo**: This is your personal workspace.
2.  **Follow Daily**: Open the folder for the day (e.g., `Day_01`).
3.  **Read NOTES.md**: Understand the mental model first.
4.  **Run Code**: Execute the simulated examples (`go run ...`).
5.  **Modify**: Experiment with the code to verify your understanding.

Let's build reliable systems!
