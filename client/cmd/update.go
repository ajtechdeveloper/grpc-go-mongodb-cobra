package cmd

import (
	"context"
	"fmt"

	employeepb "github.com/ajtechdeveloper/grpc-go-mongodb-cobra/proto"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Find a Employee by ID",
	Long: `Find a employee by MongoDB Unique identifier.
	
	If no employee post is found for the ID it will return a 'Not Found' error`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the flags from CLI
		id, err := cmd.Flags().GetString("id")
		name, err := cmd.Flags().GetString("name")
		department, err := cmd.Flags().GetString("department")
		salary, err := cmd.Flags().GetInt32("salary")
		if err != nil {
			return err
		}
		employee := &employeepb.Employee{
			Id:         id,
			Name:       name,
			Department: department,
			Salary:  salary,
		}
		// Create UpdateEmployeeRequest
		res, err := client.UpdateEmployee(
			context.TODO(),
			&employeepb.UpdateEmployeeRequest{
				Employee: employee,
			},
		)
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	},
}

func init() {
	updateCmd.Flags().StringP("id", "i", "", "The id of the employee")
	updateCmd.Flags().StringP("name", "n", "", "Add an name")
	updateCmd.Flags().StringP("department", "d", "", "A department for the employee")
	updateCmd.Flags().Int32P("salary", "s", 1, "The salary for the employee")
	updateCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(updateCmd)
}
