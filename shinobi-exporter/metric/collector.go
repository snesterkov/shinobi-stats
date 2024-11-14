package metric

import (
	"context"

	"github.com/pavelgopanenko/shinobi-exporter/shinobi"
	"github.com/prometheus/client_golang/prometheus"
)

type Server interface {
	Monitors(ctx context.Context, group shinobi.Group) ([]shinobi.Monitor, error)
}

type ServerCollector struct {
	ctx    context.Context
	server Server
	groups []shinobi.Group

	monitorsTotal       *prometheus.Desc
	monitorsStatusInfo  *prometheus.Desc
	monitorsErrorsTotal *prometheus.Desc
}

func NewServerCollector(ctx context.Context, server Server, groups ...shinobi.Group) *ServerCollector {
	const ns = "shinobi_"

	return &ServerCollector{
		ctx:    ctx,
		server: server,
		groups: groups,

		monitorsTotal: prometheus.NewDesc(
			ns+"monitors_total",
			"The total count of all known monitors.",
			[]string{"group"}, nil,
		),
		monitorsStatusInfo: prometheus.NewDesc(
			ns+"monitors_status_info",
			"Monitor status information. Codes: 0-undefined, 1-starting, 2-watching, 5-stopped, 7-died",
			[]string{"group", "mid", "name", "mode", "streams"}, nil,
//			[]string{"group", "mid", "name", "status", "mode"}, nil,
		),
		monitorsErrorsTotal: prometheus.NewDesc(
			ns+"monitors_error_total",
			"The total errors count.",
			nil, nil,
		),
	}
}

func (c *ServerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.monitorsTotal
}

func (c *ServerCollector) Collect(ch chan<- prometheus.Metric) {
	var errCount int

	for _, group := range c.groups {
		// @todo handle errors
		monitors, err := c.server.Monitors(c.ctx, group)
		if err != nil {
			errCount++

			continue
		}

		ch <- prometheus.MustNewConstMetric(
			c.monitorsTotal,
			prometheus.GaugeValue,
			float64(len(monitors)),
			string(group),
		)

		for _, monitor := range monitors {
			ch <- prometheus.MustNewConstMetric(
				c.monitorsStatusInfo,
				prometheus.GaugeValue,
				float64(monitor.Code),
				string(group),
				string(monitor.MID),
				monitor.Name,
//				string(monitor.Status),
				string(monitor.Mode),
                                string(monitor.Streams[0]),
			)
		}
	}

	ch <- prometheus.MustNewConstMetric(
		c.monitorsErrorsTotal,
		prometheus.GaugeValue,
		float64(errCount),
	)
}
