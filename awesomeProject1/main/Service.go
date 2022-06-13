package main

import (
	"errors"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
	"sort"
	"strings"
)

type Service struct {
	store *ConfigurationStore
}

//Pravi jednu konfiguraciju
func (ts *Service) createConfigurationHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeConfigurationBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := ts.store.PostConfig(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, post)
}

///Vraca sve konfiguracije
func (ts *Service) getAllConfigurationsHandler(w http.ResponseWriter, req *http.Request) {
	allTasks, err := ts.store.GetAllConfigurations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, allTasks)
}

///Vraca konfiguraciju preko id-a
func (ts *Service) getConfigByIDHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, ok := ts.store.GetConfigurationById(id)
	if ok != nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

//Vraca konfiguraciju preko id-a i verzije
func (ts *Service) getConfigByIDVersionHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	verzija := mux.Vars(req)["version"]
	task, ok := ts.store.GetConfigByIdVersion(id, verzija)
	if ok != nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

//Brise konfiguraciju
func (ts *Service) delConfigurationHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	verzija := mux.Vars(req)["version"]
	konf, err := ts.store.DeleteConfig(id, verzija)
	if err != nil {
		errors.New("konfiguracija nije pronadjena, stim nije obrisana!")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, konf)
}

///Dodaje novu verziju konfiguracije
func (ts *Service) addConfigVersionHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	rt, err := decodeConfigurationBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := mux.Vars(req)["id"]
	rt.Id = id
	config, err := ts.store.AddNewConfigVersion(rt)
	if err != nil {
		http.Error(w, "Ta verzija vec postoji!", http.StatusBadRequest)
	}
	renderJSON(w, config)

}

////
////Konfiguracije
////

///Grupe

///Vraca sve grupe
func (ts *Service) getAllGroupHandler(w http.ResponseWriter, req *http.Request) {
	allTasks, err := ts.store.GetAllGroups()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, allTasks)
}

///Pravi jednu grupu
func (ts *Service) createGroupHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeGroupBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	group, err := ts.store.PostGroup(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, group)
}

///Brisanje grupe
func (ts *Service) delGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	group, err := ts.store.DeleteGroup(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, group)
}

//Nadji grupu preko id-a
func (ts *Service) getGroupByIdHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	group, err := ts.store.GetGroupById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, group)
}

//Nadji grupu preko id-a i verzije
func (ts *Service) getGroupByIdVersionHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	verzija := mux.Vars(req)["version"]
	group, ok := ts.store.GetGroupByIdVersion(id, verzija)
	if ok != nil {
		err := errors.New("Nije pronadjeno")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, group)
}

///dodaje novu verziju
func (cs *Service) addNewGroupVersionHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	rt, err := decodeGroupBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	rt.Id = id
	rt.Version = version
	group, err := cs.store.AddNewGroupVersion(rt)
	renderJSON(w, group)

}

//Dodaje novu konfiguraciju u grupu
func (ts *Service) UpdateGroupWithNewHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	_, err := ts.store.DeleteGroup(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeGroupBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id2 := mux.Vars(req)["id"]
	version2 := mux.Vars(req)["version"]
	rt.Id = id2
	rt.Version = version2

	nova, err := ts.store.UpdateGroup(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, nova)
}

///Proveri ovde
////Nadji grupu preko labela
func (ts *Service) getGroupLabelHandler(w http.ResponseWriter, req *http.Request) {

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	label := mux.Vars(req)["label"]
	list := strings.Split(label, ";")
	sort.Strings(list)
	sortedLabel := ""
	for _, v := range list {
		sortedLabel += v + ";"
	}
	sortedLabel = sortedLabel[:len(sortedLabel)-1]
	returnConfigs, error := ts.store.GetGroupByLabel(id, version, sortedLabel)

	if error != nil {
		renderJSON(w, "Doslo je do greske!Nije pronadjena")
	}
	renderJSON(w, returnConfigs)
}
