# GoDaddy DNS Provider Test Script
$BaseURL = "http://localhost:8080/api/v1/dns/providers"
$Headers = @{
    "Content-Type" = "application/json"
}

Write-Host "=== GoDaddy DNS Provider Test ===" -ForegroundColor Green

# 1. Test connection with invalid credentials
Write-Host "`n1. Testing GoDaddy connection (invalid credentials)..." -ForegroundColor Yellow

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
    Write-Host "Response: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Cyan
} catch {
    Write-Host "Request failed (expected due to invalid credentials): $($_.Exception.Message)" -ForegroundColor Yellow
}

# 2. Create GoDaddy provider
Write-Host "`n2. Creating GoDaddy provider..." -ForegroundColor Yellow

$createBody = @{
    name = "Test GoDaddy"
    type = "godaddy"
    credentials = @{
        api_key = "test-key-123"
        api_secret = "test-secret-456"
    }
    remark = "Test provider"
}

$createJson = $createBody | ConvertTo-Json -Depth 3

try {
    $response = Invoke-RestMethod -Uri $BaseURL -Method POST -Headers $Headers -Body $createJson
    Write-Host "Success: Created provider" -ForegroundColor Green
    Write-Host "Response: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Cyan
} catch {
    Write-Host "Failed to create provider: $($_.Exception.Message)" -ForegroundColor Red
}

# 3. Get provider list
Write-Host "`n3. Getting provider list..." -ForegroundColor Yellow

try {
    $response = Invoke-RestMethod -Uri $BaseURL -Method GET -Headers $Headers
    Write-Host "Success: Got provider list" -ForegroundColor Green
    Write-Host "Response: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Cyan
} catch {
    Write-Host "Failed to get provider list: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== Test Complete ===" -ForegroundColor Green 
