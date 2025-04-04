package main

import (
	"fmt"
	"os"

	myviper "github.com/ARtorias742/viper"
	"github.com/spf13/pflag"
)

func main() {
	v := myviper.New()

	// Configure Viper
	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath(".")

	// Set some default values
	v.Set("port", 8080)
	v.Set("name", "myapp")
	v.Set("debug", false)

	// Define command-line flags
	flags := pflag.NewFlagSet("example", pflag.ExitOnError)
	flags.String("port", "", "Port to run the application on")
	flags.String("name", "", "Name of the application")
	flags.Bool("debug", false, "Enable debug mode")
	flags.Parse(os.Args[1:])

	// Bind flags to Viper
	v.BindPFlags(flags)

	// Enable environment variables with a prefix
	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	// Set an environment variable for testing
	os.Setenv("APP_PORT", "9090")

	// Write config to file (optional)
	if err := v.WriteConfig(); err != nil {
		fmt.Println("Error writing config:", err)
	}

	// Read config from file
	if err := v.ReadInConfig(); err != nil {
		fmt.Println("Error reading config:", err)
	} else {
		fmt.Println("Config loaded successfully")
	}

	// Access values
	fmt.Println("Port:", v.Get("port")) // Should pick up "APP_PORT"
	fmt.Println("Name:", v.GetString("name"))
	fmt.Println("Debug:", v.Get("debug"))
}
