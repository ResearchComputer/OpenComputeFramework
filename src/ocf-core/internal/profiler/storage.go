package profiler

import (
	"fmt"
	"ocfcore/internal/common"
	"ocfcore/internal/common/structs"
	"time"

	"github.com/nakabonne/tstorage"
)

type TStorage struct {
	Client tstorage.Storage
}

var StorageClient TStorage

func NewStorageClient() TStorage {
	if StorageClient.Client == nil {
		StorageClient.Client, _ = tstorage.NewStorage(
			tstorage.WithTimestampPrecision(tstorage.Seconds),
			tstorage.WithDataPath("./data"),
		)
	}
	return StorageClient
}

func AddPoint(host string, metric string, timestamp int64, value float64) {
	labels := []tstorage.Label{
		{Name: "host", Value: host},
	}
	err := NewStorageClient().Client.InsertRows([]tstorage.Row{
		{
			Metric:    metric,
			Labels:    labels,
			DataPoint: tstorage.DataPoint{Timestamp: timestamp, Value: value},
		},
	})
	if err != nil {
		common.Logger.Error(err)
	}
}

func QueryPoints(start int64, end int64, metric string, host string) []*tstorage.DataPoint {
	labels := []tstorage.Label{
		{Name: "host", Value: host},
	}
	rows, _ := NewStorageClient().Client.Select(
		metric,
		labels,
		start,
		end,
	)
	return rows
}

func AggregateAverageUtilization(host string, duration time.Duration) float64 {
	// current timestamp
	end := time.Now().Unix()
	// start timestamp = current timestamp - duration
	start := end - int64(duration.Seconds())
	labels := []tstorage.Label{
		{Name: "host", Value: host},
	}
	rows, _ := NewStorageClient().Client.Select(
		"GPU Utilization",
		labels,
		start,
		end,
	)
	var sum float64
	for _, row := range rows {
		sum += row.Value
	}
	return sum / float64(len(rows))
}

func QueryCardSummary(host string) (structs.CardMetrics, error) {
	var metrics structs.CardMetrics
	end := time.Now().Unix()
	start := end - int64(30*time.Second.Seconds())
	labels := []tstorage.Label{
		{Name: "host", Value: host},
	}
	var err error
	keywords := []string{"GPU Utilization", "Power Usage", "Used Memory", "Available Memory"}
	for _, keyword := range keywords {
		rows, _ := NewStorageClient().Client.Select(
			keyword,
			labels,
			start,
			end,
		)
		if len(rows) > 0 {
			if keyword == "GPU Utilization" {
				metrics.GPUUtilization = rows[len(rows)-1].Value
			}
			if keyword == "Power Usage" {
				metrics.PowerUsage = rows[len(rows)-1].Value
			}
			if keyword == "Used Memory" {
				metrics.UsedMemory = rows[len(rows)-1].Value
			}
			if keyword == "Available Memory" {
				metrics.AvailableMemory = rows[len(rows)-1].Value
			}
		} else {
			err = fmt.Errorf("no data for %s", keyword)
		}
	}
	return metrics, err
}
