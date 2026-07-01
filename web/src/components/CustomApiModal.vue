<template>
  <Transition name="fade">
    <div v-if="show" class="custom-api-overlay" @mousedown="onOverlayMouseDown" @mouseup="onOverlayMouseUp">
      <div class="custom-api-panel" @mousedown.stop>
        <div class="modal-header">
          <h3>自定义 API</h3>
          <p class="modal-hint">填入自定义 API 配置，可无限自由提问。配置仅保存在本地，服务器不存任何数据。</p>
        </div>

        <div class="modal-body">
          <div v-if="testResult" class="test-result" :class="{ success: testResult.ok, fail: !testResult.ok }">
            <p class="test-message">{{ testResult.message }}</p>
            <p v-if="testResult.reply" class="test-reply">返回内容：{{ testResult.reply }}</p>
          </div>

          <div class="form-group">
            <label class="field-label">API 格式</label>
            <div class="format-selector">
              <button
                type="button"
                class="format-chip"
                :class="{ selected: form.format === 'openai' }"
                @click="form.format = 'openai'"
              >
                OpenAI
              </button>
              <button
                type="button"
                class="format-chip"
                :class="{ selected: form.format === 'anthropic' }"
                @click="form.format = 'anthropic'"
              >
                Anthropic
              </button>
            </div>
          </div>

          <div class="form-group">
            <label class="field-label">API Base URL</label>
            <input
              v-model="form.base_url"
              type="text"
              class="auth-input"
              :placeholder="form.format === 'anthropic' ? 'https://api.anthropic.com' : 'https://api.openai.com/v1'"
            />
          </div>

          <div class="form-group">
            <label class="field-label">API Key</label>
            <input
              v-model="form.api_key"
              type="password"
              class="auth-input"
              placeholder="sk-..."
              autocomplete="off"
            />
          </div>

          <div class="form-group">
            <label class="field-label">Model</label>
            <input
              v-model="form.model"
              type="text"
              class="auth-input"
              :placeholder="form.format === 'anthropic' ? 'claude-3-5-sonnet-20241022' : 'gpt-4o-mini'"
            />
          </div>

          <!-- 阶段参数 -->
          <div class="stage-section">
            <div class="stage-section-title">
              <span>阶段参数（可选）</span>
              <span
                class="help-icon"
                @mouseenter="showStageHelp = true"
                @mouseleave="showStageHelp = false"
                @touchstart="onStageHelpTouchStart"
                @touchend="onStageHelpTouchEnd"
              >?</span>
              <div v-if="showStageHelp" class="help-tooltip">
                <p><b>阶段参数说明，留空则使用服务器默认值。</b></p>
                <p><b>意图澄清</b>：分析用户提问意图，判断是否需要追问。通常输出较短，温度偏低以保证稳定。</p>
                <p><b>牌阵分析</b>：分析牌面含义与位置关系，输出结构化解读。需要一定长度和创造性。</p>
                <p><b>综合建议</b>：基于牌阵分析给出专业建议和行动指引，不重复单卡解析。</p>
                <p><b>角色翻译</b>：将专业建议转换为选定角色风格的最终输出，只做格式转换不额外思考。</p>
                <p><b>max_tokens</b>：模型生成的最大 token 数。值越大回复越长，但消耗更多额度。</p>
                <p><b>temperature</b>：控制生成随机性。0 更确定保守，1 更随机发散，通常 0.3-0.8。</p>
                <p><b>说明</b>：默认参数未必有最好效果，只是服务器权衡成本填入的数值。</p>
              </div>
            </div>

            <div v-for="stage in stages" :key="stage.key" class="stage-block">
              <div class="stage-header">
                <span class="stage-name">{{ stage.label }}</span>
              </div>
              <div class="stage-fields">
                <div class="stage-field">
                  <label class="micro-label">max_tokens</label>
                  <input
                    v-model.number="form.stage_params[stage.key].max_tokens"
                    type="number"
                    min="0"
                    class="auth-input stage-input"
                    :placeholder="String(stage.defaultTokens)"
                  />
                </div>
                <div class="stage-field">
                  <label class="micro-label">temperature</label>
                  <input
                    v-model.number="form.stage_params[stage.key].temperature"
                    type="number"
                    min="0"
                    max="2"
                    step="0.1"
                    class="auth-input stage-input"
                    :placeholder="stage.defaultTemp.toFixed(1)"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button
            type="button"
            class="test-btn"
            :disabled="testing || saving"
            @click="handleTest"
          >
            {{ testing ? '测试中…' : '测试连接' }}
          </button>
          <button
            type="button"
            class="clear-btn"
            v-if="isEnabled"
            :disabled="testing || saving"
            @click="handleClear"
          >
            清除配置
          </button>
          <button
            type="button"
            class="save-btn"
            :disabled="testing || saving || !canSave || !testPassed"
            @click="handleSave"
          >
            {{ saving ? '保存中…' : (testPassed ? '保存' : '需先通过测试') }}
          </button>
          <button type="button" class="back-btn" @click="show = false">关闭</button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue';
