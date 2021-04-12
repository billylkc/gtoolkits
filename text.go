package gtoolkits

import (
	"context"
	"fmt"
	"log"
	"time"

	"gitlab.com/billyla/gtoolkits/protos/toolkits"
	"google.golang.org/grpc"
)

const (
	ADDRESS = "localhost:50052"
	TIMEOUT = 2
)

func getRPCConnection(server string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(server,
		grpc.WithTimeout(time.Duration(TIMEOUT)*time.Second),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return conn, fmt.Errorf("%s: %v", "grpc server unreachable", err)
	}
	return conn, nil
}

// GetSummary gets a summary from the text with n sentence
func GetSummary(s string, n int) (string, error) {
	conn, err := getRPCConnection(ADDRESS)
	if err != nil {
		return "", err
	}

	c := toolkits.NewTextClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT*time.Second)
	defer cancel()
	r, err := c.ExtractSummary(ctx, &toolkits.SummaryRequest{Text: s, Count: int64(n)})
	if err != nil {
		return "", err
	}
	return r.GetText(), nil
}

// GetKeywords gets the keywords from a text using_unsupervisored automatic keyword extraction method (YAKE)
func GetKeywords(s string, n int) ([]string, error) {
	var keywords []string
	conn, err := getRPCConnection(ADDRESS)
	if err != nil {
		return keywords, err
	}
	c := toolkits.NewTextClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.ExtractKeywords(ctx, &toolkits.KeywordRequest{Text: s})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	keywords = r.GetText()
	if len(keywords) >= n {
		keywords = keywords[0:n]
	}
	return keywords, nil
}

func GetTfIdf(doc []string, n int) ([]TFRecord, error) {
	var result []TFRecord

	conn, err := getRPCConnection(ADDRESS)
	if err != nil {
		return result, err
	}
	c := toolkits.NewTextClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.ExtractTfIdf(ctx, &toolkits.TFRequest{Text: doc})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// Convert to package struct before export
	records := r.Records
	for _, rec := range records {
		var r TFRecord
		r = TFRecord{
			Keyword: rec.Text,
			Weight:  float64(rec.Score),
		}
		result = append(result, r)
	}
	if len(result) >= n {
		result = result[0:n]
	}
	return result, nil
}
