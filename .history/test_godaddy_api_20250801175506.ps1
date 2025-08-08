# GoDaddy DNS提供商API测试脚本 (PowerShell)
$BaseURL = "http://localhost:8080/api/v1/dns/providers"
$Headers = @{
    "Content-Type" = "application/json"
    "Access-Token" = "test-token"
}

Write-Host "=== GoDaddy DNS提供商API测试 ===" -ForegroundColor Green

# 1. 测试获取提供商列表
Write-Host "`n1. 测试获取提供商列表..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri $BaseURL -Method GET -Headers $Headers
    Write-Host "✓ 获取提供商列表成功" -ForegroundColor Green
    Write-Host "响应: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ 获取提供商列表失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 2. 测试GoDaddy临时连接测试（使用无效凭证）
Write-Host "`n2. 测试GoDaddy临时连接测试（无效凭证）..." -ForegroundColor Yellow
$testData = @{
    type = "godaddy"
    credentials = @{
        api_key = "test-invalid-key"
        api_secret = "test-invalid-secret"
    }
} | ConvertTo-Json -Depth 3

try {
    $response = Invoke-RestMethod -Uri "$BaseURL/test-connection" -Method POST -Headers $Headers -Body $testData
    Write-Host "✓ 临时连接测试响应成功" -ForegroundColor Green
    Write-Host "响应: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ 临时连接测试失败: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream())
        $responseBody = $reader.ReadToEnd()
        Write-Host "错误详情: $responseBody" -ForegroundColor Red
    }
}

# 3. 测试创建GoDaddy提供商
Write-Host "`n3. 测试创建GoDaddy提供商..." -ForegroundColor Yellow
$createData = @{
    name = "测试GoDaddy DNS"
    type = "godaddy"
    credentials = @{
        api_key = "test-key-12345"
        api_secret = "test-secret-67890"
    }
    remark = "测试用GoDaddy提供商"
}
$createData = $createData | ConvertTo-Json -Depth 3

try {
    $response = Invoke-RestMethod -Uri $BaseURL -Method POST -Headers $Headers -Body $createData
    Write-Host "✓ 创建GoDaddy提供商成功" -ForegroundColor Green
    Write-Host "响应: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Cyan
    
    # 保存提供商ID用于后续测试
    $providerId = $response.data.id
    Write-Host "提供商ID: $providerId" -ForegroundColor Blue
    
    # 4. 测试已保存提供商的连接测试
    if ($providerId) {
        Write-Host "`n4. 测试已保存提供商的连接测试..." -ForegroundColor Yellow
        try {
            $testResponse = Invoke-RestMethod -Uri "$BaseURL/$providerId/test" -Method POST -Headers $Headers
            Write-Host "✓ 已保存提供商连接测试成功" -ForegroundColor Green
            Write-Host "测试结果: $($testResponse | ConvertTo-Json -Depth 3)" -ForegroundColor Cyan
        } catch {
            Write-Host "✗ 已保存提供商连接测试失败: $($_.Exception.Message)" -ForegroundColor Red
            if ($_.Exception.Response) {
                $reader = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream())
                $responseBody = $reader.ReadToEnd()
                Write-Host "错误详情: $responseBody" -ForegroundColor Red
            }
        }
        
        # 5. 测试域名同步
        Write-Host "`n5. 测试域名同步..." -ForegroundColor Yellow
        try {
            $syncResponse = Invoke-RestMethod -Uri "$BaseURL/$providerId/sync-domains" -Method POST -Headers $Headers
            Write-Host "✓ 域名同步请求成功" -ForegroundColor Green
            Write-Host "同步结果: $($syncResponse | ConvertTo-Json -Depth 3)" -ForegroundColor Cyan
        } catch {
            Write-Host "✗ 域名同步失败: $($_.Exception.Message)" -ForegroundColor Red
            if ($_.Exception.Response) {
                $reader = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream())
                $responseBody = $reader.ReadToEnd()
                Write-Host "错误详情: $responseBody" -ForegroundColor Red
            }
        }
    }
    
} catch {
    Write-Host "✗ 创建GoDaddy提供商失败: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream())
        $responseBody = $reader.ReadToEnd()
        Write-Host "错误详情: $responseBody" -ForegroundColor Red
    }
}

# 6. 再次获取提供商列表验证创建结果
Write-Host "`n6. 验证提供商列表..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri $BaseURL -Method GET -Headers $Headers
    Write-Host "✓ 获取提供商列表成功" -ForegroundColor Green
    
    # 检查是否有GoDaddy提供商
    $goDaddyProviders = $response.data.items | Where-Object { $_.type -eq "godaddy" }
    if ($goDaddyProviders) {
        Write-Host "✓ 找到 $($goDaddyProviders.Count) 个GoDaddy提供商" -ForegroundColor Green
        $goDaddyProviders | ForEach-Object {
            Write-Host "  - ID: $($_.id), 名称: $($_.name), 状态: $($_.status)" -ForegroundColor Cyan
        }
    } else {
        Write-Host "⚠ 未找到GoDaddy提供商" -ForegroundColor Yellow
    }
} catch {
    Write-Host "✗ 获取提供商列表失败: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== GoDaddy DNS测试完成 ===" -ForegroundColor Green
Write-Host "注意：由于使用的是测试凭证，连接测试应该会失败，这是正常的。" -ForegroundColor Yellow
Write-Host "请使用真实的GoDaddy API Key和Secret进行实际测试。" -ForegroundColor Yellow 
