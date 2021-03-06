package cmd

import (
	"context"
	"fmt"
	"io"

	employeepb "github.com/ajtechdeveloper/grpc-go-mongodb-cobra/proto"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Get all employees",
	RunE: func(cmd *cobra.Command, args []string) error {
		req := &employeepb.GetAllEmployeesRequest{}
		// Call GetAllEmployees that returns a stream
		stream, err := client.GetAllEmployees(context.Background(), req)
		if err != nil {
			return err
		}
		// Iterate
		for {
			res, err := stream.Recv()
			// If it is end of the stream, then break the loop
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			fmt.Println(res.GetEmployee())
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
