package xattrdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	mux "github.com/gorilla/mux"
)

type Value struct {
	Value string `json:"value"`
}

type Snapshot struct {
	Snapshot string `json:"snapshot"`
}

func ServerCreate(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintf(writer, "Failed to create.")
		return
	}
	requestContent := make(map[string]string)
	json.Unmarshal(body, &requestContent)
	key, value := requestContent["key"], requestContent["value"]
	if key == "" || value == "" || !CreateData(key, value) {
		fmt.Fprintf(writer, "Failed to create.")
	}
}

func ServerRead(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	value, err := ReadData(params["key"])
	if err != nil {
		fmt.Fprintf(writer, "Failed to read.")
		return
	}
	json.NewEncoder(writer).Encode(Value{value})
}

func ServerUpdate(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	key := params["key"]
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintf(writer, "Failed to update.")
		return
	}
	requestContent := make(map[string]string)
	json.Unmarshal(body, &requestContent)
	value := requestContent["value"]
	if value == "" || !UpdateData(key, value) {
		fmt.Fprintf(writer, "Failed to create.")
	}
}

func ServerDelete(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	if !DeleteData(params["key"]) {
		fmt.Fprintf(writer, "Failed to delete.")
	}
}

func ServerCreateSnapshot(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	snapshot := CreateSnapshot()
	json.NewEncoder(writer).Encode(Snapshot{snapshot})

}

func ServerReadSnapshot(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	snapshot := params["snapshot"]
	key := params["key"]
	value, err := ReadSnapshot(key, snapshot)
	if err != nil {
		fmt.Fprintf(writer, "Failed to read.")
		return
	}
	json.NewEncoder(writer).Encode(Value{value})
}