import { useCustomApi } from '../composables/useCustomApi.js';

const show = defineModel('show', { type: Boolean, default: false });

const { config, isEnabled, saveConfig, clearConfig, testConfig } = useCustomApi();

const form = ref(createEmptyForm());
const testing = ref(false);
const saving = ref(false);
const testResult = ref(null);
const testPassed = ref(false);
const showStageHelp = ref(false);
// 标记是否正在自动保存回写，避免 watch 触发循环
let autoSaving = false;

// 阶段定义，含默认值
const stages = reactive([
  {
    key: 'intent_clarifier',
    label: '意图澄清',
    defaultTokens: 2048,
    defaultTemp: 0.4,
  },
  {
    key: 'spread_analyst',
    label: '牌阵分析',
    defaultTokens: 8192,
    defaultTemp: 0.6,
  },
  {
    key: 'advisor',
    label: '综合建议',
    defaultTokens: 4096,
    defaultTemp: 0.5,
  },
  {
    key: 'persona',
    label: '角色翻译',
    defaultTokens: 4096,
    defaultTemp: 0.3,
  },
]);

function createEmptyForm() {
  return {
    api_key: '',
    base_url: '',
    model: '',
    format: 'openai',
    stage_params: {
      intent_clarifier: { max_tokens: null, temperature: null },
      spread_analyst: { max_tokens: null, temperature: null },
      advisor: { max_tokens: null, temperature: null },
      persona: { max_tokens: null, temperature: null },
    },
  };
}

const canSave = computed(() => {
  return (
    form.value.api_key.trim() &&
    form.value.base_url.trim() &&
    form.value.model.trim()
  );
});

watch(show, (open) => {
  if (open) {
    testResult.value = null;
    testPassed.value = false;
    // 初始化 form 时暂停自动保存，避免回写触发 watch
    autoSaving = true;
    if (config.value) {
      form.value = {
        ...createEmptyForm(),
        ...config.value,
        stage_params: {
          intent_clarifier: { max_tokens: null, temperature: null, ...(config.value.stage_params?.intent_clarifier || {}) },
          spread_analyst: { max_tokens: null, temperature: null, ...(config.value.stage_params?.spread_analyst || {}) },
          advisor: { max_tokens: null, temperature: null, ...(config.value.stage_params?.advisor || {}) },
          persona: { max_tokens: null, temperature: null, ...(config.value.stage_params?.persona || {}) },
        },
      };
    } else {
      form.value = createEmptyForm();
    }
    // 下一轮事件循环恢复自动保存
    setTimeout(() => { autoSaving = false; }, 0);
  }
});

async function handleTest() {
  testResult.value = null;
  testing.value = true;
  try {
    const result = await testConfig(form.value);
    testResult.value = result;
    testPassed.value = !!result.ok;
  } catch (err) {
    testResult.value = {
      ok: false,
      message: err.message || '测试失败，请检查网络或配置',
    };
    testPassed.value = false;
  } finally {
    testing.value = false;
  }
}

function handleSave() {
  if (!testPassed.value) return;
  saving.value = true;
  try {
    // 清理空值，避免存储无意义字段
    const cleanStageParams = {};
    for (const stage of stages) {
      const sp = form.value.stage_params[stage.key];
      const cleaned = {};
      if (sp.max_tokens && sp.max_tokens > 0) cleaned.max_tokens = sp.max_tokens;
      if (sp.temperature && sp.temperature > 0) cleaned.temperature = sp.temperature;
      if (Object.keys(cleaned).length > 0) cleanStageParams[stage.key] = cleaned;
    }
    const cfg = {
      api_key: form.value.api_key.trim(),
      base_url: form.value.base_url.trim(),
      model: form.value.model.trim(),
      format: form.value.format,
    };
    if (Object.keys(cleanStageParams).length > 0) {
      cfg.stage_params = cleanStageParams;
    }
    saveConfig(cfg);
    show.value = false;
  } finally {
    saving.value = false;
  }
}

