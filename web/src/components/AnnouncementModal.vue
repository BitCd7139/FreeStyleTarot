<template>
  <Transition name="fade">
    <div v-if="show" class="announcement-overlay">
      <div class="announcement-panel">
        <div class="announcement-header">
          <div class="announcement-badge">公告</div>
          <h3>{{ title }}</h3>
          <p class="announcement-hint">{{ hintText }}</p>
        </div>

        <div class="announcement-body">
          <div class="markdown-body announcement-content" v-html="parsedContent"></div>
        </div>

        <div class="announcement-footer">
          <button type="button" class="confirm-btn" @click="handleDismiss">我知道了</button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { computed } from 'vue';
import { marked } from 'marked';
import DOMPurify from 'dompurify';
import { useAnnouncement } from '../composables/useAnnouncement.js';
import { useAuth } from '../composables/useAuth.js';

const { show, title, content, isGuestView, isManualView, dismiss, dismissGuest, dismissManual } = useAnnouncement();
const { user } = useAuth();

const parsedContent = computed(() => {
  if (!content.value) return '';
  return DOMPurify.sanitize(marked.parse(content.value));
});

const hintText = computed(() => {
  if (isManualView.value) return '手动查看公告，关闭后不影响自动弹出记录';
  if (isGuestView.value) return '游客模式下每次进入都会展示公告';
  return '服务器重启后首次登录可见，关闭后本次运行期间不再弹出';
});

function handleDismiss() {
  if (isManualView.value) {
    dismissManual();
  } else if (isGuestView.value) {
    dismissGuest();
  } else {
    dismiss(user.value);
  }
}
</script>

<style scoped>
.announcement-overlay {
  position: fixed;
  inset: 0;
  z-index: 2000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  background: rgba(0, 0, 0, 0.72);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
}

.announcement-panel {
  width: min(520px, 100%);
  max-height: min(80vh, 714px);
  display: flex;
  flex-direction: column;
  background: rgba(30, 27, 36, 0.96);
  border: 1px solid rgba(194, 163, 95, 0.45);
  border-radius: 16px;
  box-shadow: 0 24px 64px rgba(0, 0, 0, 0.55);
  overflow: hidden;
}

.announcement-header {
  padding: 24px 24px 12px;
  border-bottom: 1px solid rgba(194, 163, 95, 0.2);
}

.announcement-badge {
  display: inline-block;
  margin-bottom: 10px;
  padding: 4px 10px;
  border-radius: 999px;
  background: rgba(194, 163, 95, 0.18);
  border: 1px solid rgba(194, 163, 95, 0.35);
  color: #e5d8b0;
  font-size: 12px;
  letter-spacing: 0.08em;
}

.announcement-header h3 {
  margin: 0 0 8px;
  color: #e5d8b0;
  font-size: 20px;
}

.announcement-hint {
  margin: 0;
  color: #a89f91;
  font-size: 12px;
  line-height: 1.5;
}

.announcement-body {
  flex: 1;
  overflow-y: auto;
  padding: 20px 24px;
  scrollbar-width: thin;
  scrollbar-color: rgba(194, 163, 95, 0.45) transparent;
}

.announcement-body::-webkit-scrollbar {
  width: 4px;
}

.announcement-body::-webkit-scrollbar-track {
  background: transparent;
}

.announcement-body::-webkit-scrollbar-thumb {
  background: rgba(194, 163, 95, 0.45);
  border: none;
  border-radius: 999px;
}

.announcement-body::-webkit-scrollbar-thumb:hover {
  background: rgba(194, 163, 95, 0.7);
}

.announcement-content {
  margin: 0;
  font-size: 14px;
  line-height: 1.75;
  color: #f0e6d2;
}

.announcement-content :deep(p) {
  margin: 0 0 0.75em;
}

.announcement-content :deep(p:last-child) {
  margin-bottom: 0;
}

.announcement-content :deep(ul),
.announcement-content :deep(ol) {
  margin: 0 0 0.75em;
  padding-left: 1.25em;
}

.announcement-content :deep(a) {
  color: #c2a35f;
}

.announcement-footer {
  padding: 16px 24px 24px;
  border-top: 1px solid rgba(194, 163, 95, 0.2);
}

.announcement-footer .confirm-btn {
  width: 100%;
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
  .announcement-overlay {
    padding: 16px;
  }

  .announcement-header,
  .announcement-body,
  .announcement-footer {
    padding-left: 18px;
    padding-right: 18px;
  }
}
</style>
