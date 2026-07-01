<template>
  <Transition name="modal-fade">
    <div v-if="show" class="modal-overlay" @mousedown="onOverlayMouseDown" @mouseup="onOverlayMouseUp">
      <div class="modal-content" @mousedown.stop>
        <div class="modal-header">
          <h3>选择你的体验模式</h3>
          <p class="modal-hint">选择后双击卡片即可确认</p>
        </div>

        <div class="modal-body">
          <div class="mode-cards">
            <button
              type="button"
              class="mode-card"
              :class="{ active: picked === 'beginner' }"
              @click="handleTap('beginner')"
            >
              <div class="mode-icon">🌱</div>
              <div class="mode-title">新手模式</div>
              <div class="mode-desc">选择牌阵 → 引导抽牌 → 自动识别阵位，适合刚接触塔罗的你。</div>
              <span v-if="picked === 'beginner' && saving" class="double-tap-badge saving">确认中…</span>
            </button>
            <button
              type="button"
              class="mode-card"
              :class="{ active: picked === 'advanced' }"
              @click="handleTap('advanced')"
            >
              <div class="mode-icon">🔮</div>
              <div class="mode-title">高级模式</div>
              <div class="mode-desc">自由摆放、自定义牌阵，完整体验所有进阶玩法。</div>
              <span v-if="picked === 'advanced' && saving" class="double-tap-badge saving">确认中…</span>
            </button>
          </div>
          <span v-if="error" class="error-text">{{ error }}</span>
        </div>

        <!-- 底部提示：可点击遮罩关闭 -->
        <p v-if="cancellable" class="dismiss-hint">或点击外部区域关闭</p>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, watch } from 'vue';
import { useAuth } from '../composables/useAuth.js';

const props = defineProps({
  cancellable: { type: Boolean, default: true },
});
const show = defineModel('show', { type: Boolean, default: false });
const emit = defineEmits(['confirmed']);

const { user, updateExperienceMode } = useAuth();
const picked = ref('beginner');
const saving = ref(false);
const error = ref('');

watch(show, (visible) => {
  if (visible) {
    picked.value = user.value?.experience_mode === 'advanced' ? 'advanced' : 'beginner';
    error.value = '';
  }
});

function handleDoubleClick(mode) {
  if (saving.value) return;
  picked.value = mode;
  confirm();
}

// 移动端 dblclick 不可靠，手动检测同一卡片上的快速双击（兼容桌面与移动）
let lastTapTime = 0;
let lastTapMode = '';
const DOUBLE_TAP_INTERVAL = 300;

function handleTap(mode) {
  picked.value = mode;
  const now = Date.now();
  if (lastTapMode === mode && now - lastTapTime < DOUBLE_TAP_INTERVAL) {
    lastTapTime = 0;
    lastTapMode = '';
    handleDoubleClick(mode);
  } else {
    lastTapTime = now;
    lastTapMode = mode;
  }
}

async function confirm() {
  if (!picked.value || saving.value) return;
  saving.value = true;
  error.value = '';
  try {
    await updateExperienceMode(picked.value);
    emit('confirmed', picked.value);
    show.value = false;
  } catch (err) {
    error.value = err.message || '保存失败';
  } finally {
    saving.value = false;
  }
}

let mouseDownOutside = false;
function onOverlayMouseDown() {
  mouseDownOutside = true;
}
function onOverlayMouseUp() {
  if (mouseDownOutside && props.cancellable) show.value = false;
  mouseDownOutside = false;
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 16px;
  overflow: hidden;
}

.modal-content {
  height: fit-content;
  width: 100%;
  max-width: 480px;
  max-height: calc(100vh - 32px);
  background: linear-gradient(180deg, #2a2438 0%, #1a1626 100%);
  border: 1px solid rgba(194, 163, 95, 0.4);
  border-radius: 16px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.6);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal-header {
  padding: 18px 22px 8px;
  text-align: center;
  flex-shrink: 0;
}

.modal-header h3 {
  margin: 0;
  color: #e5d8b0;
  font-size: 19px;
}

.modal-hint {
  margin: 4px 0 0;
  color: #a89f91;
  font-size: 12px;
}

.modal-body {
  padding: 10px 22px 6px;
  flex-shrink: 1;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
}

.mode-cards {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.mode-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 16px 12px 12px;
  border-radius: 12px;
  border: 1px solid rgba(194, 163, 95, 0.25);
  background: rgba(0, 0, 0, 0.25);
  color: #e5d8b0;
  cursor: pointer;
  text-align: center;
  transition: all 0.2s ease;
  position: relative;
  user-select: none;
  -webkit-user-select: none;
  -webkit-tap-highlight-color: transparent;
  touch-action: manipulation;
}

.mode-card:hover {
  border-color: rgba(194, 163, 95, 0.6);
  background: rgba(194, 163, 95, 0.08);
}

.mode-card.active {
  border-color: #c2a35f;
  background: rgba(194, 163, 95, 0.15);
  box-shadow: 0 0 0 2px rgba(194, 163, 95, 0.4) inset;
}

.mode-icon {
  font-size: 30px;
  line-height: 1;
}

.mode-title {
  font-size: 16px;
  font-weight: 600;
  color: #e5d8b0;
}

.mode-desc {
  font-size: 12px;
  color: #a89f91;
  line-height: 1.45;
}

.double-tap-badge {
  display: inline-block;
  margin-top: 2px;
  padding: 3px 10px;
  border-radius: 20px;
  background: rgba(194, 163, 95, 0.2);
  color: #c2a35f;
  font-size: 11px;
  font-weight: 500;
  letter-spacing: 0.5px;
  animation: pulse-badge 2s ease-in-out infinite;
  border: 1px solid rgba(194, 163, 95, 0.35);
}

.double-tap-badge.saving {
  animation: none;
  opacity: 0.7;
  background: rgba(194, 163, 95, 0.1);
  color: #a89f91;
  border-color: rgba(194, 163, 95, 0.2);
}

@keyframes pulse-badge {
  0%, 100% {
    opacity: 0.75;
    transform: scale(1);
  }
  50% {
    opacity: 1;
    transform: scale(1.04);
  }
}

.error-text {
  display: block;
  color: #ff4d4f;
  font-size: 12px;
  margin-top: 10px;
  text-align: center;
}

.dismiss-hint {
  text-align: center;
  color: #6b6378;
  font-size: 11px;
  margin: 0;
  padding: 0 22px 14px;
  flex-shrink: 0;
}

@media (max-width: 600px) {
  .modal-overlay {
    align-items: center;
    justify-content: center;
    padding: 12px;
  }
  .modal-content {
    height: fit-content;
    max-width: 100%;
    max-height: calc(100vh - 24px);
  }
  .mode-cards {
    grid-template-columns: 1fr;
    gap: 10px;
  }
  .mode-card {
    padding: 14px 14px 12px;
    gap: 5px;
  }
  .mode-icon {
    font-size: 26px;
  }
  .mode-title {
    font-size: 15px;
  }
  .mode-desc {
    font-size: 11px;
  }
  .modal-body {
    padding: 16px 16px 10px;
  }
  .dismiss-hint {
    padding: 0 16px 12px;
  }
  .double-tap-badge {
    font-size: 10px;
    padding: 3px 8px;
  }
}
</style>