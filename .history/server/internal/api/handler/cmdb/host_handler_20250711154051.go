package cmdb

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"path/filepath"
	"strings"
	"time"

	"api-server/internal/model/cmdb"
	model "api-server/internal/model/cmdb"
	"api-server/internal/response/response"
	svc "api-server/internal/service/cmdb"

	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// 主机管理相关接口处理器
// 负责主机的增删改查等HTTP请求处理
type createHostRequest struct {
	// createHostRequest 主机创建/更新请求体结构体
	ProviderID    *uint          `json:"provider_id"`
	InstanceID    string         `json:"instance_id" binding:"required"`
	Name          string         `json:"name" binding:"required"`
	ResourceType  string         `json:"resource_type"`
	Region        string         `json:"region"`
	Username      string         `json:"username"`
	Password      string         `json:"password"`
	PublicIP      datatypes.JSON `json:"public_ip"`
	PrivateIP     datatypes.JSON `json:"private_ip"`
	Configuration datatypes.JSON `json:"configuration"`
	OS            string         `json:"os"`
	Status        string         `json:"status"`
	ExpiredAt     *time.Time     `json:"expired_at"`
	Remark        string         `json:"remark"`
}

// HostHandler 主机管理接口处理器
type HostHandler struct {
	// HostHandler 创建主机处理器实例
	hostSvc *svc.HostService
}

// NewHostHandler 创建主机处理器实例
func NewHostHandler(hostSvc *svc.HostService) *HostHandler {
	return &HostHandler{hostSvc: hostSvc}
}

// CreateHost 创建主机接口
func (h *HostHandler) CreateHost(c *gin.Context) {
	var req createHostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, err.Error())
		return
	}

	host := model.Host{
		ProviderID:    req.ProviderID,
		InstanceID:    req.InstanceID,
		Name:          req.Name,
		ResourceType:  req.ResourceType,
		Region:        req.Region,
		Username:      req.Username,
		Password:      req.Password,
		PublicIP:      req.PublicIP,
		PrivateIP:     req.PrivateIP,
		Configuration: req.Configuration,
		OS:            req.OS,
		Status:        req.Status,
		ExpiredAt:     req.ExpiredAt,
		Remark:        req.Remark,
	}

	if err := h.hostSvc.CreateHost(&host); err != nil {
		response.ReturnError(c, response.INTERNAL, err.Error())
		return
	}
	response.ReturnData(c, host)
}

// UpdateHost 更新主机接口
func (h *HostHandler) UpdateHost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "Invalid host ID")
		return
	}

	var req createHostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, err.Error())
		return
	}

	host := model.Host{
		Model:         gorm.Model{ID: uint(id)},
		ProviderID:    req.ProviderID,
		InstanceID:    req.InstanceID,
		Name:          req.Name,
		ResourceType:  req.ResourceType,
		Region:        req.Region,
		Username:      req.Username,
		Password:      req.Password,
		PublicIP:      req.PublicIP,
		PrivateIP:     req.PrivateIP,
		Configuration: req.Configuration,
		OS:            req.OS,
		Status:        req.Status,
		ExpiredAt:     req.ExpiredAt,
		Remark:        req.Remark,
	}

	if err := h.hostSvc.UpdateHost(&host); err != nil {
		response.ReturnError(c, response.INTERNAL, err.Error())
		return
	}
	response.ReturnData(c, host)
}

// GetHost 查询单个主机接口
func (h *HostHandler) GetHost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "Invalid host ID")
		return
	}

	host, err := h.hostSvc.GetHost(uint(id))
	if err != nil {
		response.ReturnError(c, response.NOT_FOUND, err.Error())
		return
	}
	response.ReturnData(c, host)
}

// ListHostsRequest 查询主机列表请求
type ListHostsRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
	Keyword  string `form:"keyword"`
	Status   string `form:"status"`
	GroupID  *uint  `form:"group_id"`
	Region   string `form:"region"`
}

