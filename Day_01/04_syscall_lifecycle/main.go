package main

import (
	"fmt"
	"os"
	"time"
)

// 04_syscall_lifecycle.go
// -----------------------------------------------------------------------------
// This file demonstrates the interactions between YOUR CODE and the OS KERNEL.
//
// MENTAL MODEL:
// Your program runs in "User Space" (restricted access).
// To do ANYTHING real (write file, network, screen, sleep), it must ask the OS.
// These requests are called "SYSCALLS" (System Calls).
//
// 1. CPU Work -> Happens entirely in your program (User Space).
// 2. I/O Work -> Happens by switching to the OS (Kernel Space).
// -----------------------------------------------------------------------------

func main() {
	fmt.Println("=== Syscall & OS Interaction Simulation ===")

	// 1. Pure User Space (CPU Binding)
	// The OS doesn't care about this. It's just math.
	x := 10 * 10
	fmt.Printf("1. Math calculated in User Space: %d\n", x)

	// 2. Getting System Info (Syscall: getpid)
	// We ask the Kernel: "Who am I?"
	pid := os.Getpid()
	fmt.Printf("2. Syscall (GetPID): My Process ID is %d\n", pid)

	// 3. Writing to Console (Syscall: write)
	// 'fmt.Println' is wrappers around the 'write' syscall.
	// We are asking the OS to send bytes to the standard output device (screen).
	os.Stdout.WriteString("3. Syscall (Write): Hello from Kernel Space!\n")

	// 4. Time Management (Syscall: nanosleep/Sleep)
	// We ask the Kernel: "Please take the CPU away from me for 1 second."
	fmt.Println("4. Syscall (Sleep): Asking OS to suspend me...")
	time.Sleep(1 * time.Second)
	fmt.Println("   ...I am back!")
}
