<template>
  <Transition name="fade">
    <div v-if="showModal" class="modal-overlay">
      <div class="modal-content answer-modal">
        
        <div class="modal-header">
          <h3>命运的指引</h3>
          <p class="modal-hint">塔罗的解答已在此显现</p>
        </div>
        
        <!-- 
          修改点1：将 scroll 监听和 ref 移到最外层 modal-body。
          现在牌阵和文字都在同一个滚动区域里，向下滚动时牌阵会被推走。
        -->
        <div class="modal-body" ref="scrollContainer" @scroll="handleScroll">
          
          <div class="scroll-content" ref="captureArea">
            <!-- 现场渲染的迷你嵌套牌阵 -->
            <div class="stage-container">
              <MiniTarot 
                :drawnCards="drawnCards" 
                :cardWidth="cardWidth" 
                :cardHeight="cardHeight" 
                :containerWidth="720" 
                :containerHeight="240" 
              />
            </div>

            <!-- 解析展示区域 -->
            <div class="answer-box">
                <div class="markdown-body" v-html="parsedAnswer"></div>
            </div>
          </div>

        </div>

        <div class="modal-footer">
          <button class="btn-secondary" @click="saveScreenshot" :disabled="isStreaming">
            <i class="icon-camera"></i> 保存截图
          </button>
          
          <button class="btn-secondary" @click="copyAnswer">
            <i class="icon-copy"></i> 复制解析
          </button>
          
          <button class="btn-primary" @click="closeModal" :disabled="isStreaming">
            {{ isStreaming ? '正在解析...' : '完成并关闭' }}
          </button>
        </div>

      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, watch, computed, nextTick } from 'vue';
import MiniTarot from './MiniTarot.vue';
import { marked } from 'marked';
import DOMPurify from 'dompurify';
import html2canvas from 'html2canvas'; 

const showModal = defineModel('showModal');
const props = defineProps({
  fullAnswer: { type: String, default: '' },
  drawnCards: { type: Array, default: () => [] },
  cardWidth: { type: Number, default: 120 },
  cardHeight: { type: Number, default: 210 },
  isStreaming: { type: Boolean, default: false } 
});

// 指向新的统排滚动容器
const scrollContainer = ref(null);
const captureArea = ref(null);
const userScrolledUp = ref(false); 

const parsedAnswer = computed(() => {
  const rawText = props.isStreaming 
    ? props.fullAnswer + '<span class="cursor">|</span>' 
    : props.fullAnswer;

  const rawHtml = marked.parse(rawText);
  return DOMPurify.sanitize(rawHtml, { ADD_ATTR: ['class'] });
});

watch(() => props.fullAnswer, async () => {
  if (!scrollContainer.value) return;
  
  if (!userScrolledUp.value) {
    await nextTick();
    scrollToBottom();
  }
});

const handleScroll = (e) => {
  const { scrollTop, scrollHeight, clientHeight } = e.target;
  if (scrollHeight - scrollTop - clientHeight > 50) {
    userScrolledUp.value = true;
  } else {
    userScrolledUp.value = false; 
  }
};

const scrollToBottom = () => {
  if (scrollContainer.value) {
    scrollContainer.value.scrollTop = scrollContainer.value.scrollHeight;
  }
};

const copyAnswer = async () => {
  try {
    const tempDiv = document.createElement('div');
    tempDiv.innerHTML = parsedAnswer.value;
    const plainText = tempDiv.innerText || tempDiv.textContent;
    
    await navigator.clipboard.writeText(plainText);
    alert('解析内容已复制到剪贴板'); 
  } catch (err) {
    console.error('复制失败', err);
  }
};

const saveScreenshot = async () => {
  const originalElement = document.querySelector('.answer-modal');
  if (!originalElement) return;

  try {
    const clone = originalElement.cloneNode(true);
    
    Object.assign(clone.style, {
      position: 'absolute',
      top: '-9999px',
      left: '0',
      width: originalElement.offsetWidth + 'px', 
      height: 'auto',       
      maxHeight: 'none',    
      overflow: 'visible',  
    });

    // 截图修复：由于现在滚动条在 modal-body 上，需要把克隆体的 modal-body 完全展开
    const clonedBody = clone.querySelector('.modal-body');
    if (clonedBody) {
      Object.assign(clonedBody.style, {
        height: 'auto',      
        maxHeight: 'none',   
        overflow: 'visible' 
      });
    }

    const clonedFooter = clone.querySelector('.modal-footer');
    if (clonedFooter) clonedFooter.style.display = 'none';

    document.body.appendChild(clone);

    const canvas = await html2canvas(clone, {
      backgroundColor: '#151518', 
      scale: 2,                  
      useCORS: true,             
      logging: false,
      windowHeight: clone.scrollHeight 
    });

    document.body.removeChild(clone);

    const url = canvas.toDataURL('image/png');
    const a = document.createElement('a');
    a.href = url;
    a.download = `塔罗解析-${new Date().toLocaleDateString()}.png`;
    a.click();
    
  } catch (err) {
    console.error('长截图生成失败:', err);
    alert('截图生成失败，请重试');
  }
};