// ListHosts 查询主机列表接口
// @Summary 查询主机列表
// @Description 分页查询主机列表，支持关键字搜索和状态过滤
// @Tags CMDB-主机管理
// @Accept json
// @Produce json
// @Param page query int true "页码"
// @Param page_size query int true "每页数量"
// @Param keyword query string false "搜索关键字"
// @Param status query string false "主机状态"
// @Param group_id query int false "主机组ID"
// @Param region query string false "区域"
// @Success 200 {object} response.Response{data=[]model.Host}
// @Router /api/v1/cmdb/hosts [get]
func (h *HostHandler) ListHosts(c *gin.Context) {
	var req ListHostsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, err.Error())
		return
	}

	// 构建查询参数
	params := map[string]interface{}{
		"page":      req.Page,
		"page_size": req.PageSize,
	}
	if req.Keyword != "" {
		params["keyword"] = req.Keyword
	}
	if req.Status != "" {
		params["status"] = req.Status
	}
	if req.GroupID != nil {
		params["group_id"] = *req.GroupID
	}
	if req.Region != "" {
		params["region"] = req.Region
	}

	// 查询总数
	total, err := h.hostSvc.CountHosts(params)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, err.Error())
		return
	}

	// 查询列表
	hosts, err := h.hostSvc.ListHosts(params)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, err.Error())
		return
	}

	response.ReturnDataWithCount(c, int(total), hosts)
}

// DeleteHost 删除主机接口
func (h *HostHandler) DeleteHost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "Invalid host ID")
		return
	}

	if err := h.hostSvc.DeleteHost(uint(id)); err != nil {
		response.ReturnError(c, response.INTERNAL, err.Error())
		return
	}
	response.ReturnSuccess(c)
}

// BatchImportResult 批量导入结果
type BatchImportResult struct {
	Total     int      `json:"total"`      // 总记录数
	Success   int      `json:"success"`    // 成功数
	Failed    int      `json:"failed"`     // 失败数
	FailedMsg []string `json:"failed_msg"` // 失败原因
}

// BatchImportHosts 主机批量导入接口
// @Summary 主机批量导入
// @Description 通过Excel/CSV文件批量导入主机资产
// @Tags CMDB-主机管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "主机Excel/CSV文件"
// @Success 200 {object} response.Response{data=BatchImportResult} "导入结果"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/cmdb/hosts/batch_import [post]
func (h *HostHandler) BatchImportHosts(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请选择要上传的文件")
		return
	}

	f, err := file.Open()
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "文件打开失败")
		return
	}
	defer f.Close()

	ext := strings.ToLower(filepath.Ext(file.Filename))
	var hosts []model.Host
	var result BatchImportResult

	switch ext {
	case ".csv":
		hosts, result = h.parseCSVFile(f)
	case ".xlsx":
		hosts, result = h.parseExcelFile(f)
	default:
		response.ReturnError(c, response.INVALID_ARGUMENT, "不支持的文件格式，请使用CSV或Excel文件")
		return
	}

	if len(hosts) > 0 {
		// 批量创建主机
		if err := h.hostSvc.BatchCreateHosts(hosts); err != nil {
			result.Failed += len(hosts)
			result.FailedMsg = append(result.FailedMsg, "批量创建失败: "+err.Error())
		} else {
			result.Success += len(hosts)
		}
	}

	response.ReturnData(c, result)
}

// parseCSVFile 解析CSV文件
func (h *HostHandler) parseCSVFile(f io.Reader) ([]model.Host, BatchImportResult) {
	var hosts []model.Host
	var result BatchImportResult

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		result.FailedMsg = append(result.FailedMsg, "CSV解析失败: "+err.Error())
		return hosts, result
	}

	if len(records) < 2 { // 至少需要标题行和一行数据
		result.FailedMsg = append(result.FailedMsg, "文件内容为空")
		return hosts, result
	}

	// 验证标题行
	headers := records[0]
	if !h.validateHeaders(headers) {
		result.FailedMsg = append(result.FailedMsg, "文件格式错误，请使用正确的模板")
		return hosts, result
	}

	result.Total = len(records) - 1 // 减去标题行
	// 解析数据行
	for i, record := range records[1:] {
		host, err := h.parseHostRecord(record, headers)
		if err != nil {
			result.Failed++
			result.FailedMsg = append(result.FailedMsg, fmt.Sprintf("第%d行解析失败: %s", i+2, err.Error()))
			continue
		}
		hosts = append(hosts, *host)
	}

	return hosts, result
}

