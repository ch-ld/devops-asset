package cmdb

import (
	repo "api-server/internal/repository/cmdb"
	"api-server/internal/response/response"
	svc "api-server/internal/service/cmdb"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// DashboardHandler 仪表盘处理器
type DashboardHandler struct {
	hostRepo       *repo.HostRepository
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
		response.ReturnError(c, response.INTERNAL, "获取主机总数失败")
		return
	}

	// 获取运行中主机数量
	runningCount, err := h.hostRepo.GetRunningHostCount()
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取运行中主机数量失败")
		return
	}

	// 获取停止状态主机数量
	stoppedCount, err := h.hostRepo.CountByConditions(map[string]interface{}{"status": "stopped"})
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取停止状态主机数量失败")
		return
	}

	// 获取异常状态主机数量
	errorCount, err := h.hostRepo.CountByConditions(map[string]interface{}{"status": "error"})
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取异常状态主机数量失败")
		return
	}

	// 计算运行中主机百分比
	var runningPercent float64
	if totalCount > 0 {
		runningPercent = float64(runningCount) / float64(totalCount) * 100
	}

	// 构建响应数据
	data := gin.H{
		"total":           totalCount,
		"running":         runningCount,
		"stopped":         stoppedCount,
		"error":           errorCount,
		"running_percent": runningPercent,
	}

	response.ReturnData(c, data)
}

// HandleGetHostGroups 获取主机组统计信息
func (h *DashboardHandler) HandleGetHostGroups(c *gin.Context) {
	// 获取所有主机组
	groups, err := h.hostRepo.GetHostGroupsWithStats()
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取主机组统计失败")
		return
	}

	response.ReturnData(c, groups)
}

// HandleGetRecentHosts 获取最近添加的主机列表
func (h *DashboardHandler) HandleGetRecentHosts(c *gin.Context) {
	limitParam := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	// 获取最近添加的主机
	hosts, err := h.hostRepo.GetRecentHosts(limit)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取最近主机列表失败")
		return
	}

	response.ReturnData(c, hosts)
}

// HandleGetDistribution 获取主机分布统计
func (h *DashboardHandler) HandleGetDistribution(c *gin.Context) {
	// 按提供商类型统计
	providerCounts, err := h.hostRepo.CountByProviderType()
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取云提供商统计失败")
		return
	}

	// 按区域统计
	regionCounts, err := h.hostRepo.CountByRegion()
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取区域统计失败")
		return
	}

	// 按状态统计
	statusCounts, err := h.hostRepo.CountByStatus()
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取状态统计失败")
		return
	}

	// 按操作系统统计
	osCounts, err := h.hostRepo.CountByOS()
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取操作系统统计失败")
		return
	}

	// 构建响应数据
	data := gin.H{
		"provider": formatForChart(providerCounts),
		"region":   formatForChart(regionCounts),
		"status":   formatForChart(statusCounts),
		"os":       formatForChart(osCounts),
	}

	response.ReturnData(c, data)
}

// HandleGetResourceUsage 获取资源使用情况
func (h *DashboardHandler) HandleGetResourceUsage(c *gin.Context) {
	// 获取主机资源使用情况
	usage, err := h.metricsService.GetHostsResourceUsage()
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取资源使用情况失败")
		return
	}

	response.ReturnData(c, usage)
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
		response.ReturnError(c, response.INTERNAL, "获取历史趋势数据失败")
		return
	}

	response.ReturnData(c, history)
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
		response.ReturnError(c, response.INTERNAL, "获取告警信息失败")
		return
	}

	// 获取异常状态的主机告警
	abnormalAlerts, err := h.alertService.CheckAbnormalHosts()
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取异常状态告警失败")
		return
	}

	// 合并告警
	allAlerts := append(alerts, abnormalAlerts...)

	response.ReturnData(c, allAlerts)
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
		dashboardRouter.GET("/host_groups", h.HandleGetHostGroups)
		dashboardRouter.GET("/recent_hosts", h.HandleGetRecentHosts)
	}
}
