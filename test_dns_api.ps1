# DNS提供商API测试脚本 (PowerShell)
$BaseURL = "http://localhost:8080/api/v1/dns/providers"
$Headers = @{
    "Content-Type" = "application/json"
    "Access-Token" = "test-token"
}

Write-Host "=== DNS提供商API测试 ===" -ForegroundColor Green

# 1. 测试获取提供商列表
Write-Host "`n1. 测试获取提供商列表..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri $BaseURL -Method GET -Headers $Headers
    Write-Host "✓ 获取提供商列表成功" -ForegroundColor Green
    Write-Host "响应: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ 获取提供商列表失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 2. 测试临时连接测试
Write-Host "`n2. 测试临时连接测试..." -ForegroundColor Yellow
$testData = @{
    type = "aliyun"
    credentials = @{
        access_key_id = "test-key-id"
        access_key_secret = "test-key-secret"
        region = "cn-hangzhou"
    }
} | ConvertTo-Json -Depth 3

try {
    $response = Invoke-RestMethod -Uri "$BaseURL/test-connection" -Method POST -Headers $Headers -Body $testData
    Write-Host "✓ 临时连接测试成功" -ForegroundColor Green
    Write-Host "响应: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ 临时连接测试失败: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream())
        $responseBody = $reader.ReadToEnd()
        Write-Host "错误详情: $responseBody" -ForegroundColor Red
    }
}

# 3. 测试创建提供商
Write-Host "`n3. 测试创建提供商..." -ForegroundColor Yellow
$createData = @{
    name = "测试阿里云DNS"
    type = "aliyun"
    credentials = @{
        access_key_id = "test-key-id"
        access_key_secret = "test-key-secret"
        region = "cn-hangzhou"
    }
    remark = "测试用提供商"
} | ConvertTo-Json -Depth 3

try {
    $response = Invoke-RestMethod -Uri $BaseURL -Method POST -Headers $Headers -Body $createData
    Write-Host "✓ 创建提供商成功" -ForegroundColor Green
    Write-Host "响应: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ 创建提供商失败: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream())
        $responseBody = $reader.ReadToEnd()
        Write-Host "错误详情: $responseBody" -ForegroundColor Red
    }
}

Write-Host "`n=== 测试完成 ===" -ForegroundColor Green 
