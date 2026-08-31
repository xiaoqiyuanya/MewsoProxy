package sshclient

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type Config struct {
	Host       string
	Port       int
	User       string
	Password   string
	PrivateKey string
	Timeout    time.Duration
}

type Client struct {
	cfg    Config
	client *ssh.Client
}

func New(cfg Config) (*Client, error) {
	if cfg.Host == "" || cfg.User == "" {
		return nil, fmt.Errorf("SSH 地址或用户名缺失")
	}
	if cfg.Port == 0 {
		cfg.Port = 22
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	var auths []ssh.AuthMethod
	if cfg.PrivateKey != "" {
		signer, err := parsePrivateKey([]byte(cfg.PrivateKey))
		if err != nil {
			return nil, err
		}
		auths = append(auths, ssh.PublicKeys(signer))
	}
	if cfg.Password != "" {
		auths = append(auths, ssh.Password(cfg.Password))
	}
	if len(auths) == 0 {
		return nil, fmt.Errorf("SSH 认证信息缺失")
	}
	sshCfg := &ssh.ClientConfig{
		User:            cfg.User,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         cfg.Timeout,
	}
	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	client, err := ssh.Dial("tcp", addr, sshCfg)
	if err != nil {
		return nil, err
	}
	return &Client{cfg: cfg, client: client}, nil
}

func parsePrivateKey(pemBytes []byte) (ssh.Signer, error) {
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("解析 SSH 私钥失败: %w", err)
	}
	return signer, nil
}

func (c *Client) Close() error { return c.client.Close() }

func (c *Client) Run(ctx context.Context, cmd string) (string, string, error) {
	var stdout, stderr bytes.Buffer
	session, err := c.client.NewSession()
	if err != nil {
		return "", "", err
	}
	defer session.Close()
	session.Stdout = &stdout
	session.Stderr = &stderr
	if err := session.Run(cmd); err != nil {
		return stdout.String(), stderr.String(), err
	}
	return stdout.String(), stderr.String(), nil
}

func (c *Client) RunWithLog(ctx context.Context, cmd string, onLine func(string)) error {
	session, err := c.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		return err
	}
	if err := session.Start(cmd); err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go scanLines(stdout, onLine, &wg)
	go scanLines(stderr, onLine, &wg)
	err = session.Wait()
	wg.Wait()
	return err
}

func scanLines(r io.Reader, onLine func(string), wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if onLine != nil {
			onLine(scanner.Text())
		}
	}
}
