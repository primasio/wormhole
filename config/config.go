/*
 * Copyright 2018 Primas Lab Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"errors"
	"github.com/spf13/viper"
)

var config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string, configPath *string) error {
	var err error
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName(env)

	if configPath != nil {
		v.AddConfigPath(*configPath)
	} else {
		v.AddConfigPath("config/")
	}

	err = v.ReadInConfig()
	if err != nil {
		return errors.New("error on parsing configuration file")
	}
	config = v

	return nil
}

func GetConfig() *viper.Viper {
	return config
}
