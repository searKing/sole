package viper

import (
	"encoding/json"
	"reflect"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Unmarshal persists and returns the latest config viper proto
func Unmarshal(v proto.Message) error {
	// config file -> ViperProto
	if err := viper.Unmarshal(v, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.TagName = "json" // trick of protobuf, which generates json tag only
		decoderConfig.WeaklyTypedInput = true
		decoderConfig.DecodeHook = anyToProtoMessageHookFunc(v)
	}); err != nil {
		logrus.WithError(err).Warnf("ignore config file changed")
		return err
	}

	// persist using config
	return PersistConfig()
}

// anyToProtoMessageHookFunc returns a DecodeHookFunc that converts
// root struct to config.ViperProto.
// Trick of protobuf, which generates json tag only
func anyToProtoMessageHookFunc(v proto.Message) mapstructure.DecodeHookFunc {
	return func(src reflect.Type, dst reflect.Type, data interface{}) (interface{}, error) {
		// Convert it by parsing
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		// apply protobuf check
		if err := jsonpb.UnmarshalString(string(dataBytes), v); err != nil {
			return data, err
		}
		return v, nil
	}
}
