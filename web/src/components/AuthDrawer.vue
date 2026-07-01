<template>
  <Transition name="slide-right">
    <div v-if="show" class="modal-overlay auth-overlay">
      <div class="modal-content">
        <div class="modal-header">
          <h3>登录 FreeStyleTarot</h3>
          <p class="modal-hint">{{ headerHint }}</p>
        </div>

        <div class="modal-body">
          <div class="mode-tabs">
            <button
              type="button"
              class="mode-tab"
              :class="{ active: mode === 'password' }"
              @click="switchMode('password')"
            >
              密码登录
            </button>
            <button
              type="button"
              class="mode-tab"
              :class="{ active: mode === 'code' }"
              @click="switchMode('code')"
            >
              验证码登录
            </button>
          </div>

          <!-- 密码登录 / 注册 / 找回密码 -->
          <div v-if="mode === 'password'" class="form-group">
            <!-- 找回密码 -->
            <div v-if="showForgotPassword">
              <div v-if="forgotStep === 1">
                <label class="field-label">邮箱</label>
                <input
                  v-model="email"
                  type="email"
                  class="auth-input"
                  placeholder="your@email.com"
                  :disabled="loading"
                  @keyup.enter="handleSendForgotCode"
                />
                <span v-if="error" class="error-text">{{ error }}</span>
                <button
                  class="confirm-btn auth-btn"
                  :disabled="sendingCode || cooldown > 0"
                  @click="handleSendForgotCode"
                >
                  <template v-if="sendingCode">发送中...</template>
                  <template v-else-if="cooldown > 0">{{ cooldown }}s 后可重发</template>
                  <template v-else>获取验证码</template>
                </button>
                <button type="button" class="cancel-btn auth-back" :disabled="loading" @click="cancelForgotPassword">
                  返回登录
                </button>
              </div>
              <div v-else>
                <p class="sent-hint">验证码已发送至 {{ email }}，30 分钟内有效</p>
                <label class="field-label">验证码</label>
                <div class="code-input-row">
                  <input
                    v-model="code"
                    type="text"
                    class="auth-input code-input"
                    placeholder="6 位数字"
                    maxlength="6"
                    inputmode="numeric"
                    :disabled="loading"
                  />
                  <button
                    type="button"
                    class="resend-code-btn"
                    :disabled="sendingCode || loading || cooldown > 0"
                    @click="handleSendForgotCode"
                  >
                    <template v-if="sendingCode">发送中...</template>
                    <template v-else-if="cooldown > 0">{{ cooldown }}s</template>
                    <template v-else>重新获取</template>
                  </button>
                </div>
                <label class="field-label">新密码</label>
                <input
                  v-model="password"
                  type="password"
                  class="auth-input"
                  placeholder="6-20 位，含数字/字母/符号中至少两种"
                  :disabled="loading"
                />
                <p class="password-hint">
                  密码需 6-20 位，且包含数字、字母、符号中的至少两种
                </p>
                <label class="field-label">确认新密码</label>
                <input
                  v-model="confirmPassword"
                  type="password"
                  class="auth-input"
                  placeholder="再次输入新密码"
                  :disabled="loading"
                  @keyup.enter="handleResetPassword"
                />
                <span v-if="error" class="error-text">{{ error }}</span>
                <button class="confirm-btn auth-btn" :disabled="loading" @click="handleResetPassword">
                  {{ loading ? '重置中...' : '重置密码' }}
                </button>
                <button type="button" class="cancel-btn auth-back" :disabled="loading" @click="forgotStep = 1">
                  返回修改邮箱
                </button>
              </div>
            </div>

            <!-- 登录 / 注册 -->
            <template v-else>
            <label class="field-label">邮箱</label>
            <input
              v-model="email"
              type="email"
              class="auth-input"
              placeholder="your@email.com"
              :disabled="loading"
            />
            <label class="field-label">密码</label>
            <input
              v-model="password"
              type="password"
              class="auth-input"
              :placeholder="passwordPlaceholder"
              :disabled="loading"
              @keyup.enter="isRegister ? handleRegister() : handlePasswordLogin()"
            />
            <p v-if="isRegister" class="password-hint">
              密码需 6-20 位，且包含数字、字母、符号中的至少两种
            </p>
            <div v-if="!isRegister" class="forgot-password-row">
              <button
                type="button"
                class="forgot-password-link"
                :disabled="loading"
                @click="startForgotPassword"
              >
                忘记密码？
              </button>
            </div>
            <template v-if="isRegister">
              <label class="field-label">确认密码</label>
              <input
                v-model="confirmPassword"
                type="password"
                class="auth-input"
                placeholder="再次输入密码"
                :disabled="loading"
              />
              <label class="field-label">昵称（可选，2-32 字）</label>
              <input
                v-model="nickname"
                type="text"
                class="auth-input"
                placeholder="留空则使用默认昵称"
                maxlength="32"
                :disabled="loading"
              />
              <label class="field-label">邮箱验证码</label>
              <div class="code-input-row">
                <input
                  v-model="registerCode"
                  type="text"
                  class="auth-input code-input"
                  placeholder="6 位数字"
                  maxlength="6"
                  inputmode="numeric"
                  :disabled="loading"
                />
                <button
                  type="button"
                  class="resend-code-btn"
                  :disabled="sendingCode || loading || cooldown > 0"
                  @click="handleSendRegisterCode"
                >
                  <template v-if="sendingCode">发送中...</template>
                  <template v-else-if="cooldown > 0">{{ cooldown }}s</template>
                  <template v-else>获取验证码</template>
                </button>
              </div>
              <p class="password-hint">验证码 30 分钟内有效，请查收邮件</p>
            </template>
            <span v-if="error" class="error-text">{{ error }}</span>
            <p v-if="successMessage" class="success-text">{{ successMessage }}</p>
            <button
              class="confirm-btn auth-btn"
              :disabled="loading"
              @click="isRegister ? handleRegister() : handlePasswordLogin()"
            >
              {{ loading ? '处理中...' : (isRegister ? '注册' : '登录') }}
            </button>
            <button type="button" class="cancel-btn auth-back" :disabled="loading" @click="toggleRegister">
              {{ isRegister ? '已有账号？去登录' : '没有账号？去注册' }}
            </button>
            </template>
          </div>

          <!-- 验证码登录 -->
          <div v-else class="form-group">
            <p v-if="skipVerifyHint" class="stub-hint">{{ skipVerifyHint }}</p>
            <div v-if="codeStep === 1">
              <label class="field-label">邮箱</label>
              <input
                v-model="email"
                type="email"
                class="auth-input"
                placeholder="your@email.com"
                :disabled="loading"
                @keyup.enter="handleSendCode"
              />
              <span v-if="error" class="error-text">{{ error }}</span>
              <button
                class="confirm-btn auth-btn"
                :disabled="sendingCode || cooldown > 0"
                @click="handleSendCode"
              >
                <template v-if="sendingCode">发送中...</template>
                <template v-else-if="cooldown > 0">{{ cooldown }}s 后可重发</template>
                <template v-else>获取验证码</template>
              </button>
            </div>
            <div v-else-if="codeStep === 2">
              <p class="sent-hint">验证码已发送至 {{ email }}，30 分钟内有效</p>
              <label class="field-label">验证码</label>
              <div class="code-input-row">
                <input
                  v-model="code"
                  type="text"
                  class="auth-input code-input"
                  placeholder="6 位数字"
                  maxlength="6"
                  inputmode="numeric"
                  :disabled="loading"
                  @keyup.enter="handleVerifyCode"
                />
                <button
                  type="button"
                  class="resend-code-btn"
                  :disabled="sendingCode || loading || cooldown > 0"
                  @click="handleSendCode"
                >
                  <template v-if="sendingCode">发送中...</template>
                  <template v-else-if="cooldown > 0">{{ cooldown }}s</template>
                  <template v-else>重新获取</template>
                </button>
              </div>
              <span v-if="error" class="error-text">{{ error }}</span>
              <button class="confirm-btn auth-btn" :disabled="loading" @click="handleVerifyCode">
                {{ loading ? '验证中...' : '验证并继续' }}
              </button>
              <button class="cancel-btn auth-back" :disabled="loading" @click="codeStep = 1">
                返回修改邮箱
              </button>
            </div>
            <div v-else>
              <p class="sent-hint">验证码已通过，请完善账号信息</p>
              <label class="field-label">昵称（2-32 字）</label>
              <input
                v-model="nickname"
                type="text"
                class="auth-input"
                placeholder="请输入昵称"
                maxlength="32"
                :disabled="loading"
              />
              <label class="field-label">密码</label>
              <input
                v-model="password"
                type="password"
                class="auth-input"
                placeholder="6-20 位，含数字/字母/符号中至少两种"
                :disabled="loading"
              />
              <p class="password-hint">
                密码需 6-20 位，且包含数字、字母、符号中的至少两种
              </p>
              <label class="field-label">确认密码</label>
              <input
                v-model="confirmPassword"
                type="password"
                class="auth-input"
                placeholder="再次输入密码"
                :disabled="loading"
                @keyup.enter="handleCompleteCodeSignup"
              />
              <span v-if="error" class="error-text">{{ error }}</span>
              <button class="confirm-btn auth-btn" :disabled="loading" @click="handleCompleteCodeSignup">
                {{ loading ? '注册中...' : '完成注册' }}
              </button>
              <button class="cancel-btn auth-back" :disabled="loading" @click="backToInitialLogin">
                返回登录
              </button>
            </div>
          </div>
        </div>

        <div class="modal-footer auth-footer">
          <button type="button" class="cancel-btn auth-close-btn" @click="handleFooterBack">返回</button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, computed, onUnmounted, watch } from 'vue';
