package main

import (
	"bytes"
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

type (
	// SnippetID is the ID of a config snippet
	SnippetID string
	// BaseConfigContent is the content of a base config
	BaseConfigContent string
	// Selector is a map of selector labels
	Selector map[string]string
)

type Namespace struct {
	BaseConfig BaseConfigContent      `json:"base_config" yaml:"base_config"`
	Snippets   map[SnippetID]*Snippet `json:"snippets" yaml:"snippets"`
}

type BaseConfig struct {
	Config string `json:"config" yaml:"config"`
}

// Snippet is a snippet of configuration for a specific selectors.
type Snippet struct {
	// Config is the snippet of config to be included.
	Config string `json:"config" yaml:"config"`
	// Selector is map to label the snippet.
	Selector Selector `json:"selector" yaml:"selector"`
}

const configPath = "cfg/"
const baseFilename = "base.yaml"
const snipsPath = "snips/"
const apiPath = "agent-management/api/config/v1/namespace"

var APIHost = os.Getenv("AGENT_MANAGEMENT_HOST")
var APIToken = os.Getenv("AGENT_MANAGEMENT_TOKEN")

func main() {
	files, err := ioutil.ReadDir(configPath)
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
	nsPath := filepath.Join(configPath, name)
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
	if err != nil {
		log.Fatal(err)
	}

	snips := make(map[SnippetID]*Snippet, len(snipsFiles))
	for _, file := range snipsFiles {
		if !file.IsDir() {
			snipBuf, err := os.ReadFile(filepath.Join(nsSnipsPath, file.Name()))
			if err != nil {
				log.Fatal(err)
			}
			id := SnippetID(strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())))
			var snip Snippet
			err = yaml.Unmarshal(snipBuf, &snip)
			if err != nil {
				log.Fatal(err)
			}

			snips[id] = &snip
		}
	}

	ns := Namespace{
		BaseConfig: BaseConfigContent(base.Config),
		Snippets:   snips,
	}

	payload, err := yaml.Marshal(ns)
	if err != nil {
		log.Fatal(err)
	}

	uri, err := url.JoinPath("http://", APIHost, apiPath, name)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("PUT", uri, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(uri)

	req.Header.Add("Authorization", fmt.Sprintf("Basic %v", APIToken))
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done!")
}
