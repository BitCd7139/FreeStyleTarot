<template>
    <div 
      class="card-item"
      :class="{ 'dragging': isActive, 'resizing': isGlobalResizing }"
      :style="{ 
        left: card.x + 'px', 
        top: card.y + 'px',
        width: width + 'px',
        height: height + 'px',
        zIndex: isActive ? 1000 : card.order 
      }"
      @mousedown.stop="$emit('drag-start', card, $event)"
    >
      <!-- 卡牌图片 -->
      <div class="card-inner">
        <img 
          :src="card.isRevealed ? getCardUrl(card.name) : getCardUrl('back')" 
          :class="{ 'reversed': card.isReversed && card.isRevealed }"
          class="tarot-img"
          draggable="false"
        />
        <div v-if="card.isRevealed" class="order-badge">{{ card.order }}</div>
      </div>
  
      <!-- 四角缩放手柄 -->
      <div class="resizer nw" @mousedown.stop="$emit('resize-start', $event)"></div>
      <div class="resizer ne" @mousedown.stop="$emit('resize-start', $event)"></div>
      <div class="resizer sw" @mousedown.stop="$emit('resize-start', $event)"></div>
      <div class="resizer se" @mousedown.stop="$emit('resize-start', $event)"></div>
    </div>
  </template>
  
  <script setup>
  import { defineProps, defineEmits } from 'vue';
  
  const props = defineProps({
    card: {
      type: Object,
      required: true
    },
    width: {
      type: Number,
      required: true
    },
    height: {
      type: Number,
      required: true
    },
    isActive: {
      type: Boolean,
      default: false
    },
    isGlobalResizing: {
      type: Boolean,
      default: false
    }
  });
  
  defineEmits(['drag-start', 'resize-start']);
  
  // 处理图片 URL
  const getCardUrl = (name) => {
    if (name === 'back') return new URL(`../assets/tarots/back.jpeg`, import.meta.url).href;
    const ext = 'jpeg';
    return new URL(`../assets/tarots/${name}.${ext}`, import.meta.url).href;
  };
  </script>
  
  <style scoped>
/* ==================== 卡牌专属样式 ==================== */
.card-item {
  position: absolute;
  cursor: grab;
  background: transparent; 
  transition: transform 0.2s cubic-bezier(0.25, 0.46, 0.45, 0.94);
}

.card-item.dragging {
  cursor: grabbing;
  transform: scale(1.02);
  filter: brightness(1.1);
  z-index: 999;
}

.card-inner {
  position: relative;
  width: 100%;
  height: 100%;
  border-radius: 6px; /* 圆角稍微收敛一点 */
  box-shadow: 0 6px 20px rgba(0,0,0,0.6);
  overflow: visible;
}

.tarot-img {
  width: 100%;
  height: 100%;
  border-radius: 6px;
  object-fit: 100% 100%; 
  display: block;
  pointer-events: none;
  /* border: 1px solid rgba(255, 255, 255, 0.1); */
}

.reversed {
  transform: rotate(180deg);
}

/* 序号标签也稍微变小一点，匹配更小的缩放手柄 */
.order-badge {
  position: absolute;
  top: -8px;
  right: -8px;
  background: linear-gradient(135deg, #ffd700, #b8860b);
  color: #000;
  border-radius: 50%;
  width: 22px;
  height: 22px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: bold;
  border: 1px solid #000;
  box-shadow: 0 2px 4px rgba(0,0,0,0.5);
}

/* 缩放手柄样式 */
.resizer {
  position: absolute;
  /* 尺寸从 16px 缩小到 10px */
  width: 10px;
  height: 10px;
  background: rgba(255, 255, 255, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.6);
  border-radius: 50%;
  z-index: 10;
  opacity: 0;
  transition: opacity 0.3s, background 0.2s, transform 0.2s;
}

/* 鼠标悬浮在卡牌上时显示手柄 */
.card-item:hover .resizer {
  opacity: 1;
}

/* 鼠标悬浮在手柄上时的反馈 */
.resizer:hover {
  background: #fff;
  box-shadow: 0 0 8px #fff;
  transform: scale(1.2);
}

/* 调整手柄位置偏移（因为尺寸变成了 10px，所以向外偏移 5px 居中在角上） */
.nw { top: -5px; left: -5px; cursor: nwse-resize; }
.ne { top: -5px; right: -5px; cursor: nesw-resize; }
.sw { bottom: -5px; left: -5px; cursor: nesw-resize; }
.se { bottom: -5px; right: -5px; cursor: nwse-resize; }
</style>