<template>
  <div class="beginner-page">
    <div class="beginner-header">
      <h2>选择一个牌阵</h2>
      <p class="beginner-hint">不同牌阵适合不同问题，挑一个开始吧</p>
      <a class="advanced-link" @click="goAdvanced">我想自己摆，点击进入高级模式</a>
    </div>

    <div class="spread-grid">
      <div
        v-for="spread in catalog"
        :key="spread.id"
        class="spread-card"
        :class="{ active: selectedId === spread.id }"
        @click="selectedId = spread.id"
      >
        <div class="spread-preview">
          <MiniTarot
            :drawn-cards="spread.preview.cards"
            :card-width="previewCardW"
            :card-height="previewCardH"
            :fixed="true"
            :fixed-stage-width="spread.preview.stageWidth"
            :fixed-stage-height="spread.preview.stageHeight"
            :connections="spread.preview.connections"
          />
        </div>
        <div class="spread-info">
          <div class="spread-name">{{ spread.name }}</div>
          <div class="spread-meta">{{ spread.cardCount }} 张牌</div>
          <div class="spread-desc">{{ spread.description }}</div>
          <div class="spread-usecase">{{ spread.useCase }}</div>
        </div>
      </div>
    </div>

    <div class="beginner-footer">
      <button type="button" class="start-btn" :disabled="!selectedId" @click="start">
        开始抽牌
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import MiniTarot from './MiniTarot.vue';
import { getSpreadCatalog } from '../spread/index.js';
import { slotsToFixedPreviewCards, getSpreadBBox } from '../utils/spreadLayout.js';
import { useTarotFlow } from '../composables/useTarotFlow.js';
import { useAuth } from '../composables/useAuth.js';

// 卡牌尺寸：不超过屏幕宽高的 1/3
const screenThirdW = window.innerWidth / 3;
const screenThirdH = window.innerHeight / 3;
const previewCardW = Math.min(40, screenThirdW);
const previewCardH = Math.min(68, screenThirdH);

// 统一舞台尺寸：取所有牌阵最大 bbox + padding，保证容器大小强制一致
const allCatalog = getSpreadCatalog();
const padding = 14;
let maxBboxW = 0;
let maxBboxH = 0;
for (const spread of allCatalog) {
  const { bboxW, bboxH } = getSpreadBBox(spread.slots, previewCardW, previewCardH);
  if (bboxW > maxBboxW) maxBboxW = bboxW;
  if (bboxH > maxBboxH) maxBboxH = bboxH;
}
const previewStageW = maxBboxW + padding * 2;
const previewStageH = maxBboxH + padding * 2;

function buildPreview(spread) {
  const cards = slotsToFixedPreviewCards(
    spread.slots,
    previewCardW,
    previewCardH,
    previewStageW,
    previewStageH
  );
  const connections = (spread.connections || []).map(([a, b]) => [
    `fixed-preview-${a}`,
    `fixed-preview-${b}`,
  ]);
  return { cards, stageWidth: previewStageW, stageHeight: previewStageH, connections };
}

const catalog = allCatalog
  .map(spread => ({ ...spread, preview: buildPreview(spread) }))
  .sort((a, b) => a.cardCount - b.cardCount);
const { selectSpread } = useTarotFlow();
const { updateExperienceMode } = useAuth();
const selectedId = ref('');

function start() {
  const spread = catalog.find(s => s.id === selectedId.value);
  if (!spread) return;
  selectSpread(spread);
}

async function goAdvanced() {
  try {
    await updateExperienceMode('advanced');
  } catch (e) {
    // 忽略错误：user.experience_mode 更新后 isBeginner 自动变 false，App.vue 会切回 TarotMain
  }
}
</script>

<style scoped>
.beginner-page {
  min-height: 100vh;
  padding: 24px 20px 120px;
  box-sizing: border-box;
  background: #433843;
}

.beginner-header {
  text-align: center;
  margin-bottom: 24px;
}

.beginner-header h2 {
  margin: 0;
  color: #e5d8b0;
  font-size: 24px;
}

.beginner-hint {
  margin: 8px 0 0;
  color: #a89f91;
  font-size: 14px;
}

.advanced-link {
  display: inline-block;
  margin-top: 12px;
  color: #6b6358;
  font-size: 12px;
  cursor: pointer;
  text-decoration: underline;
  transition: color 0.2s ease;
}

.advanced-link:hover {
  color: #a89f91;
}

.spread-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  max-width: 960px;
  margin: 0 auto;
}

.spread-card {
  background: linear-gradient(180deg, #2a2438 0%, #1a1626 100%);
  border: 1px solid rgba(194, 163, 95, 0.25);
  border-radius: 14px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.spread-card:hover {
  border-color: rgba(194, 163, 95, 0.6);
  transform: translateY(-2px);
}

.spread-card.active {
  border-color: #c2a35f;
  box-shadow: 0 0 0 2px rgba(194, 163, 95, 0.5) inset, 0 6px 20px rgba(0, 0, 0, 0.4);
}

.spread-preview {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 10px;
  padding: 12px;
  min-height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.spread-info {
  text-align: center;
}

.spread-name {
  color: #e5d8b0;
  font-size: 17px;
  font-weight: 600;
  margin-bottom: 4px;
}

.spread-meta {
  color: #c2a35f;
  font-size: 12px;
  margin-bottom: 8px;
}

.spread-desc {
  color: #a89f91;
  font-size: 13px;
  line-height: 1.5;
}

.spread-usecase {
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px dashed rgba(194, 163, 95, 0.25);
  color: #c9b890;
  font-size: 12px;
  line-height: 1.5;
  text-align: left;
}

.spread-usecase::before {
  content: "使用场景：";
  color: #c2a35f;
  font-weight: 600;
}

.beginner-footer {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 16px 20px calc(16px + env(safe-area-inset-bottom));
  background: linear-gradient(180deg, rgba(26, 22, 38, 0) 0%, rgba(26, 22, 38, 0.95) 30%);
  display: flex;
  justify-content: center;
}

.start-btn {
  width: 100%;
  max-width: 420px;
  padding: 14px;
  background: linear-gradient(180deg, #c2a35f 0%, #a88a45 100%);
  border: none;
  color: #1a1626;
  border-radius: 10px;
  cursor: pointer;
  font-size: 16px;
  font-weight: 600;
}

.start-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

@media (max-width: 600px) {
  .spread-grid {
    grid-template-columns: 1fr;
  }
}
</style>
