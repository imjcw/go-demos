package main

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// Serve the index.html file at the root path
	router.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})
	router.GET("/events", func(c *gin.Context) {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		clientChan := make(chan string)

		go func() {
			for i := 0; ; i++ {
				message := fmt.Sprintf("Message %d at %s", i, time.Now().Format(time.RFC3339))
				clientChan <- message
				time.Sleep(1 * time.Second)
			}
		}()

		c.Stream(func(w io.Writer) bool {
			msg, open := <-clientChan
			if !open {
				close(clientChan)
				return false
			}
			fmt.Fprintf(w, "data: %s\n\n", msg)
			return true
		})
	})

	router.Run(":8080")
}
