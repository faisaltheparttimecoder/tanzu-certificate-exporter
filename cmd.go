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
	SkipSsl        bool
	CACertFile     string
}

// Defaults
const (
	defaultInterval   = 86400 // 1 day
	defaultPort       = 8080
	envDebug          = "DEBUG"
	envSslSkip        = "SKIP_SSL_VALIDATION"
	envInterval       = "INTERVAL"
	envPort           = "PORT"
	envOpsManUrl      = "OPSMAN_URL"
	envOpsManUserName = "OPSMAN_USERNAME"
	envOpsManPassword = "OPSMAN_PASSWORD"
	envEnvironment    = "ENVIRONMENT"
	envCaCertFile     = "CACERTFILE"
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
		viper.GetBool(envDebug), "enable verbose or debug logging. Environment Variable: "+envDebug)
	rootCmd.PersistentFlags().BoolVarP(&cmdOptions.SkipSsl, "skip-ssl-validation", "k",
		viper.GetBool(envSslSkip), "skip validating certificate. Environment Variable: "+envSslSkip)
	rootCmd.PersistentFlags().IntVarP(&cmdOptions.Interval, "interval", "i",
		viper.GetInt(envInterval), "scrapping interval in seconds. Environment Variable: "+envInterval)
	rootCmd.PersistentFlags().IntVarP(&cmdOptions.Port, "port", "p",
		viper.GetInt(envPort), "port number to start the web server. Environment Variable: "+envPort)
	rootCmd.PersistentFlags().StringVarP(&cmdOptions.OpsManHostname, "opsman-address", "a",
		viper.GetString(envOpsManUrl),
		"[required] provide the hostname or IP address of the ops manager url. Environment Variable: "+envOpsManUrl)
	rootCmd.PersistentFlags().StringVarP(&cmdOptions.OpsManUsername, "opsman-username", "u",
		viper.GetString(envOpsManUserName),
		"[required] provide the username to connect to ops manager. Environment Variable: "+envOpsManUserName)
	rootCmd.PersistentFlags().StringVarP(&cmdOptions.OpsManPassword, "opsman-password", "w",
		viper.GetString(envOpsManPassword),
		"[required] provide the password to connect to ops manager. Environment Variable: "+envOpsManPassword)
	rootCmd.PersistentFlags().StringVarP(&cmdOptions.Environment, "environment", "e",
		viper.GetString(envEnvironment),
		"[required] provide the environment name for this foundation. Environment Variable: "+envEnvironment)
	rootCmd.PersistentFlags().StringVarP(&cmdOptions.CACertFile, "ca-cert-file", "c",
		viper.GetString(envCaCertFile),
		"[required if skip ssl is false] provide the environment name for this foundation. Environment Variable: "+envCaCertFile)
}
