<template>
    <div class="mini-tarot-section" ref="sectionRef">
      <div v-if="!fixed" class="markdown-body" v-html="renderedTitle"></div>
      
      <div 
        class="mini-tarot-stage"
        :style="{ 
          width: layoutState.stageWidth + 'px',
          height: layoutState.stageHeight + 'px'
        }"
      >
        <div
          class="cards-wrapper"
          :style="{
            transform: `scale(${layoutState.scale})`,
            transformOrigin: 'top left',
            width: layoutState.spreadWidth + 'px',
            height: layoutState.spreadHeight + 'px'
          }"
        >
          <svg
            v-if="fixed && connectionLines.length"
            class="connections-svg"
            :width="layoutState.spreadWidth"
            :height="layoutState.spreadHeight"
            aria-hidden="true"
          >
            <line
              v-for="(line, i) in connectionLines"
              :key="i"
              :x1="line.x1"
              :y1="line.y1"
              :x2="line.x2"
              :y2="line.y2"
              stroke="rgba(194, 163, 95, 0.55)"
              stroke-width="1"
              stroke-linecap="round"
            />
          </svg>
          <div
            v-for="card in drawnCards"
            :key="card.id"
            class="static-card"
            :style="{
              left: layoutState.relativePositions[card.id]?.x + 'px',
              top: layoutState.relativePositions[card.id]?.y + 'px',
              width: cardWidth + 'px',
              height: cardHeight + 'px',
              zIndex: card.order
            }"
          >
            <img
              v-if="card.name && card.name !== 'back'"
              :src="getCardUrl(card.name)"
              :class="{ 'is-reversed': card.isReversed }"
              alt="tarot"
            />
            <CardBack v-else :class="{ 'is-reversed': card.isReversed }" />
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue';
import { marked } from 'marked';
import DOMPurify from 'dompurify';
import CardBack from './CardBack.vue';
  
  const props = defineProps({
    drawnCards: { type: Array, default: () => [] },
    cardWidth: { type: Number, default: 120 },
    cardHeight: { type: Number, default: 210 },
    fixed: { type: Boolean, default: false },
    fixedStageWidth: { type: Number, default: 220 },
    fixedStageHeight: { type: Number, default: 220  },
    // fixed 模式下的卡牌关系连接线：[[cardIdA, cardIdB], ...]
    connections: { type: Array, default: () => [] }
  });
  
  const sectionRef = ref(null);
  const parentWidth = ref(400);
  
  const updateDimensions = () => {
    if (sectionRef.value) {
      const parent = sectionRef.value.parentElement;
      if (parent) {
        parentWidth.value = parent.offsetWidth || 400;
      }
    }
  };
  
  onMounted(() => {
    nextTick(updateDimensions);
    window.addEventListener('resize', updateDimensions);
  });
  
  onUnmounted(() => {
    window.removeEventListener('resize', updateDimensions);
  });
  
  const renderedTitle = computed(() => {
    const rawHtml = marked.parse('## 🌈 牌阵展示');
    const styledHtml = `<div style="color: rgb(224, 224, 224);">${rawHtml}</div>`;
    return DOMPurify.sanitize(styledHtml);
  });
  
  const layoutState = computed(() => {
    // fixed 模式：固定舞台 + 固定卡牌大小 + 居中（坐标已由 slotsToFixedPreviewCards 算好）
    if (props.fixed) {
      const stageW = props.fixedStageWidth;
      const stageH = props.fixedStageHeight;
      if (!props.drawnCards || props.drawnCards.length === 0) {
        return { stageWidth: stageW, stageHeight: stageH, spreadWidth: 0, spreadHeight: 0, scale: 1, relativePositions: {} };
      }
      const relativePositions = {};
      props.drawnCards.forEach(c => {
        relativePositions[c.id] = { x: c.x, y: c.y };
      });
      return { stageWidth: stageW, stageHeight: stageH, spreadWidth: stageW, spreadHeight: stageH, scale: 1, relativePositions };
    }
    if (!props.drawnCards || props.drawnCards.length === 0) {
      return {
        stageWidth: 300,
        stageHeight: 200,
        spreadWidth: 300,
        spreadHeight: 200,
        scale: 0.5,
        relativePositions: {}
      };
    }
    
    let minX = Infinity, minY = Infinity;
    let maxX = -Infinity, maxY = -Infinity;
    
    props.drawnCards.forEach(c => {
      if (c.x < minX) minX = c.x;
      if (c.y < minY) minY = c.y;
      if (c.x + props.cardWidth > maxX) maxX = c.x + props.cardWidth;
      if (c.y + props.cardHeight > maxY) maxY = c.y + props.cardHeight;
    });
    
    const spreadWidth = maxX - minX;
    const spreadHeight = maxY - minY;
    
    const availableWidth = parentWidth.value - 20;
    const scaleByWidth = availableWidth / (spreadWidth || 1);
    
    const isNarrowScreen = window.innerWidth < 768;
    const maxStageHeight = window.innerHeight * (isNarrowScreen ? 0.35 : 0.6);
    const scaleByHeight = maxStageHeight / (spreadHeight || 1);
    
    let scale = Math.min(scaleByWidth, scaleByHeight);
    if (scale > 1) scale = 1;
    if (scale < 0.15) scale = 0.15;
    
    const stageWidth = Math.max(spreadWidth * scale, 150);
    const stageHeight = Math.max(spreadHeight * scale, 120);
    
    const relativePositions = {};
    props.drawnCards.forEach(c => {
      relativePositions[c.id] = {
        x: c.x - minX,
        y: c.y - minY
      };
    });
    
    return {
      stageWidth,
      stageHeight,
      spreadWidth,
      spreadHeight,
      scale,
      relativePositions
    };
  });
  
  const getCardUrl = (name) => {
    return new URL(`../assets/tarots/${name}.jpeg`, import.meta.url).href;
  };

  // fixed 模式下的关系连接线坐标（基于 relativePositions 计算卡牌中心点）
  const connectionLines = computed(() => {
    if (!props.fixed || !Array.isArray(props.connections) || props.connections.length === 0) return [];
    const pos = layoutState.value.relativePositions || {};
    const lines = [];
    for (const pair of props.connections) {
      if (!Array.isArray(pair) || pair.length !== 2) continue;
      const a = pos[pair[0]];
      const b = pos[pair[1]];
      if (!a || !b) continue;
      lines.push({
        x1: a.x + props.cardWidth / 2,
        y1: a.y + props.cardHeight / 2,
        x2: b.x + props.cardWidth / 2,
        y2: b.y + props.cardHeight / 2,
      });
    }
    return lines;
  });
  </script>
  
  <style scoped>
.mini-tarot-section {
  width: 100%;
  margin-bottom: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.markdown-body {
  width: 100%;
  color: #fff;
}

/* 牌阵标题：去除下划线，复用全局 h2 样式 */
.markdown-body :deep(h2) {
  border-bottom: none;
  padding-bottom: 0;
}

.mini-tarot-stage {
  background: transparent;
  border-radius: 0;
  border: none;
  position: relative;
  overflow: visible;
  margin: 0;
  padding: 0;
}

.cards-wrapper {
  position: absolute;
  top: 0;
  left: 0;
}

.connections-svg {
  position: absolute;
  top: 0;
  left: 0;
  z-index: 0;
  pointer-events: none;
}

.static-card {
  position: absolute;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.5);
  background: #111;
}

.static-card img {
  width: 100%;
  height: 100%;
  border-radius: 8px;
  object-fit: cover;
}

.static-card :deep(.card-back-svg) {
  border-radius: 8px;
}

.static-card img.is-reversed {
  transform: rotate(180deg);
}

@media (max-width: 768px) {
  .mini-tarot-stage {
    min-height: 120px;
  }
}
</style>