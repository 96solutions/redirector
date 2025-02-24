package logger

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

// OpenSearchWriter is a custom type that implements the Write method/interface
// for logging directly to OpenSearch, without the help of Logstash.
type OpenSearchWriter struct {
	Client *opensearch.Client
	index  string
}

// NewOpenSearchWriter creates an instance of OpenSearchWriter which can be used for logging.
func NewOpenSearchWriter(host, port, indexName, user, pass string) (*OpenSearchWriter, error) {
	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: []string{"https://" + host + ":" + port},
		Username:  user,
		Password:  pass,
	})

	return &OpenSearchWriter{Client: client, index: indexName}, err
}

// Write method to write logs directly to OpenSearch.
func (ow *OpenSearchWriter) Write(p []byte) (n int, err error) {
	req := opensearchapi.IndexRequest{
		Index: ow.index,
		Body:  strings.NewReader(string(p)),
	}

	insertResponse, err := req.Do(context.Background(), ow.Client)
	if err != nil {
		return 0, err
	}

	defer insertResponse.Body.Close()

	return len(p), nil
}
