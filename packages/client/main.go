package client

import (
	"fmt"
	"crypto/tls"

	irc "github.com/fluffle/goirc/client"
)

type ClientOptions struct {
	Channels []string
}

type Client struct {
	Conn          *irc.Conn
	ClientOptions ClientOptions
}

func Create(server string, port int, nickname string, options ClientOptions,
            ssl bool, allowInsecure bool) *Client {
	config := irc.NewConfig(nickname)
	config.SSL = ssl
	if ssl {
		config.SSLConfig = &tls.Config{}
		if allowInsecure {
			config.SSLConfig.InsecureSkipVerify = true
		} else {
			config.SSLConfig.ServerName = server
		}
	}
	config.Server = fmt.Sprintf("%s:%d", server, port)
	config.NewNick = func(n string) string { return n + "^" }

	client := irc.Client(config)
	client.EnableStateTracking()

	client.HandleFunc(irc.CONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			for _, channel := range options.Channels {
				conn.Join(channel)
			}
		})

	return &Client{
		Conn:          client,
		ClientOptions: options,
	}
}

func (c *Client) Connect() error {
	if err := c.Conn.Connect(); err != nil {
		return err
	}

	return nil
}
