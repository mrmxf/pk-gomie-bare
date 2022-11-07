package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config = viper.Viper

var cfg *Config

/** entry point for configuration
 *
 */
func GetConfig(forceConfigName ...string) *Config {
	if cfg != nil {
		return cfg
	}

	//load any .env file - ignore consequences
	_ = godotenv.Load(".env")

	thisExecutable, err := os.Executable()
	if err != nil {
		panic(err)
	}

	//just take the basename and remove the extension
	exeBase := filepath.Base(thisExecutable)
	exeName := exeBase[:len(exeBase)-len(filepath.Ext(exeBase))]

	configPaths := []string{
		".",
		"$HOME",
		filepath.Join("$HOME", "."+exeName),
		exeName,
	}

	cfg = viper.New()
	if len(forceConfigName) > 0 {
		cfg.SetConfigName(forceConfigName[0])
	} else {
		cfg.SetConfigName("config-" + exeName) // name of config file (without extension)
	}
	cfg.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name

	pathStrings := "["
	for _, s := range configPaths {
		pathStrings += "\"" + s + "\",  "
		cfg.AddConfigPath(s)
	}

	err = cfg.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		log.Printf("fatal error in config file: %v", cfg.ConfigFileUsed())
		log.Fatalf("fatal error: %v", err)
	}
	cfg.AutomaticEnv()
	cfg.BindEnv("LIMELM_API_KEY")

	// Now that config is loaded we can safely fire up the logger (which needs config)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix(cfg.GetString("app_name"))

	return cfg
}
