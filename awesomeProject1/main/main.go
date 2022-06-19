package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.StrictSlash(true)

	store, err := New()
	if err != nil {
		log.Fatal(err)
	}

	server := Service{
		store: store,
	}

	router.HandleFunc("/group/{id}/{version}", server.delGroupHandler).Methods("DELETE")

	router.HandleFunc("/config/{id}/{version}", server.delConfigurationHandler).Methods("DELETE")

	router.HandleFunc("/config", server.getAllConfigurationsHandler).Methods("GET")

	router.HandleFunc("/config/{id}", server.getConfigByIDHandler).Methods("GET")

	router.HandleFunc("/config/{id}/{version}", server.getConfigByIDVersionHandler).Methods("GET")

	router.HandleFunc("/group", server.getAllGroupHandler).Methods("GET")

	router.HandleFunc("/group/{id}/{version}", server.getGroupByIdVersionHandler).Methods("GET")

	router.HandleFunc("/group/{id}", server.getGroupByIdHandler).Methods("GET")

	router.HandleFunc("/group/{id}/{version}/{label}", server.getGroupLabelHandler).Methods("GET")

	router.HandleFunc("/config", server.createConfigurationHandler).Methods("POST")

	router.HandleFunc("/group", server.createGroupHandler).Methods("POST")

	router.HandleFunc("/config/{id}", server.addConfigVersionHandler).Methods("POST")

	router.HandleFunc("/group/{id}/{version}", server.addNewGroupVersionHandler).Methods("POST")

	router.HandleFunc("/group/{id}/{version}", server.UpdateGroupWithNewHandler).Methods("PUT")

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

	<-quit

	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
