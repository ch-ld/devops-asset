package cmdb

import (
	"strconv"
	"time"

	"api-server/internal/response/response"
	svc "api-server/internal/service/cmdb"

	"github.com/gin-gonic/gin"
)

// HostMetricsHandler 主机监控指标处理器
type HostMetricsHandler struct {
	metricsService *svc.MetricsService
	hostService    *svc.HostService
}

// NewHostMetricsHandler 创建主机监控指标处理器
func NewHostMetricsHandler(metricsService *svc.MetricsService, hostService *svc.HostService) *HostMetricsHandler {
	return &HostMetricsHandler{
		metricsService: metricsService,
		hostService:    hostService,
	}
}

// GetHostMetrics 获取主机最新指标
// @Summary 获取主机最新指标
// @Description 获取指定主机的最新CPU、内存、磁盘等监控指标
// @Tags CMDB-主机监控
// @Accept json
// @Produce json
// @Param id path int true "主机ID"
// @Success 200 {object} response.Response{data=svc.HostMetrics}
// @Router /api/v1/cmdb/hosts/{id}/metrics [get]
func (h *HostMetricsHandler) GetHostMetrics(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的主机ID")
		return
	}

	// 获取主机信息，验证权限
	host, err := h.hostService.GetHost(uint(id))
	if err != nil {
		response.ReturnError(c, response.NOT_FOUND, "主机不存在")
		return
	}

	// 获取最新指标
	metrics, err := h.metricsService.GetLatestMetrics(uint(id))
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取主机指标失败: "+err.Error())
		return
	}

	response.ReturnData(c, metrics)
}

// GetHostMetricsHistory 获取主机历史指标
// @Summary 获取主机历史指标
// @Description 获取指定主机的历史监控指标数据
// @Tags CMDB-主机监控
// @Accept json
// @Produce json
// @Param id path int true "主机ID"
// @Param period query string false "时间段(last_hour, last_day, last_week, last_month)" default(last_day)
// @Param start_time query string false "开始时间(ISO8601格式)"
// @Param end_time query string false "结束时间(ISO8601格式)"
// @Param metric_type query string false "指标类型(cpu, memory, disk, network, all)" default(all)
// @Success 200 {object} response.Response{data=[]svc.HostMetrics}
// @Router /api/v1/cmdb/hosts/{id}/metrics/history [get]
func (h *HostMetricsHandler) GetHostMetricsHistory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的主机ID")
		return
	}

	// 获取主机信息，验证权限
	host, err := h.hostService.GetHost(uint(id))
	if err != nil {
		response.ReturnError(c, response.NOT_FOUND, "主机不存在")
		return
	}

	// 解析查询参数
	period := c.DefaultQuery("period", "last_day")
	metricType := c.DefaultQuery("metric_type", "all")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	var startTime, endTime time.Time
	endTime = time.Now()

	// 如果没有明确指定时间范围，则根据period参数确定
	if startTimeStr == "" || endTimeStr == "" {
		switch period {
		case "last_hour":
			startTime = endTime.Add(-1 * time.Hour)
		case "last_day":
			startTime = endTime.Add(-24 * time.Hour)
		case "last_week":
			startTime = endTime.Add(-7 * 24 * time.Hour)
		case "last_month":
			startTime = endTime.Add(-30 * 24 * time.Hour)
		default:
			startTime = endTime.Add(-24 * time.Hour)
		}
	} else {
		// 解析用户指定的时间范围
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			response.ReturnError(c, response.INVALID_ARGUMENT, "无效的开始时间格式")
			return
		}

		if endTimeStr != "" {
			endTime, err = time.Parse(time.RFC3339, endTimeStr)
			if err != nil {
				response.ReturnError(c, response.INVALID_ARGUMENT, "无效的结束时间格式")
				return
			}
		}
	}

	// 获取历史指标
	history, err := h.metricsService.GetHostMetricsHistory(uint(id), startTime, endTime, metricType)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取主机历史指标失败: "+err.Error())
		return
	}

	response.ReturnData(c, history)
}

// GetHostsOverallMetrics 获取所有主机的汇总指标
// @Summary 获取所有主机的汇总指标
// @Description 获取所有主机的CPU、内存、磁盘等资源使用汇总情况
// @Tags CMDB-主机监控
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Router /api/v1/cmdb/hosts/metrics/overall [get]
func (h *HostMetricsHandler) GetHostsOverallMetrics(c *gin.Context) {
	// 获取所有主机的汇总指标
	overall, err := h.metricsService.GetOverallMetrics()
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取主机汇总指标失败: "+err.Error())
		return
	}

	response.ReturnData(c, overall)
}

// RegisterRoutes 注册路由
func (h *HostMetricsHandler) RegisterRoutes(r *gin.RouterGroup) {
	metricsGroup := r.Group("/hosts")
	{
		metricsGroup.GET("/:id/metrics", h.GetHostMetrics)
		metricsGroup.GET("/:id/metrics/history", h.GetHostMetricsHistory)
		metricsGroup.GET("/metrics/overall", h.GetHostsOverallMetrics)
	}
} 