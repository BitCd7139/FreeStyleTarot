/**
 * 前端提交频率限制工具
 * 逻辑：限制 30 分钟内最多提交 5 次
 */

const STORAGE_KEY = 'tarot_submit_history';
const MAX_COUNT = 5;
const WINDOW_MS = 300 * 60 * 1000; 

export const rateLimiter = {
  /**
   * 检查是否允许提交
   * @returns {Object} { allowed: boolean, message: string }
   */
  checkLimit() {
    const now = Date.now();
    const history = JSON.parse(localStorage.getItem(STORAGE_KEY) || '[]');

    // 1. 过滤掉 30 分钟之前的记录
    const recentSubmits = history.filter(timestamp => now - timestamp < WINDOW_MS);
    
    // 更新一下存储（清理掉过期的）
    localStorage.setItem(STORAGE_KEY, JSON.stringify(recentSubmits));

    if (recentSubmits.length >= MAX_COUNT) {
      // 计算最早的一条记录还有多久过期，提示用户等待时间
      const oldestTimestamp = recentSubmits[0];
      const waitTimeMs = WINDOW_MS - (now - oldestTimestamp);
      const waitMinutes = Math.ceil(waitTimeMs / (1000 * 60));
      
      return {
        allowed: false,
        message: `为了占卜的严谨性，请勿频繁求问。请等待 ${waitMinutes} 分钟后再试。`
      };
    }

    return { allowed: true };
  },

  /**
   * 记录一次成功的提交
   */
  recordSubmission() {
    const now = Date.now();
    const history = JSON.parse(localStorage.getItem(STORAGE_KEY) || '[]');
    history.push(now);
    // 只保留最近的 MAX_COUNT 条即可
    localStorage.setItem(STORAGE_KEY, JSON.stringify(history.slice(-MAX_COUNT)));
  }
};