// parseExcelFile 解析Excel文件
func (h *HostHandler) parseExcelFile(f io.Reader) ([]model.Host, BatchImportResult) {
	var hosts []model.Host
	var result BatchImportResult

	xlsx, err := excelize.OpenReader(f)
	if err != nil {
		result.FailedMsg = append(result.FailedMsg, "Excel解析失败: "+err.Error())
		return hosts, result
	}

	// 获取第一个工作表
	sheet := xlsx.GetSheetList()[0]
	rows, err := xlsx.GetRows(sheet)
	if err != nil {
		result.FailedMsg = append(result.FailedMsg, "Excel读取失败: "+err.Error())
		return hosts, result
	}

	if len(rows) < 2 { // 至少需要标题行和一行数据
		result.FailedMsg = append(result.FailedMsg, "文件内容为空")
		return hosts, result
	}

	// 验证标题行
	headers := rows[0]
	if !h.validateHeaders(headers) {
		result.FailedMsg = append(result.FailedMsg, "文件格式错误，请使用正确的模板")
		return hosts, result
	}

	result.Total = len(rows) - 1 // 减去标题行
	// 解析数据行
	for i, row := range rows[1:] {
		host, err := h.parseHostRecord(row, headers)
		if err != nil {
			result.Failed++
			result.FailedMsg = append(result.FailedMsg, fmt.Sprintf("第%d行解析失败: %s", i+2, err.Error()))
			continue
		}
		hosts = append(hosts, *host)
	}

	return hosts, result
}

// validateHeaders 验证文件标题行
func (h *HostHandler) validateHeaders(headers []string) bool {
	requiredHeaders := []string{
		"provider_id",
		"instance_id",
		"name",
		"resource_type",
		"region",
		"username",
		"password",
		"public_ip",
		"private_ip",
		"configuration",
		"os",
		"status",
		"expired_at",
		"remark",
	}

	headerMap := make(map[string]bool)
	for _, h := range headers {
		headerMap[strings.TrimSpace(strings.ToLower(h))] = true
	}

	for _, required := range requiredHeaders {
		if !headerMap[required] {
			return false
		}
	}
	return true
}

// parseHostRecord 解析主机记录
func (h *HostHandler) parseHostRecord(record []string, headers []string) (*model.Host, error) {
	if len(record) != len(headers) {
		return nil, fmt.Errorf("列数不匹配")
	}

	// 创建字段映射
	fields := make(map[string]string)
	for i, header := range headers {
		fields[strings.TrimSpace(strings.ToLower(header))] = record[i]
	}

	// 解析必填字段
	providerID, err := strconv.ParseUint(fields["provider_id"], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("provider_id格式错误")
	}

	if fields["instance_id"] == "" {
		return nil, fmt.Errorf("instance_id不能为空")
	}

	if fields["name"] == "" {
		return nil, fmt.Errorf("name不能为空")
	}

	// 解析JSON字段
	var publicIP, privateIP, configuration datatypes.JSON
	if fields["public_ip"] != "" {
		if err := json.Unmarshal([]byte(fields["public_ip"]), &publicIP); err != nil {
			return nil, fmt.Errorf("public_ip格式错误")
		}
	}
	if fields["private_ip"] != "" {
		if err := json.Unmarshal([]byte(fields["private_ip"]), &privateIP); err != nil {
			return nil, fmt.Errorf("private_ip格式错误")
		}
	}
	if fields["configuration"] != "" {
		if err := json.Unmarshal([]byte(fields["configuration"]), &configuration); err != nil {
			return nil, fmt.Errorf("configuration格式错误")
		}
	}

	// 解析时间字段
	var expiredAt *time.Time
	if fields["expired_at"] != "" {
		t, err := time.Parse("2006-01-02 15:04:05", fields["expired_at"])
		if err != nil {
			return nil, fmt.Errorf("expired_at格式错误，请使用YYYY-MM-DD HH:mm:ss格式")
		}
		expiredAt = &t
	}

	var providerIDPtr *uint
	if providerID != 0 {
		providerIDUint := uint(providerID)
		providerIDPtr = &providerIDUint
	}

	return &model.Host{
		ProviderID:    providerIDPtr,
		InstanceID:    fields["instance_id"],
		Name:          fields["name"],
		ResourceType:  fields["resource_type"],
		Region:        fields["region"],
		Username:      fields["username"],
		Password:      fields["password"],
		PublicIP:      publicIP,
		PrivateIP:     privateIP,
		Configuration: configuration,
		OS:            fields["os"],
		Status:        fields["status"],
		ExpiredAt:     expiredAt,
		Remark:        fields["remark"],
	}, nil
}

