package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Networks map[string]NetworkConfig `mapstructure:"networks"`
	Default  string                   `mapstructure:"default_network"`
}

type NetworkConfig struct {
	RPCURL   string `mapstructure:"rpc_url"`
	ChainID  int64  `mapstructure:"chain_id"`
	Symbol   string `mapstructure:"symbol"`
	Explorer string `mapstructure:"explorer"`
}

// LoadConfig loads the configuration from file and environment variables
func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".tokit")
	configName := "config"
	configFile := filepath.Join(configPath, configName+".yaml")

	// Create default config if it doesn't exist
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := createDefaultConfig(configPath, configFile); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func createDefaultConfig(path, file string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	viper.SetDefault("default_network", "ethereum")
	viper.SetDefault("networks.ethereum.rpc_url", "https://eth.llamarpc.com")
	viper.SetDefault("networks.ethereum.chain_id", 1)
	viper.SetDefault("networks.ethereum.symbol", "ETH")
	viper.SetDefault("networks.ethereum.explorer", "https://etherscan.io")

	viper.SetDefault("networks.arbitrum.rpc_url", "https://arb1.arbitrum.io/rpc")
	viper.SetDefault("networks.arbitrum.chain_id", 42161)
	viper.SetDefault("networks.arbitrum.symbol", "ETH")
	viper.SetDefault("networks.arbitrum.explorer", "https://arbiscan.io")

	viper.SetDefault("networks.optimism.rpc_url", "https://mainnet.optimism.io")
	viper.SetDefault("networks.optimism.chain_id", 10)
	viper.SetDefault("networks.optimism.symbol", "ETH")
	viper.SetDefault("networks.optimism.explorer", "https://optimistic.etherscan.io")

	viper.SetDefault("networks.base.rpc_url", "https://mainnet.base.org")
	viper.SetDefault("networks.base.chain_id", 8453)
	viper.SetDefault("networks.base.symbol", "ETH")
	viper.SetDefault("networks.base.explorer", "https://basescan.org")

	return viper.WriteConfigAs(file)
}