import { useAuth } from '../composables/useAuth.js';
import {
  validateEmail,
  validateRegisterPassword,
  validateNickname,
  validateRequiredNickname,
} from '../utils/authValidation.js';

const props = defineProps({
  allowDismiss: { type: Boolean, default: true },
});

const show = defineModel('show', { type: Boolean, default: false });

const {
  login,
  register,
  resetPassword,
  requestCode,
  verifyEmailCode,
  finishCodeSignup,
  applyAuthResult,
  skipVerify,
  isAuthenticated,
} = useAuth();

const mode = ref('password');
const isRegister = ref(false);
const codeStep = ref(1);
const forgotStep = ref(1);
const email = ref('');
const password = ref('');
const confirmPassword = ref('');
const code = ref('');
const registerCode = ref('');
const nickname = ref('');
const error = ref('');
const successMessage = ref('');
const showForgotPassword = ref(false);
const loading = ref(false);
const sendingCode = ref(false);
const cooldown = ref(0);
let cooldownTimer = null;

const headerHint = computed(() => {
  if (mode.value === 'password') {
    if (showForgotPassword.value) {
      return forgotStep.value === 1
        ? '输入注册邮箱，获取验证码重置密码'
        : '设置新密码完成重置';
    }
    return isRegister.value
      ? '填写邮箱与密码，并完成邮箱验证后注册'
      : '使用邮箱与密码登录或注册';
  }
  if (codeStep.value === 3) {
    return '验证码已通过，请设置昵称与密码';
  }
  return skipVerify.value
    ? '验证码登录（当前 skip_verify=true，验证已架空）'
    : '使用邮箱验证码登录，新用户验证后将设置密码';
});

