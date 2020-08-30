package main

import (
	"github.com/spf13/viper"
)

func configDefaults() {
	// viper.SetDefault("default.i2c-addr", 0x70)

	// viper.SetDefault("ads1115.i2c-addr", 0x49)
}

func configInit() {
	viper.SetConfigName("dfr0591")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./artifacts")
	viper.AddConfigPath("./test")
	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("/etc/dfr0591")
	viper.AddConfigPath("/usr/local/dfr0591/etc")
	viper.AddConfigPath("/run/secrets")
	viper.AddConfigPath(".")

	configDefaults()
}
