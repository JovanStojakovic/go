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
	//Dodaje novu konfiguraciju
	router.HandleFunc("/config", server.createConfigurationHandler).Methods("POST")
	//Dodaje novu grupu
	router.HandleFunc("/group", server.createGroupHandler).Methods("POST")
	//Dodaje novu verziju postojece konfiguracije
	router.HandleFunc("/config/{id}", server.addConfigVersionHandler).Methods("POST")
	//Dodaje novu verziju postojece grupe
	router.HandleFunc("/group/{id}/{version}", server.addNewGroupVersionHandler).Methods("POST")
	////
	////Post

	////
	///
	////GET

	//Nadji sve konfiguracije
	router.HandleFunc("/config", server.getAllConfigurationsHandler).Methods("GET")

	//Nadji konfiguraciju po id-u
	router.HandleFunc("/config/{id}", server.getConfigByIDHandler).Methods("GET")

	//Nadji konfiguraciju po id-u i verziji
	router.HandleFunc("/config/{id}/{version}", server.getConfigByIDVersionHandler).Methods("GET")

	//Nadji sve grupe
	router.HandleFunc("/group", server.getAllGroupHandler).Methods("GET")

	//Nadji grupu po id-u i verziji
	router.HandleFunc("/group/{id}/{version}", server.getGroupByIdVersionHandler).Methods("GET")

	//Nadji grupu po id-u
	router.HandleFunc("/group/{id}", server.getGroupByIdHandler).Methods("GET")

	///Nadji konfiguraciju unutar grupe preko labela - nedovrseno
	router.HandleFunc("/group/{id}/{version}/{label}", server.getGroupLabelHandler).Methods("GET") /// Ovde je problem

	///
	///
	///

	///
	///Delete
	//Obrisi grupu
	router.HandleFunc("/group/{id}/{version}", server.delGroupHandler).Methods("DELETE")

	//Obrisi konfiguraciju
	router.HandleFunc("/config/{id}/{version}", server.delConfigurationHandler).Methods("DELETE")
	//
	//PUT
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