const skipVerifyHint = computed(() =>
  mode.value === 'code'
);

const passwordPlaceholder = computed(() =>
  isRegister.value ? '6-20 位，含数字/字母/符号中至少两种' : '请输入密码'
);

function resetCodeFlow() {
  codeStep.value = 1;
  code.value = '';
  registerCode.value = '';
  nickname.value = '';
  password.value = '';
  confirmPassword.value = '';
}

function resetForgotFlow() {
  forgotStep.value = 1;
  code.value = '';
  password.value = '';
  confirmPassword.value = '';
}

function startForgotPassword() {
  showForgotPassword.value = true;
  error.value = '';
  resetForgotFlow();
}

function cancelForgotPassword() {
  showForgotPassword.value = false;
  error.value = '';
  resetForgotFlow();
}

function backToInitialLogin() {
  resetCodeFlow();
  resetForgotFlow();
  mode.value = 'password';
  isRegister.value = false;
  showForgotPassword.value = false;
  error.value = '';
}

function handleFooterBack() {
  if (mode.value === 'code' && codeStep.value === 3) {
    backToInitialLogin();
    return;
  }
  if (mode.value === 'password' && showForgotPassword.value) {
    if (forgotStep.value === 2) {
      forgotStep.value = 1;
      error.value = '';
      return;
    }
    cancelForgotPassword();
    return;
  }
  if (props.allowDismiss) {
    show.value = false;
  }
}