function handleClear() {
  autoSaving = true;
  clearConfig();
  form.value = createEmptyForm();
  testResult.value = null;
  testPassed.value = false;
  show.value = false;
  setTimeout(() => { autoSaving = false; }, 0);
}

/**
 * 自动保存 model 和阶段参数到 localStorage。
 * 仅在已存在配置（用户之前测试通过并保存过）时生效，
 * 这样用户修改这些非敏感字段后无需再点保存按钮。
 * api_key / base_url / format 仍需通过测试后手动保存。
 */
function autoSaveModelAndStageParams() {
  // 未保存过配置时不自动保存
  if (!config.value || !config.value.api_key) return;
  // 正在回写 form 时跳过
  if (autoSaving) return;

  const merged = {
    api_key: config.value.api_key,
    base_url: config.value.base_url,
    model: form.value.model.trim() || config.value.model,
    format: config.value.format,
  };

  // 清理并合并阶段参数
  const cleanStageParams = {};
  for (const stage of stages) {
    const sp = form.value.stage_params[stage.key];
    const cleaned = {};
    if (sp.max_tokens && sp.max_tokens > 0) cleaned.max_tokens = sp.max_tokens;
    if (sp.temperature && sp.temperature > 0) cleaned.temperature = sp.temperature;
    if (Object.keys(cleaned).length > 0) cleanStageParams[stage.key] = cleaned;
  }
  if (Object.keys(cleanStageParams).length > 0) {
    merged.stage_params = cleanStageParams;
  }
  saveConfig(merged);
}

// 监听 model 变化自动保存
watch(() => form.value.model, () => {
  autoSaveModelAndStageParams();
});

// 监听阶段参数变化自动保存（深度监听）
watch(() => form.value.stage_params, () => {
  autoSaveModelAndStageParams();
}, { deep: true });

// 移动端长按帮助提示
let helpTouchTimer = null;
function onStageHelpTouchStart() {
  helpTouchTimer = setTimeout(() => { showStageHelp.value = true; }, 500);
}
function onStageHelpTouchEnd() {
  if (helpTouchTimer) { clearTimeout(helpTouchTimer); helpTouchTimer = null; }
  setTimeout(() => { showStageHelp.value = false; }, 3000);
}

let mouseDownOutside = false;
function onOverlayMouseDown() {
  mouseDownOutside = true;
}
function onOverlayMouseUp() {
  if (mouseDownOutside) {
    show.value = false;
  }
  mouseDownOutside = false;
}
</script>

<style scoped>
.custom-api-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.65);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1200;
  padding: 16px;
}

.custom-api-panel {
  background: rgba(30, 27, 36, 0.96);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(194, 163, 95, 0.4);
  border-radius: 16px;
  width: min(520px, 100%);
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.6);
}

.modal-header {
  padding: 20px 25px 10px;
  border-bottom: 1px solid rgba(194, 163, 95, 0.2);
}

.modal-header h3 {
  margin: 0 0 4px;
  color: #e5d8b0;
  font-size: 18px;
}

.modal-hint {
  margin: 0;
  color: #a89f91;
  font-size: 12px;
  line-height: 1.5;
}

.modal-body {
  padding: 20px 25px;
  overflow-y: auto;
  flex: 1;
}

.form-group {
  margin-bottom: 16px;
}

.field-label {
  display: block;
  color: #a89f91;
  font-size: 13px;
  margin-bottom: 6px;
}

.auth-input {
  width: 100%;
  box-sizing: border-box;
  padding: 10px 12px;
  border-radius: 8px;
  border: 1px solid rgba(194, 163, 95, 0.3);
  background: rgba(0, 0, 0, 0.3);
  color: #e5d8b0;
  font-size: 14px;
}

.auth-input:focus {
  outline: none;
  border-color: #c2a35f;
}

/* 移除 number input 的上下箭头（spinner），仅保留文本输入 */
.auth-input[type="number"]::-webkit-outer-spin-button,
.auth-input[type="number"]::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}
.auth-input[type="number"] {
  -moz-appearance: textfield;
  appearance: textfield;
}

.format-selector {
  display: flex;
  gap: 10px;
}

.format-chip {
  flex: 1;
  padding: 10px 14px;
  border-radius: 10px;
  border: 1px solid rgba(194, 163, 95, 0.35);
  background: rgba(194, 163, 95, 0.06);
  color: #e5d8b0;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.format-chip:hover {
  background: rgba(194, 163, 95, 0.12);
  border-color: rgba(194, 163, 95, 0.55);
}

