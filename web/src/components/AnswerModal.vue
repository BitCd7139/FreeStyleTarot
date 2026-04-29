<template>
    <Transition name="fade">
      <div v-if="showModal" class="modal-overlay">
        <div class="modal-content answer-modal">
          <div class="modal-header">
            <h3>命运的指引</h3>
            <p class="modal-hint">塔罗的解答已在此显现</p>
          </div>
          
          <div class="modal-body">
            <!-- 现场渲染的迷你嵌套牌阵 -->
            <MiniTarot 
              :drawnCards="drawnCards" 
              :cardWidth="cardWidth" 
              :cardHeight="cardHeight" 
              :containerWidth="720" 
              :containerHeight="240" 
            />
  
            <!-- 解析展示区域 -->
            <div class="answer-box">
                <div class="markdown-body" v-html="parsedAnswer"></div>
            </div>
          </div>
  
          <div class="modal-footer">
            <button 
              class="confirm-btn" 
              @click="closeModal" 
              :disabled="isTyping"
            >
              {{ isTyping ? '正在倾听...' : '完成并关闭' }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </template>
  
  <script setup>
  import { ref, watch, computed, defineModel } from 'vue';
  import MiniTarot from './MiniTarot.vue';
  import { marked } from 'marked';
  import DOMPurify from 'dompurify';
  
  const showModal = defineModel('showModal');
  const props = defineProps({
    fullAnswer: { type: String, default: '' },
    // 接收来自主界面的牌阵数据
    drawnCards: { type: Array, default: () => [] },
    cardWidth: { type: Number, default: 120 },
    cardHeight: { type: Number, default: 210 }
  });
  
  // 1. 动态计算牌阵的边界与缩放比例 (核心逻辑)
  const layoutState = computed(() => {
    if (!props.drawnCards || props.drawnCards.length === 0) {
      return { scale: 1, minX: 0, minY: 0, sWidth: 0, sHeight: 0 };
    }
    
    let minX = Infinity, minY = Infinity;
    let maxX = -Infinity, maxY = -Infinity;
    
    // 找出所有卡牌组成的最外围边界
    props.drawnCards.forEach(c => {
      if (c.x < minX) minX = c.x;
      if (c.y < minY) minY = c.y;
      if (c.x + props.cardWidth > maxX) maxX = c.x + props.cardWidth;
      if (c.y + props.cardHeight > maxY) maxY = c.y + props.cardHeight;
    });
    
    const sWidth = maxX - minX;
    const sHeight = maxY - minY;
    
    // 定义预览容器的大致宽高 (对应CSS中的 .embedded-stage)
    const containerWidth = 720; 
    const containerHeight = 240; 
    
    // 计算缩放比例 (乘以 0.9 是为了留出 10% 的内边距)
    const scaleX = containerWidth / sWidth;
    const scaleY = containerHeight / sHeight;
    let scale = Math.min(scaleX, scaleY) * 0.9; 
    if (scale > 1) scale = 1; // 如果原本就很小，就不放大，保持原大小
  
    return { scale, minX, minY, sWidth, sHeight };
  });
  
  // 获取图片的工具函数
  const getCardUrl = (name) => {
    if (!name || name === 'back') return new URL(`../assets/tarots/back.jpeg`, import.meta.url).href;
    const ext = 'jpeg';
    return new URL(`../assets/tarots/${name}.${ext}`, import.meta.url).href;
  };
  
  // 2. 流式输出状态与逻辑
  const displayedText = ref('');
  const isTyping = ref(false);
  let typingInterval = null;
  
  // === 新增：实时解析 Markdown 并处理光标 ===
  const parsedAnswer = computed(() => {
    // 巧妙处理：把光标的 HTML 标签直接拼接到当前文本末尾，再一起进行 Markdown 解析
    // 这样光标就会自然地跟在当前段落的最后一个字后面，不会发生错位
    const rawText = isTyping.value 
      ? displayedText.value + '<span class="cursor">|</span>' 
      : displayedText.value;
  
    // 将 Markdown 转换为 HTML
    const rawHtml = marked.parse(rawText);
    // 净化 HTML 防止 XSS，同时允许我们自定义的 cursor class
    return DOMPurify.sanitize(rawHtml, { ADD_ATTR: ['class'] });
  });

  watch(() => props.fullAnswer, (newVal) => {
  if (!newVal) return;
  displayedText.value = '';
  isTyping.value = true;
  let currentIndex = 0;
  
  if (typingInterval) clearInterval(typingInterval);

  typingInterval = setInterval(() => {
    if (currentIndex < newVal.length) {
      displayedText.value += newVal[currentIndex];
      currentIndex++;
      // 自动滚动到底部 (可选，提升体验)
      scrollToBottom();
    } else {
      clearInterval(typingInterval);
      isTyping.value = false;
    }
  }, 30);
}, { immediate: true });

const closeModal = () => {
  showModal.value = false;
};

// 自动滚动的辅助函数
const scrollToBottom = () => {
  const box = document.querySelector('.answer-box');
  if (box) {
    box.scrollTop = box.scrollHeight;
  }
};
</script>
  
  <style scoped>
  /* 背景蒙版 */
  .modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.75);
    backdrop-filter: blur(8px);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 2000;
  }
  
  /* 弹窗主体变得更宽，以容纳牌阵和文字 */
  .modal-content.answer-modal {
    background: #151518;
    padding: 30px 40px;
    border-radius: 24px;
    width: 800px; 
    max-width: 95vw;
    border: 1px solid #333;
    text-align: center;
    box-shadow: 0 10px 40px rgba(0,0,0,0.8);
    display: flex;
    flex-direction: column;
    gap: 20px;
  }
  
  .modal-header h3 { margin: 0; font-size: 24px; }
  .modal-hint { color: #666; font-size: 14px; margin-top: 5px; }
  
  /* === 迷你牌阵舞台样式 === */
  .embedded-stage {
    width: 100%;
    height: 240px; /* 预览区高度固定 */
    background: #0a0a0c;
    border-radius: 16px;
    border: 1px solid #222;
    position: relative;
    overflow: hidden;
  }
  
  .spread-wrapper {
    position: absolute;
    top: 50%;
    left: 50%;
    /* 核心：动态宽高 + 居中偏移 + 整体缩放 都在 template 的 inline-style 中计算过了 */
    transform-origin: center center;
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
    width: 20px;
    height: 20px;
    border-radius: 50%;
    font-size: 12px;
    font-weight: bold;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 2px 4px rgba(0,0,0,0.5);
  }
  
  /* === 解析文本区域 === */
  .answer-box {
    width: 100%;
    height: 250px;
    background: #0a0a0c;
    color: #ddd;
    border: 1px solid #222;
    border-radius: 16px;
    padding: 20px;
    font-size: 15px;
    line-height: 1.8;
    text-align: left;
    overflow-y: auto;
    white-space: pre-wrap;
  }
  
  .cursor {
    display: inline-block;
    width: 2px;
    animation: blink 1s step-end infinite;
    color: #fff;
    margin-left: 2px;
  }
  
  @keyframes blink { 0%, 100% { opacity: 1; } 50% { opacity: 0; } }
  
  /* 底部按钮 */
  .confirm-btn { 
    background: #fff; 
    color: #000; 
    padding: 12px 40px;
    border-radius: 12px;
    cursor: pointer;
    font-weight: bold;
    border: none;
    font-size: 16px;
    transition: all 0.3s;
  }
  .confirm-btn:disabled { opacity: 0.5; cursor: not-allowed; }
  .confirm-btn:not(:disabled):hover { transform: scale(1.05); }
  
  /* 动画 */
  .fade-enter-active, .fade-leave-active { transition: opacity 0.3s; }
  .fade-enter-from, .fade-leave-to { opacity: 0; }
  </style>