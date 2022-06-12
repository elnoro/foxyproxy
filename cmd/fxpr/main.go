package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"foxyproxy/cmd/fxpr/app"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("missing command")
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup

	a := app.New()

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
	case "test-server":
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := a.RunTestServer(ctx)
			if err != nil {
				log.Println(err)
			}
		}()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("Exiting...")
	cancel()
	wg.Wait()
	log.Println("Cleanup done!")
}
