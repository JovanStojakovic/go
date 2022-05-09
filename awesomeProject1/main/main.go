package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.StrictSlash(true)

	server := Service{
		data: map[string][]*Config{},
	}
	router.HandleFunc("/post/", server.createPostHandler).Methods("POST")
	router.HandleFunc("/posts/", server.getAllHandler).Methods("GET")
	router.HandleFunc("/post/{id}/", server.getPostHandler).Methods("GET")
	router.HandleFunc("/post/{id}", server.delPostHandler).Methods("DELETE")

	// start server
	srv := &http.Server{Addr: "0.0.0.0:8080", Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

}
