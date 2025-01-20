package get

import (
	"encoding/json"
	"fmt"
	"github.com/divyansh42/tkn-results/pkg/cli"
	"github.com/divyansh42/tkn-results/pkg/records"
	"github.com/jonboulle/clockwork"
	"github.com/spf13/cobra"
	"github.com/tektoncd/cli/pkg/formatted"
	pb "github.com/tektoncd/results/proto/v1alpha2/results_go_proto"
	"log"
	"text/tabwriter"
	"text/template"

	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

const trListTemplate = `{{- $trl := len .TaskRuns.Items -}}{{- if eq $trl 0 -}}
No TaskRuns found
{{ else -}}
{{- if not $.NoHeaders -}}
{{- if $.AllNamespaces -}}
NAMESPACE	NAME	UID	STARTED	DURATION	STATUS
{{ else -}}
NAME	UID	STARTED	DURATION	STATUS
{{ end -}}
{{- end -}}
{{- range $_, $tr := .TaskRuns.Items }}{{- if $tr }}{{- if $.AllNamespaces -}}
{{ $tr.Namespace }}	{{ $tr.Name }} {{ $tr.UID }}	{{ formatAge $tr.Status.StartTime $.Time }}	{{ formatDuration $tr.Status.StartTime $tr.Status.CompletionTime }}	{{ formatCondition $tr.Status.Conditions }}
{{ else -}}
{{ $tr.Name }}	{{ $tr.UID }}	{{ formatAge $tr.Status.StartTime $.Time }}	{{ formatDuration $tr.Status.StartTime $tr.Status.CompletionTime }}	{{ formatCondition $tr.Status.Conditions }}
{{ end -}}{{- end -}}{{- end -}}
{{- end -}}`

func TrCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "tr",
		Short: "Get all TaskRuns in the namespace",
		RunE: func(cmd *cobra.Command, args []string) error {
			namespace, _ := cmd.Flags().GetString("namespace")
			server, _ := cmd.Flags().GetString("server")
			token, _ := cmd.Flags().GetString("token")
			records, err := records.GetRecords(namespace, server, token)
			if err != nil {
				return err
			}

			trs := getTaskRuns(records)

			stream := &cli.Stream{
				Out: cmd.OutOrStdout(),
				Err: cmd.OutOrStderr(),
			}

			if trs != nil {
				err = printFormattedTr(stream, trs, clockwork.NewRealClock(), false, false)
			}
			if err != nil {
				return fmt.Errorf("failed to print TaskRuns: %v", err)
			}

			return nil
		},
	}
}

func getTaskRuns(records []*pb.Record) *v1.TaskRunList {
	var pipelineRuns = new(v1.TaskRunList)

	for _, record := range records {
		var tr v1.TaskRun
		if record.Data.Type == "tekton.dev/v1.TaskRun" {
			if err := json.Unmarshal(record.Data.Value, &tr); err != nil {
				log.Fatalf("Failed to unmarshal JSON: %v", err)
			}
			pipelineRuns.Items = append(pipelineRuns.Items, tr)
		}
	}

	return pipelineRuns
}

func printFormattedTr(s *cli.Stream, trs *v1.TaskRunList, c clockwork.Clock, allnamespaces bool, noheaders bool) error {
	var data = struct {
		TaskRuns      *v1.TaskRunList
		Time          clockwork.Clock
		AllNamespaces bool
		NoHeaders     bool
	}{
		TaskRuns:      trs,
		Time:          c,
		AllNamespaces: allnamespaces,
		NoHeaders:     noheaders,
	}

	funcMap := template.FuncMap{
		"formatAge":       formatted.Age,
		"formatDuration":  formatted.Duration,
		"formatCondition": formatted.Condition,
	}

	w := tabwriter.NewWriter(s.Out, 0, 5, 3, ' ', tabwriter.TabIndent)
	t := template.Must(template.New("List TaskRuns").Funcs(funcMap).Parse(trListTemplate))

	err := t.Execute(w, data)
	if err != nil {
		return err
	}

	return w.Flush()
}
