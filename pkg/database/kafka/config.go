// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kafka

import (
	"time"

	"github.com/go-playground/validator/v10"
	logrus_ "github.com/searKing/golang/third_party/github.com/sirupsen/logrus"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	viper_ "github.com/searKing/golang/third_party/github.com/spf13/viper"
)

type Config struct {
	Proto     Kafka
	Validator *validator.Validate

	viper     *viper.Viper
	viperKeys []string
}

type completedConfig struct {
	*Config

	// for Complete Only
	completeError error
}

type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// NewConfig returns a Config struct with the default values
func NewConfig() *Config {
	return &Config{}
}

// NewViperConfig returns a Config struct with the global viper instance
// key representing a sub tree of this instance.
// NewViperConfig is case-insensitive for a key.
func NewViperConfig(v *viper.Viper, keys ...string) *Config {
	c := NewConfig()
	c.viper = v
	c.viperKeys = keys
	return c
}

// Validate checks Config.
func (c *completedConfig) Validate() error {
	return c.Validator.Struct(c)
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to ApplyOptions, do that first. It's mutating the receiver.
// ApplyOptions is called inside.
func (c *Config) Complete() CompletedConfig {
	if c.viper != nil {
		err := viper_.UnmarshalKeys(c.viperKeys, &c.Proto)
		if err != nil {
			return CompletedConfig{&completedConfig{completeError: err}}
		}
	}
	if c.Validator == nil {
		c.Validator = validator.New()
	}
	return CompletedConfig{&completedConfig{Config: c}}
}

// NewClient creates a new client which a high-level API to interact with kafka brokers.
func (c completedConfig) NewClient() (*kafka.Client, error) {
	if c.completeError != nil {
		return nil, c.completeError
	}
	err := c.Validate()
	if err != nil {
		return nil, err
	}
	return &kafka.Client{
		Addr: kafka.TCP(c.Proto.GetAddrs()...),
	}, nil
}

// NewWriter creates a new writer which provides the implementation of a producer of kafka messages
// that automatically distributes messages across partitions of a single topic
// using a configurable balancing policy.
func (c completedConfig) NewWriter() (*kafka.Writer, error) {
	if c.completeError != nil {
		return nil, c.completeError
	}
	return &kafka.Writer{
		Addr:         kafka.TCP(c.Proto.GetAddrs()...),
		Logger:       logrus_.AsStdLogger(logrus.StandardLogger(), logrus.InfoLevel, "", 0),
		ErrorLogger:  logrus_.AsStdLogger(logrus.StandardLogger(), logrus.ErrorLevel, "", 0),
		Balancer:     &kafka.Hash{},
		BatchSize:    100,
		BatchTimeout: 100 * time.Millisecond,
	}, nil
}

// NewReaderConfig creates a new writer that is a configuration object used to create new instances of Reader.
func (c completedConfig) NewReaderConfig() (*kafka.Reader, error) {
	if c.completeError != nil {
		return nil, c.completeError
	}
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     c.Proto.GetAddrs(),
		Logger:      logrus_.AsStdLogger(logrus.StandardLogger(), logrus.InfoLevel, "", 0),
		ErrorLogger: logrus_.AsStdLogger(logrus.StandardLogger(), logrus.ErrorLevel, "", 0),
		MaxWait:     100 * time.Millisecond,
		MinBytes:    1,
	}), nil
}
