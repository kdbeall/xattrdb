package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kdbeall/xattrdb/internal/server"
)

func main() {
	server.SetPath("/tmp/xattrdb")
	server.SetShards(2)
	server.CreateShards()
	router := mux.NewRouter()
	router.HandleFunc("/xattrdb", server.ServerCreate).Methods("POST")
	router.HandleFunc("/xattrdb/{key}", server.ServerRead).Methods("GET")
	router.HandleFunc("/xattrdb/{key}", server.ServerUpdate).Methods("PUT")
	router.HandleFunc("/xattrdb/{key}", server.ServerDelete).Methods("DELETE")
	router.HandleFunc("/xattrdb/snapshots", server.ServerCreateSnapshot).Methods("POST")
	router.HandleFunc("/xattrdb/snapshots/{snapshot}/{key}", server.ServerReadSnapshot).Methods("GET")
	http.ListenAndServe(":8000", router)
}
