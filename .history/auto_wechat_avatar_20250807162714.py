# -*- coding: utf-8 -*-
"""
微信头像自动更换器 - UI自动化版本
通过模拟鼠标点击和键盘操作来自动更换微信头像

功能：
1. 每天自动根据星期几更换对应头像
2. 自动打开微信设置页面
3. 自动上传新头像
4. 支持定时任务

注意：
- 需要微信桌面版已登录
- 运行时请不要操作鼠标键盘
- 建议在固定分辨率下使用
"""

import datetime
import os
import time
import schedule
from pathlib import Path
import pyautogui
import pygetwindow as gw
import cv2
import numpy as np

# 禁用pyautogui的安全机制
pyautogui.FAILSAFE = False

# =============== 配置 ===============
AVATAR_FOLDER = "avatars"
CHANGE_TIME = "08:00"  # 每天8点自动更换

# 工作日头像映射
WEEKDAY_AVATARS = {
    0: "monday.png",     # 周一
    1: "tuesday.png",    # 周二
    2: "wednesday.png",  # 周三
    3: "thursday.png",   # 周四
    4: "friday.png",     # 周五
}

# UI元素识别配置
WECHAT_WINDOW_TITLES = ["微信", "WeChat"]
CLICK_DELAY = 1.0  # 点击间隔
CONFIDENCE = 0.8   # 图像识别置信度

class WeChatAvatarChanger:
    def __init__(self):
        self.setup_folders()
        self.wechat_window = None
    
    def setup_folders(self):
        """初始化文件夹"""
        Path(AVATAR_FOLDER).mkdir(exist_ok=True)
    
    def log(self, message):
        """日志输出"""
        timestamp = datetime.datetime.now().strftime("%H:%M:%S")
        print(f"[{timestamp}] {message}")
    
    def find_wechat_window(self):
        """查找微信窗口"""
        for title in WECHAT_WINDOW_TITLES:
            try:
                windows = gw.getWindowsWithTitle(title)
                if windows:
                    self.wechat_window = windows[0]
                    self.log(f"找到微信窗口: {title}")
                    return True
            except:
                continue
        return False
    
    def activate_wechat(self):
        """激活微信窗口"""
        if not self.wechat_window:
            if not self.find_wechat_window():
                self.log("❌ 未找到微信窗口，请确保微信已启动并登录")
                return False
        
        try:
            # 激活窗口
            self.wechat_window.activate()
            time.sleep(1)
            
            # 确保窗口在前台
            if self.wechat_window.isMinimized:
                self.wechat_window.restore()
            
            self.log("✅ 微信窗口已激活")
            return True
        except Exception as e:
            self.log(f"❌ 激活微信窗口失败: {e}")
            return False
    
    def click_with_retry(self, x, y, retries=3):
        """带重试的点击"""
        for i in range(retries):
            try:
                pyautogui.click(x, y)
                time.sleep(CLICK_DELAY)
                return True
            except Exception as e:
                self.log(f"点击失败 ({i+1}/{retries}): {e}")
                time.sleep(0.5)
        return False
    
    def find_and_click_image(self, image_path, confidence=CONFIDENCE):
        """查找并点击图像"""
        try:
            location = pyautogui.locateOnScreen(image_path, confidence=confidence)
            if location:
                center = pyautogui.center(location)
                self.click_with_retry(center.x, center.y)
                return True
        except Exception as e:
            self.log(f"查找图像失败 {image_path}: {e}")
        return False
    
    def open_settings_by_menu(self):
        """通过菜单打开设置"""
        try:
            # 方法1: 尝试点击设置按钮（通常在左下角）
            # 这里需要根据具体的微信界面调整坐标
            window_rect = (self.wechat_window.left, self.wechat_window.top, 
                          self.wechat_window.width, self.wechat_window.height)
            
            # 点击左下角设置区域
            settings_x = window_rect[0] + 50
            settings_y = window_rect[1] + window_rect[3] - 50
            
            self.click_with_retry(settings_x, settings_y)
            time.sleep(1)
            
            # 方法2: 使用快捷键
            pyautogui.hotkey('ctrl', ',')  # 微信设置快捷键
            time.sleep(2)
            
            self.log("✅ 尝试打开设置页面")
            return True
            
        except Exception as e:
            self.log(f"❌ 打开设置失败: {e}")
            return False
    
    def click_avatar_area(self):
        """点击头像区域"""
        try:
            # 在设置页面查找头像位置
            # 通常头像在设置页面的左上角区域
            time.sleep(2)
            
            # 获取当前活动窗口
            current_window = gw.getActiveWindow()
            if current_window:
                # 点击设置窗口左上角头像区域
                avatar_x = current_window.left + 100
                avatar_y = current_window.top + 150
                
                self.click_with_retry(avatar_x, avatar_y)
                time.sleep(1)
                
                self.log("✅ 点击头像区域")
                return True
        except Exception as e:
            self.log(f"❌ 点击头像失败: {e}")
        return False
    
    def upload_new_avatar(self, avatar_path):
        """上传新头像"""
        try:
            # 等待文件选择对话框出现
            time.sleep(2)
            
            # 输入文件路径
            pyautogui.write(str(Path(avatar_path).absolute()))
            time.sleep(0.5)
            
            # 按回车确认
            pyautogui.press('enter')
            time.sleep(2)
            
            # 如果有裁剪界面，点击确认
            # 通常会有"确定"或"完成"按钮
            pyautogui.press('enter')  # 或者点击确定按钮
            time.sleep(1)
            
            self.log("✅ 头像上传完成")
            return True
            
        except Exception as e:
            self.log(f"❌ 上传头像失败: {e}")
            return False
    
    def change_avatar_auto(self):
        """自动更换头像的完整流程"""
        today = datetime.datetime.now()
        weekday = today.weekday()
        
        # 只在工作日更换
        if weekday not in WEEKDAY_AVATARS:
            self.log("今天不是工作日，跳过头像更换")
            return
        
        avatar_file = WEEKDAY_AVATARS[weekday]
        avatar_path = Path(AVATAR_FOLDER) / avatar_file
        
        if not avatar_path.exists():
            self.log(f"❌ 头像文件不存在: {avatar_path}")
            return
        
        weekday_names = ["周一", "周二", "周三", "周四", "周五"]
        self.log(f"🚀 开始更换{weekday_names[weekday]}头像...")
        
        try:
            # 1. 激活微信窗口
            if not self.activate_wechat():
                return
            
            # 2. 打开设置
            if not self.open_settings_by_menu():
                return
            
            # 3. 点击头像
            if not self.click_avatar_area():
                return
            
            # 4. 上传新头像
            if not self.upload_new_avatar(avatar_path):
                return
            
            self.log(f"🎉 {weekday_names[weekday]}头像更换成功！")
            
        except Exception as e:
            self.log(f"❌ 更换头像过程中出错: {e}")
    
    def manual_test(self):
        """手动测试模式"""
        self.log("🧪 进入手动测试模式")
        self.log("请确保微信已打开并登录")
        input("按回车开始测试...")
        
        self.change_avatar_auto()

