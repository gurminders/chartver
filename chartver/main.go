package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"gopkg.in/yaml.v2"
)

var repoURL string
var count int

func main() {
	flag.StringVar(&repoURL, "repo", "https://svl-artifactory.juniper.net/atom-helm", "URL of Helm Chart repo")
	flag.IntVar(&count, "count", 5, "Limit number of versions of chart to display")
	flag.Parse()

	index, err := getIndex()
	if err != nil {
		log.Fatal(err)
	}

	if len(flag.Args()) > 0 {
		for _, chartName := range flag.Args() {
			printChartVersion(index, chartName, true)
		}
	} else {
		printChartNames(index)
	}
}

func getIndex() (*Index, error) {
	indexURL := repoURL + "/index.yaml"
	log.Println("Getting and parsing repo index ", indexURL)

	res, err := http.Get(indexURL)
	if err != nil {
		return nil, err
	}

	// read the yaml into memory
	log.Println("Parsing index")

	index := Index{}
	err = yaml.NewDecoder(res.Body).Decode(&index)
	if err != nil {
		return nil, err
	}

	log.Println("Get complete")

	return &index, nil
}

func printChartNames(index *Index) {
	var charts []string
	for key, _ := range index.Entries {
		charts = append(charts, key)
	}

	fmt.Printf("%-30s %-30s %-30s Created\n", "Chart", "Version", "AppVersion")
	fmt.Println("----------------------------------------------------------------------------------------------------------------------------------")
	count = 1; // only print the latest

	sort.Strings(charts)
	for _, name := range charts {
		printChartVersion(index, name, false)
	}
}

func printChartVersion(index *Index, chartName string, headers bool) {
	entries := index.Entries[chartName]
	if entries == nil {
		fmt.Println("Chart ", chartName, " not found in repository")
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Created.After(entries[j].Created)
	})

	if headers {
		fmt.Println("Chart: ", chartName)
		fmt.Printf("%-30s %-30s Created\n", "Version", "AppVersion")
		fmt.Println("----------------------------------------------------------------------------------------------------")
	}

	for idx, entry := range entries {
		if idx+1 > count {
			break
		}

		if headers {
			fmt.Printf("%-30s %-30s %v\n", entry.Version, entry.AppVersion, entry.Created)
		} else {
			fmt.Printf("%-30s %-30s %-30s %v\n", chartName, entry.Version, entry.AppVersion, entry.Created)
		}
	}

	if headers {
		fmt.Println("----------------------------------------------------------------------------------------------------\n")
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
