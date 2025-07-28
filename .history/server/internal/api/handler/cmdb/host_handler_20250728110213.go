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
	GroupID       *uint          `json:"group_id"` // 添加主机组ID字段
}

// HostHandler 主机管理接口处理器
type HostHandler struct {
	// HostHandler 创建主机处理器实例
	hostSvc  *svc.HostService
	groupSvc *svc.HostGroupService
}

// NewHostHandler 创建主机处理器实例
func NewHostHandler(hostSvc *svc.HostService, groupSvc *svc.HostGroupService) *HostHandler {
	return &HostHandler{
		hostSvc:  hostSvc,
		groupSvc: groupSvc,
	}
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
		GroupID:       req.GroupID, // 添加主机组ID映射
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
		GroupID:       req.GroupID, // 添加主机组ID映射
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
	Name     string `form:"name"`
	IP       string `form:"ip"`
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
	if req.Name != "" {
		params["name"] = req.Name
	}
	if req.IP != "" {
		params["ip"] = req.IP
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
		// 处理主机组并批量创建主机
		if err := h.processHostGroupsAndCreateHosts(hosts, &result); err != nil {
			result.Failed += len(hosts)
			result.FailedMsg = append(result.FailedMsg, "批量创建失败: "+err.Error())
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
	// 中文字段名到英文字段名的映射
	chineseToEnglishMap := map[string]string{
		"云厂商id": "provider_id",
		"实例id":  "instance_id",
		"主机名称":  "name",
		"主机组":   "group_name",
		"资源类型":  "resource_type",
		"地域":    "region",
		"用户名":   "username",
		"密码":    "password",
		"公网ip":  "public_ip",
		"私网ip":  "private_ip",
		"配置信息":  "configuration",
		"操作系统":  "os",
		"状态":    "status",
		"过期时间":  "expired_at",
		"备注":    "remark",
	}

	requiredHeaders := []string{
		"provider_id",
		"instance_id",
		"name",
		"group_name",
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
		normalizedHeader := strings.TrimSpace(strings.ToLower(h))
		// 如果是中文字段名，转换为英文字段名
		if englishName, exists := chineseToEnglishMap[normalizedHeader]; exists {
			headerMap[englishName] = true
		} else {
			headerMap[normalizedHeader] = true
		}
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

	// 中文字段名到英文字段名的映射
	chineseToEnglishMap := map[string]string{
		"云厂商id": "provider_id",
		"实例id":  "instance_id",
		"主机名称":  "name",
		"主机组":   "group_name",
		"资源类型":  "resource_type",
		"地域":    "region",
		"用户名":   "username",
		"密码":    "password",
		"公网ip":  "public_ip",
		"私网ip":  "private_ip",
		"配置信息":  "configuration",
		"操作系统":  "os",
		"状态":    "status",
		"过期时间":  "expired_at",
		"备注":    "remark",
	}

	// 创建字段映射
	fields := make(map[string]string)
	for i, header := range headers {
		normalizedHeader := strings.TrimSpace(strings.ToLower(header))
		// 如果是中文字段名，转换为英文字段名
		if englishName, exists := chineseToEnglishMap[normalizedHeader]; exists {
			fields[englishName] = record[i]
		} else {
			fields[normalizedHeader] = record[i]
		}
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

	// 处理主机组名称，临时存储在ExtraFields中
	var extraFields datatypes.JSON
	if fields["group_name"] != "" {
		extraFieldsMap := map[string]interface{}{
			"_temp_group_name": fields["group_name"],
		}
		extraFieldsBytes, _ := json.Marshal(extraFieldsMap)
		extraFields = datatypes.JSON(extraFieldsBytes)
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
		ExtraFields:   extraFields,
		Remark:        fields["remark"],
	}, nil
}

// BatchExportHosts 主机批量导出接口
// @Summary 主机批量导出
// @Description 导出主机清单为Excel/CSV文件，支持多条件筛选和字段选择
// @Tags 主机管理
// @Accept json
// @Produce application/octet-stream
// @Param format query string false "导出格式（excel/csv），默认excel"
// @Param scope query string false "导出范围（all/group/current），默认all"
// @Param group_id query int false "主机组ID（当scope=group时必填）"
// @Param fields query string false "导出字段，逗号分隔"
// @Param name query string false "主机名称筛选（当scope=current时使用）"
// @Param status query string false "状态筛选（当scope=current时使用）"
// @Param region query string false "区域筛选（当scope=current时使用）"
// @Param provider query string false "提供商筛选（当scope=current时使用）"
// @Success 200 {file} file "导出文件"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /cmdb/hosts/batch_export [get]
func (h *HostHandler) BatchExportHosts(c *gin.Context) {
	format := c.DefaultQuery("format", "excel")
	scope := c.DefaultQuery("scope", "all")
	fieldsParam := c.DefaultQuery("fields", "")

	// 解析导出字段
	var fields []string
	if fieldsParam != "" {
		fields = strings.Split(fieldsParam, ",")
	} else {
		// 默认字段
		fields = []string{"name", "instance_id", "status", "public_ip", "private_ip", "os", "region"}
	}

	// 构建查询参数
	params := map[string]interface{}{}

	switch scope {
	case "group":
		groupIDStr := c.Query("group_id")
		if groupIDStr == "" {
			response.ReturnError(c, response.INVALID_ARGUMENT, "导出指定主机组时必须提供group_id参数")
			return
		}
		groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
		if err != nil {
			response.ReturnError(c, response.INVALID_ARGUMENT, "主机组ID格式错误")
			return
		}
		params["group_id"] = uint(groupID)

	case "current":
		// 使用当前筛选条件
		if name := c.Query("name"); name != "" {
			params["name"] = name
		}
		if status := c.Query("status"); status != "" {
			params["status"] = status
		}
		if region := c.Query("region"); region != "" {
			params["region"] = region
		}
		if provider := c.Query("provider"); provider != "" {
			params["provider_type"] = provider
		}

	case "all":
		// 导出所有主机，不添加筛选条件

	default:
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的导出范围，支持：all/group/current")
		return
	}

	// 获取主机列表
	hosts, err := h.hostSvc.ListHosts(params)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, err.Error())
		return
	}

	// 生成文件名
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	var filename string
	switch scope {
	case "group":
		filename = fmt.Sprintf("主机列表_主机组_%s.%s", timestamp, format)
	case "current":
		filename = fmt.Sprintf("主机列表_筛选结果_%s.%s", timestamp, format)
	default:
		filename = fmt.Sprintf("主机列表_全部_%s.%s", timestamp, format)
	}

	if format == "csv" {
		h.exportToCSV(c, hosts, fields, filename)
	} else {
		h.exportToExcel(c, hosts, fields, filename)
	}
}

// exportToCSV 导出为CSV格式
func (h *HostHandler) exportToCSV(c *gin.Context, hosts []cmdb.Host, fields []string, filename string) {
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// 添加BOM以支持中文
	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})

	w := csv.NewWriter(c.Writer)
	defer w.Flush()

	// 写入表头
	headers := h.getFieldHeaders(fields)
	w.Write(headers)

	// 写入数据行
	for _, host := range hosts {
		row := h.getHostRowData(host, fields)
		w.Write(row)
	}
}

