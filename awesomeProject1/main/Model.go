package main

type Config struct {
	Entries map[string]string `json:"entries"`
	Id      string            `json:"id"`
	Version string            `json:"version"`
}

type Group struct {
	Configs []map[string]string `json:"entries"`
	Id      string              `json:"id"`
	Version string              `json:"version"`
}
