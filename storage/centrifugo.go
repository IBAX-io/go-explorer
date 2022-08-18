/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package storage

import "github.com/centrifugal/gocent"

var publisher *gocent.Client

type CentrifugoConfig struct {
	Enable bool   `yaml:"enable"`
	Secret string `yaml:"secret"`
	URL    string `yaml:"url"`
	Socket string `yaml:"socket"`
	Key    string `yaml:"key"`
}

func (c *CentrifugoConfig) Init() error {
	if c.Enable {
		publisher = gocent.New(gocent.Config{
			Addr: c.URL,
			Key:  c.Key,
		})
	}
	return nil
}

func (c *CentrifugoConfig) Conn() *gocent.Client {
	return publisher
}

func (l *CentrifugoConfig) Close() error {
	return nil
}
