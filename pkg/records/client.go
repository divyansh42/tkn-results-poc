package records

import (
	"crypto/tls"
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"log"
	"net/http"

	pb "github.com/tektoncd/results/proto/v1alpha2/results_go_proto"
)

func GetRecords(namespace, server, token string) ([]*pb.Record, error) {
	url := fmt.Sprintf("%s/apis/results.tekton.dev/v1alpha2/parents/%s/results/-/records", server, namespace)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	client := getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch records: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var records pb.ListRecordsResponse
	if err := protojson.Unmarshal(body, &records); err != nil {
		log.Fatalf("Failed to decode record: %v", err)
	}
	return records.Records, nil
}

func getHTTPClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: transport}
	//return http.DefaultClient
}
