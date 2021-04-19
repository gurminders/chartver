package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"gopkg.in/yaml.v2"
)

const defaultRepo = "https://svl-artifactory.juniper.net/atom-helm"

func main() {
	var chartName string
	if len(os.Args) > 1 {
		chartName = os.Args[1]
	}

	repoURL := defaultRepo
	if len(os.Args) > 2 {
		repoURL = os.Args[2]
	}

	index, err := getIndex(repoURL)
	if err != nil {
		log.Fatal(err)
	}

	if chartName != "" {
		printChartVersion(index, os.Args[1])
	} else {
		printChartNames(index)
	}
}

func getIndex(repoURL string) (*Index, error) {
	indexURL := repoURL + "/index.yaml"
	log.Println("Fetching ", indexURL)

	res, err := http.Get(indexURL)
	if err != nil {
		return nil, err
	}

	log.Println("Fetch complete")

	// read the yaml into memory
	log.Println("Parsing index")

	index := Index{}
	err = yaml.NewDecoder(res.Body).Decode(&index)
	if err != nil {
		return nil, err
	}

	log.Println("Parsing complete")

	return &index, nil
}

func printChartNames(index *Index) {
	var charts []string
	for key, _ := range index.Entries {
		charts = append(charts, key)
	}

	sort.Strings(charts)
	for _, name := range charts {
		fmt.Println(name)
	}
}

func printChartVersion(index *Index, chartName string) {
	entries := index.Entries[chartName]
	if entries == nil {
		fmt.Println("Chart ", chartName, " not found in repository")
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Created.After(entries[j].Created)
	})

	fmt.Printf("%-30s %-30s Created\n", "Version", "AppVersion")
	fmt.Println("----------------------------------------------------------------------------------------------------")

	for _, entry := range entries {
		fmt.Printf("%-30s %-30s %v\n", entry.Version, entry.AppVersion, entry.Created)
	}
}

type Index struct {
	Entries map[string][]ChartEntry `yaml:"entries"`
}

type ChartEntry struct {
	Name       string    `yaml:"name"`
	Created    time.Time `yaml:"created"`
	AppVersion string    `yaml:"appVersion"`
	Version    string    `yaml:version`
}
