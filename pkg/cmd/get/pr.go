package get

import (
	"encoding/json"
	"fmt"
	"github.com/divyansh42/tkn-results/pkg/cli"
	"github.com/divyansh42/tkn-results/pkg/cmd/config"
	"github.com/jonboulle/clockwork"
	"github.com/spf13/cobra"
	"github.com/tektoncd/cli/pkg/formatted"
	pb "github.com/tektoncd/results/proto/v1alpha2/results_go_proto"
	"log"
	"text/tabwriter"
	"text/template"

	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

const prListTemplate = `{{- $prl := len .PipelineRuns.Items -}}{{- if eq $prl 0 -}}
No PipelineRuns found
{{ else -}}
{{- if not $.NoHeaders -}}
{{- if $.AllNamespaces -}}
NAMESPACE	NAME	UID	STARTED	DURATION	STATUS
{{ else -}}
NAME	UID	STARTED	DURATION	STATUS
{{ end -}}
{{- end -}}
{{- range $_, $pr := .PipelineRuns.Items }}{{- if $pr }}{{- if $.AllNamespaces -}}
{{ $pr.Namespace }}	{{ $pr.Name }} {{ $pr.UID }}	{{ formatAge $pr.Status.StartTime $.Time }}	{{ formatDuration $pr.Status.StartTime $pr.Status.CompletionTime }}	{{ formatCondition $pr.Status.Conditions }}
{{ else -}}
{{ $pr.Name }}	{{ $pr.UID }}	{{ formatAge $pr.Status.StartTime $.Time }}	{{ formatDuration $pr.Status.StartTime $pr.Status.CompletionTime }}	{{ formatCondition $pr.Status.Conditions }}
{{ end -}}{{- end -}}{{- end -}}
{{- end -}}`

type Options struct {
	Client    *cli.RESTClient
	Namespace string
}

func PrCommand() *cobra.Command {
	o := &Options{}
	c := &cobra.Command{
		Use:     "pr",
		Short:   "Get all PipelineRuns in the namespace",
		PreRunE: o.PreRun,
		RunE:    o.Run,
		//RunE: func(cmd *cobra.Command, args []string) error {
		//	namespace, _ := cmd.Flags().GetString("namespace")
		//	server, _ := cmd.Flags().GetString("server")
		//	token, _ := cmd.Flags().GetString("token")
		//	records, err := records.GetRecords(namespace, server, token)
		//	if err != nil {
		//		return err
		//	}
		//
		//	prs := getPipelineRuns(records)
		//
		//	stream := &cli.Stream{
		//		Out: cmd.OutOrStdout(),
		//		Err: cmd.OutOrStderr(),
		//	}
		//
		//	if prs != nil {
		//		err = printFormattedPr(stream, prs, clockwork.NewRealClock(), false, false)
		//	}
		//	if err != nil {
		//		return fmt.Errorf("failed to print PipelineRuns: %v", err)
		//	}
		//
		//	return nil
		//},
	}
	return c
}

func (o *Options) PreRun(_ *cobra.Command, args []string) (err error) {
	c, err := config.NewConfig()
	if err != nil {
		return err
	}
	o.Client, err = cli.NewRESTClient(c.Get())
	if err != nil {
		return err
	}
	return nil
}

func (o *Options) Run(cmd *cobra.Command, args []string) (err error) {
	namespace, _ := cmd.Flags().GetString("namespace")
	records, err := o.Client.GetRecords(namespace)
	if err != nil {
		return err
	}

	prs := getPipelineRuns(records)

	stream := &cli.Stream{
		Out: cmd.OutOrStdout(),
		Err: cmd.OutOrStderr(),
	}

	if prs != nil {
		err = printFormattedPr(stream, prs, clockwork.NewRealClock(), false, false)
	}
	if err != nil {
		return fmt.Errorf("failed to print PipelineRuns: %v", err)
	}

	return nil
}
func getPipelineRuns(records []*pb.Record) *v1.PipelineRunList {
	var pipelineRuns = new(v1.PipelineRunList)

	for _, record := range records {
		var pr v1.PipelineRun
		if record.Data.Type == "tekton.dev/v1.PipelineRun" || record.Data.Type == "tekton.dev/v1beta1.PipelineRun" {
			if err := json.Unmarshal(record.Data.Value, &pr); err != nil {
				log.Fatalf("Failed to unmarshal JSON: %v", err)
			}
			pipelineRuns.Items = append(pipelineRuns.Items, pr)
		}
	}

	return pipelineRuns
}

func printFormattedPr(s *cli.Stream, prs *v1.PipelineRunList, c clockwork.Clock, allnamespaces bool, noheaders bool) error {
	var data = struct {
		PipelineRuns  *v1.PipelineRunList
		Time          clockwork.Clock
		AllNamespaces bool
		NoHeaders     bool
	}{
		PipelineRuns:  prs,
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
	t := template.Must(template.New("List PipelineRuns").Funcs(funcMap).Parse(prListTemplate))

	err := t.Execute(w, data)
	if err != nil {
		return err
	}

	return w.Flush()
}
