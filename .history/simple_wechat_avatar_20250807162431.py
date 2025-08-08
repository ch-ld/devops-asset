# -*- coding: utf-8 -*-
"""
简易微信头像自动更换器 (Windows版)
每天自动生成对应的头像并提醒更换

注意：由于微信不提供官方API，本脚本采用安全方案：
1. 根据星期几生成对应头像
2. 自动打开头像文件夹
3. 您手动将头像拖拽到微信设置中即可

使用方法：
1. 将您的5张头像文件放在 avatars 文件夹中，命名为：
   - monday.png (周一)
   - tuesday.png (周二) 
   - wednesday.png (周三)
   - thursday.png (周四)
   - friday.png (周五)
2. 运行脚本：python simple_wechat_avatar.py
3. 程序会每天早上8点提醒您更换头像
"""

import datetime
import os
import shutil
import time
import schedule
from pathlib import Path
import subprocess

# =============== 配置 ===============
AVATAR_FOLDER = "avatars"  # 头像文件夹
OUTPUT_FOLDER = "current"  # 输出文件夹
REMINDER_TIME = "08:00"    # 提醒时间

# 工作日头像文件映射
WEEKDAY_AVATARS = {
    0: "monday.png",     # 周一
    1: "tuesday.png",    # 周二
    2: "wednesday.png",  # 周三
    3: "thursday.png",   # 周四
    4: "friday.png",     # 周五
}

def setup_folders():
    """初始化文件夹"""
    Path(AVATAR_FOLDER).mkdir(exist_ok=True)
    Path(OUTPUT_FOLDER).mkdir(exist_ok=True)
    
    # 检查头像文件是否存在
    missing_files = []
    for day, filename in WEEKDAY_AVATARS.items():
        file_path = Path(AVATAR_FOLDER) / filename
        if not file_path.exists():
            missing_files.append(filename)
    
    if missing_files:
        print("⚠️  请将以下头像文件放入 avatars 文件夹：")
        for filename in missing_files:
            print(f"   - {filename}")
        print("\n💡 您也可以使用任意图片文件，重命名即可")
        return False
    
    return True

def change_avatar():
    """更换头像"""
    today = datetime.datetime.now()
    weekday = today.weekday()
    
    # 只在工作日更换头像
    if weekday not in WEEKDAY_AVATARS:
        print(f"今天是{'周末' if weekday >= 5 else '周日'}，不需要更换头像")
        return
    
    # 获取今天的头像文件
    avatar_file = WEEKDAY_AVATARS[weekday]
    source_path = Path(AVATAR_FOLDER) / avatar_file
    target_path = Path(OUTPUT_FOLDER) / "today_avatar.png"
    
    if not source_path.exists():
        print(f"❌ 头像文件不存在: {source_path}")
        return
    
    try:
        # 复制头像到输出文件夹
        shutil.copy2(source_path, target_path)
        
        weekday_names = ["周一", "周二", "周三", "周四", "周五"]
        print(f"✅ 已准备好{weekday_names[weekday]}的头像: {target_path}")
        
        # 自动打开文件夹
        os.startfile(OUTPUT_FOLDER)
        
        # 显示简单提醒
        print("\n" + "="*50)
        print("🎯 微信头像更换提醒")
        print("="*50)
        print(f"📅 今天是{weekday_names[weekday]}，新头像已准备好！")
        print(f"📁 头像位置: {target_path.absolute()}")
        print("\n📱 更换步骤：")
        print("   1. 打开微信 -> 设置 -> 个人信息")
        print("   2. 点击头像")
        print("   3. 选择刚才打开的文件夹中的 today_avatar.png")
        print("="*50)
        
    except Exception as e:
        print(f"❌ 复制头像失败: {e}")

def main():
    """主函数"""
    print("🚀 微信头像自动更换器启动中...")
    print(f"⏰ 每天 {REMINDER_TIME} 自动提醒更换头像")
    print("💼 只在工作日(周一到周五)提醒")
    
    # 初始化
    if not setup_folders():
        input("\n按回车键退出...")
        return
    
    # 立即执行一次（测试）
    print("\n🧪 测试运行...")
    change_avatar()
    
    # 设置定时任务
    schedule.every().day.at(REMINDER_TIME).do(change_avatar)
    
    print(f"\n✅ 定时任务已设置！")
    print("🔄 程序正在后台运行...")
    print("❌ 按 Ctrl+C 退出程序\n")
    
    # 运行定时任务
    try:
        while True:
            schedule.run_pending()
            time.sleep(60)  # 每分钟检查一次
    except KeyboardInterrupt:
        print("\n👋 程序已退出")

if __name__ == "__main__":
    main() 
