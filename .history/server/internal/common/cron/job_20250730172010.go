// 定时任务Job定义
// 负责主机资源同步、告警等具体任务的实现
package cron

import (
	"api-server/internal/service/cmdb"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RegionSyncJob 区域同步任务
type RegionSyncJob struct {
	regionService *cmdb.RegionService
}

// NewRegionSyncJob 创建区域同步任务
func NewRegionSyncJob(db *gorm.DB) *RegionSyncJob {
	return &RegionSyncJob{
		regionService: cmdb.NewRegionService(db),
	}
}

// Run 执行区域同步任务
func (j *RegionSyncJob) Run() {
	zap.L().Info("开始执行区域同步定时任务")

	err := j.regionService.SyncAllProviderRegions()
	if err != nil {
		zap.L().Error("区域同步定时任务执行失败", zap.Error(err))
		return
	}

	zap.L().Info("区域同步定时任务执行完成")
}