// exportToExcel 导出为Excel格式
func (h *HostHandler) exportToExcel(c *gin.Context, hosts []cmdb.Host, fields []string, filename string) {
	xl := excelize.NewFile()
	sheet := xl.GetSheetName(0)

	// 设置表头
	headers := h.getFieldHeaders(fields)
	xl.SetSheetRow(sheet, "A1", &headers)

	// 设置表头样式
	headerStyle, _ := xl.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 12},
		Fill: &excelize.Fill{Type: "pattern", Color: []string{"#E6F3FF"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
	})

	// 应用表头样式
	endCol := string(rune('A' + len(headers) - 1))
	xl.SetCellStyle(sheet, "A1", endCol+"1", headerStyle)

	// 写入数据行
	for i, host := range hosts {
		row := h.getHostRowData(host, fields)
		xl.SetSheetRow(sheet, fmt.Sprintf("A%d", i+2), &row)
	}

	// 自动调整列宽
	for i := 0; i < len(headers); i++ {
		col := string(rune('A' + i))
		xl.SetColWidth(sheet, col, col, 15)
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	xl.Write(c.Writer)
}

// getFieldHeaders 获取字段对应的中文表头
func (h *HostHandler) getFieldHeaders(fields []string) []string {
	fieldMap := map[string]string{
		"name":          "主机名称",
		"instance_id":   "实例ID",
		"status":        "状态",
		"public_ip":     "公网IP",
		"private_ip":    "私网IP",
		"os":            "操作系统",
		"region":        "区域",
		"configuration": "配置规格",
		"username":      "用户名",
		"provider_type": "提供商类型",
		"resource_type": "资源类型",
		"group_name":    "主机组",
		"provider_name": "云账号",
		"tags":          "标签",
		"expired_at":    "过期时间",
		"remark":        "备注",
		"created_at":    "创建时间",
		"updated_at":    "更新时间",
	}

	headers := make([]string, len(fields))
	for i, field := range fields {
		if header, exists := fieldMap[field]; exists {
			headers[i] = header
		} else {
			headers[i] = field
		}
	}
	return headers
}