// BatchExportHosts 主机批量导出接口
// @Summary 主机批量导出
// @Description 导出主机清单为Excel/CSV文件，支持多条件筛选
// @Tags 主机管理
// @Accept json
// @Produce application/octet-stream
// @Param format query string false "导出格式（excel/csv），默认excel"
// @Param name query string false "主机名称筛选"
// @Param group query string false "分组筛选"
// @Param tag query string false "标签筛选"
// @Success 200 {file} file "导出文件"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /cmdb/hosts/batch_export [get]
func (h *HostHandler) BatchExportHosts(c *gin.Context) {
	format := c.DefaultQuery("format", "excel")
	name := c.Query("name")
	group := c.Query("group")
	tag := c.Query("tag")
	// TODO: 支持更多筛选条件

	params := map[string]interface{}{}
	if name != "" {
		params["name"] = name
	}
	if group != "" {
		params["group"] = group
	}
	if tag != "" {
		params["tag"] = tag
	}

	hosts, err := h.hostSvc.ListHosts(params)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, err.Error())
		return
	}

	if format == "csv" {
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment; filename=hosts.csv")
		w := csv.NewWriter(c.Writer)
		w.Write([]string{"主机名称", "实例ID", "区域"}) // TODO: 其他字段
		for _, host := range hosts {
			w.Write([]string{host.Name, host.InstanceID, host.Region})
		}
		w.Flush()
		return
	}

	// 默认导出Excel
	xl := excelize.NewFile()
	sheet := xl.GetSheetName(0)
	xl.SetSheetRow(sheet, "A1", &[]string{"主机名称", "实例ID", "区域"}) // TODO: 其他字段
	for i, host := range hosts {
		row := []string{host.Name, host.InstanceID, host.Region}
		xl.SetSheetRow(sheet, fmt.Sprintf("A%d", i+2), &row)
	}
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=hosts.xlsx")
	xl.Write(c.Writer)
}

// BatchDeleteHosts 主机批量删除接口
// @Summary 主机批量删除
// @Description 批量删除主机
// @Tags 主机管理
// @Accept json
// @Produce json
// @Param ids body []uint true "主机ID列表"
// @Success 200 {object} response.ResponseData "删除结果"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /cmdb/hosts/batch_delete [delete]
func (h *HostHandler) BatchDeleteHosts(c *gin.Context) {
	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil || len(ids) == 0 {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请传入主机ID列表")
		return
	}
	if err := h.hostSvc.BatchDeleteHosts(ids); err != nil {
		response.ReturnError(c, response.INTERNAL, err.Error())
		return
	}
	response.ReturnSuccess(c)
}

