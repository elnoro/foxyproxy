package app

import (
	"flag"
	"github.com/peterbourgon/ff/v3"
	"log"
	"os"
)

const defaultProxyPort = 1337

type config struct {
	DoToken     string `json:"do_token"`
	Fingerprint string `json:"fingerprint"`
	Port        int    `json:"port,omitempty"`
}

func loadConfig() (config, error) {
	paths := []ff.Option{
		ff.WithAllowMissingConfigFile(true),
		ff.WithConfigFileParser(ff.JSONParser),
		ff.WithConfigFile(getConfigFile()),
	}

	fs := flag.NewFlagSet("default", flag.ContinueOnError)
	c := config{}

	fs.StringVar(&c.Fingerprint, "fingerprint", "", "Digital Ocean fingerprint")
	fs.StringVar(&c.DoToken, "do_token", "", "Digital Ocean token")
	fs.IntVar(&c.Port, "port", defaultProxyPort, "Digital Ocean token")

	// TODO move to main
	err := ff.Parse(fs, os.Args[2:], paths...)
	log.Println(c)

	return c, err
}

func getConfigFile() string {
	_, err := os.Stat("config.json")
	if err == nil {
		return "config.json"
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return homeDir + "/.config/fxpr/config.json"
}
