/**
 * 前端提问频率限制（设备级 localStorage，多账号共用）
 * 免费用户每小时 1 次；VIP 跳过。权威校验仍以服务端为准。
 */

const STORAGE_KEY = 'fst_device_last_predict_at';
const PREDICT_INTERVAL_MS = 60 * 60 * 1000;

export function isVipUser(user) {
  if (!user?.tier || user.tier !== 'vip') return false;
  if (!user.vip_expires_at) return false;
  return new Date(user.vip_expires_at) > new Date();
}

function getLastPredictAt() {
  const raw = localStorage.getItem(STORAGE_KEY);
  if (!raw) return null;
  const ts = Number(raw);
  return Number.isNaN(ts) ? null : ts;
}

export const rateLimiter = {
  checkLimit({ isVip = false } = {}) {
    if (isVip) {
      return { allowed: true, submitInfo: 'Result' };
    }

    const lastAt = getLastPredictAt();
    if (lastAt === null) {
      return { allowed: true, submitInfo: 'Result' };
    }

    const elapsed = Date.now() - lastAt;
    if (elapsed >= PREDICT_INTERVAL_MS) {
      return { allowed: true, submitInfo: 'Result' };
    }

    const remainMins = Math.max(1, Math.ceil((PREDICT_INTERVAL_MS - elapsed) / 60000));
    return {
      allowed: false,
      submitInfo: 'LimitReached',
      message: `提问次数已达上限，请约 ${remainMins} 分钟后再来吧`,
    };
  },

  recordSubmission() {
    localStorage.setItem(STORAGE_KEY, String(Date.now()));
  },
};
