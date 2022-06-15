package droplets

import (
	"context"
	"fmt"
	"time"

	"github.com/digitalocean/godo"
)

type SimpleClient struct {
	client      *godo.Client
	fingerPrint string
	waitTimeout time.Duration
}

type Server struct {
	Id       int
	PublicIP string
}

func NewSimpleClient(token string, fingerPrint string, waitTimeout time.Duration) *SimpleClient {
	return &SimpleClient{
		client:      godo.NewFromToken(token),
		fingerPrint: fingerPrint,
		waitTimeout: waitTimeout,
	}
}

func (s *SimpleClient) StartDroplet(ctx context.Context, tagPrefix string) (Server, error) {
	dropletName := time.Now().Format(tagPrefix + "-2006-01-02-15-04-05")

	createRequest := &godo.DropletCreateRequest{
		Name:   dropletName,
		Region: "ams3",
		Size:   "s-1vcpu-1gb",
		Image: godo.DropletCreateImage{
			Slug: "ubuntu-20-04-x64",
		},
		SSHKeys: []godo.DropletCreateSSHKey{
			{Fingerprint: s.fingerPrint},
		},
		IPv6: true,
		Tags: []string{tagPrefix},
	}

	droplet, _, err := s.client.Droplets.Create(ctx, createRequest)
	if err != nil {
		return Server{}, fmt.Errorf("creating droplet, %w", err)
	}

	server := Server{Id: droplet.ID}

	timeout, cancel := context.WithTimeout(ctx, s.waitTimeout)
	defer cancel()
	activeDroplet, err := waitForDroplet(timeout, s.client, droplet.ID)
	if err != nil {
		return server, fmt.Errorf("waiting for droplet to become active, %w", err)
	}
	publicIP, err := activeDroplet.PublicIPv4()
	if err != nil {
		return server, fmt.Errorf("getting public ip, %w", err)
	}

	server.PublicIP = publicIP

	return server, nil
}

func (s *SimpleClient) DeleteDroplet(ctx context.Context, dropletID int) error {
	_, err := s.client.Droplets.Delete(ctx, dropletID)
	if err != nil {
		return fmt.Errorf("deleting droplet, %w", err)
	}
	return nil
}

func waitForDroplet(ctx context.Context, client *godo.Client, dropletID int) (*godo.Droplet, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			droplet, _, err := client.Droplets.Get(ctx, dropletID)
			if err != nil {
				return nil, err
			}
			if droplet.Status == "active" {
				return droplet, nil
			}
			time.Sleep(10 * time.Second)
		}
	}
}