// getHostRowData 获取主机行数据
func (h *HostHandler) getHostRowData(host cmdb.Host, fields []string) []string {
	row := make([]string, len(fields))

	for i, field := range fields {
		switch field {
		case "name":
			row[i] = host.Name
		case "instance_id":
			row[i] = host.InstanceID
		case "status":
			row[i] = h.getStatusText(host.Status)
		case "public_ip":
			row[i] = h.formatIPList(host.PublicIP)
		case "private_ip":
			row[i] = h.formatIPList(host.PrivateIP)
		case "os":
			row[i] = host.OS
		case "region":
			row[i] = host.Region
		case "configuration":
			row[i] = h.formatConfiguration(host.Configuration)
		case "username":
			row[i] = host.Username
		case "provider_type":
			row[i] = h.getProviderTypeText(host.ProviderType)
		case "resource_type":
			row[i] = h.getResourceTypeText(host.ResourceType)
		case "group_name":
			if host.Group != nil {
				row[i] = host.Group.Name
			} else {
				row[i] = "未分组"
			}
		case "provider_name":
			if host.Provider != nil {
				row[i] = host.Provider.Name
			} else {
				row[i] = "手动添加"
			}
		case "tags":
			row[i] = h.formatTags(host.Tags)
		case "expired_at":
			if host.ExpiredAt != nil {
				row[i] = host.ExpiredAt.Format("2006-01-02 15:04:05")
			} else {
				row[i] = "永不过期"
			}
		case "remark":
			row[i] = host.Remark
		case "created_at":
			row[i] = host.CreatedAt.Format("2006-01-02 15:04:05")
		case "updated_at":
			row[i] = host.UpdatedAt.Format("2006-01-02 15:04:05")
		default:
			row[i] = ""
		}
	}

	return row
}

// 格式化辅助方法

// getStatusText 获取状态的中文描述
func (h *HostHandler) getStatusText(status string) string {
	statusMap := map[string]string{
		"running":  "运行中",
		"stopped":  "已停止",
		"starting": "启动中",
		"stopping": "停止中",
		"pending":  "待定",
		"unknown":  "未知",
	}
	if text, exists := statusMap[status]; exists {
		return text
	}
	return status
}