function switchMode(next) {
  mode.value = next;
  error.value = '';
  successMessage.value = '';
  showForgotPassword.value = false;
  resetForgotFlow();
  resetCodeFlow();
}

function toggleRegister() {
  isRegister.value = !isRegister.value;
  showForgotPassword.value = false;
  registerCode.value = '';
  error.value = '';
  successMessage.value = '';
}

function startCooldown(seconds = 60) {
  cooldown.value = seconds;
  clearInterval(cooldownTimer);
  cooldownTimer = setInterval(() => {
    cooldown.value -= 1;
    if (cooldown.value <= 0) {
      clearInterval(cooldownTimer);
      cooldownTimer = null;
    }
  }, 1000);
}

watch(isAuthenticated, (authed) => {
  if (authed && show.value) {
    show.value = false;
  }
});

async function handlePasswordLogin() {
  error.value = '';
  successMessage.value = '';
  const emailResult = validateEmail(email.value);
  if (!emailResult.ok) {
    error.value = emailResult.message;
    return;
  }
  if (!password.value) {
    error.value = '请输入密码';
    return;
  }

  loading.value = true;
  try {
    await login(emailResult.value, password.value);
  } catch (err) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}

async function sendCodeWithCooldown(emailAddr) {
  await requestCode(emailAddr);
  email.value = emailAddr;
  startCooldown(60);
}

async function handleSendRegisterCode() {
  if (cooldown.value > 0) return;

  error.value = '';
  const emailResult = validateEmail(email.value);
  if (!emailResult.ok) {
    error.value = emailResult.message;
    return;
  }

  sendingCode.value = true;
  try {
    await sendCodeWithCooldown(emailResult.value);
    registerCode.value = '';
  } catch (err) {
    error.value = err.message;
    if (err.status === 429) {
      const match = err.message.match(/(\d+)\s*秒/);
      if (match) startCooldown(parseInt(match[1], 10));
    }
  } finally {
    sendingCode.value = false;
  }
}

async function handleSendForgotCode() {
  if (cooldown.value > 0) return;

  error.value = '';
  const emailResult = validateEmail(email.value);
  if (!emailResult.ok) {
    error.value = emailResult.message;
    return;
  }

  sendingCode.value = true;
  try {
    await sendCodeWithCooldown(emailResult.value);
    forgotStep.value = 2;
    code.value = '';
  } catch (err) {
    error.value = err.message;
    if (err.status === 429) {
      const match = err.message.match(/(\d+)\s*秒/);
      if (match) startCooldown(parseInt(match[1], 10));
    }
  } finally {
    sendingCode.value = false;
  }
}

