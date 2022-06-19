package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"os"
	"reflect"
	"strings"
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

///find by id and vers
func (ps *ConfigurationStore) GetConfigByIdVersion(id string, verzija string) (*Config, error) {
	kv := ps.cli.KV()
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

///find by id config
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

///Add new config version
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

///Add config
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

///find all config
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

//Update
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

///Find by id and version
func (ps *ConfigurationStore) GetGroupByIdVersion(id string, version string) (*Group, error) {
	kv := ps.cli.KV()
	sid := constructKeyGroupVersion(id, version)
	pair, _, err := kv.Get(sid, nil)
	if pair == nil {
		return nil, errors.New("Group not found!")
	}

	konfGrupa := &Group{}
	err = json.Unmarshal(pair.Value, konfGrupa)
	if err != nil {
		return nil, err
	}
	return konfGrupa, nil
}

//get by label
func (ps *ConfigurationStore) GetGroupByLabel(id string, version string, label string) ([]*ConfigurationInGroup, error) {
	kv := ps.cli.KV()
	listConfigs := []*ConfigurationInGroup{}

	sid := constructKeyGroupVersion(id, version)
	pair, _, err := kv.Get(sid, nil)
	if pair == nil {
		return nil, err
	}

	labellist := strings.Split(label, ";")
	labeldb := make(map[string]string)
	for _, labela := range labellist {
		deolabele := strings.Split(labela, ":")
		if deolabele != nil {
			labeldb[deolabele[0]] = deolabele[1]
		}
	}

	confGroup := &Group{}
	err = json.Unmarshal(pair.Value, confGroup)
	if err != nil {
		return nil, err
	}

	for _, config := range confGroup.Configs {
		if len(config.Entries) == len(labeldb) {
			if reflect.DeepEqual(config.Entries, labeldb) {
				listConfigs = append(listConfigs, config)
			}
		}
	}

	return listConfigs, nil
}

//add new version
func (ps *ConfigurationStore) AddNewGroupVersion(group *Group) (*Group, error) {
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

//post group
func (ps *ConfigurationStore) PostGroup(group *Group) (*Group, error) {
	kv := ps.cli.KV()

	sid, rid := generateGroupKey(group.Version)
	group.Id = rid

	data, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	pairs, _, err := kv.Get(sid, nil)
	if pairs != nil {
		return nil, errors.New("There is already configuaration with id like that! ")
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return group, nil
}

///Find by id
func (ps *ConfigurationStore) GetGroupById(id string) ([]*Group, error) {
	kv := ps.cli.KV()
	sid := constructKeyGroupId(id)
	data, _, err := kv.List(sid, nil)
	if err != nil {
		return nil, err

	}
	GroupList := []*Group{}

	for _, pair := range data {
		grupa := &Group{}
		err = json.Unmarshal(pair.Value, grupa)
		if err != nil {
			return nil, err
		}
		GroupList = append(GroupList, grupa)

	}
	return GroupList, nil

}

//delete group
func (ps *ConfigurationStore) DeleteGroup(id string, verzija string) (map[string]string, error) {
	kv := ps.cli.KV()
	data, _, err := kv.List(constructKeyGroupVersion(id, verzija), nil)
	if err != nil || data == nil {
		return nil, errors.New("Cannot find that group!")
	} else {
		_, greska := kv.Delete(constructKeyGroupVersion(id, verzija), nil)
		if greska != nil {
			return nil, greska
		}

		return map[string]string{"Deleted:": id}, nil
	}
}

///delete config
func (ps *ConfigurationStore) DeleteConfig(id string, verzija string) (map[string]string, error) {
	kv := ps.cli.KV()
	pair, _, greska := kv.Get(constructKeyVersionConfigs(id, verzija), nil)
	if greska != nil || pair == nil {
		return nil, errors.New("No configuration like that!Cannot delete!")
	} else {
		data, err := kv.Delete(constructKeyVersionConfigs(id, verzija), nil)
		if err != nil || data == nil {
			fmt.Println("Config not deleted!")
			return nil, err

		}

		return map[string]string{"Deleted: ": id}, nil
	}
}
