package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"foxyproxy/cmd/fxpr/app"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup

	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[1] {
	case "proxy":
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := a.RunProxy(ctx)
			if err != nil {
				log.Fatal(err)
			}
		}()
	case "test":
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := a.RunTestServer(ctx)
			if err != nil {
				log.Fatal(err)
			}
		}()
	case "help":
		showHelp()
	default:
		fmt.Println("Invalid command.")
		showHelp()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("Exiting...")
	cancel()
	wg.Wait()
	log.Println("Cleanup done!")
}

func showHelp() {
	fmt.Print("fxpr is a CLI tool to quickly spin up and destroy DigitalOcean servers" +
		"\n\nUsage:\n  fxpr [command]\n\nAvailable Commands:" +
		"\n  proxy         Start a droplet and an SSH tunnel on localhost. Hit Ctrl-C to destroy the droplet" +
		"\n  test          Start a droplet you can SSH into. Hit Ctrl-C to destroy the droplet" +
		"\n")
}
