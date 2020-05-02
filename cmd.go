package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/url"
)

// Global Parameter
var (
	cmdOptions Command
)

// Root command line options
type Command struct {
	Debug          bool
	Environment    string
	Port           int
	OpsManHostname string
	OpsManUsername string
	OpsManPassword string
	Interval       int
	SkipSsl		   bool
	CACertFile	   string
}

// Defaults
const (
	defaultInterval = 86400 // 1 day
	defaultPort     = 8080
)

// The root commands.
var rootCmd = &cobra.Command{
	Use:   fmt.Sprintf("%s", programName),
	Short: "VmWare Tanzu Certificate Exporter",
	Long: "This application is designed to extract the certificate information from " +
		"vmware tanzu operation manager and for prometheus to scrape.",
	Version: programVersion,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Before running any command setup the logger log level
		initLogger(cmdOptions.Debug)

		// Define defaults for arguments that are missing or error out for which there
		// is no error
		setDefaultsOrErrorIfMissing()
	},
	Run: func(cmd *cobra.Command, args []string) {
		startHttpServer()
	},
}

// Set defaults of values of fields that are missed by the user
func setDefaultsOrErrorIfMissing() {
	suffixText := "is a required flag and its missing"
	// TODO: no need for interval set the grafana dashboard to 1 day to refresh
	if cmdOptions.Interval == 0 { // Interval
		cmdOptions.Interval = defaultInterval
	}
	if cmdOptions.Port == 0 { // Port
		cmdOptions.Port = defaultPort
	}
	if cmdOptions.OpsManHostname == "" { // ops man url
		Fatalf("Operation manager URL %s", suffixText)
	}
	if cmdOptions.OpsManUsername == "" { // ops man username
		Fatalf("Operation manager Username %s", suffixText)
	}
	if cmdOptions.OpsManPassword == "" { // ops man password
		Fatalf("Operation manager Password %s", suffixText)
	}
	if cmdOptions.Environment == "" { // environment
		Fatalf("Environment %s", suffixText)
	}
	if !cmdOptions.SkipSsl && cmdOptions.CACertFile == "" { // ca cert file is missing
		Fatalf("CA cert file parameter %s", suffixText)
	}
	if cmdOptions.OpsManHostname != "" { // validate ops man URL
		u, err := url.ParseRequestURI(cmdOptions.OpsManHostname)
		if err != nil {
			Fatalf("Invalid URL (%s), err: %v", cmdOptions.OpsManHostname, err)
		}
		cmdOptions.OpsManHostname = u.Host
	}
}

// Initialize the cobra command line
func init() {
	// Load the environment variable using viper
	viper.AutomaticEnv()

	// Root command flags
	rootCmd.PersistentFlags().BoolVarP(&cmdOptions.Debug, "debug", "d",
		viper.GetBool("DEBUG"), "enable verbose or debug logging. Environment Variable: DEBUG")
	rootCmd.PersistentFlags().BoolVarP(&cmdOptions.SkipSsl, "skip-ssl-validation", "k",
		viper.GetBool("SKIP_SSL_VALIDATION"), "skip validating certificate. Environment Variable: SKIP_SSL_VALIDATION")
	rootCmd.PersistentFlags().IntVarP(&cmdOptions.Interval, "interval", "i",
		viper.GetInt("INTERVAL"), "scrapping interval in seconds. Environment Variable: INTERVAL")
	rootCmd.PersistentFlags().IntVarP(&cmdOptions.Port, "port", "p",
		viper.GetInt("PORT"), "port number to start the web server. Environment Variable: PORT")
	rootCmd.PersistentFlags().StringVarP(&cmdOptions.OpsManHostname, "opsman-address", "a",
		viper.GetString("OPSMAN_URL"),
		"(required) provide the hostname or IP address of the ops manager url. Environment Variable: OPSMAN_URL")
	rootCmd.PersistentFlags().StringVarP(&cmdOptions.OpsManUsername, "opsman-username", "u",
		viper.GetString("OPSMAN_USERNAME"),
		"(required) provide the username to connect to ops manager. Environment Variable: OPSMAN_USERNAME")
	rootCmd.PersistentFlags().StringVarP(&cmdOptions.OpsManPassword, "opsman-password", "w",
		viper.GetString("OPSMAN_PASSWORD"),
		"(required) provide the password to connect to ops manager. Environment Variable: OPSMAN_PASSWORD")
	rootCmd.PersistentFlags().StringVarP(&cmdOptions.Environment, "environment", "e",
		viper.GetString("ENVIRONMENT"),
		"(required) provide the environment name for this foundation. Environment Variable: ENVIRONMENT")
	rootCmd.PersistentFlags().StringVarP(&cmdOptions.CACertFile, "ca-cert-file", "c",
		viper.GetString("CACERTFILE"),
		"provide the environment name for this foundation. Environment Variable: CACERTFILE")
}
