package repository

import (
	"strconv"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func NewSystemRepository() *SystemRepository {
	return &SystemRepository{}
}

type SystemRepository struct {
}

func (r *SystemRepository) GetMetrics() []Metric {
	var metrics []Metric

	virtualMemory, err := mem.VirtualMemory()
	if err == nil {
		metrics = append(metrics,
			Metric{"TotalMemory", virtualMemory.Total},
			Metric{"FreeMemory", virtualMemory.Free},
		)
	}

	if cpuUsages, err := cpu.Percent(0, true); err == nil {
		for i, usage := range cpuUsages {
			metrics = append(metrics, Metric{
				Name:  "CPUutilization" + strconv.Itoa(i+1),
				Value: uint64(usage * 100),
			})
		}
	}

	return metrics
}