async function handleRegister() {
  error.value = '';
  const emailResult = validateEmail(email.value);
  if (!emailResult.ok) {
    error.value = emailResult.message;
    return;
  }
  if (!password.value) {
    error.value = '请输入密码';
    return;
  }

  const passwordResult = validateRegisterPassword(password.value);
  if (!passwordResult.ok) {
    error.value = passwordResult.message;
    return;
  }
  if (password.value !== confirmPassword.value) {
    error.value = '两次输入的密码不一致';
    return;
  }

  const nicknameResult = validateNickname(nickname.value);
  if (!nicknameResult.ok) {
    error.value = nicknameResult.message;
    return;
  }

  const trimmedCode = registerCode.value.trim() || '000000';
  if (!skipVerify.value && trimmedCode.length !== 6) {
    error.value = '请输入 6 位邮箱验证码';
    return;
  }

  loading.value = true;
  try {
    await register(emailResult.value, password.value, nickname.value.trim(), trimmedCode);
  } catch (err) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}

async function handleResetPassword() {
  error.value = '';
  const emailResult = validateEmail(email.value);
  if (!emailResult.ok) {
    error.value = emailResult.message;
    return;
  }

  const trimmedCode = code.value.trim() || '000000';
  if (!skipVerify.value && trimmedCode.length !== 6) {
    error.value = '请输入 6 位验证码';
    return;
  }
  if (!password.value) {
    error.value = '请输入新密码';
    return;
  }

  const passwordResult = validateRegisterPassword(password.value);
  if (!passwordResult.ok) {
    error.value = passwordResult.message;
    return;
  }
  if (password.value !== confirmPassword.value) {
    error.value = '两次输入的密码不一致';
    return;
  }

  loading.value = true;
  try {
    await resetPassword(emailResult.value, trimmedCode, password.value, confirmPassword.value);
    cancelForgotPassword();
    password.value = '';
    error.value = '';
    successMessage.value = '密码已重置，请使用新密码登录';
  } catch (err) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}

async function handleSendCode() {
  if (cooldown.value > 0) return;

  error.value = '';
  const emailResult = validateEmail(email.value);
  if (!emailResult.ok) {
    error.value = emailResult.message;
    return;
  }

  sendingCode.value = true;
  try {
    await sendCodeWithCooldown(emailResult.value);
    codeStep.value = 2;
    code.value = '';
  } catch (err) {
    error.value = err.message;
    if (err.status === 429) {
      const match = err.message.match(/(\d+)\s*秒/);
      if (match) startCooldown(parseInt(match[1], 10));
    }
  } finally {
    sendingCode.value = false;
  }
}

async function handleVerifyCode() {
  error.value = '';
  const emailResult = validateEmail(email.value);
  if (!emailResult.ok) {
    error.value = emailResult.message;
    return;
  }
  const trimmedCode = code.value.trim() || '000000';
  if (!skipVerify.value && trimmedCode.length !== 6) {
    error.value = '请输入 6 位验证码';
    return;
  }

  loading.value = true;
  try {
    const data = await verifyEmailCode(emailResult.value, trimmedCode);
    if (data.needs_setup) {
      nickname.value = '';
      password.value = '';
      confirmPassword.value = '';
      codeStep.value = 3;
      return;
    }
    await applyAuthResult(data);
  } catch (err) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}

async function handleCompleteCodeSignup() {
  error.value = '';
  const emailResult = validateEmail(email.value);
  if (!emailResult.ok) {
    error.value = emailResult.message;
    return;
  }

  const nicknameResult = validateRequiredNickname(nickname.value);
  if (!nicknameResult.ok) {
    error.value = nicknameResult.message;
    return;
  }
  if (!password.value) {
    error.value = '请输入密码';
    return;
  }

  const passwordResult = validateRegisterPassword(password.value);
  if (!passwordResult.ok) {
    error.value = passwordResult.message;
    return;
  }
  if (password.value !== confirmPassword.value) {
    error.value = '两次输入的密码不一致';
    return;
  }

  loading.value = true;
  try {
    await finishCodeSignup(
      emailResult.value,
      nicknameResult.value,
      password.value,
      confirmPassword.value
    );
  } catch (err) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}

