package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
)

func decodeBody(r io.Reader) ([]*Config, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt []*Config
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return rt, nil
}

func decodeConfigurationBody(r io.Reader) (*Config, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt Config
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func decodeGroupBody(r io.Reader) (*Group, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt Group
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func createId() string {
	return uuid.New().String()
}

const (
	configurationsId     = "configs/%s"
	configurationVersion = "configs/%s/%s"
	allConfigs           = "config"

	allGroups     = "group"
	groupsId      = "groups/%s"
	groupsVersion = "groups/%s/%s"
	groupsLabel   = "groups/%s/%s/%s/"
)

func generateKeyConfiguration(verzija string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(configurationVersion, id, verzija), id
}

func constructKeyIDConfigurations(id string) string {
	return fmt.Sprintf(configurationsId, id)
}

func constructKeyVersionConfigs(id string, verzija string) string {
	return fmt.Sprintf(configurationVersion, id, verzija)
}

//Konfiguracije
///

///Grupe
func generateKeyGroup(verzija string, labele string) (string, string) { /////Proveri ovde da li ide mapa stringova ili samo string
	id := uuid.New().String()
	return fmt.Sprintf(groupsLabel, id, verzija, labele), id
}

func constructKeyGroupId(id string) string {
	return fmt.Sprintf(groupsId, id)
}

func constructKeyGroupVersion(id string, verzija string) string {
	return fmt.Sprintf(groupsVersion, id, verzija)
}

func constructKeyGroupLabels(id string, verzija string, labela string) string {
	return fmt.Sprintf(groupsLabel, id, verzija, labela)
}
