// Go program to illustrate how to
// find the capacity of the channel

package main

import "fmt"

// Main function
func main() {

	// Creating a channel
	// Using make() function
	ch := make(chan string, 5)
	ch <- "GFG"
	ch <- "gfg"
	ch <- "Geeks"
	ch <- "Geeks for Geeks"

	// Finding the capacity of the channel
	// Using cap() function
	fmt.Printf("Capacity of the channel is: %d, Length ofo the channel is : %d .", cap(ch), len(ch))
}
