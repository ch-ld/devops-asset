@echo off
chcp 65001 >nul
echo ==========================================
echo 微信头像自动更换器 - 安装脚本
echo ==========================================
echo.

echo 正在安装Python依赖...
pip install -r requirements.txt

echo.
echo 正在创建头像文件夹...
if not exist "avatars" mkdir "avatars"
if not exist "current" mkdir "current"

echo.
echo ==========================================
echo 安装完成！
echo ==========================================
echo.
echo 下一步：
echo 1. 将您的5张头像文件放入 avatars 文件夹
echo 2. 文件命名为：
echo    - monday.png (周一头像)
echo    - tuesday.png (周二头像)
echo    - wednesday.png (周三头像)
echo    - thursday.png (周四头像)
echo    - friday.png (周五头像)
echo 3. 双击 run.bat 启动程序
echo.
pause 
