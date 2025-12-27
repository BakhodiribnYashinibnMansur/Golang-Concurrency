// Go program to illustrate how to
// find the capacity of the channel

package main

import (
	"fmt"
	"time"
)

// Main function
func main() {

	// Creating a channel
	// Using make() function
	ch := make(chan string)
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println(<-ch)
		fmt.Print(<-ch)
	}()

	ch <- "GFG"
	ch <- "WTF"
	fmt.Printf("\n Capacity of the channel is: %d, Length ofo the channel is : %d .", cap(ch), len(ch))

	fmt.Printf("\n Capacity of the channel is: %d, Length ofo the channel is : %d .", cap(ch), len(ch))
}
