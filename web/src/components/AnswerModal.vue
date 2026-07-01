<template>
  <Transition name="fade">
    <div v-if="showModal" class="modal-overlay" @mousedown="onOverlayMouseDown" @mouseup="onOverlayMouseUp">
      <div class="modal-content answer-modal" @mousedown.stop>
        
        <div class="modal-header">
          <h3> FreeStyleTarot | 塔罗解析 </h3>
        </div>
        
        <!--
          现在牌阵和文字都在同一个滚动区域里，向下滚动时牌阵会被推走。
        -->
        <div class="modal-body" ref="scrollContainer" @scroll="handleScroll">
          
          <div class="scroll-content" ref="captureArea">
            <!-- 现场渲染的迷你嵌套牌阵 - now integrated without borders -->
            <div class="stage-container-integrated">
              <MiniTarot 
                :drawnCards="drawnCards" 
                :cardWidth="cardWidth" 
                :cardHeight="cardHeight" 
              />
            </div>

            <!-- 解析展示区域 -->
            <div class="answer-box">
              <template v-if="displayBlocks.length">
                <div
                  v-for="(block, index) in displayBlocks"
                  :key="block.key"
                  class="stream-block"
                >
                  <div v-if="block.loading" class="phase-loading">
                    <span class="phase-spinner" aria-hidden="true"></span>
                    <span>{{ block.label }}</span>
                  </div>
                  <div
                    v-if="block.text"
                    class="markdown-body"
                    v-html="renderBlock(block, index)"
                  ></div>
                </div>
              </template>
              <div v-else class="markdown-body" v-html="parsedLegacyAnswer"></div>
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
  streamBlocks: { type: Array, default: () => [] },
  drawnCards: { type: Array, default: () => [] },
  cardWidth: { type: Number, default: 120 },
  cardHeight: { type: Number, default: 210 },
  isStreaming: { type: Boolean, default: false },
});

const scrollContainer = ref(null);
const captureArea = ref(null);
const userScrolledUp = ref(false);

const displayBlocks = computed(() => props.streamBlocks.filter((b) => b.loading || b.text));

const hasActiveLoading = computed(() => props.streamBlocks.some((b) => b.loading));

function renderMarkdown(text) {
  // Fix double bullet prefixes: "- - text" → "- text", "  - - text" → "  - text"
  const cleaned = text.replace(/^(\s*)- - /gm, '$1- ');
  const rawHtml = marked.parse(cleaned);
  return DOMPurify.sanitize(rawHtml, { ADD_ATTR: ['class'] });
}

function renderBlock(block, index) {
  const isLast = index === displayBlocks.value.length - 1;
  const showCursor = props.isStreaming && isLast && !hasActiveLoading.value && !block.loading;
  const rawText = showCursor ? block.text + '<span class="cursor">|</span>' : block.text;
  return renderMarkdown(rawText);
}

const parsedLegacyAnswer = computed(() => {
  const showCursor = props.isStreaming && !props.fullAnswer;
  const rawText = showCursor
    ? props.fullAnswer + '<span class="cursor">|</span>'
    : props.fullAnswer;
  return renderMarkdown(rawText);
});

const copySourceHtml = computed(() => {
  if (displayBlocks.value.length) {
    return displayBlocks.value.map((block, index) => renderBlock(block, index)).join('');
  }
  return parsedLegacyAnswer.value;
});

watch([() => props.fullAnswer, () => props.streamBlocks], async () => {
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
    tempDiv.innerHTML = copySourceHtml.value;
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

let mouseDownOutside = false;
const onOverlayMouseDown = () => {
  mouseDownOutside = true;
};
const onOverlayMouseUp = () => {
  if (mouseDownOutside) {
    showModal.value = false;
  }
  mouseDownOutside = false;
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
  padding: 0;
}

/* 弹窗主体 */
.modal-content.answer-modal {
  background: #151518;
  padding: 0 30px;
  border-radius: 0;
  width: 100%;
  max-width: 900px;
  height: 100%;
  max-height: 100vh;
  border: none;
  border-left: 1px solid rgba(212, 175, 55, 0.4);
  border-right: 1px solid rgba(212, 175, 55, 0.4);
  box-shadow: 0 0 50px rgba(0,0,0,0.9), -20px 0 30px rgba(0,0,0,0.6), 20px 0 30px rgba(0,0,0,0.6);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.modal-header {
  text-align: center;
  margin-bottom: 20px;
  padding-top: 20px;
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

/* 牌阵容器 - integrated without borders */
.stage-container-integrated {
  width: 100%;
  background: transparent;
  border: none;
  border-radius: 0;
  padding: 0;
  margin: 0 0 16px 0;
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: visible;
  flex: 0 0 auto;
  height: auto;
}

/* 解析文本区域随内容自然撑开 */
.answer-box {
  background: transparent;
  color: #e0e0e0;
  border: none;
  border-radius: 0;
  padding: 0;
  font-size: 16px;
  line-height: 1.8;
  word-break: break-word; 
  flex: 0 0 auto;
  height: auto;
}

.stream-block + .stream-block {
  margin-top: 16px;
}

.phase-loading {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 16px;
  padding: 12px 16px;
  color: #D4AF37;
  font-size: 15px;
  background: rgba(212, 175, 55, 0.06);
  border-radius: 8px;
  border: 1px solid rgba(212, 175, 55, 0.2);
}

.phase-spinner {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(212, 175, 55, 0.25);
  border-top-color: #D4AF37;
  border-radius: 50%;
  animation: phase-spin 0.8s linear infinite;
  flex-shrink: 0;
}

@keyframes phase-spin {
  to { transform: rotate(360deg); }
}

.markdown-body :deep(.cursor) {
  color: #D4AF37;
  animation: blink 1s step-end infinite;
}

@keyframes blink {
  50% { opacity: 0; }
}

.modal-footer {
  flex-shrink: 0;
  display: flex;
  width: 100%;
  gap: 12px;
  margin-top: 20px;
  padding-top: 20px;
  padding-bottom: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  background: #151518;
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
  /* 1. 蒙版无留白，全覆盖 */
  .modal-overlay {
    padding: 0;
  }

  /* 2. 弹窗全覆盖，仅保留左右细线 */
  .modal-content.answer-modal {
    padding: 0 12px;
    width: 100%;
    height: 100%;
    max-width: 100%;
    max-height: 100vh;
    border-radius: 0;
    border: none;
    border-left: 1px solid rgba(212, 175, 55, 0.4);
    border-right: 1px solid rgba(212, 175, 55, 0.4);
  }

  .modal-header {
    margin-bottom: 12px;
  }
  .modal-header h3 { font-size: 18px; }
  
  .scroll-content { gap: 12px; }

  /* 3. 牌阵高度改为自适应缩小，不再傻大个 */
  .stage-container-integrated {
    height: auto !important;
  }

  /* 4. 文字框大幅缩减 Padding，留空间给文字本身 */
  .answer-box {
    padding: 0;
    font-size: 14px;
    line-height: 1.6;
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