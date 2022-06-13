package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"os"
	"sort"
)

type ConfigurationStore struct {
	cli *api.Client
}

func New() (*ConfigurationStore, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConfigurationStore{
		cli: client,
	}, nil
}

///Nadji sve konfiguracije --- ok
func (ps *ConfigurationStore) GetAllConfigurations() ([]*Config, error) {
	kv := ps.cli.KV()
	data, _, err := kv.List(allConfigs, nil)
	if err != nil {
		return nil, err
	}

	configurations := []*Config{}
	for _, pair := range data {
		post := &Config{}
		err = json.Unmarshal(pair.Value, post)
		if err != nil {
			return nil, err
		}
		configurations = append(configurations, post)
	}

	return configurations, nil
}

///Nadji konfiguraciju preko id-a i verzije ---ok
func (ps *ConfigurationStore) GetConfigByIdVersion(id string, verzija string) (*Config, error) {
	kv := ps.cli.KV()
	//sid := constructKeyVersionConfigs(id, verzija)
	//fmt.Println("OVO JE ONO STO TRAZIMO: " + sid)
	pair, _, err := kv.Get(constructKeyVersionConfigs(id, verzija), nil)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(pair.Value, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

///Nadji sve verzije jedne konfiguracije preko id-a ---ok
func (ps *ConfigurationStore) GetConfigurationById(id string) ([]*Config, error) {
	kv := ps.cli.KV()
	sid := constructKeyIDConfigurations(id)
	data, _, err := kv.List(sid, nil)
	if err != nil {
		return nil, err

	}
	configList := []*Config{}

	for _, pair := range data {
		config := &Config{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		configList = append(configList, config)

	}
	return configList, nil

}

///Pravi novu verziju konfiguracije ---ok
func (ps *ConfigurationStore) AddNewConfigVersion(config *Config) (*Config, error) {
	kv := ps.cli.KV()
	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	sid := constructKeyVersionConfigs(config.Id, config.Version)

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}
	return config, nil
}

///Dodaj konfiguraciju ---ok
func (ps *ConfigurationStore) PostConfig(configuration *Config) (*Config, error) {
	kv := ps.cli.KV()

	sid, rid := generateKeyConfiguration(configuration.Version)
	configuration.Id = rid

	data, err := json.Marshal(configuration)
	if err != nil {
		return nil, err
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return configuration, nil
}

///Obrisi konfiguraciju
func (ps *ConfigurationStore) DeleteConfig(id string, verzija string) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.Delete(constructKeyVersionConfigs(id, verzija), nil)
	if err != nil {
		fmt.Println("Konfiguracija nije obrisana!")
		return nil, err

	}

	return map[string]string{"Obrisana konfiguracija sa id: ": id}, nil
}

///Konfiguracije
///
///
///
///Grupe

///Nadji sve grupe
func (ps *ConfigurationStore) GetAllGroups() ([]*Group, error) {
	kv := ps.cli.KV()
	data, _, err := kv.List(allGroups, nil)
	if err != nil {
		return nil, err
	}

	groups := []*Group{}
	for _, pair := range data {
		group := &Group{}
		err = json.Unmarshal(pair.Value, group)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

//Dodavanje nove verzije grupe
func (ps *ConfigurationStore) AddNewGroupVersion(group *Group) (*Group, error) {
	kv := ps.cli.KV()

	for _, v := range group.Configs {
		labela := ""
		listaStringova := []string{}
		for k, val := range v {
			listaStringova = append(listaStringova, k+":"+val)
		}
		sort.Strings(listaStringova)
		for _, v := range listaStringova {
			labela += v + ";"
		}
		labela = labela[:len(labela)-1]
		sid := constructKeyGroupLabels(group.Id, group.Version, labela) + uuid.New().String()

		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		p := &api.KVPair{Key: sid, Value: data}
		_, err = kv.Put(p, nil)
		if err != nil {
			return nil, err
		}
	}

	return group, nil
}

///Nadji grupu preko id-a i verzije
func (ps *ConfigurationStore) GetGroupByIdVersion(id string, version string) (*Group, error) {
	kv := ps.cli.KV()

	data, _, err := kv.List(constructKeyGroupVersion(id, version), nil)
	if err != nil || data == nil {
		return nil, errors.New("Ne postoji ta grupa!")
	}

	configs := []map[string]string{}
	for _, pair := range data {
		config := &map[string]string{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, *config)
	}

	group := &Group{configs, id, version}

	return group, nil
}

///Nadji grupu po id-u
func (ps *ConfigurationStore) GetGroupById(id string) ([]*Group, error) {
	kv := ps.cli.KV()
	sid := constructKeyGroupId(id)
	data, _, err := kv.List(sid, nil)
	if err != nil {
		return nil, err

	}
	groupList := []*Group{}

	for _, pair := range data {
		group := &Group{}
		err = json.Unmarshal(pair.Value, group)
		if err != nil {
			return nil, err
		}
		groupList = append(groupList, group)

	}
	return groupList, nil

}

func (ps *ConfigurationStore) GetGroupByLabel(id string, version string, label string) ([]map[string]string, error) {
	kv := ps.cli.KV()
	data, _, err := kv.List(constructKeyGroupLabels(id, version, label), nil)
	if err != nil {
		return nil, err
	}

	configs := []map[string]string{}
	for _, pair := range data {
		config := &map[string]string{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, *config)
	}

	return configs, nil
}

//Brisanje grupe u bazi
func (ps *ConfigurationStore) DeleteGroup(id string, verzija string) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.DeleteTree(constructKeyGroupVersion(id, verzija), nil)
	if err != nil {
		fmt.Println("Grupa nije obrisana!")
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}

//Pravljenje grupe u bazi
func (ps *ConfigurationStore) PostGroup(group *Group) (*Group, error) {
	kv := ps.cli.KV()

	idGrupe := uuid.New().String()

	for _, v := range group.Configs {
		labela := ""
		listaStringova := []string{}
		for k, val := range v {
			listaStringova = append(listaStringova, k+":"+val)
		}
		sort.Strings(listaStringova)
		for _, v := range listaStringova {
			labela += v + ";"
		}
		labela = labela[:len(labela)-1]
		sid := constructKeyGroupLabels(idGrupe, group.Version, labela) + uuid.New().String()
		group.Id = idGrupe

		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		p := &api.KVPair{Key: sid, Value: data}
		_, err = kv.Put(p, nil)
		if err != nil {
			return nil, err
		}
	}

	return group, nil
}

//Izmena postojece grupe
func (ps *ConfigurationStore) UpdateGroup(group *Group) (*Group, error) {
	kv := ps.cli.KV()
	data, err := json.Marshal(group)

	sid := constructKeyGroupVersion(group.Id, group.Version)

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}
	return group, nil
}
