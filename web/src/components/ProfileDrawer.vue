<template>
  <Transition name="slide-right">
    <div v-if="show" class="modal-overlay" @mousedown="onOverlayMouseDown" @mouseup="onOverlayMouseUp">
      <div class="modal-content" @mousedown.stop>
        <div class="modal-header">
          <button type="button" class="header-back-btn" @click="show = false" aria-label="返回">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M15 18l-6-6 6-6"/>
            </svg>
          </button>
          <div class="header-title-group">
            <h3>个人资料</h3>
            <p class="modal-hint">管理你的账号信息</p>
          </div>
        </div>

        <div class="modal-body">
          <div class="profile-field">
            <span class="field-label">邮箱</span>
            <span class="field-value">{{ user?.email || '—' }}</span>
          </div>

          <div class="profile-field">
            <span class="field-label">昵称</span>
            <div class="nickname-row">
              <input
                v-model="editNickname"
                type="text"
                class="auth-input"
                maxlength="32"
                :disabled="saving"
              />
              <button class="btn-mini btn-yes nickname-action-btn" :disabled="saving" @click="saveNickname">
                {{ saving ? '...' : '保存' }}
              </button>
            </div>
            <span v-if="nickError" class="error-text">{{ nickError }}</span>
          </div>

          <div class="profile-field">
            <span class="field-label">体验偏好</span>
            <div class="nickname-row">
              <span class="field-value" style="flex:1">{{ experienceLabel }}</span>
              <button type="button" class="btn-mini btn-yes nickname-action-btn" @click="showExperienceModal = true">切换</button>
            </div>
          </div>

          <div class="profile-field">
            <div class="status-balance-row">
              <div class="status-balance-col">
                <span class="field-label">会员状态</span>
                <span class="field-value">{{ vipLabel }}</span>
              </div>
              <div v-if="!customApiEnabled" class="status-balance-col">
                <span class="field-label">余额</span>
                <span class="field-value">{{ balanceLabel }}</span>
              </div>
            </div>
          </div>

          <div v-if="!customApiEnabled" class="profile-field">
            <span class="field-label">提问配额</span>
            <span class="field-value">{{ quotaLabel }}</span>
          </div>

          <ProfilePredictHistory />
        </div>

        <div class="modal-footer profile-footer">
          <button type="button" class="custom-api-btn" @click="handleOpenCustomApi">
            自定义 API
            <span v-if="customApiEnabled" class="api-badge">已启用</span>
          </button>
          <button type="button" class="announcement-btn" @click="handleViewAnnouncement">查看公告</button>
          <button class="logout-btn" @click="handleLogout">退出登录</button>
        </div>
      </div>
    </div>
  </Transition>

  <CustomApiModal v-model:show="showCustomApi" />
<ExperiencePreferenceModal v-model:show="showExperienceModal" :cancellable="true" />
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { useAuth } from '../composables/useAuth.js';
import { useAnnouncement } from '../composables/useAnnouncement.js';
import { usePredictHistory } from '../composables/usePredictHistory.js';
import { useCustomApi } from '../composables/useCustomApi.js';
import ProfilePredictHistory from './ProfilePredictHistory.vue';
import CustomApiModal from './CustomApiModal.vue';
import ExperiencePreferenceModal from './ExperiencePreferenceModal.vue';

const show = defineModel('show', { type: Boolean, default: false });

const { user, canPredict, nextPredictAt, updateNickname, logout } = useAuth();
const { openManually } = useAnnouncement();
const { resetHistory } = usePredictHistory();
const { isEnabled: customApiEnabled } = useCustomApi();

const editNickname = ref('');
const saving = ref(false);
const nickError = ref('');
const showCustomApi = ref(false);
const showExperienceModal = ref(false);

watch(
  () => show.value,
  (visible) => {
    if (visible && user.value) {
      editNickname.value = user.value.nickname || '';
      nickError.value = '';
    }
  }
);

watch(user, (u) => {
  if (u) editNickname.value = u.nickname || '';
});

const vipLabel = computed(() => {
  if (!user.value) return '—';
  if (customApiEnabled.value) return 'API 用户';
  if (user.value.tier === 'vip') {
    const exp = user.value.vip_expires_at;
    if (exp) {
      return `VIP（至 ${new Date(exp).toLocaleDateString()}）`;
    }
    return 'VIP';
  }
  return '免费用户';
});

const balanceLabel = computed(() => user.value?.balance || '0.00 CNY');

const experienceLabel = computed(() => {
  const mode = user.value?.experience_mode;
  if (mode === 'beginner') return '新手模式';
  if (mode === 'advanced') return '高级模式';
  return '未设置';
});

