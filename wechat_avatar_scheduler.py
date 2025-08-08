# -*- coding: utf-8 -*-
"""
微信头像智能管理器
功能：根据星期几自动生成对应头像，并提供提醒功能

注意：由于微信官方不提供头像更换API，且使用第三方库存在封号风险，
本脚本采用安全的方案：生成头像文件并提供提醒功能。

使用方法：
1. 准备7张不同的头像模板
2. 运行脚本，它会根据当天是星期几生成对应的头像
3. 脚本会弹出提醒，您可以手动将生成的头像设置为微信头像
4. 也可以配置自动打开文件夹，方便快速更换
"""

import datetime
import os
import shutil
import time
import logging
from pathlib import Path
import schedule
from PIL import Image, ImageDraw, ImageFont
import tkinter as tk
from tkinter import messagebox
import subprocess
import sys

# =============== 配置区域 ======================

# 头像模板配置（7天对应的头像）
AVATAR_TEMPLATES = {
    0: r"D:\avatars\templates\monday.png",      # 周一
    1: r"D:\avatars\templates\tuesday.png",     # 周二  
    2: r"D:\avatars\templates\wednesday.png",   # 周三
    3: r"D:\avatars\templates\thursday.png",    # 周四
    4: r"D:\avatars\templates\friday.png",      # 周五
    5: r"D:\avatars\templates\saturday.png",    # 周六
    6: r"D:\avatars\templates\sunday.png",      # 周日
}

# 输出配置
OUTPUT_DIR = r"D:\avatars\daily"
CURRENT_AVATAR_NAME = "current_avatar.png"

# 时间配置
CHANGE_TIME = "07:00"  # 每天7点提醒更换头像

# 提醒配置
ENABLE_POPUP_REMINDER = True    # 是否弹出提醒窗口
ENABLE_FOLDER_AUTO_OPEN = True  # 是否自动打开头像文件夹
ENABLE_SOUND_REMINDER = False   # 是否播放提醒音（需要安装额外依赖）

# 高级功能配置
ENABLE_WATERMARK = True         # 是否添加日期水印
WATERMARK_POSITION = "bottom"   # 水印位置: top, bottom, center
WATERMARK_COLOR = (255, 255, 255, 180)  # 水印颜色 (R,G,B,Alpha)

# =============== 核心功能 ======================

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('avatar_scheduler.log', encoding='utf-8'),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger(__name__)

