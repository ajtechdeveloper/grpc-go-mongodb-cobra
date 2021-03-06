package cmd

import (
	"context"
	"fmt"

	employeepb "github.com/ajtechdeveloper/grpc-go-mongodb-cobra/proto"
	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Find a Employee by ID",
	Long: `Find a employee by MongoDB Unique identifier.
	
	If no employee is found for the ID it will return a 'Not Found' error`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			return err
		}
		req := &employeepb.GetEmployeeRequest{
			Id: id,
		}
		res, err := client.GetEmployee(context.Background(), req)
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	},
}

func init() {
	readCmd.Flags().StringP("id", "i", "", "The id of the employee")
	readCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(readCmd)
}
