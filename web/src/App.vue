<template>
  <div v-if="authLoading" class="auth-loading">加载中...</div>
  <div v-else-if="forceLogin && !isAuthenticated" class="app-backdrop">
    <AuthDrawer v-model:show="forceAuthOpen" :allow-dismiss="false" />
  </div>
  <template v-else>
    <TarotBeginner v-if="showBeginner" />
    <TarotMain v-else />
    <AnnouncementModal />
  </template>
  <ExperiencePreferenceModal
    v-model:show="showExperiencePref"
    :cancellable="false"
  />
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue';
import TarotMain from './components/TarotMain.vue';
import TarotBeginner from './components/TarotBeginner.vue';
import AuthDrawer from './components/AuthDrawer.vue';
import AnnouncementModal from './components/AnnouncementModal.vue';
import ExperiencePreferenceModal from './components/ExperiencePreferenceModal.vue';
import { useAuth } from './composables/useAuth.js';
import { useAnnouncement } from './composables/useAnnouncement.js';
import { useTarotFlow } from './composables/useTarotFlow.js';

const { isAuthenticated, forceLogin, initAuth, user } = useAuth();
const { checkAndShow, checkAndShowForGuest } = useAnnouncement();
const { isBeginner, selectedSpread } = useTarotFlow();
const authLoading = ref(true);
const forceAuthOpen = ref(true);

// 新手用户未选阵时进入选阵页；游客与高级用户、已选阵的新手进入 TarotMain
const showBeginner = computed(() => isBeginner.value && !selectedSpread.value);

// 登录用户若未设置体验偏好，弹窗引导选择
const showExperiencePref = ref(false);
// 跟踪认证状态：仅在「未登录 → 已登录」的新鲜登录转换时弹窗，缓存会话恢复不弹
let wasAuthenticated = isAuthenticated.value;
watch(
  [isAuthenticated, user],
  () => {
    const isAuthed = isAuthenticated.value;
    if (!wasAuthenticated && isAuthed && user.value && !user.value.experience_mode) {
      showExperiencePref.value = true;
    }
    wasAuthenticated = isAuthed;
  }
);

function maybeShowAnnouncement() {
  if (authLoading.value) return;
  if (isAuthenticated.value && user.value) {
    checkAndShow(user.value);
  } else if (!forceLogin.value) {
    checkAndShowForGuest();
  }
}

watch(
  [isAuthenticated, user, authLoading],
  () => {
    maybeShowAnnouncement();
  }
);

onMounted(async () => {
  try {
    await initAuth();
  } catch {
    // 配置或 token 校验失败时仍允许进入
  } finally {
    authLoading.value = false;
    // 缓存会话恢复不算新鲜登录，同步状态避免 watch 误触发偏好弹窗
    wasAuthenticated = isAuthenticated.value;
  }
});
</script>

<style>
.app-backdrop {
  width: 100vw;
  height: 100vh;
  background: #433843;
}

.auth-loading {
  width: 100vw;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #433843;
  color: #a89f91;
  font-size: 16px;
}
</style>
