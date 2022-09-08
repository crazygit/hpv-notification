package cmd

import (
	"fmt"
	"github.com/crazygit/hpv-notification/config"
	"github.com/crazygit/hpv-notification/internal/dal"
	"github.com/crazygit/hpv-notification/internal/util"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hpv-notification",
	Short: "HPV notification",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("Execute error: %v\n", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig, util.ConfigLog, dal.InitDBInstance)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is %s/config.yaml)", util.RootDir()))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// load default config
		viper.AddConfigPath(util.RootDir())
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		rootCmd.Println("Using config file:", viper.ConfigFileUsed())
	}
	err := viper.Unmarshal(&config.AppConfig, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.ErrorUnused = true
		decoderConfig.ErrorUnset = true
	})
	if err != nil {
		log.Fatalf("Unmarshal app config failed, error: %v\n", err)
	}
}
