package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"foxyproxy/cmd/fxpr/app"
)

var version = "development"
var date = "0000-00-00"

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

	fs := flag.NewFlagSet("flag", flag.ExitOnError)
	d := fs.Bool("d", false, "run process in the background")

	switch os.Args[1] {
	case "proxy":
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := a.RunProxy(ctx)
			if err != nil {
				log.Println(err)
			}
		}()
	case "test":
		err := fs.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}

		// daemon mode
		if *d {
			_, err := a.StartTestServer(ctx)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("don't forget to remove the server!")

			return
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := a.RunTestServer(ctx)
			if err != nil {
				log.Println(err)
			}
		}()
	case "list":
		err := a.ListDroplets(ctx)
		if err != nil {
			log.Fatal(err)
		}
		return
	case "help":
		showHelp()
		return
	default:
		fmt.Println("Invalid command.")
		showHelp()
		return
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
	fmt.Print("fxpr is a CLI tool to quickly spin up and destroy DigitalOcean servers"+
		"\nVersion: "+version+", built on "+date,
		"\n\nUsage:\n  fxpr [command]\n\nAvailable Commands:"+
			"\n  proxy         Start a droplet and an SSH tunnel on localhost. Hit Ctrl-C to destroy the droplet"+
			"\n  test          Start a droplet you can SSH into. Hit Ctrl-C to destroy the droplet"+
			"\n  list          Shows the list of existing droplets"+
			"\n")
}
