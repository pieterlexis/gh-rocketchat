package main

import (
	"flag"
	"fmt"
	"github.com/pieterlexis/gh-rocketchat/config"
	"github.com/pieterlexis/gh-rocketchat/processor"
	"github.com/pieterlexis/gh-rocketchat/version"
	"net/http"
	"os"

	"github.com/pieterlexis/gh-rocketchat/receiver"
	log "github.com/sirupsen/logrus"
)

func indexOf(element string, data []string) (int) {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1    //not found.
}

var UsagePre = fmt.Sprintf(`GitHub to Rocket.Chat Webhook Translator [%s]

This program ingests GitHub/GitLab/Gogs/Gitea webhooks and translates them
to messages that can be displayed by rocket.chat.
`, version.Version)

var configFile string
var debug bool

func init() {
	flag.StringVar(&configFile, "config", "gh-rocketchat.yaml",
		"Configuration file to load")
	flag.BoolVar(&debug, "debug", false, "Enable debug-level logging")
}

func usage() {
	fmt.Fprint(os.Stderr, UsagePre)
	fmt.Fprint(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}


func main() {
	formatter := log.TextFormatter{
		DisableTimestamp:       true,
		DisableLevelTruncation: true}
	log.SetFormatter(&formatter)

	flag.CommandLine.Usage = usage
	flag.Parse()
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Infof("Starting up!")

	config, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatalf("Could not parse config file %s: %+v", "", err)
	}

	if len(config.Hooks) == 0 {
		log.Fatal("No hooks configured!")
	}
	var hookNames []string
	for _, hook := range config.Hooks {
		if len(hook.Name) == 0 {
			log.Fatal("Hook without a name found, can not continue")
		}
		if indexOf(hook.Name, hookNames) >= 0 {
			log.Fatalf("More than one hook with name '%s' in the configuration", hook.Name)
		}
		hookNames = append(hookNames, hook.Name)

		if len(hook.Destination) == 0 {
			log.Fatalf("Destination for hook '%s' not set", hook.Name)
		}

		if len(hook.Endpoint) == 0 {
			log.Fatalf("Endpoint for hook '%s' not set", hook.Name)
		}

		var ghPayloadChan = make(chan interface{})
		rec, err := receiver.NewReceiver(hook.Name, hook.Secret)
		if err != nil {
			log.Fatalf("Could not initiate webhook receiver: %v", err)
		}
		http.Handle(hook.Endpoint, rec.Handle(ghPayloadChan))
		processor.RunProcessor(ghPayloadChan, hook.Destination, hook.Name)
	}

	log.Infof("All receivers registered and processors started")
	http.ListenAndServe(config.ListenAddress, nil)
}
