<template>
    <div class="mini-tarot-stage" :style="{ height: containerHeight + 'px' }">
      <div 
        class="spread-wrapper" 
        :style="{
          width: layoutState.sWidth + 'px',
          height: layoutState.sHeight + 'px',
          transform: `translate(-50%, -50%) scale(${layoutState.scale})`
        }"
      >
        <div 
          v-for="card in drawnCards" 
          :key="card.id"
          class="static-card"
          :style="{
            left: (card.x - layoutState.minX) + 'px',
            top: (card.y - layoutState.minY) + 'px',
            width: cardWidth + 'px',
            height: cardHeight + 'px',
            zIndex: card.order
          }"
        >
          <img 
            :src="getCardUrl(card.name)" 
            :class="{ 'is-reversed': card.isReversed }" 
            alt="tarot" 
          />
          <div class="card-badge">{{ card.order }}</div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup>
  import { computed } from 'vue';
  
  const props = defineProps({
    drawnCards: { type: Array, default: () => [] },
    cardWidth: { type: Number, default: 120 },
    cardHeight: { type: Number, default: 210 },
    // 容器的预估宽度，用来计算缩放比
    containerWidth: { type: Number, default: 400 },
    // 容器的高度，由父组件决定预览区有多高
    containerHeight: { type: Number, default: 200 }
  });
  
  const layoutState = computed(() => {
    if (!props.drawnCards || props.drawnCards.length === 0) {
      return { scale: 1, minX: 0, minY: 0, sWidth: 0, sHeight: 0 };
    }
    
    let minX = Infinity, minY = Infinity;
    let maxX = -Infinity, maxY = -Infinity;
    
    props.drawnCards.forEach(c => {
      if (c.x < minX) minX = c.x;
      if (c.y < minY) minY = c.y;
      if (c.x + props.cardWidth > maxX) maxX = c.x + props.cardWidth;
      if (c.y + props.cardHeight > maxY) maxY = c.y + props.cardHeight;
    });
    
    const sWidth = maxX - minX;
    const sHeight = maxY - minY;
    
    const scaleX = props.containerWidth / (sWidth || 1);
    const scaleY = props.containerHeight / (sHeight || 1);
    // 乘以 0.85 留出安全内边距
    let scale = Math.min(scaleX, scaleY) * 0.85; 
    if (scale > 1) scale = 1; 
  
    return { scale, minX, minY, sWidth, sHeight };
  });
  
  const getCardUrl = (name) => {
    if (!name || name === 'back') return new URL(`../assets/tarots/back.jpeg`, import.meta.url).href;
    const ext = 'jpeg';
    return new URL(`../assets/tarots/${name}.${ext}`, import.meta.url).href;
  };
  </script>
  
  <style scoped>
  .mini-tarot-stage {
    width: 100%;
    background: #0a0a0c;
    border-radius: 12px;
    border: 1px solid #222;
    position: relative;
    overflow: hidden;
    margin-bottom: 20px;
  }
  
  .spread-wrapper {
    position: absolute;
    top: 50%;
    left: 50%;
    transform-origin: center center;
    transition: transform 0.3s ease;
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
  
  .static-card img.is-reversed {
    transform: rotate(180deg);
  }
  
  .card-badge {
    position: absolute;
    bottom: -10px;
    left: 50%;
    transform: translateX(-50%);
    background: #fff;
    color: #000;
    width: 22px;
    height: 22px;
    border-radius: 50%;
    font-size: 13px;
    font-weight: bold;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 2px 4px rgba(0,0,0,0.5);
    z-index: 10;
  }
  </style>