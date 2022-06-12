package sshproxy

import (
	"fmt"
	"os/exec"
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
	portStr := fmt.Sprintf("%d", port)
	sshCmd := "ssh -D " + portStr + " -C -q -N root@" + serverIp
	cmd := exec.Command("sh", "-c", sshCmd)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return &Proxy{cmd: cmd}, nil
}
