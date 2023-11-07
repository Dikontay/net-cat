package main

import (
	"fmt"
	"sync"
)

func main() {
	// Create a channel for sending integers.
	ch := make(chan int)

	// Create a WaitGroup to wait for goroutines to finish.
	var wg sync.WaitGroup

	// Start four goroutines.
	for i := 0; i < 4; i++ {
		wg.Add(1) // Increment the WaitGroup counter.
		go func(i int) {
			defer wg.Done() // Decrement the counter when the goroutine completes.
			ch <- i         // Send 'i' to the channel.
		}(i)
	}

	// Start a goroutine to close the channel once all sends are done.
	go func() {
		wg.Wait() // Wait for all goroutines to finish.
		close(ch) // Close the channel.
	}()

	// Receive values from the channel and print them.
	for val := range ch {
		fmt.Println(val)
	}
}
