package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"foxyproxy/droplets"
	"foxyproxy/proxy"
)

type App struct {
	config config
	client *droplets.SimpleClient
}

func New() (*App, error) {
	a := &App{}
	err := a.init()
	if err != nil {
		return nil, fmt.Errorf("initialization error, %w", err)
	}

	return a, nil
}

func (a *App) RunProxy(ctx context.Context) error {
	s, err := a.startDroplet(ctx, "proxy")
	if err != nil {
		return err
	}

	log.Println("staring proxy...")
	pr, err := proxy.StartProxy(s.PublicIP, 1337)
	if err != nil {
		log.Println("cannot start proxy", err)
		time.Sleep(30 * time.Second)
		a.deleteDroplet(ctx, s)

		return err
	}
	log.Println("proxy started: ", pr)

	<-ctx.Done()
	a.deleteDroplet(context.Background(), s)
	err = pr.Stop()
	if err != nil {
		log.Println("cannot stop proxy", err)
	}

	return nil
}

func (a *App) RunTestServer(ctx context.Context) error {
	s, err := a.startDroplet(ctx, "test")
	if err != nil {
		return err
	}
	log.Println("Run ssh root@" + s.PublicIP)

	<-ctx.Done()
	a.deleteDroplet(context.Background(), s)

	return nil
}

func (a *App) ListDroplets(ctx context.Context) error {
	servers, err := a.client.List(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Servers:")
	for _, server := range servers {
		fmt.Println(server.Id, server.Name, server.PublicIP)
	}

	return nil
}

func (a *App) init() error {
	c, err := loadConfig()
	if err != nil {
		return err
	}
	a.config = c
	a.client = droplets.NewSimpleClient(c.DoToken, c.Fingerprint, time.Minute)

	return nil
}

func (a *App) deleteDroplet(ctx context.Context, s droplets.Server) {
	log.Println("deleting droplet with id", s.Id)
	if err := a.client.DeleteDroplet(ctx, s.Id); err != nil {
		log.Println("cannot delete droplet", err)
	}
}

func (a *App) startDroplet(ctx context.Context, tagPrefix string) (droplets.Server, error) {
	log.Println("starting droplet...")
	s, err := a.client.StartDroplet(ctx, tagPrefix)
	if err != nil {
		log.Println("cannot start droplet", err)
		if s.Id != 0 {
			a.deleteDroplet(ctx, s)
		}

		return s, err
	}
	log.Println("droplet is active!")

	return s, nil
}
