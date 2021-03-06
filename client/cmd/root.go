package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	employeepb "github.com/ajtechdeveloper/grpc-go-mongodb-cobra/proto"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var cfgFile string

var client employeepb.EmployeeServiceClient
var requestCtx context.Context
var requestOpts grpc.DialOption

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "employeeclient",
	Short: "a gRPC client to communicate with the EmployeeService server",
	Long: `a gRPC client to communicate with the EmployeeService server.
	You can use this client to create and read employees.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.employeeclient.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	fmt.Println("Starting Employee Service Client")
	// Establish the context to timeout if the server does not respond
	requestCtx, _ = context.WithTimeout(context.Background(), 20*time.Second)
	// Establish the insecure grpc options (no TLS)
	requestOpts = grpc.WithInsecure()
	// Dial the server, this returns a client connection
	conn, err := grpc.Dial("localhost:50051", requestOpts)
	if err != nil {
		log.Fatalf("Unable to establish client connection to localhost:50051: %v", err)
	}

	// Instantiate the EmployeeServiceClient with our client connection to the server
	client = employeepb.NewEmployeeServiceClient(conn)
}

// initConfig reads in config file and ENV variables if set.
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

		viper.AddConfigPath(home)
		viper.SetConfigName(".employeeclient")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, then read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
