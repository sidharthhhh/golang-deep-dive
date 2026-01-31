# Day 2: Variables, Types, and Memory Layout

## 1. Variables and Declaration

Go is statically typed, meaning variables have a fixed type known at compile time.

### Declaration Styles

1.  **Standard Declaration (`var`)**:
    ```go
    var age int = 30
    var name string // Defaults to "" (Zero Value)
    ```
    Use this when you need to declare a variable without an initial value (to use later) or inside `var()` blocks.

2.  **Short Declaration (`:=`)**:
    ```go
    count := 10 // Type inferred as int
    ```
    **Rule**: Can ONLY be used inside functions.
    **Gotcha**: It creates a *new* variable. Be careful of shadowing!

### Zero Values (Crucial Concept)
Go **never** leaves memory uninitialized. If you don't assign a value, it gets the "Zero Value":
- `int` -> `0`
- `float` -> `0.0`
- `bool` -> `false`
- `string` -> `""` (Empty string)
- `pointer`, `interface`, `slice`, `map`, `channel` -> `nil`

---

## 2. The Type System

Go is **Strongly Typed**. It does not do implicit type conversion (magic).

```go
var a int = 10
var b int64 = 10
// a == b // COMPILE ERROR! Types are different.
```

### Type Conversion
You must explicitly convert types:
```go
if a == int(b) { ... }
```
Note: This is "Conversion", not casting. It creates a new value of the target type.

### Custom Types
You can define your own types based on primitives:
```go
type UserID int
```
`UserID` is a distinct type from `int`. You cannot mix them without conversion. This adds safety.

---

## 3. Memory Layout (Intro)

### RAM Model
Think of memory as a massive array of bytes, each with a unique address (hexadecimal number like `0xc0000140b0`).

### Stack vs Heap
- **Stack**: Fast, structured, automatic cleanup. Functions store their local variables here. When function returns, stack frame is popped (memory "freed" instantly).
- **Heap**: Slower, unstructured, manual cleanup (handled by Garbage Collector in Go). Used for data that needs to live *allocated* longer than the function that created it (Escaping).

### Variable Addresses
Every variable has a location in memory.
- `&variable` -> Gives you the memory address (pointer).