onUnmounted(() => {
  clearInterval(cooldownTimer);
});
</script>

<style scoped>
.auth-overlay {
  pointer-events: auto;
  background: rgba(0, 0, 0, 0.5);
}

.mode-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
}

.mode-tab {
  flex: 1;
  padding: 10px;
  border-radius: 8px;
  border: 1px solid rgba(194, 163, 95, 0.3);
  background: transparent;
  color: #a89f91;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.mode-tab.active {
  background: rgba(194, 163, 95, 0.15);
  border-color: #c2a35f;
  color: #e5d8b0;
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
  padding: 12px;
  border-radius: 8px;
  border: 1px solid rgba(194, 163, 95, 0.3);
  background: rgba(0, 0, 0, 0.3);
  color: #e5d8b0;
  font-size: 14px;
  margin-bottom: 16px;
}

.auth-input:focus {
  outline: none;
  border-color: #c2a35f;
}

.code-input-row {
  display: flex;
  gap: 10px;
  align-items: stretch;
  margin-bottom: 16px;
}

.code-input-row .code-input {
  flex: 1;
  min-width: 0;
  margin-bottom: 0;
}

.resend-code-btn {
  flex-shrink: 0;
  min-width: 88px;
  padding: 0 12px;
  border-radius: 8px;
  border: 1px solid rgba(194, 163, 95, 0.45);
  background: rgba(194, 163, 95, 0.12);
  color: #e5d8b0;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.resend-code-btn:hover:not(:disabled) {
  background: rgba(194, 163, 95, 0.22);
  border-color: #c2a35f;
}

.resend-code-btn:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

.auth-btn {
  width: 100%;
  margin-top: 8px;
}

.auth-back {
  width: 100%;
  margin-top: 12px;
  padding: 10px;
  background: transparent;
  border: 1px solid rgba(194, 163, 95, 0.3);
  color: #c2a35f;
  border-radius: 8px;
  cursor: pointer;
}

.sent-hint,
.stub-hint {
  color: #a89f91;
  font-size: 13px;
  margin: 0 0 16px;
}

.stub-hint {
  padding: 10px 12px;
  background: rgba(194, 163, 95, 0.08);
  border-radius: 8px;
  border: 1px dashed rgba(194, 163, 95, 0.3);
}

.error-text {
  display: block;
  color: #ff4d4f;
  font-size: 13px;
  margin: -8px 0 12px;
  line-height: 1.5;
}

.success-text {
  display: block;
  color: #52c41a;
  font-size: 13px;
  margin: -8px 0 12px;
  line-height: 1.5;
}

.forgot-password-row {
  display: flex;
  justify-content: flex-end;
  margin: -8px 0 16px;
}

.forgot-password-link {
  padding: 0;
  border: none;
  background: none;
  color: #c2a35f;
  font-size: 13px;
  cursor: pointer;
  text-decoration: underline;
  text-underline-offset: 3px;
}

.forgot-password-link:hover:not(:disabled) {
  color: #e5d8b0;
}

.forgot-password-link:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.password-hint {
  margin: -8px 0 16px;
  color: #a89f91;
  font-size: 12px;
  line-height: 1.5;
}

.auth-footer {
  display: none;
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

  .auth-footer {
    display: block;
    padding: 16px 15px 24px;
    border-top: 1px solid rgba(194, 163, 95, 0.2);
  }

  .auth-close-btn {
    width: 100%;
    padding: 12px;
    background: transparent;
    border: 1px solid rgba(194, 163, 95, 0.3);
    color: #c2a35f;
    border-radius: 8px;
    cursor: pointer;
    font-size: 15px;
  }
}
</style>
