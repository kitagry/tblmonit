package cmd

import (
	"fmt"
	"os"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var cfg tmConfig

type tmConfig struct {
	timeZone string
}

var verbose, debug bool // for verbose and debug output

// rootCmd represents the root command
var rootCmd = &cobra.Command{
	Use:   "tblmonit",
	Short: "Monitoring tool for Bigquery tables",
	Long:  `Monitoring BigQuery table metadata to ensure the data pipeline jobs are correctly worked.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".tblmonit" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tblmonit")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Failed to read Config File", viper.ConfigFileUsed(), err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Println("Failed to read Config File", viper.ConfigFileUsed(), err)
		os.Exit(1)
	}

	loadTimezone()
	logOutput()
}

func loadTimezone() {
	var err error
	var loc *time.Location
	if cfg.timeZone != "" {
		loc, err = time.LoadLocation(cfg.timeZone)
		if err != nil {
			fmt.Println("Failed to load location from config file", cfg.timeZone)
		}
		time.Local = loc
	}
	time.Local = time.UTC
}

func logOutput() {
	zerolog.SetGlobalLevel(zerolog.Disabled) // default: quiet mode
	switch {
	case verbose:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case debug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tblmonit.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// for log output
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable varbose log output")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug log output")
}
