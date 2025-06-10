package cli

import (
	"crypto/tls"
	"fmt"
	pb "github.com/tektoncd/results/proto/v1alpha2/results_go_proto"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"k8s.io/client-go/transport"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const (
	pathPrefix = "parents"
)

type RESTClient struct {
	url    *url.URL
	client *http.Client
}

type Config struct {
	URL       *url.URL
	Timeout   time.Duration
	Transport *transport.Config
}

// NewRESTClient creates a new REST client.
func NewRESTClient(c *Config) (*RESTClient, error) {
	c.URL.Path = path.Join(c.URL.Path, pathPrefix)
	rt, err := transport.New(c.Transport)
	if err != nil {
		return nil, err
	}

	return &RESTClient{
		url: c.URL,
		client: &http.Client{
			Transport: rt,
			Timeout:   c.Timeout,
		},
	}, nil
}

func (c *RESTClient) GetRecords(ns string) ([]*pb.Record, error) {
	pageToken := ""
	orderBy := "create_time desc"
	url := fmt.Sprintf("%s/%s/results/-/records?filter=data_type==PIPELINE_RUN&page_size=100&page_token=%s&order_by=%s", c.url.String(), ns, pageToken, url.QueryEscape(orderBy))
	fmt.Println("Requesting URL:", url)
	req, err := http.NewRequest("GET", url, nil)
	//u := c.url.JoinPath(fmt.Sprintf("%s/results/-/records?page_size=100&page_token=%s", ns, pageToken))
	//url := fmt.Sprintf("%s/apis/results.tekton.dev/v1alpha2/parents/%s/results/-/records", server, namespace)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	//req.Header.Set("Authorization", "Bearer "+token)
	//client := getHTTPClient()
	resp, err := c.client.Do(req)
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

func FetchTrLogs(server, token, recordName string) ([]byte, error) {
	// API end point is apis/results.tekton.dev/v1alpha2/parents/<namespace>/results/<result-id>/logs/<record-id>
	recordName = strings.Replace(recordName, "records", "logs", -1)
	url := fmt.Sprintf("%s/apis/results.tekton.dev/v1alpha3/parents/%s", server, recordName)
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
		return nil, fmt.Errorf("failed to fetch logs: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getHTTPClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: transport}
	//return http.DefaultClient
}
