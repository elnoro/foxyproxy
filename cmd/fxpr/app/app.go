package app

import (
	"context"
	"log"
	"os"
	"time"

	"foxyproxy/doapi"
	"foxyproxy/sshproxy"
)

type App struct {
}

func New() *App {
	return &App{}
}

func (a *App) RunProxy(ctx context.Context) error {
	token := os.Getenv("DO_TOKEN")
	if token == "" {
		log.Fatal("DO_TOKEN is not set")
	}

	sc := doapi.NewSimpleClient(token, getSshKey(), time.Minute)

	s, err := sc.StartDroplet(ctx, "proxy")
	if err != nil {
		log.Println("cannot start droplet", err)
		if s.Id != 0 {
			deleteDroplet(ctx, sc, s)
		}

		return err
	}

	proxy, err := sshproxy.StartProxy(s.PublicIP, 1337)
	if err != nil {
		log.Println("cannot start proxy", err)
		deleteDroplet(ctx, sc, s)

		return err
	}
	log.Println("proxy started: ", proxy)

	<-ctx.Done()
	deleteDroplet(context.Background(), sc, s)
	err = proxy.Stop()
	if err != nil {
		log.Println("cannot stop proxy", err)
	}

	return nil
}

func (a *App) RunTestServer(ctx context.Context) error {
	token := os.Getenv("DO_TOKEN")
	if token == "" {
		log.Fatal("DO_TOKEN is not set")
	}

	sc := doapi.NewSimpleClient(token, getSshKey(), time.Minute)

	s, err := sc.StartDroplet(ctx, "proxy")
	if err != nil {
		log.Println("cannot start droplet", err)
		if s.Id != 0 {
			deleteDroplet(ctx, sc, s)
		}

		return err
	}

	log.Println("Droplet started! Run ssh root@" + s.PublicIP)

	<-ctx.Done()
	deleteDroplet(context.Background(), sc, s)

	return nil
}

func deleteDroplet(ctx context.Context, sc *doapi.SimpleClient, s doapi.Server) {
	log.Println("deleting droplet with id", s.Id)
	if err := sc.DeleteDroplet(ctx, s.Id); err != nil {
		log.Println("cannot delete droplet", err)
	}
}

func getSshKey() string {
	return "76:4f:ef:75:a5:5d:35:7a:fa:50:dd:41:85:d1:34:41"
}