// getProviderTypeText 获取提供商类型的中文描述
func (h *HostHandler) getProviderTypeText(providerType string) string {
	providerMap := map[string]string{
		"aliyun":  "阿里云",
		"aws":     "亚马逊云",
		"tencent": "腾讯云",
		"huawei":  "华为云",
		"manual":  "手动添加",
	}
	if text, exists := providerMap[providerType]; exists {
		return text
	}
	return providerType
}

// getResourceTypeText 获取资源类型的中文描述
func (h *HostHandler) getResourceTypeText(resourceType string) string {
	resourceMap := map[string]string{
		"ecs":     "云服务器ECS",
		"ec2":     "EC2实例",
		"cvm":     "云服务器CVM",
		"bcc":     "云服务器BCC",
		"unknown": "未知类型",
	}
	if text, exists := resourceMap[resourceType]; exists {
		return text
	}
	return resourceType
}

// formatIPList 格式化IP列表
func (h *HostHandler) formatIPList(ipData datatypes.JSON) string {
	if ipData == nil {
		return ""
	}

	var ips []string
	if err := json.Unmarshal(ipData, &ips); err != nil {
		return string(ipData)
	}

	return strings.Join(ips, ", ")
}

// formatConfiguration 格式化配置信息
func (h *HostHandler) formatConfiguration(configData datatypes.JSON) string {
	if configData == nil {
		return ""
	}

	var config map[string]interface{}
	if err := json.Unmarshal(configData, &config); err != nil {
		return string(configData)
	}

	var parts []string
	if cpu, ok := config["cpu"]; ok {
		parts = append(parts, fmt.Sprintf("CPU: %v核", cpu))
	}
	if memory, ok := config["memory"]; ok {
		parts = append(parts, fmt.Sprintf("内存: %vGB", memory))
	}
	if disk, ok := config["disk"]; ok {
		parts = append(parts, fmt.Sprintf("磁盘: %vGB", disk))
	}
	if instanceType, ok := config["instance_type"]; ok {
		parts = append(parts, fmt.Sprintf("规格: %v", instanceType))
	}

	if len(parts) > 0 {
		return strings.Join(parts, ", ")
	}

	return string(configData)
}