const quotaLabel = computed(() => {
  if (canPredict.value) return '可以提问';
  if (nextPredictAt.value) {
    const mins = Math.ceil((new Date(nextPredictAt.value) - Date.now()) / 60000);
    return `冷却中，约 ${mins} 分钟后可再次提问`;
  }
  return '冷却中';
});

async function saveNickname() {
  nickError.value = '';
  const trimmed = editNickname.value.trim();
  if (trimmed.length < 2 || trimmed.length > 32) {
    nickError.value = '昵称长度需在 2-32 字符之间';
    return;
  }
  saving.value = true;
  try {
    await updateNickname(trimmed);
  } catch (err) {
    nickError.value = err.message;
  } finally {
    saving.value = false;
  }
}

function handleViewAnnouncement() {
  openManually();
}

function handleOpenCustomApi() {
  showCustomApi.value = true;
}

function handleLogout() {
  logout();
  resetHistory();
  show.value = false;
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
.modal-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 18px 20px 12px;
  border-bottom: 1px solid rgba(194, 163, 95, 0.2);
}

.header-back-btn {
  flex-shrink: 0;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: 1px solid rgba(194, 163, 95, 0.3);
  border-radius: 50%;
  color: #c2a35f;
  cursor: pointer;
  padding: 0;
  transition: background 0.2s, border-color 0.2s;
}

.header-back-btn:hover {
  background: rgba(194, 163, 95, 0.12);
  border-color: rgba(194, 163, 95, 0.6);
}

.header-back-btn svg {
  width: 18px;
  height: 18px;
}

.header-title-group {
  margin-left: auto;
  text-align: right;
}

.modal-header h3 {
  margin: 0;
  color: #e5d8b0;
  font-size: 18px;
}

.modal-hint {
  margin: 4px 0 0;
  color: #a89f91;
  font-size: 12px;
}

.profile-field {
  margin-bottom: 20px;
}

.field-label {
  display: block;
  color: #a89f91;
  font-size: 13px;
  margin-bottom: 6px;
}

.field-value {
  color: #e5d8b0;
  font-size: 15px;
}

.status-balance-row {
  display: flex;
  gap: 24px;
}

.status-balance-col {
  flex: 1;
  min-width: 0;
}

.nickname-row {
  display: flex;
  gap: 8px;
  align-items: center;
}

.nickname-action-btn {
  min-width: 64px;
}

.auth-input {
  flex: 1;
  box-sizing: border-box;
  padding: 10px 12px;
  border-radius: 8px;
  border: 1px solid rgba(194, 163, 95, 0.3);
  background: rgba(0, 0, 0, 0.3);
  color: #e5d8b0;
  font-size: 14px;
}

.announcement-btn {
  width: 100%;
  padding: 12px;
  background: rgba(194, 163, 95, 0.12);
  border: 1px solid rgba(194, 163, 95, 0.4);
  color: #e5d8b0;
  border-radius: 8px;
  cursor: pointer;
  font-size: 15px;
}

.announcement-btn:hover {
  background: rgba(194, 163, 95, 0.2);
}

.custom-api-btn {
  width: 100%;
  padding: 12px;
  background: rgba(82, 196, 26, 0.12);
  border: 1px solid rgba(82, 196, 26, 0.4);
  color: #b7eb8f;
  border-radius: 8px;
  cursor: pointer;
  font-size: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.custom-api-btn:hover {
  background: rgba(82, 196, 26, 0.2);
}

.api-badge {
  display: inline-block;
  padding: 2px 8px;
  background: rgba(82, 196, 26, 0.25);
  border: 1px solid rgba(82, 196, 26, 0.5);
  border-radius: 10px;
  font-size: 11px;
  font-weight: 500;
}

.logout-btn {
  width: 100%;
  margin-top: 0;
  padding: 12px;
  background: rgba(255, 77, 79, 0.15);
  border: 1px solid rgba(255, 77, 79, 0.4);
  color: #ff7875;
  border-radius: 8px;
  cursor: pointer;
  font-size: 15px;
}

.logout-btn:hover {
  background: rgba(255, 77, 79, 0.25);
}

.profile-footer {
  padding: 16px 25px 24px;
  border-top: 1px solid rgba(194, 163, 95, 0.2);
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.error-text {
  display: block;
  color: #ff4d4f;
  font-size: 13px;
  margin-top: 4px;
}

@media (max-width: 600px) {
  .modal-overlay {
    padding: 0;
    align-items: stretch;
    justify-content: stretch;
  }

  .modal-content {
    width: 100%;
    min-width: 0;
    height: 100%;
    max-height: 100vh;
    border-radius: 0;
    border-left: none;
    border-top: none;
  }

  .modal-body {
    padding: 15px;
  }

  .profile-footer {
    padding: 16px 15px 24px;
  }
}
</style>
