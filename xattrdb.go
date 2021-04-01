package xattrdb

import (
	"net/http"

	mux "github.com/gorilla/mux"
)

func Init(path string, shards int) {
	SetPath(path)
	SetShards(shards)
	CreateShards()
	router := mux.NewRouter()
	router.HandleFunc("/xattrdb", ServerCreate).Methods("POST")
	router.HandleFunc("/xattrdb/{key}", ServerRead).Methods("GET")
	router.HandleFunc("/xattrdb/{key}", ServerUpdate).Methods("PUT")
	router.HandleFunc("/xattrdb/{key}", ServerDelete).Methods("DELETE")
	http.ListenAndServe(":8000", router)
}
