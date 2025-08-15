package dns

import (
    "fmt"

    "api-server/internal/response/response"
    svc "api-server/internal/service/dns"

    "github.com/gin-gonic/gin"
)

// RecordStatHandler 记录统计相关处理器
type RecordStatHandler struct {
    recordService *svc.RecordService
}

func NewRecordStatHandler(recordService *svc.RecordService) *RecordStatHandler {
    return &RecordStatHandler{recordService: recordService}
}

// CountByDomainIDs 批量按域名ID统计解析记录数量
// @Summary 批量统计域名解析记录数
// @Tags DNS记录
// @Accept json
// @Produce json
// @Param body body object true "{\"domain_ids\":[1,2,3]}"
// @Success 200 {object} response.Response{data=map[string]int64}
// @Router /api/v1/dns/records/count-by-domain [post]
func (h *RecordStatHandler) CountByDomainIDs(c *gin.Context) {
    var req struct{
        DomainIDs []uint `json:"domain_ids"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        response.ReturnError(c, response.INVALID_ARGUMENT, "请求参数错误: "+err.Error())
        return
    }
    if len(req.DomainIDs) == 0 {
        response.ReturnData(c, map[string]int64{})
        return
    }

    tenantIDVal, ok := c.Get("tenant_id")
    if !ok {
        response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
        return
    }
    tenantID := tenantIDVal.(uint)

    counts, err := h.recordService.RecordCountsByDomainIDs(req.DomainIDs, tenantID)
    if err != nil {
        response.ReturnError(c, response.INTERNAL, "统计失败: "+err.Error())
        return
    }

    // 转换 key 为字符串，便于前端以对象读取
    out := make(map[string]int64, len(counts))
    for k, v := range counts {
        out[fmt.Sprintf("%d", k)] = v
    }

    response.ReturnData(c, out)
}

