<template>
  <Transition name="fade">
    <div v-if="show" class="picker-overlay" @click.self="close">
      <div class="picker-panel">
        <div class="picker-header">
          <h3>自选卡牌</h3>
          <p class="picker-hint">搜索卡牌名称、花色或数字，选择正逆位后确认添加</p>
          <button type="button" class="close-btn" @click="close" aria-label="关闭">×</button>
        </div>

        <div class="search-bar">
          <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8"></circle>
            <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
          </svg>
          <input
            ref="searchInput"
            v-model="searchQuery"
            type="search"
            placeholder="搜索：愚者、0、大阿卡罗纳、宝剑、11侍从、14国王..."
            autocomplete="off"
          />
        </div>

        <div class="picker-body">
          <div v-if="filteredCards.length === 0" class="empty-hint">
            未找到匹配的卡牌，试试「0愚人」「大阿卡罗纳13」「宝剑11侍从」
          </div>
          <div v-else class="card-grid">
            <button
              v-for="card in filteredCards"
              :key="card.id"
              type="button"
              class="card-option"
              :class="{ selected: selectedCardId === card.id }"
              @click="selectCard(card.id)"
            >
              <img :src="getCardImageUrl(card.id)" :alt="card.cnName" draggable="false" />
              <span class="card-name">{{ card.cnName }}</span>
            </button>
          </div>
        </div>

        <div class="picker-footer">
          <div v-if="selectedCardId" class="preview-row">
            <img
              :src="getCardImageUrl(selectedCardId)"
              :class="{ 'is-reversed': isReversed }"
              alt="preview"
              draggable="false"
            />
            <div class="preview-info">
              <span class="preview-name">{{ getName(selectedCardId) }}</span>
              <div class="orientation-toggle">
                <button
                  type="button"
                  :class="{ active: !isReversed }"
                  @click="isReversed = false"
                >正位</button>
                <button
                  type="button"
                  :class="{ active: isReversed }"
                  @click="isReversed = true"
                >逆位</button>
              </div>
            </div>
          </div>
          <div v-else class="preview-placeholder">请先选择一张卡牌</div>

          <div class="action-btns">
            <button type="button" class="btn-cancel" @click="close">取消</button>
            <button
              type="button"
              class="btn-confirm"
              :disabled="!selectedCardId"
              @click="confirm"
            >确认添加</button>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue';
import { searchCards, getName, getCardImageUrl } from '../utils/cardInfo.js';

const props = defineProps({
  show: Boolean,
});

const emit = defineEmits(['update:show', 'confirm']);

const searchQuery = ref('');
const selectedCardId = ref(null);
const isReversed = ref(false);
const searchInput = ref(null);

const filteredCards = computed(() => searchCards(searchQuery.value));

watch(() => props.show, async (visible) => {
  if (visible) {
    searchQuery.value = '';
    selectedCardId.value = null;
    isReversed.value = false;
    await nextTick();
    searchInput.value?.focus();
  }
});

function selectCard(id) {
  selectedCardId.value = id;
}

function close() {
  emit('update:show', false);
}

function confirm() {
  if (!selectedCardId.value) return;
  emit('confirm', {
    name: selectedCardId.value,
    isReversed: isReversed.value,
  });
  emit('update:show', false);
}
</script>

<style scoped>
.picker-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.65);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  padding: 20px;
}

.picker-panel {
  background: rgba(30, 27, 36, 0.95);
  backdrop-filter: blur(16px);
  border: 1px solid rgba(194, 163, 95, 0.35);
  border-radius: 16px;
  width: min(560px, 100%);
  max-height: min(85vh, 720px);
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.6);
  overflow: hidden;
}

.picker-header {
  padding: 22px 24px 12px;
  border-bottom: 1px solid rgba(194, 163, 95, 0.2);
  position: relative;
}

.picker-header h3 {
  margin: 0 0 6px;
  color: #e5d8b0;
  font-size: 18px;
}

.picker-hint {
  margin: 0;
  font-size: 12px;
  color: #a89f91;
}

.close-btn {
  position: absolute;
  top: 16px;
  right: 16px;
  background: none;
  border: none;
  color: #a89f91;
  font-size: 24px;
  line-height: 1;
  cursor: pointer;
  padding: 4px 8px;
  transition: color 0.2s;
}

.close-btn:hover {
  color: #e5d8b0;
}

.search-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 16px 20px 0;
  padding: 10px 14px;
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(194, 163, 95, 0.35);
  border-radius: 10px;
  color: #a89f91;
}

.search-bar input {
  flex: 1;
  background: transparent;
  border: none;
  color: #fff;
  font-size: 14px;
  outline: none;
}

.search-bar input::placeholder {
  color: #6b6358;
}

.picker-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px 20px;
  min-height: 200px;
}

.picker-body::-webkit-scrollbar {
  width: 5px;
}

.picker-body::-webkit-scrollbar-thumb {
  background: #c2a35f;
  border-radius: 3px;
}

.empty-hint {
  text-align: center;
  color: #a89f91;
  font-size: 14px;
  padding: 40px 20px;
}

.card-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(72px, 1fr));
  gap: 10px;
}

.card-option {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 8px 4px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(194, 163, 95, 0.15);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.card-option:hover {
  background: rgba(194, 163, 95, 0.12);
  border-color: rgba(194, 163, 95, 0.4);
}

.card-option.selected {
  background: rgba(194, 163, 95, 0.2);
  border-color: #c2a35f;
  box-shadow: 0 0 12px rgba(194, 163, 95, 0.25);
}

.card-option img {
  width: 48px;
  height: 80px;
  object-fit: cover;
  border-radius: 4px;
}

.card-name {
  font-size: 10px;
  color: #e5d8b0;
  text-align: center;
  line-height: 1.2;
  word-break: break-all;
}

.picker-footer {
  padding: 16px 20px 20px;
  border-top: 1px solid rgba(194, 163, 95, 0.2);
}

.preview-row {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.preview-row img {
  width: 56px;
  height: 94px;
  object-fit: cover;
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.5);
  transition: transform 0.3s;
}

.preview-row img.is-reversed {
  transform: rotate(180deg);
}

.preview-name {
  display: block;
  color: #e5d8b0;
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 10px;
}

.orientation-toggle {
  display: flex;
  gap: 8px;
}

.orientation-toggle button {
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid rgba(194, 163, 95, 0.4);
  background: transparent;
  color: #a89f91;
}

.orientation-toggle button.active {
  background: #c2a35f;
  color: #1e1b24;
  border-color: #c2a35f;
}

.preview-placeholder {
  color: #6b6358;
  font-size: 13px;
  margin-bottom: 16px;
  text-align: center;
}

.action-btns {
  display: flex;
  gap: 12px;
}

.btn-cancel,
.btn-confirm {
  flex: 1;
  padding: 12px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-cancel {
  background: transparent;
  border: 1px solid rgba(194, 163, 95, 0.4);
  color: #c2a35f;
}

.btn-cancel:hover {
  background: rgba(194, 163, 95, 0.1);
}

.btn-confirm {
  background: #c2a35f;
  border: 1px solid #c2a35f;
  color: #1e1b24;
}

.btn-confirm:hover:not(:disabled) {
  background: #d4af37;
  box-shadow: 0 0 12px rgba(194, 163, 95, 0.4);
}

.btn-confirm:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