class AvatarManager:
    def __init__(self):
        self.ensure_directories()
    
    def ensure_directories(self):
        """确保必要的目录存在"""
        Path(OUTPUT_DIR).mkdir(parents=True, exist_ok=True)
        
        # 检查模板文件是否存在
        missing_templates = []
        for day, template_path in AVATAR_TEMPLATES.items():
            if not Path(template_path).exists():
                missing_templates.append((day, template_path))
        
        if missing_templates:
            logger.warning("以下头像模板文件不存在:")
            for day, path in missing_templates:
                weekday_name = ['周一', '周二', '周三', '周四', '周五', '周六', '周日'][day]
                logger.warning(f"  {weekday_name}: {path}")
            
            # 创建示例模板
            self.create_sample_templates()
    
    def create_sample_templates(self):
        """创建示例头像模板"""
        logger.info("正在创建示例头像模板...")
        
        weekday_names = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
        colors = [
            (255, 99, 99),   # 红色 - 周一
            (99, 255, 99),   # 绿色 - 周二  
            (99, 99, 255),   # 蓝色 - 周三
            (255, 255, 99),  # 黄色 - 周四
            (255, 99, 255),  # 紫色 - 周五
            (99, 255, 255),  # 青色 - 周六
            (255, 165, 0),   # 橙色 - 周日
        ]
        
        for day, template_path in AVATAR_TEMPLATES.items():
            if not Path(template_path).exists():
                # 确保模板目录存在
                Path(template_path).parent.mkdir(parents=True, exist_ok=True)
                
                # 创建示例头像
                img = Image.new('RGB', (400, 400), colors[day])
                draw = ImageDraw.Draw(img)
                
                # 添加文字
                try:
                    font = ImageFont.truetype("arial.ttf", 60)
                except:
                    font = ImageFont.load_default()
                
                text = weekday_names[day]
                # 计算文字位置（居中）
                bbox = draw.textbbox((0, 0), text, font=font)
                text_width = bbox[2] - bbox[0]
                text_height = bbox[3] - bbox[1]
                x = (400 - text_width) // 2
                y = (400 - text_height) // 2
                
                draw.text((x, y), text, fill=(255, 255, 255), font=font)
                img.save(template_path)
                logger.info(f"已创建示例模板: {template_path}")
    
    def add_watermark(self, image, text):
        """为图片添加日期水印"""
        if not ENABLE_WATERMARK:
            return image
        
        draw = ImageDraw.Draw(image)
        
        try:
            font = ImageFont.truetype("arial.ttf", 24)
        except:
            font = ImageFont.load_default()
        
        # 计算水印位置
        bbox = draw.textbbox((0, 0), text, font=font)
        text_width = bbox[2] - bbox[0]
        text_height = bbox[3] - bbox[1]
        
        img_width, img_height = image.size
        
        if WATERMARK_POSITION == "top":
            x, y = 10, 10
        elif WATERMARK_POSITION == "bottom":
            x, y = img_width - text_width - 10, img_height - text_height - 10
        else:  # center
            x, y = (img_width - text_width) // 2, (img_height - text_height) // 2
        
        # 添加半透明背景
        overlay = Image.new('RGBA', image.size, (0, 0, 0, 0))
        overlay_draw = ImageDraw.Draw(overlay)
        
        # 绘制背景矩形
        overlay_draw.rectangle([x-5, y-5, x+text_width+5, y+text_height+5], 
                             fill=(0, 0, 0, 100))
        
        # 绘制文字
        overlay_draw.text((x, y), text, fill=WATERMARK_COLOR, font=font)
        
        # 合并图层
        if image.mode != 'RGBA':
            image = image.convert('RGBA')
        
        return Image.alpha_composite(image, overlay).convert('RGB')
    
    def generate_daily_avatar(self):
        """根据当天星期几生成对应头像"""
        today = datetime.datetime.now()
        weekday = today.weekday()  # 0=Monday, 6=Sunday
        weekday_names = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
        
        template_path = AVATAR_TEMPLATES.get(weekday)
        if not template_path or not Path(template_path).exists():
            logger.error(f"今天({weekday_names[weekday]})的头像模板不存在: {template_path}")
            return None
        
        try:
            # 加载模板图片
            img = Image.open(template_path)
            
            # 添加日期水印
            date_text = today.strftime("%Y-%m-%d")
            img = self.add_watermark(img, date_text)
            
            # 保存当前头像
            output_path = Path(OUTPUT_DIR) / CURRENT_AVATAR_NAME
            img.save(output_path)
            
            logger.info(f"已生成今天({weekday_names[weekday]})的头像: {output_path}")
            return output_path
            
        except Exception as e:
            logger.error(f"生成头像时出错: {e}")
            return None
    
    def show_reminder(self, avatar_path):
        """显示更换头像提醒"""
        if not ENABLE_POPUP_REMINDER:
            return
        
        try:
            root = tk.Tk()
            root.withdraw()  # 隐藏主窗口
            
            weekday_names = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
            today_name = weekday_names[datetime.datetime.now().weekday()]
            
            message = f"""
🎯 该更换微信头像啦！

今天是{today_name}，新头像已准备好：
{avatar_path}

点击"确定"打开头像文件夹
点击"取消"跳过此次提醒
            """.strip()
            
            result = messagebox.askokcancel("微信头像提醒", message)
            
            if result and ENABLE_FOLDER_AUTO_OPEN:
                self.open_avatar_folder()
            
            root.destroy()
            
        except Exception as e:
            logger.error(f"显示提醒窗口时出错: {e}")
    
    def open_avatar_folder(self):
        """打开头像文件夹"""
        try:
            if sys.platform == "win32":
                os.startfile(OUTPUT_DIR)
            elif sys.platform == "darwin":  # macOS
                subprocess.run(["open", OUTPUT_DIR])
            else:  # Linux
                subprocess.run(["xdg-open", OUTPUT_DIR])
            logger.info(f"已打开头像文件夹: {OUTPUT_DIR}")
        except Exception as e:
            logger.error(f"打开文件夹时出错: {e}")
    
    def daily_task(self):
        """每日任务：生成头像并提醒"""
        logger.info("开始执行每日头像更换任务...")
        
        avatar_path = self.generate_daily_avatar()
        if avatar_path:
            self.show_reminder(avatar_path)
            logger.info("每日头像更换任务完成")
        else:
            logger.error("每日头像更换任务失败")

def main():
    """主函数"""
    print("=" * 50)
    print("微信头像智能管理器")
    print("=" * 50)
    print(f"配置信息:")
    print(f"  提醒时间: {CHANGE_TIME}")
    print(f"  输出目录: {OUTPUT_DIR}")
    print(f"  弹窗提醒: {'开启' if ENABLE_POPUP_REMINDER else '关闭'}")
    print(f"  自动打开文件夹: {'开启' if ENABLE_FOLDER_AUTO_OPEN else '关闭'}")
    print(f"  日期水印: {'开启' if ENABLE_WATERMARK else '关闭'}")
    print("=" * 50)
    
    # 初始化头像管理器
    avatar_manager = AvatarManager()
    
    # 立即执行一次（用于测试）
    print("正在生成今天的头像...")
    avatar_manager.daily_task()
    
    # 设置定时任务
    schedule.every().day.at(CHANGE_TIME).do(avatar_manager.daily_task)
    
    print(f"\n✅ 定时任务已启动！")
    print(f"📅 每天 {CHANGE_TIME} 将自动生成新头像并提醒您更换")
    print("🔄 程序正在后台运行，按 Ctrl+C 退出...")
    
    # 主循环
    try:
        while True:
            schedule.run_pending()
            time.sleep(30)  # 每30秒检查一次
    except KeyboardInterrupt:
        print("\n👋 程序已退出")
        logger.info("程序被用户手动退出")

if __name__ == "__main__":
    main() 
