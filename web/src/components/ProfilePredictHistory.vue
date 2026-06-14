<template>
  <div class="profile-field history-section">
    <button type="button" class="history-section-toggle" @click="toggleSection">
      <span class="field-label">{{ sectionTitle }}</span>
      <span class="section-chevron">{{ sectionOpen ? '▾' : '▸' }}</span>
    </button>

    <div v-if="sectionOpen" class="history-panel">
      <p v-if="loading" class="field-value history-status">加载中...</p>
      <p v-else-if="error" class="error-text">{{ error }}</p>
      <p v-else-if="!items.length" class="field-value history-status">暂无提问记录</p>

      <div v-else class="history-list">
        <button
          v-for="item in items"
          :key="item.id"
          type="button"
          class="history-row"
          @click="openItem(item)"
        >
          <span class="history-row-inner">
            <span class="history-dot" aria-hidden="true"></span>
            <span class="field-value history-question-text">{{ item.question }}</span>
          </span>
        </button>
      </div>
    </div>

    <Teleport to="body">
      <AnswerModal
        v-model:showModal="showAnswerModal"
        :fullAnswer="selectedAnswer"
        :drawnCards="selectedCards"
        :cardWidth="selectedCardWidth"
        :cardHeight="selectedCardHeight"
        :isStreaming="false"
        phaseLabel=""
      />
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import AnswerModal from './AnswerModal.vue';
import { usePredictHistory } from '../composables/usePredictHistory.js';

const { items, limit, loading, error, loaded, loadHistory, refreshHistory } = usePredictHistory();

const sectionOpen = ref(false);
const showAnswerModal = ref(false);
const selectedAnswer = ref('');
const selectedCards = ref([]);
const selectedCardWidth = ref(120);
const selectedCardHeight = ref(210);

const sectionTitle = computed(() => {
  if (!loaded.value) return '提问历史';
  const n = items.value.length;
  return `提问历史（当前展示 ${n}/${limit.value} 条）`;
});

async function toggleSection() {
  sectionOpen.value = !sectionOpen.value;
  if (sectionOpen.value && !loaded.value) {
    await loadHistory();
  }
}

function toDrawnCards(item) {
  return (item.cards || []).map((c, i) => ({
    id: i,
    order: c.order,
    name: c.name,
    x: c.x,
    y: c.y,
    isReversed: c.orientation === 'reversed',
    meaning: c.meaning || '',
  }));
}

function openItem(item) {
  selectedAnswer.value = item.answer || '';
  selectedCards.value = toDrawnCards(item);
  selectedCardWidth.value = item.cardSize?.width || 120;
  selectedCardHeight.value = item.cardSize?.height || 210;
  showAnswerModal.value = true;
}

defineExpose({ refreshHistory });
</script>

<style scoped>
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

.history-section {
  margin-bottom: 0;
}

.history-section-toggle {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 0;
  background: transparent;
  border: none;
  cursor: pointer;
  text-align: left;
}

.history-section-toggle .field-label {
  margin-bottom: 0;
}

.section-chevron {
  color: #c2a35f;
  font-size: 12px;
  flex-shrink: 0;
}

.history-panel {
  margin-top: 8px;
}

.history-status {
  margin: 0;
}

.history-list {
  display: flex;
  flex-direction: column;
  max-height: 320px;
  overflow-y: auto;
  border-top: 1px solid rgba(194, 163, 95, 0.2);
}

.history-row {
  width: 100%;
  display: flex;
  align-items: center;
  padding: 12px 4px;
  background: transparent;
  border: none;
  border-bottom: 1px solid rgba(194, 163, 95, 0.15);
  cursor: pointer;
  text-align: left;
}

.history-row-inner {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
  width: 100%;
}

.history-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #c2a35f;
  flex-shrink: 0;
}

.history-row:last-child {
  border-bottom: none;
}

.history-row:hover {
  background: rgba(194, 163, 95, 0.08);
}

.history-row:hover .field-value {
  color: #f0e6c8;
}

.history-question-text {
  flex: 1;
  min-width: 0;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.error-text {
  display: block;
  color: #ff4d4f;
  font-size: 13px;
  margin: 0;
}
</style>
