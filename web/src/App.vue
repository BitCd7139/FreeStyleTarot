<template>
  <div v-if="authLoading" class="auth-loading">加载中...</div>
  <div v-else-if="forceLogin && !isAuthenticated" class="app-backdrop">
    <AuthDrawer v-model:show="forceAuthOpen" :allow-dismiss="false" />
  </div>
  <template v-else>
    <TarotMain />
    <AnnouncementModal />
  </template>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue';
import TarotMain from './components/TarotMain.vue';
import AuthDrawer from './components/AuthDrawer.vue';
import AnnouncementModal from './components/AnnouncementModal.vue';
import { useAuth } from './composables/useAuth.js';
import { useAnnouncement } from './composables/useAnnouncement.js';

const { isAuthenticated, forceLogin, initAuth, user } = useAuth();
const { checkAndShow, checkAndShowForGuest } = useAnnouncement();
const authLoading = ref(true);
const forceAuthOpen = ref(true);

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
