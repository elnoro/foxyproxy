package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var errConfigNotFound = errors.New("no valid config found")

type config struct {
	DoToken     string `json:"do_token"`
	Fingerprint string `json:"fingerprint"`
}

func loadConfig() (config, error) {
	var paths []string
	homeDir, err := os.UserHomeDir()
	if err == nil {
		paths = append(paths, homeDir+"/.config/fxpr/config.json")
	}
	paths = append(paths, "config.json")
	for _, v := range paths {
		b, err := os.ReadFile(v)
		if err != nil {
			continue
		}
		c := config{}
		err = json.Unmarshal(b, &c)
		if err != nil {
			return c, fmt.Errorf("malformed config, %w", err)
		}

		return c, nil
	}

	return config{}, errConfigNotFound
}
