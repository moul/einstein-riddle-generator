package einsteinsriddleapp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/moul/einstein-riddle-generator"
)

func init() {
	http.HandleFunc("/", handler)
}

type RiddleRequest struct {
	Options map[string]bool `json:"Options,omitempty"`
}

type RiddleResponse struct {
	Facts     []string `json:"facts"`
	Questions []string `json:"questions"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var riddleRequest RiddleRequest
	err := decoder.Decode(&riddleRequest)
	if err != nil {
		// FIXME
	}
	//if err != nil {
	//panic(err)
	//}

	options := einsteinriddle.Options{}
	/*
		if RiddleRequest.Options["Size"] {
			options.Size = RiddleRequest.Options["Size"]
		}
	*/
	generator := einsteinriddle.NewGenerator(options)

	// Shazam
	err = generator.Shazam()
	//if err != nil {
	//	panic(err)
	//}

	// Print map
	//generator.Show()

	// Print riddle
	var response RiddleResponse
	response.Facts = make([]string, 0)
	response.Questions = make([]string, 0)
	for _, group := range generator.Pickeds {
		response.Facts = append(response.Facts, generator.GroupString(group))
	}
	for _, item := range generator.Missings() {
		response.Questions = append(response.Questions, fmt.Sprintf("Where is %s ?", item.Name()))
	}
	b, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "applicaton/json")
	w.Write(b)
}
