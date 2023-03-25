package main

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/compute/metadata"
	"github.com/enescakir/emoji"
)

type Payload struct {
	PodName        string `json:"pod_name"`
	Timestamp      string `json:"timestamp"`
	Zone           string `json:"zone,omitempty"`
	ProjectId      string `json:"project_id,omitempty"`
	InstanceId     string `json:"gce_instance_id,omitempty"`
	ServiceAccount string `json:"gce_service_account,omitempty"`
	PodNameEmoji   string `json:"pod_name_emoji"`
}

// laziness and creating a global payload
var payload Payload

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
	//log.Printf(projectId + "\n")
	if gceErr != nil {
		log.Println("Unable to capture GCE metadata")

	} else {
		payload.ProjectId = projectId
		payload.Zone, _ = metadata.Zone()
		payload.InstanceId, _ = metadata.InstanceID()
		payload.ServiceAccount, _ = metadata.InstanceAttributeValue("serviceAccounts/default/email") // doesn't work for now
	}
	return payload
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	// else condition
	log.Printf("Environment variable %s not found\n", key)
	return fallback
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	log.Printf("got / request\n")

	// update timestamp
	payload.Timestamp = time.Now().UTC().String()
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
	// initialize payload
	payload = generatePayload()
	http.HandleFunc("/", getRoot)

	err := http.ListenAndServe(":"+getEnv("PORT", "8080"), nil)

	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("server closed\n")
	} else if err != nil {
		log.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