def setup_environment():
    """环境检查和设置"""
    print("🔧 检查运行环境...")
    
    # 检查头像文件
    missing_files = []
    for day, filename in WEEKDAY_AVATARS.items():
        file_path = Path(AVATAR_FOLDER) / filename
        if not file_path.exists():
            missing_files.append(filename)
    
    if missing_files:
        print("⚠️  请将以下头像文件放入 avatars 文件夹：")
        for filename in missing_files:
            print(f"   - {filename}")
        return False
    
    print("✅ 环境检查完成")
    return True

def main():
    """主函数"""
    print("=" * 50)
    print("🤖 微信头像自动更换器 (UI自动化版)")
    print("=" * 50)
    
    if not setup_environment():
        input("按回车退出...")
        return
    
    changer = WeChatAvatarChanger()
    
    # 询问运行模式
    print("\n选择运行模式:")
    print("1. 立即测试 (测试当前功能)")
    print("2. 定时运行 (每天自动运行)")
    
    choice = input("请选择 (1 或 2): ").strip()
    
    if choice == "1":
        changer.manual_test()
    elif choice == "2":
        # 设置定时任务
        schedule.every().day.at(CHANGE_TIME).do(changer.change_avatar_auto)
        
        print(f"\n✅ 定时任务已设置！每天 {CHANGE_TIME} 自动更换头像")
        print("🔄 程序正在后台运行...")
        print("❌ 按 Ctrl+C 退出\n")
        
        try:
            while True:
                schedule.run_pending()
                time.sleep(60)
        except KeyboardInterrupt:
            print("\n👋 程序已退出")
    else:
        print("❌ 无效选择")

if __name__ == "__main__":
    main() 
