package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Selector map[string]string

type Namespace struct {
	BaseConfig string              `json:"base_config" yaml:"base_config"`
	Snippets   map[string]*Snippet `json:"snippets" yaml:"snippets"`
}

type BaseConfig struct {
	Config string `json:"config" yaml:"config"`
}

type Snippet struct {
	Config   string   `json:"config" yaml:"config"`
	Selector Selector `json:"selector" yaml:"selector"`
}

const baseFilename = "base.yaml"
const snipsPath = "snips/"
const apiPath = "agent-management/api/config/v1/namespaces"
const apiScheme = "https"

var APIHost = os.Getenv("AGENT_MANAGEMENT_HOST")
var APIUsername = os.Getenv("AGENT_MANAGEMENT_USERNAME")
var APIPassword = os.Getenv("AGENT_MANAGEMENT_PASSWORD")
var ConfigPath = os.Getenv("CONFIG_PATH")

func main() {
	files, err := ioutil.ReadDir(ConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			processNamespace(file.Name())
		}
	}
}

func processNamespace(name string) {
	fmt.Printf("Processing namespace \"%v\"...", name)

	nsPath := filepath.Join(ConfigPath, name)
	nsSnipsPath := filepath.Join(nsPath, snipsPath)

	baseBuf, err := os.ReadFile(filepath.Join(nsPath, baseFilename))
	if err != nil {
		log.Fatal(err)
	}
	var base BaseConfig
	err = yaml.Unmarshal(baseBuf, &base)
	if err != nil {
		log.Fatal(err)
	}

	snipsFiles, err := ioutil.ReadDir(nsSnipsPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}

	snips := make(map[string]*Snippet, len(snipsFiles))
	for _, file := range snipsFiles {
		if !file.IsDir() {
			snipBuf, err := os.ReadFile(filepath.Join(nsSnipsPath, file.Name()))
			if err != nil {
				log.Fatal(err)
			}
			id := string(strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())))
			var snip Snippet
			err = yaml.Unmarshal(snipBuf, &snip)
			if err != nil {
				log.Fatal(err)
			}

			snips[id] = &snip
		}
	}

	ns := Namespace{
		BaseConfig: base.Config,
		Snippets:   snips,
	}

	payload, err := yaml.Marshal(ns)
	if err != nil {
		log.Fatal(err)
	}

	uri, err := url.JoinPath(fmt.Sprintf("%v://", apiScheme), APIHost, apiPath, name)
	if err != nil {
		log.Fatal(err)
	}

	token := fmt.Sprintf("%v:%v", APIUsername, APIPassword)
	encodedToken := base64.StdEncoding.EncodeToString([]byte(token))

	req, err := http.NewRequest("PUT", uri, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Basic %v", encodedToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(fmt.Errorf("something went wrong, aborting: Error: %v", err))
	}
	if resp.StatusCode != 202 {
		log.Fatal(fmt.Errorf("something went wrong, aborting: Status Code %v", resp.StatusCode))
	}
	fmt.Printf(" Done\n")
}