// BatchUpdateHosts 主机批量编辑接口
// @Summary 主机批量编辑
// @Description 批量修改主机属性
// @Tags 主机管理
// @Accept json
// @Produce json
// @Param hosts body []model.Host true "主机列表（需包含ID）"
// @Success 200 {object} response.ResponseData "编辑结果"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /cmdb/hosts/batch_update [put]
func (h *HostHandler) BatchUpdateHosts(c *gin.Context) {
	var hosts []model.Host
	if err := c.ShouldBindJSON(&hosts); err != nil || len(hosts) == 0 {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请传入主机列表")
		return
	}
	if err := h.hostSvc.BatchUpdateHosts(hosts); err != nil {
		response.ReturnError(c, response.INTERNAL, err.Error())
		return
	}
	response.ReturnSuccess(c)
}

// BatchAssignHosts 主机批量分配接口
// @Summary 主机批量分配
// @Description 批量分配分组、标签、负责人等
// @Tags 主机管理
// @Accept json
// @Produce json
// @Param req body svc.BatchAssignRequest true "批量分配请求体"
// @Success 200 {object} response.ResponseData "分配结果"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /cmdb/hosts/batch_assign [put]
func (h *HostHandler) BatchAssignHosts(c *gin.Context) {
	var req svc.BatchAssignRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请传入主机ID列表")
		return
	}
	if err := h.hostSvc.BatchAssignHosts(req); err != nil {
		response.ReturnError(c, response.INTERNAL, err.Error())
		return
	}
	response.ReturnSuccess(c)
}

// BatchLifecycleHosts 主机生命周期管理接口
// @Summary 主机生命周期管理
// @Description 批量到期、下线、回收等操作
// @Tags 主机管理
// @Accept json
// @Produce json
// @Param req body svc.BatchLifecycleRequest true "生命周期操作请求体"
// @Success 200 {object} response.ResponseData "操作结果"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /cmdb/hosts/lifecycle [put]
func (h *HostHandler) BatchLifecycleHosts(c *gin.Context) {
	var req svc.BatchLifecycleRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请传入主机ID列表")
		return
	}
	if err := h.hostSvc.BatchLifecycleHosts(req); err != nil {
		response.ReturnError(c, response.INTERNAL, err.Error())
		return
	}
	response.ReturnSuccess(c)
}

// BatchSetCustomFields 主机批量自定义字段赋值接口
// @Summary 主机批量自定义字段赋值
// @Description 批量设置主机的自定义扩展字段
// @Tags 主机管理
// @Accept json
// @Produce json
// @Param req body svc.BatchSetCustomFieldsRequest true "批量自定义字段赋值请求体"
// @Success 200 {object} response.ResponseData "操作结果"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /cmdb/hosts/custom_fields [put]
func (h *HostHandler) BatchSetCustomFields(c *gin.Context) {
	var req svc.BatchSetCustomFieldsRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 || len(req.ExtraFields) == 0 {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请传入主机ID列表和自定义字段内容")
		return
	}
	if err := h.hostSvc.BatchSetCustomFields(req); err != nil {
		response.ReturnError(c, response.INTERNAL, err.Error())
		return
	}
	response.ReturnSuccess(c)
}

// BatchChangeStatus 主机批量状态变更接口
// @Summary 主机批量状态变更
// @Description 批量变更主机状态（如上线、下线、禁用、启用等）
// @Tags 主机管理
// @Accept json
// @Produce json
// @Param req body svc.BatchChangeStatusRequest true "批量状态变更请求体"
// @Success 200 {object} response.ResponseData "操作结果"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /cmdb/hosts/batch_status [put]
func (h *HostHandler) BatchChangeStatus(c *gin.Context) {
	var req svc.BatchChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 || req.Status == "" {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请传入主机ID列表和目标状态")
		return
	}
	if err := h.hostSvc.BatchChangeStatus(req); err != nil {
		response.ReturnError(c, response.INTERNAL, err.Error())
		return
	}
	response.ReturnSuccess(c)
}

