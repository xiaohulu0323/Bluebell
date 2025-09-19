package controller

import (
	"context"
	"net/http"
	"time"
	"web-app/dao/mysql"

	"github.com/gin-gonic/gin"
)

// GetDBStatsHandler 获取数据库连接池统计信息
func GetDBStatsHandler(c *gin.Context) {
	stats := mysql.GetDBStats()

	if stats == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1000,
			"msg":  "success",
			"data": gin.H{
				"message": "使用原始连接池，无详细统计信息",
				"type":    "basic",
			},
		})
		return
	}

	// 计算连接池使用率
	connectionUsageRate := float64(0)
	if stats.MaxOpenConnections > 0 {
		connectionUsageRate = float64(stats.OpenConnections) / float64(stats.MaxOpenConnections) * 100
	}

	// 计算空闲率
	idleRate := float64(0)
	if stats.OpenConnections > 0 {
		idleRate = float64(stats.Idle) / float64(stats.OpenConnections) * 100
	}

	response := gin.H{
		"code": 1000,
		"msg":  "success",
		"data": gin.H{
			"type": "advanced",
			"connection_pool": gin.H{
				"max_open_connections":  stats.MaxOpenConnections,
				"open_connections":      stats.OpenConnections,
				"in_use":                stats.InUse,
				"idle":                  stats.Idle,
				"connection_usage_rate": connectionUsageRate,
				"idle_rate":             idleRate,
			},
			"wait_stats": gin.H{
				"wait_count":           stats.WaitCount,
				"wait_duration_ms":     float64(stats.WaitDuration.Nanoseconds()) / 1e6,
				"max_idle_closed":      stats.MaxIdleClosed,
				"max_idle_time_closed": stats.MaxIdleTimeClosed,
				"max_lifetime_closed":  stats.MaxLifetimeClosed,
			},
			"query_stats": gin.H{
				"write_query_count": stats.WriteQueryCount,
				"read_query_count":  stats.ReadQueryCount,
				"total_query_count": stats.WriteQueryCount + stats.ReadQueryCount,
				"read_write_ratio":  stats.ReadWriteRatio,
				"avg_query_time_ms": stats.AvgQueryTime,
				"slow_query_count":  stats.SlowQueryCount,
				"error_count":       stats.ErrorCount,
			},
			"last_update_time": stats.LastUpdateTime.Format("2006-01-02 15:04:05"),
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetDBHealthHandler 数据库健康检查
func GetDBHealthHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := mysql.DBHealthCheck(ctx)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "unhealthy",
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
			"error":     err.Error(),
			"checks": gin.H{
				"database_ping": false,
			},
		})
		return
	}

	// 获取连接池统计信息进行健康评估
	stats := mysql.GetDBStats()
	healthStatus := "healthy"
	warnings := make([]string, 0)

	if stats != nil {
		// 检查连接池使用率
		if stats.MaxOpenConnections > 0 {
			usageRate := float64(stats.OpenConnections) / float64(stats.MaxOpenConnections)
			if usageRate > 0.9 {
				healthStatus = "warning"
				warnings = append(warnings, "连接池使用率过高")
			}
		}

		// 检查等待情况
		if stats.WaitCount > 100 {
			healthStatus = "warning"
			warnings = append(warnings, "连接池等待次数过多")
		}

		// 检查错误率
		totalQueries := stats.WriteQueryCount + stats.ReadQueryCount
		if totalQueries > 0 {
			errorRate := float64(stats.ErrorCount) / float64(totalQueries)
			if errorRate > 0.05 { // 错误率超过5%
				healthStatus = "warning"
				warnings = append(warnings, "数据库错误率过高")
			}
		}

		// 检查慢查询
		if totalQueries > 0 {
			slowQueryRate := float64(stats.SlowQueryCount) / float64(totalQueries)
			if slowQueryRate > 0.1 { // 慢查询率超过10%
				healthStatus = "warning"
				warnings = append(warnings, "慢查询比例过高")
			}
		}
	}

	statusCode := http.StatusOK
	if healthStatus == "warning" {
		statusCode = http.StatusOK // 警告状态仍返回200，但在响应体中标明
	}

	response := gin.H{
		"status":    healthStatus,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"checks": gin.H{
			"database_ping": true,
		},
	}

	if len(warnings) > 0 {
		response["warnings"] = warnings
	}

	if stats != nil {
		response["metrics"] = gin.H{
			"connection_usage_rate": func() float64 {
				if stats.MaxOpenConnections > 0 {
					return float64(stats.OpenConnections) / float64(stats.MaxOpenConnections) * 100
				}
				return 0
			}(),
			"wait_count":  stats.WaitCount,
			"error_count": stats.ErrorCount,
			"slow_query_rate": func() float64 {
				total := stats.WriteQueryCount + stats.ReadQueryCount
				if total > 0 {
					return float64(stats.SlowQueryCount) / float64(total) * 100
				}
				return 0
			}(),
		}
	}

	c.JSON(statusCode, response)
}

// OptimizeDBPoolHandler 动态优化数据库连接池（仅开发环境）
func OptimizeDBPoolHandler(c *gin.Context) {
	// 这是一个演示接口，实际生产环境中连接池参数应该通过配置文件管理
	stats := mysql.GetDBStats()

	if stats == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4000,
			"msg":  "当前使用原始连接池，无法动态优化",
		})
		return
	}

	// 基于当前统计信息提供优化建议
	suggestions := make([]string, 0)

	// 连接池使用率分析
	if stats.MaxOpenConnections > 0 {
		usageRate := float64(stats.OpenConnections) / float64(stats.MaxOpenConnections)
		if usageRate > 0.8 {
			suggestions = append(suggestions, "建议增加max_open_conns，当前使用率较高")
		} else if usageRate < 0.3 {
			suggestions = append(suggestions, "可以适当减少max_open_conns，当前使用率较低")
		}
	}

	// 空闲连接分析
	if stats.OpenConnections > 0 {
		idleRate := float64(stats.Idle) / float64(stats.OpenConnections)
		if idleRate > 0.7 {
			suggestions = append(suggestions, "空闲连接过多，建议减少max_idle_conns")
		} else if idleRate < 0.2 && stats.WaitCount > 0 {
			suggestions = append(suggestions, "空闲连接不足且出现等待，建议增加max_idle_conns")
		}
	}

	// 等待情况分析
	if stats.WaitCount > 50 {
		suggestions = append(suggestions, "连接等待较多，建议增加连接池大小")
	}

	// 慢查询分析
	totalQueries := stats.WriteQueryCount + stats.ReadQueryCount
	if totalQueries > 0 {
		slowQueryRate := float64(stats.SlowQueryCount) / float64(totalQueries)
		if slowQueryRate > 0.05 {
			suggestions = append(suggestions, "慢查询较多，建议检查SQL语句和索引")
		}
	}

	if len(suggestions) == 0 {
		suggestions = append(suggestions, "当前连接池配置良好，无需调整")
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1000,
		"msg":  "success",
		"data": gin.H{
			"current_stats": stats,
			"suggestions":   suggestions,
			"optimization_tips": []string{
				"max_open_conns: 建议设置为CPU核心数的2-4倍",
				"max_idle_conns: 建议设置为max_open_conns的1/4到1/2",
				"conn_max_lifetime: 建议设置为30-60分钟",
				"conn_max_idle_time: 建议设置为15-30分钟",
			},
		},
	})
}
