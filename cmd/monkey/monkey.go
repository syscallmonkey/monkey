package main

import "fmt"

var (
	Version, Build string
)

func main() {
	fmt.Printf("Version %s, build %s\n", Version, Build)
	fmt.Println("Hello, world.")
}