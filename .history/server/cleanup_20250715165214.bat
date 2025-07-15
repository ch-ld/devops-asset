@echo off
REM 清理冗余main.go文件

echo 正在备份冗余文件...
if not exist backup mkdir backup
if not exist backup\cmd mkdir backup\cmd
if not exist backup\cmd\app mkdir backup\cmd\app
if not exist backup\cmd\server mkdir backup\cmd\server

REM 备份文件
copy cmd\app\main.go backup\cmd\app\ /Y
copy cmd\server\main.go backup\cmd\server\ /Y

echo 正在删除冗余文件...
del cmd\app\main.go /F /Q
del cmd\server\main.go /F /Q

echo 清理完成！冗余文件已备份到 backup 目录。
echo 现在项目只有一个统一的入口点：cmd\main.go

echo 启动服务的方法：
echo   1. 直接运行: go run cmd\main.go --dev
echo   2. 编译后运行: build.bat 然后运行 server.bat 
