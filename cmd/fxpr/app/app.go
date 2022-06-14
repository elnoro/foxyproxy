package app

import (
	"context"
	"log"
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
	sc := newClient()
	log.Println("starting droplet...")
	s, err := sc.StartDroplet(ctx, "proxy")
	if err != nil {
		log.Println("cannot start droplet", err)
		if s.Id != 0 {
			deleteDroplet(ctx, sc, s)
		}

		return err
	}
	log.Println("droplet is active!")

	log.Println("staring proxy...")
	proxy, err := sshproxy.StartProxy(s.PublicIP, 1337)
	if err != nil {
		log.Println("cannot start proxy", err)
		time.Sleep(30 * time.Second)
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
	sc := newClient()

	log.Println("starting droplet...")
	s, err := sc.StartDroplet(ctx, "test-server")
	if err != nil {
		log.Println("cannot start droplet", err)
		if s.Id != 0 {
			deleteDroplet(ctx, sc, s)
		}

		return err
	}

	log.Println("droplet is active! Run ssh root@" + s.PublicIP)

	<-ctx.Done()
	deleteDroplet(context.Background(), sc, s)

	return nil
}

func newClient() *doapi.SimpleClient {
	config, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	sc := doapi.NewSimpleClient(config.DoToken, config.Fingerprint, time.Minute)

	return sc
}

func deleteDroplet(ctx context.Context, sc *doapi.SimpleClient, s doapi.Server) {
	log.Println("deleting droplet with id", s.Id)
	if err := sc.DeleteDroplet(ctx, s.Id); err != nil {
		log.Println("cannot delete droplet", err)
	}
}
