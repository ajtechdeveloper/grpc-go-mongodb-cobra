package cmd

import (
	"context"
	"fmt"

	employeepb "github.com/ajtechdeveloper/grpc-go-mongodb-cobra/proto"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new employee",
	Long: `Create a new employee on the server through gRPC. 
	
	A employee post requires an Name, Department and Salary.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, err := cmd.Flags().GetString("name")
		department, err := cmd.Flags().GetString("department")
		salary, err := cmd.Flags().GetInt32("salary")
		if err != nil {
			return err
		}
		employee := &employeepb.Employee{
			Name:       name,
			Department: department,
			Salary:  salary,
		}
		res, err := client.CreateEmployee(
			context.TODO(),
			&employeepb.CreateEmployeeRequest{
				Employee: employee,
			},
		)
		if err != nil {
			return err
		}
		fmt.Printf("Employee created: %s\n", res.Employee.Id)
		return nil
	},
}

func init() {
	createCmd.Flags().StringP("name", "n", "", "Add an name")
	createCmd.Flags().StringP("department", "d", "", "A department for the employee")
	createCmd.Flags().Int32P("salary", "s", 1, "The salary for the employee")
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("department")
	createCmd.MarkFlagRequired("salary")
	rootCmd.AddCommand(createCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
