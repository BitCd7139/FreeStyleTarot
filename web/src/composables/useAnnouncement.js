import { ref } from 'vue';
import { fetchAnnouncement } from '../utils/announcementApi.js';

const SEEN_KEY_PREFIX = 'fst_announcement_seen_boot:';

const show = ref(false);
const title = ref('系统公告');
const content = ref('');
const isGuestView = ref(false);
const isManualView = ref(false);
let currentBootId = '';

function seenKey(userId) {
  return `${SEEN_KEY_PREFIX}${userId}`;
}

function hasSeen(userId, bootId) {
  return localStorage.getItem(seenKey(userId)) === bootId;
}

function markSeen(userId, bootId) {
  localStorage.setItem(seenKey(userId), bootId);
}

async function loadAnnouncement() {
  const data = await fetchAnnouncement();
  if (!data.enabled || !data.boot_id || !data.content?.trim()) return null;
  return data;
}

export function useAnnouncement() {
  async function checkAndShow(user) {
    if (!user?.id) return;

    try {
      const data = await loadAnnouncement();
      if (!data) {
        show.value = false;
        return;
      }
      if (hasSeen(user.id, data.boot_id)) {
        show.value = false;
        return;
      }

      currentBootId = data.boot_id;
      isGuestView.value = false;
      title.value = data.title?.trim() || '系统公告';
      content.value = data.content;
      show.value = true;
    } catch {
      // 公告拉取失败时不阻断主流程
    }
  }

  async function checkAndShowForGuest() {
    try {
      const data = await loadAnnouncement();
      if (!data) return;

      currentBootId = data.boot_id;
      isGuestView.value = true;
      title.value = data.title?.trim() || '系统公告';
      content.value = data.content;
      show.value = true;
    } catch {
      // 公告拉取失败时不阻断主流程
    }
  }

  function dismiss(user) {
    if (user?.id && currentBootId) {
      markSeen(user.id, currentBootId);
    }
    show.value = false;
  }

  function dismissGuest() {
    show.value = false;
  }

  async function openManually() {
    try {
      const data = await loadAnnouncement();
      if (!data) return;

      currentBootId = data.boot_id;
      isGuestView.value = false;
      isManualView.value = true;
      title.value = data.title?.trim() || '系统公告';
      content.value = data.content;
      show.value = true;
    } catch {
      // 公告拉取失败时不阻断主流程
    }
  }

  function dismissManual() {
    isManualView.value = false;
    show.value = false;
  }

  return {
    show,
    title,
    content,
    isGuestView,
    isManualView,
    checkAndShow,
    checkAndShowForGuest,
    openManually,
    dismiss,
    dismissGuest,
    dismissManual,
  };
}
