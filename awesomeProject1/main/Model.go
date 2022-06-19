package main

type Config struct {
	Entries map[string]string `json:"entries"`
	Id      string            `json:"id"`
	Version string            `json:"version"`
}

type Group struct {
	Configs []*ConfigurationInGroup `json:"configs"`
	Id      string                  `json:"id"`
	Version string                  `json:"version"`
}

type ConfigurationInGroup struct {
	Entries map[string]string `json:"labele"`
	Labele  map[string]string `json:"entries"`
}
