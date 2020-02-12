package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/SimonRichardson/juju-inspect/rules"
	"gopkg.in/yaml.v2"
)

func main() {
	files := os.Args[1:]
	if len(files) == 0 {
		log.Fatal("expected at least on file")
	}
	allRules := []Rule{
		rules.NewRaftRule(),
		rules.NewMongoRule(),
	}
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}

		row1, _, err := bufio.NewReader(f).ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		_, err = f.Seek(int64(len(row1)), io.SeekStart)
		if err != nil {
			log.Fatal(err)
		}

		var report rules.Report
		if err := yaml.NewDecoder(f).Decode(&report); err != nil {
			log.Fatal(err)
		}

		// TODO (pick a better name somehow - agent.name?)
		name := report.Manifolds["agent"].Report.Agent
		for _, rule := range allRules {
			rule.Run(name, report)
		}
	}

	fmt.Println("")
	fmt.Println("Analysis of Engine Report:")
	fmt.Println("")
	for _, rule := range allRules {
		fmt.Println(rule.Summary())
		fmt.Println("\t", rule.Analyse())
	}
}

type Rule interface {
	Run(string, rules.Report)
	Summary() string
	Analyse() string
}