const closeModal = () => {
  showModal.value = false;
};
</script>

<style scoped>
* {
  box-sizing: border-box;
}

/* 背景蒙版 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.85); 
  backdrop-filter: blur(10px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 2000;
  padding: 20px; 
}

/* 弹窗主体 */
.modal-content.answer-modal {
  background: #151518;
  padding: 30px;
  border-radius: 20px;
  width: 900px;
  max-width: 100%; 
  height: 90vh;
  max-height: 1200px;
  border: 1px solid rgba(212, 175, 55, 0.4);
  box-shadow: 0 10px 50px rgba(0,0,0,0.9), 0 0 20px rgba(212, 175, 55, 0.1);
  display: flex;
  flex-direction: column;
  overflow: hidden; 
}

.modal-header {
  text-align: center;
  margin-bottom: 20px;
  flex-shrink: 0;
}

.modal-header h3 {
  margin: 0;
  font-size: 24px;
  color: #D4AF37;
  letter-spacing: 2px;
}
.modal-hint {
  color: #888;
  font-size: 14px;
  margin-top: 5px;
}

/* Modal Body */
.modal-body {
  flex: 1; 
  overflow-y: auto; /* 开启整体垂直滚动 */
  overflow-x: hidden;
  min-height: 0; 
  padding-right: 8px; /* 留给滚动条的空间 */
}

/* 统排内容包裹层 */
.scroll-content {
  display: flex;
  flex-direction: column;
  gap: 20px; /* 牌阵和文字的间距 */
}

/* 牌阵容器 */
.stage-container {
  flex-shrink: 0;
  width: 100%;
  background: #0a0a0c;
  border-radius: 12px;
  border: 1px solid rgba(212, 175, 55, 0.3);
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow-x: hidden;
  flex: 0 0 auto;
  height: auto;
}

/* 解析文本区域随内容自然撑开 */
.answer-box {
  background: #0a0a0c;
  color: #e0e0e0;
  border: 1px solid rgba(212, 175, 55, 0.3);
  border-radius: 12px;
  padding: 24px;
  font-size: 16px;
  line-height: 1.8;
  word-break: break-word; 
  flex: 0 0 auto;
  height: auto;
}

.modal-footer {
  flex-shrink: 0;
  display: flex;
  width: 100%; 
  gap: 12px;   
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-footer button {
  flex: 1;               
  display: flex;         
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px 0;       
  border-radius: 10px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;   
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.05);
  color: #D4AF37;
  border: 1px solid rgba(212, 175, 55, 0.5);
}

.btn-primary {
  background: linear-gradient(135deg, #D4AF37, #AA7700);
  color: #000;
  border: none;
}

button:not(:disabled):hover {
  transform: translateY(-2px);
  filter: brightness(1.2);
}

button:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  filter: grayscale(1);
}

/* 滚动条美化移到了 modal-body 上 */
.modal-body::-webkit-scrollbar { width: 6px; }
.modal-body::-webkit-scrollbar-thumb { background: rgba(212, 175, 55, 0.3); border-radius: 3px; }

.fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }


/* === 针对移动端的深度压缩优化 === */
@media (max-width: 768px) {
  /* 1. 极致压缩蒙版空白 */
  .modal-overlay {
    padding: 8px; 
  }

  /* 2. 缩小弹窗内边距，让屏幕每一寸都用在刀刃上 */
  .modal-content.answer-modal {
    padding: 16px 12px; 
    height: 96vh;  
    border-radius: 16px;
  }

  .modal-header {
    margin-bottom: 12px;
  }
  .modal-header h3 { font-size: 18px; }
  
  .scroll-content { gap: 12px; }

  /* 3. 牌阵高度改为自适应缩小，不再傻大个 */
  .stage-container {
    height: auto !important; 
    min-height: 140px; /* 保底高度 */
    aspect-ratio: 2 / 1; /* 保持宽高比例缩小 */
  }

  /* 4. 文字框大幅缩减 Padding，留空间给文字本身 */
  .answer-box {
    padding: 12px 10px; /* 紧凑边距 */
    font-size: 14px;
    line-height: 1.6;
    /* 可选：去掉边框，减少视觉阻碍 */
    border: none;
    background: transparent;
  }

  /* 5. 按钮底部适配 */
  .modal-footer {
    flex-wrap: wrap; 
    gap: 8px;
    margin-top: 10px;
    padding-top: 10px;
  }
  .modal-footer button {
    padding: 12px 4px;
    font-size: 14px;
  }
}
</style>