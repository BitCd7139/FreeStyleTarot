<!-- components/TrashCan.vue -->
<template>
    <div 
      class="trash-container" 
      :class="{ 
        'near-trash': isHovered, 
        'is-pressing': isPressing 
      }"
      @mousedown="startPress"
      @mouseup="cancelPress"
      @mouseleave="cancelPress"
      @touchstart="startPress" 
      @touchend="cancelPress"
    >
      <div class="icon-wrapper">
        <!-- 进度条圆环 -->
        <svg class="progress-ring" viewBox="0 0 100 100">
          <circle 
            class="progress-ring__circle" 
            stroke="white" 
            stroke-width="4" 
            fill="transparent" 
            r="45" cx="50" cy="50"
          />
        </svg>
  
        <!-- 垃圾桶图标 -->
        <svg class="trash-icon" viewBox="0 0 24 24" width="32" height="32" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="3 6 5 6 21 6"></polyline>
          <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
          <line x1="10" y1="11" x2="10" y2="17"></line>
          <line x1="14" y1="11" x2="14" y2="17"></line>
        </svg>
      </div>
      <span class="trash-label">{{ isPressing ? '正在清空...' : (isHovered ? '释放以删除' : '长按清空全部') }}</span>
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue';
  
  const props = defineProps({
    isHovered: Boolean
  });
  
  const emit = defineEmits(['clear-all']);
  
  const isPressing = ref(false);
  let pressTimer = null;
  const LONG_PRESS_TIME = 800; // 长按800毫秒触发
  
  const startPress = (e) => {
    if (props.isHovered) return;
    
    isPressing.value = true;
    pressTimer = setTimeout(() => {
      emit('clear-all');
      isPressing.value = false;
    }, LONG_PRESS_TIME);
  };
  
  const cancelPress = () => {
    isPressing.value = false;
    if (pressTimer) clearTimeout(pressTimer);
  };
  </script>
  
  <style scoped>
  .trash-container {
    position: absolute;
    top: 30px;
    right: 30px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    color: rgba(255, 255, 255, 0.4);
    transition: all 0.3s ease;
    z-index: 1000;
    padding: 20px;
    cursor: pointer;
    user-select: none;
    pointer-events: auto; 
  }
  
  .icon-wrapper {
    position: relative;
    background: rgba(255, 255, 255, 0.05);
    padding: 12px;
    border-radius: 50%;
    border: 1px solid rgba(255, 255, 255, 0.1);
    transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  /* 进度条圆环样式 */
  .progress-ring {
    position: absolute;
    width: 100%;
    height: 100%;
    transform: rotate(-90deg);
    pointer-events: none;
    opacity: 0;
    transition: opacity 0.2s;
  }
  
  .progress-ring__circle {
    stroke-dasharray: 283; /* 2 * PI * r (45) */
    stroke-dashoffset: 283;
    transition: stroke-dashoffset 0.8s linear; /* 时间与 LONG_PRESS_TIME 一致 */
  }
  
  /* 长按时的状态 */
  .is-pressing .progress-ring {
    opacity: 1;
  }
  .is-pressing .progress-ring__circle {
    stroke-dashoffset: 0;
    stroke: #ff4d4f;
  }
  .is-pressing .icon-wrapper {
    transform: scale(0.9);
    background: rgba(255, 77, 79, 0.1);
  }
  
  /* 靠近时的状态 (拖拽卡牌中) */
  .near-trash {
    color: #ff4d4f;
  }
  .near-trash .icon-wrapper {
    background: rgba(255, 77, 79, 0.2);
    border-color: #ff4d4f;
    transform: scale(1.25);
    box-shadow: 0 0 25px rgba(255, 77, 79, 0.3);
  }
  
  .trash-label {
    font-size: 12px;
    opacity: 0;
    transition: opacity 0.3s;
  }
  
  .near-trash .trash-label, .is-pressing .trash-label {
    opacity: 1;
  }
  </style>