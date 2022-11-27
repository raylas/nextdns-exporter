package main

import (
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/raylas/nextdns-exporter/internal/api"
	"github.com/raylas/nextdns-exporter/internal/util"
)

type exporter struct {
	profile, apiKey string

	// Totaled metrics.
	totalQueries        *prometheus.Desc
	totalAllowedQueries *prometheus.Desc
	totalBlockedQueries *prometheus.Desc

	// Detailed metrics.
	blockedQueries     *prometheus.Desc
	deviceQueries      *prometheus.Desc
	protocolQueries    *prometheus.Desc
	typeQueries        *prometheus.Desc
	ipVersionQueries   *prometheus.Desc
	dnssecQueries      *prometheus.Desc
	encryptedQueries   *prometheus.Desc
	destinationQueries *prometheus.Desc
}

func newExporter(profile, apiKey string) *exporter {
	return &exporter{
		profile: profile,
		apiKey:  apiKey,

		// Totaled metrics.
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

		// Detailed metrics.
		blockedQueries: prometheus.NewDesc(
			prometheus.BuildFQName(util.Namespace, "blocked", "queries"),
			"Number of blocked queries per domain.",
			[]string{"profile", "domain", "root", "tracker"}, nil,
		),
		deviceQueries: prometheus.NewDesc(
			prometheus.BuildFQName(util.Namespace, "device", "queries"),
			"Number of queries per device.",
			[]string{"profile", "id", "name", "model", "local_ip"}, nil,
		),
		protocolQueries: prometheus.NewDesc(
			prometheus.BuildFQName(util.Namespace, "protocol", "queries"),
			"Number of queries per protocol.",
			[]string{"profile", "protocol"}, nil,
		),
		typeQueries: prometheus.NewDesc(
			prometheus.BuildFQName(util.Namespace, "type", "queries"),
			"Number of queries per type.",
			[]string{"profile", "type", "name"}, nil,
		),
		ipVersionQueries: prometheus.NewDesc(
			prometheus.BuildFQName(util.Namespace, "ip_version", "queries"),
			"Number of queries per IP version.",
			[]string{"profile", "version"}, nil,
		),
		dnssecQueries: prometheus.NewDesc(
			prometheus.BuildFQName(util.Namespace, "dnssec", "queries"),
			"Number of DNSSEC and non-DNSSEC queries.",
			[]string{"profile", "validated"}, nil,
		),
		encryptedQueries: prometheus.NewDesc(
			prometheus.BuildFQName(util.Namespace, "encrypted", "queries"),
			"Number of encrypted and unencrypted queries.",
			[]string{"profile", "encrypted"}, nil,
		),
		destinationQueries: prometheus.NewDesc(
			prometheus.BuildFQName(util.Namespace, "destination", "queries"),
			"Number of queries per geographic destination.",
			[]string{"profile", "code", "name"}, nil,
		),
	}
}

func (e *exporter) Describe(ch chan<- *prometheus.Desc) {
	// Totaled metrics.
	ch <- e.totalQueries
	ch <- e.totalAllowedQueries
	ch <- e.totalBlockedQueries

	// Detailed metrics.
	ch <- e.blockedQueries
	ch <- e.deviceQueries
	ch <- e.protocolQueries
	ch <- e.typeQueries
	ch <- e.ipVersionQueries
	ch <- e.dnssecQueries
	ch <- e.encryptedQueries
	ch <- e.destinationQueries
}

