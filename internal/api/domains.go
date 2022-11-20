package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/raylas/nextdns-exporter/internal/util"
)

type DomainsResponse struct {
	Domains []Domain `json:"data"`
}

type Domain struct {
	Domain  string `json:"domain"`
	Root    string `json:"root"`
	Tracker string `json:"tracker"`
	Queries int    `json:"queries"`
}

type DomainsMetrics struct {
	BlockedDomains []DomainMetric
}

type DomainMetric struct {
	Domain  string
	Root    string
	Tracker string
	Queries float64
}

func (c Client) CollectDomains(profile, apiKey string) (*DomainsMetrics, error) {
	domainsURL := fmt.Sprintf("%s/profiles/%s/analytics/domains", c.url, profile)

	domainsResponse := DomainsResponse{}
	metrics := DomainsMetrics{}

	req, err := http.NewRequest("GET", domainsURL, nil)
	if err != nil {
		util.Log.Error("error creating request", "error", err)
		return nil, err
	}
	req.Header.Set("X-Api-Key", apiKey)
	req.URL.RawQuery = url.Values{
		"from":   {util.FilterFrom},
		"limit":  {util.ResultLimit},
		"status": {"blocked"},
	}.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		util.Log.Error("error making request", "error", err)
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		util.Log.Error("error reading response body", "error", err)
		return nil, err
	}

	err = json.Unmarshal(body, &domainsResponse)
	if err != nil {
		util.Log.Error("error unmarshalling response body", "error", err)
		return nil, err
	}

	for _, domain := range domainsResponse.Domains {
		domain := DomainMetric{
			Domain:  domain.Domain,
			Root:    domain.Root,
			Tracker: domain.Tracker,
			Queries: float64(domain.Queries),
		}

		metrics.BlockedDomains = append(metrics.BlockedDomains, domain)
	}

	return &metrics, nil
}
