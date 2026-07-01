const API_BASE = '/api';
// 使用独立的 storage key，与登录 cookie 完全隔离。
// localStorage 仅在本设备本浏览器生效，不会随 HTTP 请求发送到服务端。
const STORAGE_KEY = 'fst_custom_api_config';

import { ref, computed } from 'vue';

const config = ref(loadConfig());

function loadConfig() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return null;
    const parsed = JSON.parse(raw);
    if (!parsed || !parsed.api_key) return null;
    return parsed;
  } catch {
    return null;
  }
}

function persistConfig(cfg) {
  config.value = cfg;
  if (cfg) {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(cfg));
  } else {
    localStorage.removeItem(STORAGE_KEY);
  }
}

export function useCustomApi() {
  const isEnabled = computed(() => !!config.value && !!config.value.api_key);

  function getConfig() {
    return config.value;
  }

  /**
   * 保存配置到设备 localStorage。
   * 调用方应先通过 testConfig 验证配置可用后再保存。
   */
  function saveConfig(cfg) {
    persistConfig(cfg);
  }

  function clearConfig() {
    persistConfig(null);
  }

  /**
   * 测试自定义 API 配置。
   * 返回 { ok, message, reply } 或抛出网络错误。
   */
  async function testConfig(cfg) {
    const response = await fetch(`${API_BASE}/custom-api/test`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        api_key: cfg.api_key,
        base_url: cfg.base_url,
        model: cfg.model,
        format: cfg.format || 'openai',
        stage_params: cfg.stage_params || undefined,
      }),
    });
    const data = await response.json().catch(() => ({}));
    if (!response.ok) {
      const err = new Error(data.error || data.message || `测试失败 (HTTP ${response.status})`);
      err.status = response.status;
      err.data = data;
      throw err;
    }
    return data;
  }

  return {
    config,
    isEnabled,
    getConfig,
    saveConfig,
    clearConfig,
    testConfig,
  };
}
