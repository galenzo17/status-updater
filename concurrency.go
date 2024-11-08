package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

var (
	activeConnections int64
	connectionLock    sync.Mutex
	wg                sync.WaitGroup
)

// incrementConnections increments the active connection count in a thread-safe manner
func incrementConnections() {
	connectionLock.Lock()
	defer connectionLock.Unlock()
	activeConnections++
}

// decrementConnections decrements the active connection count in a thread-safe manner
func decrementConnections() {
	connectionLock.Lock()
	defer connectionLock.Unlock()
	activeConnections--
}

// handleRequest simulates handling a request with random processing time
func handleRequest(w http.ResponseWriter, r *http.Request) {
	defer wg.Done()
	incrementConnections()
	defer decrementConnections()

	// Simulate processing time
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	// Print active connections
	fmt.Fprintf(w, "Active connections: %d\n", activeConnections)
}

// memoryMonitor logs the memory usage at intervals
func memoryMonitor() {
	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		log.Printf("Alloc = %v MiB, TotalAlloc = %v MiB, Sys = %v MiB, NumGC = %v\n",
			m.Alloc/1024/1024, m.TotalAlloc/1024/1024, m.Sys/1024/1024, m.NumGC)
		time.Sleep(5 * time.Second)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Start memory monitoring in a separate goroutine
	go memoryMonitor()

	// Define and start the HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wg.Add(1)
		go handleRequest(w, r)
	})

	port := ":8080"
	fmt.Printf("Server starting on port %s...\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}

	// Wait for all requests to finish before exiting
	wg.Wait()
	fmt.Println("All connections have been handled.")
}