.format-chip.selected {
  background: rgba(194, 163, 95, 0.22);
  border-color: #c2a35f;
  box-shadow: 0 0 12px rgba(194, 163, 95, 0.25);
}

/* 阶段参数 */
.stage-section {
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid rgba(194, 163, 95, 0.15);
}

.stage-section-title {
  position: relative;
  display: flex;
  align-items: center;
  gap: 8px;
  color: #c2a35f;
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 12px;
}

.stage-block {
  margin-bottom: 14px;
  padding: 12px;
  border-radius: 10px;
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(194, 163, 95, 0.15);
}

.stage-header {
  position: relative;
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 8px;
}

.stage-name {
  color: #e5d8b0;
  font-size: 13px;
  font-weight: 500;
}

.stage-fields {
  display: flex;
  gap: 10px;
}

.stage-field {
  flex: 1;
}

.stage-input {
  padding: 8px 10px;
  font-size: 13px;
}

.micro-label {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #a89f91;
  font-size: 11px;
  margin-bottom: 4px;
  position: relative;
}

/* 问号帮助图标 */
.help-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: rgba(194, 163, 95, 0.2);
  color: #c2a35f;
  font-size: 11px;
  font-weight: 700;
  cursor: help;
  user-select: none;
  flex-shrink: 0;
}

.help-tooltip {
  position: absolute;
  top: 100%;
  right: 0;
  z-index: 10;
  background: rgba(20, 18, 24, 0.98);
  border: 1px solid rgba(194, 163, 95, 0.4);
  border-radius: 8px;
  padding: 10px 14px;
  color: #d4c9a8;
  font-size: 12px;
  line-height: 1.6;
  max-width: 340px;
  white-space: normal;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.5);
  pointer-events: none;
}

.help-tooltip p {
  margin: 0 0 6px;
}

.help-tooltip p:last-child {
  margin-bottom: 0;
}

.help-tooltip b {
  color: #c2a35f;
}

.test-result {
  margin-bottom: 16px;
  padding: 12px;
  border-radius: 8px;
  font-size: 13px;
}

.test-result.success {
  background: rgba(82, 196, 26, 0.12);
  border: 1px solid rgba(82, 196, 26, 0.4);
  color: #b7eb8f;
}

.test-result.fail {
  background: rgba(255, 77, 79, 0.12);
  border: 1px solid rgba(255, 77, 79, 0.4);
  color: #ff7875;
}

.test-message {
  margin: 0;
  font-weight: 500;
  word-break: break-word;
}

.test-reply {
  margin: 6px 0 0;
  color: #a89f91;
  word-break: break-all;
}

.modal-footer {
  padding: 16px 25px 24px;
  border-top: 1px solid rgba(194, 163, 95, 0.2);
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.test-btn {
  width: 100%;
  padding: 12px;
  background: rgba(82, 196, 26, 0.15);
  border: 1px solid rgba(82, 196, 26, 0.4);
  color: #b7eb8f;
  border-radius: 8px;
  cursor: pointer;
  font-size: 15px;
}

.test-btn:hover:not(:disabled) {
  background: rgba(82, 196, 26, 0.25);
}

.test-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.clear-btn {
  width: 100%;
  padding: 12px;
  background: rgba(255, 77, 79, 0.15);
  border: 1px solid rgba(255, 77, 79, 0.4);
  color: #ff7875;
  border-radius: 8px;
  cursor: pointer;
  font-size: 15px;
}

.clear-btn:hover:not(:disabled) {
  background: rgba(255, 77, 79, 0.25);
}

.save-btn {
  width: 100%;
  padding: 12px;
  background: rgba(194, 163, 95, 0.2);
  border: 1px solid rgba(194, 163, 95, 0.55);
  color: #e5d8b0;
  border-radius: 8px;
  cursor: pointer;
  font-size: 15px;
}

.save-btn:hover:not(:disabled) {
  background: rgba(194, 163, 95, 0.3);
}

.save-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.back-btn {
  width: 100%;
  padding: 12px;
  background: transparent;
  border: 1px solid rgba(194, 163, 95, 0.3);
  color: #c2a35f;
  border-radius: 8px;
  cursor: pointer;
  font-size: 15px;
}

.back-btn:hover {
  background: rgba(194, 163, 95, 0.1);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@media (max-width: 600px) {
  .custom-api-overlay {
    padding: 0;
    align-items: stretch;
    justify-content: stretch;
  }

  .custom-api-panel {
    width: 100%;
    height: 100%;
    max-height: 100vh;
    border-radius: 0;
  }

  .stage-fields {
    flex-direction: column;
    gap: 8px;
  }
}
</style>
