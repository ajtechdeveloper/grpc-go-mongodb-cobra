package cmd

import (
	"context"
	"fmt"

	employeepb "github.com/ajtechdeveloper/grpc-go-mongodb-cobra/proto"
	"github.com/spf13/cobra"
)

// deleteCmd represents the read command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Employee by ID",
	Long: `Delete a employee by MongoDB Unique identifier.
	
	If no employee post is found for the ID it will return a 'Not Found' error`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			return err
		}
		req := &employeepb.DeleteEmployeeRequest{
			Id: id,
		}
		_, err = client.DeleteEmployee(context.Background(), req)
		if err != nil {
			return err
		}
		fmt.Printf("Succesfully deleted the employee with id %s\n", id)
		return nil
	},
}

func init() {
	deleteCmd.Flags().StringP("id", "i", "", "The id of the employee")
	deleteCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(deleteCmd)
}
