// 测试时间格式化
const testTimes = [
  "2025-01-13T07:30:45.123Z",
  "2025-01-13 15:30:45",
  "1970-01-01T00:00:00Z",
  "0",
  "",
  null,
  undefined,
  "1641888000000", // 2022年毫秒时间戳
  "1641888000"     // 2022年秒时间戳
]

function formatTime(time) {
  if (!time) return '-'
  try {
    let date
    
    // 处理不同的时间格式
    if (typeof time === 'string') {
      // 如果是 ISO 字符串或标准时间字符串
      if (time.includes('T') || time.includes('-')) {
        date = new Date(time)
      } else {
        // 如果是时间戳（毫秒或秒）
        const timestamp = parseInt(time)
        if (timestamp < 9999999999) {
          // 秒级时间戳，转换为毫秒
          date = new Date(timestamp * 1000)
        } else {
          // 毫秒级时间戳
          date = new Date(timestamp)
        }
      }
    } else {
      date = new Date(time)
    }
    
    // 检查时间是否有效且不是1970年
    if (isNaN(date.getTime()) || date.getFullYear() < 2020) {
      console.warn('无效的时间数据:', time)
      return '时间数据异常'
    }
    
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false
    })
  } catch (error) {
    console.error('时间格式化错误:', error, '原始数据:', time)
    return '时间格式错误'
  }
}

console.log('时间格式化测试结果:')
testTimes.forEach(time => {
  console.log(`输入: ${time} -> 输出: ${formatTime(time)}`)
})