// HostHistory 主机变更历史/审计日志接口
// @Summary 主机变更历史/审计日志
// @Description 查询指定主机的操作日志
// @Tags 主机管理
// @Accept json
// @Produce json
// @Param id query string true "主机ID或名称"
// @Success 200 {object} response.ResponseData "操作日志列表"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /cmdb/hosts/history [get]
func (h *HostHandler) HostHistory(c *gin.Context) {
	idOrName := c.Query("id")
	if idOrName == "" {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请传入主机ID或名称")
		return
	}
	// TODO: 如需操作日志功能，请补充依赖，否则此处返回空列表
	response.ReturnData(c, []interface{}{})
}

// BatchSSH 主机批量WebSSH/命令执行接口
// @Summary 主机批量WebSSH/命令执行
// @Description 对多台主机批量下发命令并收集执行结果
// @Tags 主机管理
// @Accept json
// @Produce json
// @Param req body svc.BatchSSHRequest true "批量命令执行请求体"
// @Success 200 {object} response.ResponseData "执行结果列表"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /cmdb/hosts/batch_ssh [post]
func (h *HostHandler) BatchSSH(c *gin.Context) {
	var req svc.BatchSSHRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 || req.Cmd == "" {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请传入主机ID列表和命令")
		return
	}
	results := h.hostSvc.BatchSSH(req)
	response.ReturnData(c, results)
}

