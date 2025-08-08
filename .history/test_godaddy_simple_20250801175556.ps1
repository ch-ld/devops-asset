# 简单的GoDaddy DNS提供商测试脚本
$BaseURL = "http://localhost:8080/api/v1/dns/providers"
$Headers = @{
    "Content-Type" = "application/json"
}

Write-Host "=== GoDaddy DNS提供商测试 ===" -ForegroundColor Green

# 1. 测试临时连接（无效凭证）
Write-Host "`n1. 测试GoDaddy连接（无效凭证）..." -ForegroundColor Yellow

$testBody = @{
    type = "godaddy"
    credentials = @{
        api_key = "invalid-key"
        api_secret = "invalid-secret"
    }
}

$testJson = $testBody | ConvertTo-Json -Depth 3

try {
    $response = Invoke-RestMethod -Uri "$BaseURL/test-connection" -Method POST -Headers $Headers -Body $testJson
    Write-Host "响应: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Cyan
} catch {
    Write-Host "请求失败（这是预期的，因为凭证无效）: $($_.Exception.Message)" -ForegroundColor Yellow
}

# 2. 创建GoDaddy提供商
Write-Host "`n2. 创建GoDaddy提供商..." -ForegroundColor Yellow

$createBody = @{
    name = "测试GoDaddy"
    type = "godaddy"
    credentials = @{
        api_key = "test-key-123"
        api_secret = "test-secret-456"
    }
    remark = "测试提供商"
}

$createJson = $createBody | ConvertTo-Json -Depth 3

try {
    $response = Invoke-RestMethod -Uri $BaseURL -Method POST -Headers $Headers -Body $createJson
    Write-Host "✓ 创建成功" -ForegroundColor Green
    Write-Host "响应: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ 创建失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 3. 获取提供商列表
Write-Host "`n3. 获取提供商列表..." -ForegroundColor Yellow

try {
    $response = Invoke-RestMethod -Uri $BaseURL -Method GET -Headers $Headers
    Write-Host "✓ 获取成功" -ForegroundColor Green
    Write-Host "响应: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ 获取失败: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== 测试完成 ===" -ForegroundColor Green 
