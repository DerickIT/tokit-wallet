package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Networks       map[string]NetworkConfig `mapstructure:"networks"`
	Default        string                   `mapstructure:"default_network"`
	DefaultAccount string                   `mapstructure:"default_account"`
}

type NetworkConfig struct {
	RPCURL   string `mapstructure:"rpc_url"`
	ChainID  int64  `mapstructure:"chain_id"`
	Symbol   string `mapstructure:"symbol"`
	Explorer string `mapstructure:"explorer"`
}

// LoadConfig loads the configuration from file and environment variables.
func LoadConfig() (*Config, error) {
	configPath, configFile, err := configLocation()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := createDefaultConfig(configPath, configFile); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
	}

	v := newViper(configPath)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// SaveDefaultAccount persists the preferred local account address.
func SaveDefaultAccount(address string) error {
	configPath, configFile, err := configLocation()
	if err != nil {
		return err
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := createDefaultConfig(configPath, configFile); err != nil {
			return fmt.Errorf("failed to create default config: %w", err)
		}
	}

	v := newViper(configPath)
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	v.Set("default_account", address)
	if err := v.WriteConfig(); err != nil {
		return fmt.Errorf("failed to update config file: %w", err)
	}

	return nil
}

func configLocation() (string, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}

	configPath := filepath.Join(home, ".tokit")
	configFile := filepath.Join(configPath, "config.yaml")
	return configPath, configFile, nil
}

func newViper(configPath string) *viper.Viper {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	return v
}

func createDefaultConfig(path, file string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	v := viper.New()
	v.SetConfigFile(file)
	v.SetConfigType("yaml")
	v.SetDefault("default_network", "ethereum")
	v.SetDefault("default_account", "")
	v.SetDefault("networks.ethereum.rpc_url", "https://eth.llamarpc.com")
	v.SetDefault("networks.ethereum.chain_id", 1)
	v.SetDefault("networks.ethereum.symbol", "ETH")
	v.SetDefault("networks.ethereum.explorer", "https://etherscan.io")

	v.SetDefault("networks.arbitrum.rpc_url", "https://arb1.arbitrum.io/rpc")
	v.SetDefault("networks.arbitrum.chain_id", 42161)
	v.SetDefault("networks.arbitrum.symbol", "ETH")
	v.SetDefault("networks.arbitrum.explorer", "https://arbiscan.io")

	v.SetDefault("networks.optimism.rpc_url", "https://mainnet.optimism.io")
	v.SetDefault("networks.optimism.chain_id", 10)
	v.SetDefault("networks.optimism.symbol", "ETH")
	v.SetDefault("networks.optimism.explorer", "https://optimistic.etherscan.io")

	v.SetDefault("networks.base.rpc_url", "https://mainnet.base.org")
	v.SetDefault("networks.base.chain_id", 8453)
	v.SetDefault("networks.base.symbol", "ETH")
	v.SetDefault("networks.base.explorer", "https://basescan.org")

	return v.WriteConfigAs(file)
}
