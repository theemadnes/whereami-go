package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/compute/metadata"
	"github.com/enescakir/emoji"
)

type Payload struct {
	PodName        string                 `json:"pod_name"`
	Timestamp      string                 `json:"timestamp"`
	Zone           string                 `json:"zone,omitempty"`
	ProjectId      string                 `json:"project_id,omitempty"`
	InstanceId     string                 `json:"gce_instance_id,omitempty"`
	ServiceAccount string                 `json:"gce_service_account,omitempty"`
	PodNameEmoji   string                 `json:"pod_name_emoji"`
	BackendResult  map[string]interface{} `json:"backend_result,omitempty"`
	ClusterName    string                 `json:"cluster_name,omitempty"`
}

// laziness and creating a global payload
var payload Payload

// headers to propagate
var headersToPropagate = []string{
	"x-request-id",
	"x-b3-traceid",
	"x-b3-spanid",
	"x-b3-parentspanid",
	"x-b3-sampled",
	"x-b3-flags",
	"x-ot-span-context",
	"x-cloud-trace-context",
	"traceparent",
	"grpc-trace-bin",
}

// pick a random value from a map (used for emoji assignment)
// using https://programming-idioms.org/idiom/250/pick-a-random-value-from-a-map/4435/go
func pick(m map[string]string) string {
	k := rand.Intn(len(m))
	for _, x := range m {
		if k == 0 {
			return x
		}
		k--
	}
	panic("unreachable")
}

func generatePayload() Payload {

	//p := Payload{}
	payload.PodName, _ = os.Hostname()
	payload.Timestamp = time.Now().UTC().String()
	payload.PodNameEmoji = pick(emoji.Map())
	//gceMetadataClient := metadata.NewClient(&http.Client{})

	//projectId, gceErr := gceMetadataClient.ProjectID()
	projectId, gceErr := metadata.ProjectID()
	//fmt.Printf(projectId + "\n")
	if gceErr != nil {
		fmt.Println("Unable to capture GCE metadata")

	} else {
		payload.ProjectId = projectId
		payload.Zone, _ = metadata.Zone()
		payload.InstanceId, _ = metadata.InstanceID()
		payload.ServiceAccount, _ = metadata.InstanceAttributeValue("serviceAccounts/default/email") // doesn't work for now
		payload.ClusterName, _ = metadata.InstanceAttributeValue("cluster-name")
	}
	return payload
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	// else condition
	fmt.Printf("Environment variable %s not found\n", key)
	return fallback
}

func contains(s []string, str string) bool {
	for _, v := range s {
		// using EqualFold to make string comparison case-insensitive
		if strings.EqualFold(v, str) {
			return true
		}
	}

	return false
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")

	// update timestamp
	payload.Timestamp = time.Now().UTC().String()

	// check for backend service call
	if getEnv("BACKEND_ENABLED", "") == "True" {
		if backendUrl, ok := os.LookupEnv("BACKEND_SERVICE"); ok {
			client := &http.Client{}
			req, err := http.NewRequest("GET", backendUrl, nil)
			if err != nil {
				panic(err)
			}
			// populate headers to request
			/*for _, k := range headersToPropagate {
				if r.Header.Values(k) {
					req.Header.Add(k, v)
				}
			}*/
			for k, v := range r.Header {
				//fmt.Printf("checking %s %s\n", k, v)
				//fmt.Printf("test\n")
				if contains(headersToPropagate, k) {
					//req.Header.Add(k, r.Header.Values(k))
					for _, vOther := range v {
						req.Header.Add(k, vOther)
						//fmt.Printf("%s %s\n", k, vOther)
					}
				}
			}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Call to %s failed", backendUrl)
				defer resp.Body.Close()
			} else {
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("Invalid response from %s", backendUrl)
				} else {
					var jsonRes map[string]interface{}
					err = json.Unmarshal(body, &jsonRes)
					if err != nil {
						fmt.Printf("Unable to unmarshal response from %s", backendUrl)
					} else {
						payload.BackendResult = jsonRes
					}
					//payload.BackendResult = jsonRes.Marshal()
					//marRes, _ := json.Marshal(jsonRes)
					//fmt.Println(string(marRes))

				}
			}
		}
	}
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
	// initialize payload and hack to wait for WI to come online
	time.Sleep(5 * time.Second)
	payload = generatePayload()
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	http.HandleFunc("/", getRoot)

	err := http.ListenAndServe(":"+getEnv("PORT", "8080"), nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