func (e *exporter) Collect(ch chan<- prometheus.Metric) {
	c := api.NewClient(util.BaseURL, e.profile, e.apiKey)

	// Totaled metrics.
	status, err := c.CollectStatus()
	if err != nil {
		util.Log.Error("error collecting status data", "error", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(e.totalQueries, prometheus.GaugeValue, status.TotalQueries, e.profile)
	ch <- prometheus.MustNewConstMetric(e.totalAllowedQueries, prometheus.GaugeValue, status.AllowedQueries, e.profile)
	ch <- prometheus.MustNewConstMetric(e.totalBlockedQueries, prometheus.GaugeValue, status.BlockedQueries, e.profile)

	// Detailed metrics.
	domains, err := c.CollectDomains()
	if err != nil {
		util.Log.Error("error collecting domains data", "error", err)
		return
	}
	for _, domain := range domains.BlockedDomains {
		ch <- prometheus.MustNewConstMetric(
			e.blockedQueries,
			prometheus.GaugeValue,
			domain.Queries, e.profile, domain.Domain, domain.Root, domain.Tracker,
		)
	}

	devices, err := c.CollectDevices()
	if err != nil {
		util.Log.Error("error collecting devices data", "error", err)
		return
	}
	for _, device := range devices.Devices {
		ch <- prometheus.MustNewConstMetric(
			e.deviceQueries,
			prometheus.GaugeValue,
			device.Queries, e.profile, device.ID, device.Name, device.Model, device.LocalIP,
		)
	}

	protocols, err := c.CollectProtocols()
	if err != nil {
		util.Log.Error("error collecting protocols data", "error", err)
		return
	}
	for _, protocol := range protocols.Protocols {
		ch <- prometheus.MustNewConstMetric(
			e.protocolQueries,
			prometheus.GaugeValue,
			protocol.Queries, e.profile, protocol.Protocol,
		)
	}

	queryTypes, err := c.CollectQueryTypes()
	if err != nil {
		util.Log.Error("error collecting query types data", "error", err)
		return
	}
	for _, queryType := range queryTypes.QueryTypes {
		ch <- prometheus.MustNewConstMetric(
			e.typeQueries,
			prometheus.GaugeValue,
			queryType.Queries, e.profile, queryType.Type, queryType.Name,
		)
	}

	ipVersions, err := c.CollectIPVersions()
	if err != nil {
		util.Log.Error("error collecting IP versions data", "error", err)
		return
	}
	for _, ipVersion := range ipVersions.IPVersions {
		ch <- prometheus.MustNewConstMetric(
			e.ipVersionQueries,
			prometheus.GaugeValue,
			ipVersion.Queries, e.profile, ipVersion.Version,
		)
	}

	dnssec, err := c.CollectDNSSEC()
	if err != nil {
		util.Log.Error("error collecting DNSSEC data", "error", err)
		return
	}
	for _, data := range dnssec.Data {
		ch <- prometheus.MustNewConstMetric(
			e.dnssecQueries,
			prometheus.GaugeValue,
			data.Queries, e.profile, data.Validated,
		)
	}

	encryption, err := c.CollectEncryption()
	if err != nil {
		util.Log.Error("error collecting encryption data", "error", err)
		return
	}
	for _, data := range encryption.Data {
		ch <- prometheus.MustNewConstMetric(
			e.encryptedQueries,
			prometheus.GaugeValue,
			data.Queries, e.profile, data.Encrypted,
		)
	}

	destinations, err := c.CollectDestinations()
	if err != nil {
		util.Log.Error("error collecting destinations data", "error", err)
		return
	}
	for _, destination := range destinations.Destinations {
		ch <- prometheus.MustNewConstMetric(
			e.destinationQueries,
			prometheus.GaugeValue,
			destination.Queries, e.profile, destination.Code, destination.Name,
		)
	}
}

func main() {
	exporter := newExporter(util.Profile, util.APIKey)
	prometheus.MustRegister(exporter)

	util.Log.Info("starting exporter", "version", util.Version, "port", util.Port, "path", util.MetricsPath)
	http.Handle(util.MetricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, util.MetricsPath, http.StatusMovedPermanently)
	})
	if err := http.ListenAndServe(util.Port, nil); err != nil {
		util.Log.Error("error starting exporter", "error", err)
		os.Exit(1)
	}
}