// formatTags 格式化标签
func (h *HostHandler) formatTags(tagsData datatypes.JSON) string {
	if tagsData == nil {
		return ""
	}

	var tags []string
	if err := json.Unmarshal(tagsData, &tags); err != nil {
		// 尝试解析为map格式的标签
		var tagMap map[string]string
		if err2 := json.Unmarshal(tagsData, &tagMap); err2 != nil {
			return string(tagsData)
		}

		for key, value := range tagMap {
			tags = append(tags, fmt.Sprintf("%s:%s", key, value))
		}
	}

	return strings.Join(tags, ", ")
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

// BatchLifecycleHosts 主机批量生命周期管理接口
// @Summary 主机批量生命周期管理
// @Description 批量设置主机过期时间、状态、回收等生命周期操作
// @Tags 主机管理
// @Accept json
// @Produce json
// @Param req body model.BatchLifecycleRequest true "批量生命周期请求体"
// @Success 200 {object} response.ResponseData "操作结果"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /cmdb/hosts/lifecycle [put]
func (h *HostHandler) BatchLifecycleHosts(c *gin.Context) {
	var req model.BatchLifecycleRequest
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
// @Param req body model.BatchSetCustomFieldsRequest true "批量自定义字段赋值请求体"
// @Success 200 {object} response.ResponseData "操作结果"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /cmdb/hosts/custom_fields [put]
func (h *HostHandler) BatchSetCustomFields(c *gin.Context) {
	var req model.BatchSetCustomFieldsRequest
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
// @Param req body model.BatchChangeStatusRequest true "批量状态变更请求体"
// @Success 200 {object} response.ResponseData "操作结果"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /cmdb/hosts/batch_status [put]
func (h *HostHandler) BatchChangeStatus(c *gin.Context) {
	var req model.BatchChangeStatusRequest
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
// @Param req body model.BatchSSHRequest true "批量命令执行请求体"
// @Success 200 {object} response.ResponseData "执行结果列表"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /cmdb/hosts/batch_ssh [post]
func (h *HostHandler) BatchSSH(c *gin.Context) {
	var req model.BatchSSHRequest
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

	// 调试：打印接收到的数据
	fmt.Printf("🔍 后端接收到的主机数据: %+v\n", host)
	fmt.Printf("🔍 主机组ID: %v\n", host.GroupID)

	if err := h.hostSvc.CreateManualHost(&host); err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("创建自建主机失败: %s", err.Error()))
		return
	}

	response.ReturnSuccess(c)
}

// GetHostFilterOptions 获取主机筛选选项接口
// @Summary 获取主机筛选选项
// @Description 获取主机状态、地域、云厂商等筛选选项
// @Tags CMDB-主机管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/cmdb/hosts/filter_options [get]
func (h *HostHandler) GetHostFilterOptions(c *gin.Context) {
	options, err := h.hostSvc.GetFilterOptions()
	if err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("获取筛选选项失败: %s", err.Error()))
		return
	}
	response.ReturnData(c, options)
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

// HostWithGroupName 临时结构，用于存储主机和主机组名称
type HostWithGroupName struct {
	Host      model.Host
	GroupName string
}

// processHostGroupsAndCreateHosts 处理主机组并批量创建主机
func (h *HostHandler) processHostGroupsAndCreateHosts(hosts []model.Host, result *BatchImportResult) error {
	// 收集所有需要的主机组名称
	groupNameToID := make(map[string]uint)

	for i := range hosts {
		// 从主机的ExtraFields中获取主机组名称
		if hosts[i].ExtraFields != nil {
			var extraFields map[string]interface{}
			if err := json.Unmarshal([]byte(hosts[i].ExtraFields), &extraFields); err == nil {
				if groupName, exists := extraFields["_temp_group_name"].(string); exists && groupName != "" {
					// 查找或创建主机组
					groupID, err := h.findOrCreateHostGroup(groupName)
					if err != nil {
						result.Failed++
						result.FailedMsg = append(result.FailedMsg, fmt.Sprintf("主机 %s 的主机组 %s 处理失败: %s", hosts[i].Name, groupName, err.Error()))
						continue
					}

					// 设置主机的GroupID
					hosts[i].GroupID = &groupID
					groupNameToID[groupName] = groupID

					// 清除临时字段
					delete(extraFields, "_temp_group_name")
					if len(extraFields) == 0 {
						hosts[i].ExtraFields = nil
					} else {
						extraFieldsBytes, _ := json.Marshal(extraFields)
						hosts[i].ExtraFields = datatypes.JSON(extraFieldsBytes)
					}
				}
			}
		}
	}

	// 批量创建主机
	if err := h.hostSvc.BatchCreateHosts(hosts); err != nil {
		result.Failed += len(hosts)
		result.FailedMsg = append(result.FailedMsg, "批量创建失败: "+err.Error())
		return err
	} else {
		result.Success += len(hosts)
	}

	return nil
}

// findOrCreateHostGroup 查找或创建主机组
func (h *HostHandler) findOrCreateHostGroup(groupName string) (uint, error) {
	// 首先尝试查找现有的主机组
	groups, err := h.groupSvc.ListGroups("")
	if err != nil {
		return 0, fmt.Errorf("查询主机组失败: %w", err)
	}

	// 查找同名的主机组
	for _, group := range groups {
		if group.Name == groupName {
			return group.ID, nil
		}
	}

	// 如果没有找到，创建新的主机组
	newGroup := &model.HostGroup{
		Name:        groupName,
		Description: fmt.Sprintf("批量导入时自动创建的主机组: %s", groupName),
		ParentID:    nil, // 创建为根级主机组
		Sort:        0,
	}

	if err := h.groupSvc.CreateHostGroup(newGroup); err != nil {
		return 0, fmt.Errorf("创建主机组失败: %w", err)
	}

	return newGroup.ID, nil
}
