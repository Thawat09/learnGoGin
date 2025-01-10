package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func Health(c *gin.Context) {
	startTime := time.Now()

	cpuUsage, err := cpu.Percent(0, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get CPU usage", "details": err.Error()})
		return
	}

	memUsage, err := mem.VirtualMemory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get memory usage", "details": err.Error()})
		return
	}

	diskUsage, err := disk.Usage("/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get disk usage", "details": err.Error()})
		return
	}

	responseTime := time.Since(startTime).Microseconds()

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"resource_utilization": gin.H{
			"cpu_usage":    fmt.Sprintf("%.2f%%", cpuUsage[0]),
			"memory_usage": fmt.Sprintf("%.2f%%", memUsage.UsedPercent),
			"disk_usage":   fmt.Sprintf("%.2f%%", diskUsage.UsedPercent),
		},
		"response_time": fmt.Sprintf("%dµs", responseTime),
		"last_checked":  time.Now().Format(time.RFC3339),
	})
}
