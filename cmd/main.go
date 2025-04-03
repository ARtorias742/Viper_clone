package main

import (
	"fmt"

	myviper "github.com/ARtorias742/viper"
)

func main() {
	v := myviper.New()

	// Configure Viper
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	// Set some default values
	v.Set("port", 8080)
	v.Set("name", "myapp")

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
	fmt.Println("Port:", v.Get("port"))
	fmt.Println("Name:", v.GetString("name"))
}
