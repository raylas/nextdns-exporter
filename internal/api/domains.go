package api

import (
	"encoding/json"
	"fmt"
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

func (c Client) CollectDomains() (*DomainsMetrics, error) {
	domainsURL := fmt.Sprintf("%s/profiles/%s/analytics/domains", c.url, c.profile)

	domainsResponse := DomainsResponse{}
	metrics := DomainsMetrics{}

	params := url.Values{
		"status": {"blocked"},
	}

	body, err := c.Request(domainsURL, params)
	if err != nil {
		util.Log.Error("error making request", "error", err)
		return nil, err
	}

	err = json.Unmarshal(body, &domainsResponse)
	if err != nil {
		util.Log.Error("error unmarshalling response body", "error", err)
		return nil, err
	}

	for _, domain := range domainsResponse.Domains {
		// Some entries appear not to have a root, in which case replicate the domain.
		// https://github.com/raylas/nextdns-exporter/issues/20
		if len(domain.Root) == 0 {
			domain.Root = domain.Domain
		}
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
