// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	errors_ "github.com/searKing/golang/go/errors"
	os_ "github.com/searKing/golang/go/os"
	proto_ "github.com/searKing/golang/third_party/github.com/golang/protobuf/proto"
	viper_ "github.com/searKing/sole/api/protobuf-spec/v1/viper"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Load load config from file and protos, and save to a using file
// load sequence: protos..., file, env, replace if member has been set
// that is, later cfg appeared, higher priority cfg has
func Load(cfgFile string, protos ...*viper_.ViperProto) (*viper_.ViperProto, error) {
	// read default config from protobuf
	mergeConfigFromProto(protos...)
	// read from file
	if cfgFile != "" {
		if err := mergeConfigFromFile(cfgFile); err != nil {
			err = errors.WithMessage(err, "load config proto from the file failed")
			log.Fatalf("[WARN] %s\n", err)
			return nil, err
		}
	}
	// read in environment variables that match
	mergeConfigFromENV()

	// dump config
	return dumpConfig()
}

// anyToViperProtoHookFunc returns a DecodeHookFunc that converts
// root struct to config.ViperProto.
// Trick of protobuf, which generates json tag only
func anyToViperProtoHookFunc() mapstructure.DecodeHookFunc {
	return func(src reflect.Type, dst reflect.Type, data interface{}) (interface{}, error) {
		// Convert it by parsing
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		var msg viper_.ViperProto
		// apply protobuf check
		if err := jsonpb.UnmarshalString(string(dataBytes), &msg); err != nil {
			return data, err
		}
		return &msg, nil
	}
}

// dumpConfig persists and returns the latest config viper proto
func dumpConfig() (*viper_.ViperProto, error) {
	var using viper_.ViperProto
	// config file -> ViperProto
	if err := viper.Unmarshal(&using, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.TagName = "json" // trick of protobuf, which generates json tag only
		decoderConfig.WeaklyTypedInput = true
		decoderConfig.DecodeHook = anyToViperProtoHookFunc()
	}); err != nil {
		log.Printf(`[WARN] %s`,
			errors.WithMessagef(err, "ignore config file changed"))
		return nil, err
	}

	// persist using config
	return &using, persistConfig()
}

// read from file
func mergeConfigFromFile(cfgFile string) error {
	if cfgFile != "" {
		// enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
		log.Printf("[INFO] Using config file: %s\n", cfgFile)
		if err := os_.CreateAllIfNotExist(cfgFile, 0666); err != nil {
			return fmt.Errorf("create %s %w", cfgFile, err)
		}
	}

	return viper.MergeInConfig()
}

// read from env
func mergeConfigFromENV() {
	// read in environment variables that match
	viper.AutomaticEnv()            // read in environment variables that match
	viper.SetEnvPrefix(ServiceName) // will be uppercase automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

// read from protobuf
// merge protos into viper one by one, replace if member has been set
// that is, later proto appeared, higher priority proto has
func mergeConfigFromProto(protos ...*viper_.ViperProto) {
	viper.SetConfigType("yaml")
	defer viper.SetConfigType("")
	var marshalErrs []error
	var mergeErrs []error
	for _, proto := range protos {
		protoMap, err := proto_.ToGolangMap(proto)
		if err != nil {
			marshalErrs = append(marshalErrs, err)
			continue
		}
		protoBytes, err := yaml.Marshal(protoMap)
		if err != nil {
			marshalErrs = append(marshalErrs, err)
			continue
		}

		// merge if not exist
		if err := viper.MergeConfig(bytes.NewReader(protoBytes)); err != nil {
			mergeErrs = append(mergeErrs, err)
			continue
		}
	}
	if len(marshalErrs) > 0 {
		log.Printf("[WARN] %s\n",
			errors.WithMessage(errors_.Multi(marshalErrs...), "marshal config proto failed"))
	}
	if len(mergeErrs) > 0 {
		log.Printf("[WARN] %s\n",
			errors.WithMessage(errors_.Multi(mergeErrs...), "merge config proto failed"))
	}
}
