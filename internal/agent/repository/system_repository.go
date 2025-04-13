package repository

import (
	"metrics/internal/agent/dto"
	"strconv"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func NewSystemRepository() *SystemRepository {
	return &SystemRepository{}
}

type SystemRepository struct {
}

func (r *SystemRepository) GetMetrics() []dto.Metric {
	const percent = 100
	var metrics []dto.Metric

	virtualMemory, err := mem.VirtualMemory()
	if err == nil {
		metrics = append(metrics,
			dto.Metric{Name: "TotalMemory", Value: virtualMemory.Total},
			dto.Metric{Name: "FreeMemory", Value: virtualMemory.Free},
		)
	}

	if cpuUsages, err := cpu.Percent(0, true); err == nil {
		for i, usage := range cpuUsages {
			metrics = append(metrics, dto.Metric{
				Name:  "CPUutilization" + strconv.Itoa(i+1),
				Value: uint64(usage * percent),
			})
		}
	}

	return metrics
}