// BatchSFTP 主机批量SFTP/文件分发接口
// @Summary 主机批量SFTP/文件分发
// @Description 对多台主机批量上传（分发）文件
// @Tags 主机管理
// @Accept multipart/form-data
// @Produce json
// @Param ids formData []uint true "主机ID列表"
// @Param remote_path formData string true "目标路径"
// @Param file formData file true "上传文件"
// @Success 200 {object} response.ResponseData "分发结果列表"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /cmdb/hosts/batch_sftp [post]
type SFTPResult struct {
	HostID  uint   `json:"host_id"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func (h *HostHandler) BatchSFTP(c *gin.Context) {
	idsStr := c.PostFormArray("ids")
	remotePath := c.PostForm("remote_path")
	file, err := c.FormFile("file")
	if len(idsStr) == 0 || remotePath == "" || err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请传入主机ID列表、目标路径和文件")
		return
	}
	var ids []uint
	for _, s := range idsStr {
		id, _ := strconv.Atoi(s)
		ids = append(ids, uint(id))
	}
	f, err := file.Open()
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "文件打开失败")
		return
	}
	defer f.Close()
	results := h.hostSvc.BatchSFTP(ids, remotePath, f)
	response.ReturnData(c, results)
}

// AlertHosts 主机异常/到期告警查询接口
// @Summary 主机异常/到期告警查询
// @Description 查询异常、到期、即将到期主机
// @Tags 主机管理
// @Accept json
// @Produce json
// @Param days query int false "即将到期天数，默认7天"
// @Success 200 {object} response.ResponseData "告警主机列表"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /cmdb/hosts/alert [get]
func (h *HostHandler) AlertHosts(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	alerts, err := h.hostSvc.AlertHosts(days)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, err.Error())
		return
	}
	response.ReturnData(c, alerts)
}

// SyncHosts 同步云主机接口
// @Summary 同步云主机
// @Description 从云服务商同步主机资源
// @Tags CMDB-主机管理
// @Accept json
// @Produce json
// @Param provider_id query int false "云账号ID，不传则同步所有云账号"
// @Success 200 {object} response.Response
// @Router /api/v1/cmdb/hosts/sync [post]
func (h *HostHandler) SyncHosts(c *gin.Context) {
	// 获取云账号ID
	providerIDStr := c.Query("provider_id")
	if providerIDStr != "" {
		providerID, err := strconv.ParseUint(providerIDStr, 10, 32)
		if err != nil {
			response.ReturnError(c, response.INVALID_ARGUMENT, "云账号ID格式错误")
			return
		}
		// 同步指定云账号的主机
		if err := h.hostSvc.SyncHostsFromCloud(uint(providerID)); err != nil {
			response.ReturnError(c, response.INTERNAL, fmt.Sprintf("同步主机失败: %s", err.Error()))
			return
		}
	} else {
		// 同步所有云账号的主机
		if err := h.hostSvc.SyncAllProviderHosts(); err != nil {
			response.ReturnError(c, response.INTERNAL, fmt.Sprintf("同步主机失败: %s", err.Error()))
			return
		}
	}

	response.ReturnSuccess(c)
}

// SyncHostStatus 同步主机状态接口
// @Summary 同步主机状态
// @Description 从云服务商同步指定主机的状态
// @Tags CMDB-主机管理
// @Accept json
// @Produce json
// @Param id path int true "主机ID"
// @Success 200 {object} response.Response
// @Router /api/v1/cmdb/hosts/{id}/sync_status [post]
func (h *HostHandler) SyncHostStatus(c *gin.Context) {
	// 获取主机ID
	hostID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "主机ID格式错误")
		return
	}

	// 同步主机状态
	if err := h.hostSvc.SyncHostStatus(uint(hostID)); err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("同步主机状态失败: %s", err.Error()))
		return
	}

	response.ReturnSuccess(c)
}

// CreateManualHost 创建自建主机接口
// @Summary 创建自建主机
// @Description 创建一个自建主机
// @Tags CMDB-主机管理
// @Accept json
// @Produce json
// @Param host body cmdb.Host true "主机信息"
// @Success 200 {object} response.Response
// @Router /api/v1/cmdb/hosts/manual [post]
func (h *HostHandler) CreateManualHost(c *gin.Context) {
	var host cmdb.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请求参数错误")
		return
	}

	if err := h.hostSvc.CreateManualHost(&host); err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("创建自建主机失败: %s", err.Error()))
		return
	}

	response.ReturnSuccess(c)
}

// ListGroupHosts 获取主机组下的主机列表接口
// @Summary 获取主机组下的主机列表
// @Description 获取指定主机组下的主机列表，支持分页和搜索
// @Tags CMDB-主机管理
// @Accept json
// @Produce json
// @Param group_id path int true "主机组ID"
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页数量，默认20"
// @Param keyword query string false "搜索关键字"
// @Success 200 {object} response.Response{data=response.PageData{list=[]cmdb.Host}}
// @Router /api/v1/cmdb/groups/{group_id}/hosts [get]
func (h *HostHandler) ListGroupHosts(c *gin.Context) {
	// 获取主机组ID
	groupID, err := strconv.ParseUint(c.Param("group_id"), 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "主机组ID格式错误")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	// 查询主机列表
	hosts, total, err := h.hostSvc.GetGroupHosts(uint(groupID), page, pageSize, keyword)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("获取主机列表失败: %s", err.Error()))
		return
	}

	response.ReturnDataWithCount(c, int(total), hosts)
}

// MoveHost 移动主机到指定组接口
// @Summary 移动主机到指定组
// @Description 将主机移动到指定的主机组
// @Tags CMDB-主机管理
// @Accept json
// @Produce json
// @Param id path int true "主机ID"
// @Param group_id query int false "主机组ID，不传则移出主机组"
// @Success 200 {object} response.Response
// @Router /api/v1/cmdb/hosts/{id}/move [post]
func (h *HostHandler) MoveHost(c *gin.Context) {
	// 获取主机ID
	hostID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "主机ID格式错误")
		return
	}

	// 获取主机组ID
	var groupID *uint
	if groupIDStr := c.Query("group_id"); groupIDStr != "" {
		id, err := strconv.ParseUint(groupIDStr, 10, 32)
		if err != nil {
			response.ReturnError(c, response.INVALID_ARGUMENT, "主机组ID格式错误")
			return
		}
		uid := uint(id)
		groupID = &uid
	}

	// 移动主机
	if err := h.hostSvc.MoveHost(uint(hostID), groupID); err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("移动主机失败: %s", err.Error()))
		return
	}

	response.ReturnSuccess(c)
}
