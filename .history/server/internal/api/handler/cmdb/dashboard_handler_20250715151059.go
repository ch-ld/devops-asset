package cmdb

import (
	repoModel "api-server/internal/model/cmdb"
	repo "api-server/internal/repository/cmdb"
	"api-server/internal/response/response"
	svc "api-server/internal/service/cmdb"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// DashboardHandler 仪表盘处理器
type DashboardHandler struct {
	hostRepo      *repo.HostRepository
	metricsService *svc.MetricsService
	alertService   *svc.AlertService
}

// NewDashboardHandler 创建仪表盘处理器
func NewDashboardHandler(
	hostRepo *repo.HostRepository,
	metricsService *svc.MetricsService,
	alertService *svc.AlertService,
) *DashboardHandler {
	return &DashboardHandler{
		hostRepo:       hostRepo,
		metricsService: metricsService,
		alertService:   alertService,
	}
}

// HandleGetSummary 获取仪表盘摘要数据
func (h *DashboardHandler) HandleGetSummary(c *gin.Context) {
	// 获取主机总数
	totalCount, err := h.hostRepo.GetTotalHostCount()
	if err != nil {
		response.SystemError(c, "获取主机总数失败")
		return
	}

	// 获取运行中主机数量
	runningCount, err := h.hostRepo.GetRunningHostCount()
	if err != nil {
		response.SystemError(c, "获取运行中主机数量失败")
		return
	}

	// 获取30天内过期的主机数量
	expiringCount, err := h.hostRepo.GetExpiringHostCount(30)
	if err != nil {
		response.SystemError(c, "获取即将过期主机数量失败")
		return
	}

	// 计算运行中主机百分比
	var runningPercent float64
	if totalCount > 0 {
		runningPercent = float64(runningCount) / float64(totalCount) * 100
	}

	// 计算即将过期主机百分比
	var expiringPercent float64
	if totalCount > 0 {
		expiringPercent = float64(expiringCount) / float64(totalCount) * 100
	}

	// 获取告警数量 (简化实现，实际应查询告警表)
	alertCount := 8
	alertTrend := -12.5 // 假设告警趋势下降了12.5%

	// 构建响应数据
	data := gin.H{
		"total_hosts":      totalCount,
		"hosts_trend":      5.2, // 假设主机数量增长了5.2%
		"running_hosts":    runningCount,
		"running_percent":  runningPercent,
		"expiring_hosts":   expiringCount,
		"expiring_percent": expiringPercent,
		"alert_count":      alertCount,
		"alert_trend":      alertTrend,
	}

	response.Success(c, data)
}

// HandleGetDistribution 获取主机分布统计
func (h *DashboardHandler) HandleGetDistribution(c *gin.Context) {
	// 按提供商类型统计
	providerCounts, err := h.hostRepo.CountByProviderType()
	if err != nil {
		response.SystemError(c, "获取云提供商统计失败")
		return
	}

	// 按区域统计
	regionCounts, err := h.hostRepo.CountByRegion()
	if err != nil {
		response.SystemError(c, "获取区域统计失败")
		return
	}

	// 按状态统计
	statusCounts, err := h.hostRepo.CountByStatus()
	if err != nil {
		response.SystemError(c, "获取状态统计失败")
		return
	}

	// 按操作系统统计
	osCounts, err := h.hostRepo.CountByOS()
	if err != nil {
		response.SystemError(c, "获取操作系统统计失败")
		return
	}

	// 构建响应数据
	data := gin.H{
		"provider": formatForChart(providerCounts),
		"region":   formatForChart(regionCounts),
		"status":   formatForChart(statusCounts),
		"os":       formatForChart(osCounts),
	}

	response.Success(c, data)
}

// HandleGetResourceUsage 获取资源使用情况
func (h *DashboardHandler) HandleGetResourceUsage(c *gin.Context) {
	// 获取主机资源使用情况
	usage, err := h.metricsService.GetHostsResourceUsage()
	if err != nil {
		response.SystemError(c, "获取资源使用情况失败")
		return
	}

	response.Success(c, usage)
}

// HandleGetHistoryTrend 获取历史趋势
func (h *DashboardHandler) HandleGetHistoryTrend(c *gin.Context) {
	// 获取时间范围
	daysParam := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysParam)
	if err != nil || days <= 0 || days > 365 {
		days = 30
	}

	// 计算开始和结束时间
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	// 获取历史数据
	history, err := h.hostRepo.GetHistoryByDateRange(startDate, endDate)
	if err != nil {
		response.SystemError(c, "获取历史趋势数据失败")
		return
	}

	response.Success(c, history)
}

// HandleGetAlerts 获取告警信息
func (h *DashboardHandler) HandleGetAlerts(c *gin.Context) {
	// 获取时间范围
	daysParam := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysParam)
	if err != nil || days <= 0 || days > 365 {
		days = 30
	}

	// 获取告警信息
	alerts, err := h.alertService.CheckExpiringHosts(days)
	if err != nil {
		response.SystemError(c, "获取告警信息失败")
		return
	}

	// 获取异常状态的主机告警
	abnormalAlerts, err := h.alertService.CheckAbnormalHosts()
	if err != nil {
		response.SystemError(c, "获取异常状态告警失败")
		return
	}

	// 合并告警
	allAlerts := append(alerts, abnormalAlerts...)

	response.Success(c, allAlerts)
}

// formatForChart 将统计数据格式化为图表需要的格式
func formatForChart(counts map[string]int64) []gin.H {
	var result []gin.H
	for name, count := range counts {
		result = append(result, gin.H{
			"name":  name,
			"value": count,
		})
	}
	return result
}

// RegisterRoutes 注册路由
func (h *DashboardHandler) RegisterRoutes(router *gin.RouterGroup) {
	dashboardRouter := router.Group("/dashboard")
	{
		dashboardRouter.GET("/summary", h.HandleGetSummary)
		dashboardRouter.GET("/distribution", h.HandleGetDistribution)
		dashboardRouter.GET("/resource_usage", h.HandleGetResourceUsage)
		dashboardRouter.GET("/history_trend", h.HandleGetHistoryTrend)
		dashboardRouter.GET("/alerts", h.HandleGetAlerts)
	}
}
