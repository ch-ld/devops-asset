package dns

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"api-server/internal/model/dns"
	"api-server/internal/response/response"
	svc "api-server/internal/service/dns"

	"github.com/gin-gonic/gin"
)

// RecordImportExportHandler 记录导入导出处理器
type RecordImportExportHandler struct {
	recordService *svc.RecordService
}

func NewRecordImportExportHandler(recordService *svc.RecordService) *RecordImportExportHandler {
	return &RecordImportExportHandler{recordService: recordService}
}

// ImportRecordsFromJSON 从前端JSON导入记录
// @Router /api/v1/dns/import-export/records/import [post]
func (h *RecordImportExportHandler) ImportRecordsFromJSON(c *gin.Context) {
	var req struct {
		Format  string                   `json:"format"`
		Data    []map[string]interface{} `json:"data" binding:"required"`
		Options map[string]interface{}   `json:"options"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请求参数错误: "+err.Error())
		return
	}

	userIDVal, ok := c.Get("user_id")
	if !ok {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}
	tenantIDVal, ok := c.Get("tenant_id")
	if !ok {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
		return
	}

	success := 0
	failed := 0
	failedItems := make([]string, 0)

	for i, item := range req.Data {
		// 必要字段解析
		domainID := toUint(item["domain_id"]) // 由前端传入
		name := toString(item["name"])        // 主机记录
		recType := toString(item["type"])     // 记录类型
		value := toString(item["value"])      // 记录值
		if domainID == 0 || name == "" || recType == "" || value == "" {
			failed++
			failedItems = append(failedItems, fmt.Sprintf("第%d行缺少必填字段", i+1))
			continue
		}
		rec := &dns.Record{
			DomainID:  domainID,
			Name:      name,
			Type:      recType,
			Value:     value,
			TTL:       toInt(item["ttl"], 600),
			Remark:    toString(item["remark"]),
			TenantID:  tenantIDVal.(uint),
			CreatedBy: userIDVal.(uint),
			UpdatedBy: userIDVal.(uint),
			Status:    "active",
		}
		if v := toPtrInt(item["priority"]); v != nil {
			rec.Priority = v
		}
		if v := toPtrInt(item["weight"]); v != nil {
			rec.Weight = v
		}
		if v := toPtrInt(item["port"]); v != nil {
			rec.Port = v
		}

		if err := h.recordService.CreateRecord(rec, userIDVal.(uint), c.ClientIP()); err != nil {
			failed++
			failedItems = append(failedItems, fmt.Sprintf("第%d行: %s", i+1, err.Error()))
			continue
		}
		success++
	}

	response.ReturnData(c, gin.H{
		"success":      success,
		"failed":       failed,
		"total":        success + failed,
		"failed_items": failedItems,
	})
}

// ExportRecords 导出记录为CSV（支持filters）
// @Router /api/v1/dns/import-export/records/export [post]
func (h *RecordImportExportHandler) ExportRecords(c *gin.Context) {
	var req struct {
		Format  string                 `json:"format"`
		Filters map[string]interface{} `json:"filters"`
		Fields  []string               `json:"fields"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请求参数错误: "+err.Error())
		return
	}

	tenantIDVal, ok := c.Get("tenant_id")
	if !ok {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
		return
	}

	filters := req.Filters
	if filters == nil {
		filters = map[string]interface{}{}
	}
	filters["tenant_id"] = tenantIDVal.(uint)

	// 拉全部（limit=0, offset=0）
	records, _, err := h.recordService.ListRecords(tenantIDVal.(uint), filters, 0, 0)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "导出失败: "+err.Error())
		return
	}

	fields := req.Fields
	if len(fields) == 0 {
		fields = []string{"name", "type", "value", "ttl", "priority", "weight", "port", "remark"}
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	_ = w.Write(fields)
	for _, r := range records {
		row := make([]string, 0, len(fields))
		for _, f := range fields {
			switch f {
			case "name":
				row = append(row, r.Name)
			case "type":
				row = append(row, r.Type)
			case "value":
				row = append(row, r.Value)
			case "ttl":
				row = append(row, strconv.Itoa(r.TTL))
			case "priority":
				if r.Priority != nil {
					row = append(row, strconv.Itoa(*r.Priority))
				} else {
					row = append(row, "")
				}
			case "weight":
				if r.Weight != nil {
					row = append(row, strconv.Itoa(*r.Weight))
				} else {
					row = append(row, "")
				}
			case "port":
				if r.Port != nil {
					row = append(row, strconv.Itoa(*r.Port))
				} else {
					row = append(row, "")
				}
			case "remark":
				row = append(row, r.Remark)
			default:
				row = append(row, "")
			}
		}
		_ = w.Write(row)
	}
	w.Flush()

	filename := "dns_records_export.csv"
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "text/csv", buf.Bytes())
}

// DownloadRecordCSVTemplate 下载记录CSV模板
// @Router /api/v1/dns/records/template/csv [get]
func (h *RecordImportExportHandler) DownloadRecordCSVTemplate(c *gin.Context) {
	csvContent := "name,type,value,ttl,priority,weight,port,remark\n@,A,1.1.1.1,600,,,,根域名示例"
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=\"dns_records_template.csv\"")
	c.Data(http.StatusOK, "text/csv", []byte(csvContent))
}

// DownloadRecordExcelTemplate 下载Excel模板（临时同CSV）
// @Router /api/v1/dns/records/template/excel [get]
func (h *RecordImportExportHandler) DownloadRecordExcelTemplate(c *gin.Context) {
	csvContent := "name,type,value,ttl,priority,weight,port,remark\n"
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=\"dns_records_template.csv\"")
	c.Data(http.StatusOK, "text/csv", []byte(csvContent))
}

// helpers
func toString(v interface{}) string {
	s, _ := v.(string)
	return s
}

func toInt(v interface{}, def int) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	case string:
		if n, err := strconv.Atoi(t); err == nil {
			return n
		}
	}
	return def
}

func toPtrInt(v interface{}) *int {
	switch t := v.(type) {
	case float64:
		iv := int(t)
		return &iv
	case int:
		iv := t
		return &iv
	case string:
		if n, err := strconv.Atoi(t); err == nil {
			return &n
		}
	}
	return nil
}

func toUint(v interface{}) uint {
	switch t := v.(type) {
	case float64:
		return uint(t)
	case int:
		return uint(t)
	case string:
		if n, err := strconv.Atoi(t); err == nil {
			return uint(n)
		}
	}
	return 0
}
