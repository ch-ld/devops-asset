@echo off
REM 统一编译脚本 - 用于构建DevOps Asset Management System

echo 正在编译服务器...

REM 编译主服务器
go build -o bin/server.exe cmd/main.go

if %errorlevel% neq 0 (
    echo 编译失败！请检查错误信息。
    exit /b %errorlevel%
)

echo 编译成功！服务器可执行文件已保存到 bin/server.exe

echo 使用方式：
echo   server.exe          - 以默认配置启动服务
echo   server.exe --dev    - 以开发模式启动服务
echo   server.exe --migrate - 执行数据库迁移
echo   server.exe --config path/to/config.yaml - 使用指定配置文件启动服务

echo 创建 server.bat 启动脚本...
(
echo @echo off
echo cd %%~dp0
echo bin\server.exe %%*
) > server.bat

echo 完成！ 
