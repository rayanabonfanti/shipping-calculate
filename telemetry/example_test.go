package telemetry

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

func ExampleRecordLatencyOperationA() {
	now := time.Now()
	// operation to measure
	func() {
		fmt.Println("doing operation...")
	}()

	elapsed := time.Since(now).Milliseconds()

	RecordLatencyOperationA(context.Background(), elapsed, "local")
}

func ExampleRecordMemoryHeapServer() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	RecordMemoryHeapServer(context.Background(), int64(m.HeapAlloc/1024/1024))
}

func ExampleRecordMemoryNoHeapServer() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	RecordMemoryNoHeapServer(context.Background(), int64(m.StackInuse/1024/1024))
}

func ExampleIncrementHttpRequestHandled() {
	// custom server router
	router := http.NewServeMux()

	router.HandleFunc("/GET test", func(w http.ResponseWriter, r *http.Request) {
		IncrementHttpRequestHandled(r.Context(), "GET", http.StatusOK)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	router.HandleFunc("/POST test", func(w http.ResponseWriter, r *http.Request) {
		IncrementHttpRequestHandled(r.Context(), "POST", http.StatusCreated)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("created"))
	})

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
