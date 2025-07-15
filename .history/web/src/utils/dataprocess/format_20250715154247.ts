/**
 * 数据格式化相关工具函数
 */

import { formatTimestamp as formatTimestampUtil } from '../time'

// 时间戳转时间
export function timestampToTime(timestamp: number = Date.now(), isMs: boolean = true): string {
  const date = new Date(isMs ? timestamp : timestamp * 1000)
  return date.toISOString().replace('T', ' ').slice(0, 19)
}

// 数字格式化（千位分隔符）
export function commafy(num: number): string {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

// 生成随机数
export function randomNum(min: number, max?: number): number {
  if (max === undefined) {
    max = min
    min = 0
  }
  return Math.floor(Math.random() * (max - min + 1)) + min
}

// 移除HTML标签
export function removeHtmlTags(str: string = ''): string {
  return str.replace(/<[^>]*>/g, '')
}

/**
 * 格式化时间戳
 * @param timestamp 秒级时间戳
 * @param format 格式化模板，默认为 'YYYY-MM-DD HH:mm:ss'
 * @returns 格式化后的时间字符串
 */
export function formatTimestamp(timestamp: number, format: string = 'YYYY-MM-DD HH:mm:ss'): string {
  return formatTimestampUtil(timestamp, format)
}

/**
 * 格式化JSON数据
 * @param jsonData JSON字符串或对象
 * @param indentation 缩进空格数，默认为2
 * @returns 格式化后的JSON字符串
 */
export function formatJsonData(jsonData: string | object, indentation: number = 2): string {
  try {
    // 如果传入的是字符串，先尝试解析为对象
    const obj = typeof jsonData === 'string' ? JSON.parse(jsonData) : jsonData

    // 再将对象转换为格式化的字符串
    return JSON.stringify(obj, null, indentation)
  } catch (error) {
    // 如果解析失败，返回原始字符串
    return typeof jsonData === 'string' ? jsonData : JSON.stringify(jsonData)
  }
}
