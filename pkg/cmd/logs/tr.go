package logs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/divyansh42/tkn-results/pkg/cli"
	"github.com/divyansh42/tkn-results/pkg/records"
	"github.com/spf13/cobra"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	pb "github.com/tektoncd/results/proto/v1alpha2/results_go_proto"
	"log"
)

func TrCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "tr [name]",
		Short: "Get logs of the TaskRun in the namespace",
		RunE: func(cmd *cobra.Command, args []string) error {
			namespace, _ := cmd.Flags().GetString("namespace")
			server, _ := cmd.Flags().GetString("server")
			token, _ := cmd.Flags().GetString("token")
			// TaskRun name to fetch the logs
			trName := args[0]
			recordsData, err := records.GetRecords(namespace, server, token)
			if err != nil {
				return err
			}
			recordName := getRecordIdToFetchLogs(recordsData, trName)
			if recordName == "" {
				return errors.New("No TaskRun found for " + trName)
			}

			logs, err := records.FetchTrLogs(server, token, recordName)
			if err != nil {
				return err
			}

			stream := &cli.Stream{
				Out: cmd.OutOrStdout(),
				Err: cmd.OutOrStderr(),
			}

			_, err = stream.Out.Write(logs)

			if err != nil {
				return fmt.Errorf("failed to print PipelineRuns: %v", err)
			}

			return nil
		},
	}
}

func getRecordIdToFetchLogs(records []*pb.Record, trName string) string {
	for _, record := range records {
		var tr v1.TaskRun
		if record.Data.Type == "tekton.dev/v1.TaskRun" {
			if err := json.Unmarshal(record.Data.Value, &tr); err != nil {
				log.Fatalf("Failed to unmarshal JSON: %v", err)
			}
			if tr.Name == trName {
				// Record Name with hold value like "<namespace>/results/<results-id>/records/<record-id>"
				return record.Name
			}
		}
	}
	return ""
}
