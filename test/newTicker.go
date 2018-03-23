package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second)
	done := make(chan bool)
	go func() {
		time.Sleep(5 * time.Second)
		done <- true
	}()
	for {
		select {
		case <-done:
			fmt.Println("Done!")
			ticker.Stop()
			return
		case t := <-ticker.C:
			fmt.Println("Current time: ", t)
		}
	}
}
