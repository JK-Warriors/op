// Scrape `information_schema.innodb_metrics`.

package collector

import (
	"context"
	"database/sql"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const infoSchemaInnodbMetricsQuery = `
		SELECT
		  name, subsystem, type, comment,
		  count
		  FROM information_schema.innodb_metrics
		  WHERE status = 'enabled'
		`

// Metrics descriptors.
var (
	infoSchemaBufferPageReadTotalDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, informationSchema, "innodb_metrics_buffer_page_read_total"),
		"Total number of buffer pages read total.",
		[]string{"type"}, nil,
	)
	infoSchemaBufferPageWrittenTotalDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, informationSchema, "innodb_metrics_buffer_page_written_total"),
		"Total number of buffer pages written total.",
		[]string{"type"}, nil,
	)
	infoSchemaBufferPoolPagesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, informationSchema, "innodb_metrics_buffer_pool_pages"),
		"Total number of buffer pool pages by state.",
		[]string{"state"}, nil,
	)
	infoSchemaBufferPoolPagesDirtyDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, informationSchema, "innodb_metrics_buffer_pool_dirty_pages"),
		"Total number of dirty pages in the buffer pool.",
		nil, nil,
	)
)

// Regexp for matching metric aggregations.
var (
	bufferRE     = regexp.MustCompile(`^buffer_(pool_pages)_(.*)$`)
	bufferPageRE = regexp.MustCompile(`^buffer_page_(read|written)_(.*)$`)
)

// ScrapeInnodbMetrics collects from `information_schema.innodb_metrics`.
type ScrapeInnodbMetrics struct{}

// Name of the Scraper.
func (ScrapeInnodbMetrics) Name() string {
	return informationSchema + ".innodb_metrics"
}

// Help returns additional information about Scraper.
func (ScrapeInnodbMetrics) Help() string {
	return "Collect metrics from information_schema.innodb_metrics"
}

// Version of MySQL from which scraper is available.
func (ScrapeInnodbMetrics) Version() float64 {
	return 5.6
}

// Scrape collects data.
func (ScrapeInnodbMetrics) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric) error {
	innodbMetricsRows, err := db.QueryContext(ctx, infoSchemaInnodbMetricsQuery)
	if err != nil {
		return err
	}
	defer innodbMetricsRows.Close()

	var (
		name, subsystem, metricType, comment string
		value                                float64
	)

	for innodbMetricsRows.Next() {
		if err := innodbMetricsRows.Scan(
			&name, &subsystem, &metricType, &comment, &value,
		); err != nil {
			return err
		}
		// Special handling of the "buffer_page_io" subsystem.
		if subsystem == "buffer_page_io" {
			match := bufferPageRE.FindStringSubmatch(name)
			if len(match) != 3 {
				log.Warnln("innodb_metrics subsystem buffer_page_io returned an invalid name:", name)
				continue
			}
			switch match[1] {
			case "read":
				ch <- prometheus.MustNewConstMetric(
					infoSchemaBufferPageReadTotalDesc, prometheus.CounterValue, value, match[2],
				)
			case "written":
				ch <- prometheus.MustNewConstMetric(
					infoSchemaBufferPageWrittenTotalDesc, prometheus.CounterValue, value, match[2],
				)
			}
			continue
		}
		if subsystem == "buffer" {
			match := bufferRE.FindStringSubmatch(name)
			// Many buffer subsystem metrics are not matched, fall through to generic metric.
			if match != nil {
				switch match[1] {
				case "pool_pages":
					switch match[2] {
					case "total":
						// Ignore total, it is an aggregation of the rest.
						continue
					case "dirty":
						// Dirty pages are a separate metric, not in the total.
						ch <- prometheus.MustNewConstMetric(
							infoSchemaBufferPoolPagesDirtyDesc, prometheus.GaugeValue, value,
						)
					default:
						ch <- prometheus.MustNewConstMetric(
							infoSchemaBufferPoolPagesDesc, prometheus.GaugeValue, value, match[2],
						)
					}
				}
				continue
			}
		}
		metricName := "innodb_metrics_" + subsystem + "_" + name
		// MySQL returns counters named two different ways. "counter" and "status_counter"
		// value >= 0 is necessary due to upstream bugs: http://bugs.mysql.com/bug.php?id=75966
		if (metricType == "counter" || metricType == "status_counter") && value >= 0 {
			description := prometheus.NewDesc(
				prometheus.BuildFQName(namespace, informationSchema, metricName+"_total"),
				comment, nil, nil,
			)
			ch <- prometheus.MustNewConstMetric(
				description,
				prometheus.CounterValue,
				value,
			)
		} else {
			description := prometheus.NewDesc(
				prometheus.BuildFQName(namespace, informationSchema, metricName),
				comment, nil, nil,
			)
			ch <- prometheus.MustNewConstMetric(
				description,
				prometheus.GaugeValue,
				value,
			)
		}
	}
	return nil
}
