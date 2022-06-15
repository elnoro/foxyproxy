package proxy

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"time"
)

type Proxy struct {
	cmd *exec.Cmd
}

func (p *Proxy) String() string {
	if p.cmd == nil {
		return "Proxy is not running"
	} else {
		return fmt.Sprintf("Proxy is running with PID %d. Launch command %s", p.cmd.Process.Pid, p.cmd.String())
	}
}

func (p *Proxy) Stop() error {
	err := p.cmd.Process.Kill()
	if err != nil {
		return fmt.Errorf("cannot stop proxy, %w", err)
	}

	return nil
}

func StartProxy(serverIp string, port int) (*Proxy, error) {
	err := waitForPort(serverIp)
	if err != nil {
		return nil, fmt.Errorf("issues with the server, %w", err)
	}

	portStr := fmt.Sprintf("%d", port)
	sshCmd := "ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -D " + portStr + " -C -N root@" + serverIp
	cmd := exec.Command("sh", "-c", sshCmd)
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start ssh tunnel, %w", err)
	}
	return &Proxy{cmd: cmd}, nil
}

func waitForPort(serverIP string) error {
	log.Println("waiting for ssh to become available on", serverIP)
	timeout := time.Second
	retries := 60
	for retries > 0 {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(serverIP, "22"), timeout)
		if err != nil {
			retries--
			time.Sleep(time.Second)
			continue
		}
		if conn != nil {
			err := conn.Close()
			if err != nil {
				return fmt.Errorf("connectivity issues, %w", err)
			}
			return nil
		}
	}

	return fmt.Errorf("ssh is not available on ip %s", serverIP)
}
