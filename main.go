package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/raylas/nextdns-exporter/internal/api"
	"github.com/raylas/nextdns-exporter/internal/util"
)

type exporter struct {
	profile, apiKey string

	totalQueries        *prometheus.Desc
	totalAllowedQueries *prometheus.Desc
	totalBlockedQueries *prometheus.Desc

	blockedQueries *prometheus.Desc
}

func newExporter(profile, apiKey string) *exporter {
	return &exporter{
		profile: profile,
		apiKey:  apiKey,

		totalQueries: prometheus.NewDesc(
			prometheus.BuildFQName(util.Namespace, "queries", "total"),
			"Total number of queries.",
			[]string{"profile"}, nil,
		),
		totalAllowedQueries: prometheus.NewDesc(
			prometheus.BuildFQName(util.Namespace, "allowed_queries", "total"),
			"Total number of allowed queries.",
			[]string{"profile"}, nil,
		),
		totalBlockedQueries: prometheus.NewDesc(
			prometheus.BuildFQName(util.Namespace, "blocked_queries", "total"),
			"Total number of blocked queries.",
			[]string{"profile"}, nil,
		),
		blockedQueries: prometheus.NewDesc(
			prometheus.BuildFQName(util.Namespace, "blocked", "queries"),
			"Number of blcoked queries per domain.",
			[]string{"profile", "domain", "root", "tracker"}, nil,
		),
	}
}

func (e *exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.totalQueries
	ch <- e.totalAllowedQueries
	ch <- e.totalBlockedQueries
}

func (e *exporter) Collect(ch chan<- prometheus.Metric) {
	c := api.NewClient(util.BaseURL)

	status, err := c.CollectStatus(e.profile, e.apiKey)
	if err != nil {
		util.Log.Error("error collecting status", "error", err)
		return
	}

	domains, err := c.CollectDomains(e.profile, e.apiKey)
	if err != nil {
		util.Log.Error("error collecting domains", "error", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(e.totalQueries, prometheus.GaugeValue, status.TotalQueries, e.profile)
	ch <- prometheus.MustNewConstMetric(e.totalAllowedQueries, prometheus.GaugeValue, status.AllowedQueries, e.profile)
	ch <- prometheus.MustNewConstMetric(e.totalBlockedQueries, prometheus.GaugeValue, status.BlockedQueries, e.profile)

	fmt.Println(len(domains.BlockedDomains))

	for _, domain := range domains.BlockedDomains {
		ch <- prometheus.MustNewConstMetric(
			e.blockedQueries,
			prometheus.GaugeValue,
			domain.Queries, e.profile, domain.Domain, domain.Root, domain.Tracker,
		)
	}
}

func main() {
	exporter := newExporter(util.Profile, util.APIKey)
	prometheus.MustRegister(exporter)

	util.Log.Info("starting exporter", "port", util.Port, "path", util.MetricsPath)
	http.Handle(util.MetricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, util.MetricsPath, http.StatusMovedPermanently)
	})
	if err := http.ListenAndServe(util.Port, nil); err != nil {
		util.Log.Error("error starting exporter", "error", err)
		os.Exit(1)
	}
}
