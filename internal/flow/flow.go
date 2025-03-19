package flow

import (
	"fmt"
	"time"
)

func Info(message string) {
	fmt.Printf("→ %s\n", message)
}

func Warn(message string) {
	fmt.Printf("! %s\n", message)
}

func Error(message string) {
	fmt.Printf("✗ %s\n", message)
}

func Success(message string) {
	fmt.Printf("✓ %s\n", message)
}

func Action(message string) {
	fmt.Printf("$ %s\n", message)
}

func Start() {
	fmt.Println("lazyenv starting")
}

func Done() {
	fmt.Println("\nall done")
}

func Section(title string) {
	fmt.Printf("\n%s\n", title)
}

func Progress(message string) {
	t := time.Now().Format("15:04:05")
	fmt.Printf("[%s] %s\n", t, message)
}

func FileAction(action, path string) {
	fmt.Printf("  %s: %s\n", action, path)
}
