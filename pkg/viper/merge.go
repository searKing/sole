package viper

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	errors_ "github.com/searKing/golang/go/errors"
	os_ "github.com/searKing/golang/go/os"
	proto_ "github.com/searKing/golang/third_party/github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// merge sequence: protos..., file, env, replace if member has been set
func MergeAll(cfgFile string, envPrefix string, protos ...proto.Message) error {
	// read default config from protobuf
	MergeConfigFromProtoMessages(protos...)
	// read from file
	if err := MergeConfigFromFile(cfgFile); err != nil {
		logrus.WithError(err).Fatalf("load config proto from the file")
		return err
	}

	// read in environment variables that match
	MergeConfigFromENV(envPrefix)
	return nil
}

// read from file
func MergeConfigFromFile(cfgFile string) error {
	if cfgFile != "" {
		// enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
		logrus.WithField("file", cfgFile).Info("using config file")
		file, err := os_.CreateAllIfNotExist(cfgFile)
		if err != nil {
			return fmt.Errorf("create %s %w", cfgFile, err)
		}
		defer file.Close()
	}

	return viper.MergeInConfig()
}

// read from env
func MergeConfigFromENV(envPrefix string) {
	// read in environment variables that match
	viper.AutomaticEnv()          // read in environment variables that match
	viper.SetEnvPrefix(envPrefix) // will be uppercase automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

// read from protobuf
// merge protos into viper one by one, replace if member has been set
// that is, later proto appeared, higher priority proto has
func MergeConfigFromProtoMessages(protos ...proto.Message) {
	viper.SetConfigType("yaml")
	defer viper.SetConfigType("")
	var errs []error
	for _, proto := range protos {
		protoMap, err := proto_.ToGolangMap(proto)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		protoBytes, err := yaml.Marshal(protoMap)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		// merge if not exist
		if err := viper.MergeConfig(bytes.NewReader(protoBytes)); err != nil {
			errs = append(errs, err)
			continue
		}
	}
	var err = errors_.Multi(errs...)
	if err != nil {
		logrus.WithError(err).Warnf("marshal config proto")
		return
	}
}
