package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Payload struct {
	PodName   string `json:"pod_name"`
	Timestamp string `json:"timestamp"`
	TestValue string `json:"test_value,omitempty"`
}

/*func generatePayload() *Payload {
	p := Payload{}
	p.pod_name = "test"
	p.timestamp = "test2"
	return &p
}*/

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	hostname, _ := os.Hostname() // sloppy but ignoring error

	//payload := &Payload
	payload := Payload{PodName: hostname, Timestamp: time.Now().UTC().String()}
	//p, _ := json.Marshal(&generatePayload())
	w.Header().Set("Content-Type", "application/json")
	/*p, err := json.Marshal(Payload{PodName: "test", Timestamp: "test2"})
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	io.WriteString(w, string(p))*/
	json.NewEncoder(w).Encode(payload)
}

func main() {
	http.HandleFunc("/", getRoot)

	err := http.ListenAndServe(":"+getEnv("PORT", "8080"), nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
