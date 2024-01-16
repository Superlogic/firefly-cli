// Copyright © 2021 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ethsigner

import (
	"os"

	"github.com/hyperledger/firefly-cli/pkg/types"
	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"gopkg.in/yaml.v2"
)

type FileWalletFilenamesConfig struct {
	With0xPrefix      bool   `yaml:"with0xPrefix,omitempty"`
	PrimaryExt        string `yaml:"primaryExt,omitempty"`
	PasswordExt       string `yaml:"passwordExt,omitempty"`
	PrimaryMatchRegex string `yaml:"primaryMatchRegex,omitempty"`
}

type FileWalletMetadataConfig struct {
	Format               string `yaml:"format,omitempty"`
	KeyFileProperty      string `yaml:"keyFileProperty,omitempty"`
	PasswordFileProperty string `yaml:"passwordFileProperty,omitempty"`
}

type FileWalletConfig struct {
	Path                string                     `yaml:"path,omitempty"`
	DefaultPasswordFile string                     `yaml:"defaultPasswordFile,omitempty"`
	Filenames           *FileWalletFilenamesConfig `yaml:"filenames,omitempty"`
	Metadata            *FileWalletMetadataConfig  `yaml:"metadata,omitempty"`
}

type MPCWalletConfig struct {
	Path      string                     `yaml:"path,omitempty"`
	URL       string                     `yaml:"url,omitempty"`
	Filenames *FileWalletFilenamesConfig `yaml:"filenames,omitempty"`
	Enabled   bool                       `yaml:"enabled,omitempty"`
}

type ServerConfig struct {
	Port    int    `yaml:"port,omitempty"`
	Address string `yaml:"address,omitempty"`
}

type BackendConfig struct {
	ChainID *int64 `yaml:"chainId,omitempty"`
	URL     string `yaml:"url,omitempty"`
}

type LogConfig struct {
	Level string `yaml:"level,omitempty"`
}

type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Backend    BackendConfig    `yaml:"backend"`
	FileWallet FileWalletConfig `yaml:"fileWallet"`
	MPCWallet  MPCWalletConfig  `yaml:"mpcWallet"`
	Log        LogConfig        `yaml:"log"`
}

func (e *Config) WriteConfig(filename string) error {
	configYamlBytes, _ := yaml.Marshal(e)
	return os.WriteFile(filename, configYamlBytes, 0755)
}

func GenerateSignerConfig(signerType fftypes.FFEnum, chainID int64, rpcURL string, remoteSignerURL string) *Config {
	config := &Config{
		Server: ServerConfig{
			Port:    8545,
			Address: "0.0.0.0",
		},
		Backend: BackendConfig{
			URL:     rpcURL,
			ChainID: &chainID,
		},
		Log: LogConfig{
			Level: "debug",
		},
	}

	switch signerType {
	case types.SignerTypeMPC:
		config.MPCWallet = MPCWalletConfig{
			Path: "/data/keystore",
			URL:  remoteSignerURL,
			Filenames: &FileWalletFilenamesConfig{
				PrimaryMatchRegex: "^((0x)?[0-9a-z]+).key.json$",
			},
			Enabled: true,
		}
	default:
		config.FileWallet = FileWalletConfig{
			Path: "/data/keystore",
			Filenames: &FileWalletFilenamesConfig{
				PrimaryExt: ".toml",
			},
			Metadata: &FileWalletMetadataConfig{
				KeyFileProperty:      `{{ index .signing "key-file" }}`,
				PasswordFileProperty: `{{ index .signing "password-file" }}`,
			},
		}
	}

	return config
}
