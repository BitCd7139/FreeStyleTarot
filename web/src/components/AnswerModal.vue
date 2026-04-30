<template>
  <Transition name="fade">
    <div v-if="showModal" class="modal-overlay">
      <div class="modal-content answer-modal">
        
        <div class="modal-header">
          <h3>命运的指引</h3>
          <p class="modal-hint">塔罗的解答已在此显现</p>
        </div>
        
        <!-- 增加 ref 用于截图 -->
        <div class="modal-body" ref="captureArea">
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
          <div class="answer-box" ref="answerBoxRef" @scroll="handleScroll">
              <div class="markdown-body" v-html="parsedAnswer"></div>
          </div>
        </div>

        <!-- 找到 <div class="modal-footer"> 这一块，修改为如下结构 -->
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
  // 由父组件传入，标识后端 SSE 是否还在输出
  isStreaming: { type: Boolean, default: false } 
});

const answerBoxRef = ref(null);
const captureArea = ref(null);
const userScrolledUp = ref(false); // 记录用户是否主动向上滚动了

// 解析 Markdown 与光标
const parsedAnswer = computed(() => {
  const rawText = props.isStreaming 
    ? props.fullAnswer + '<span class="cursor">|</span>' 
    : props.fullAnswer;

  const rawHtml = marked.parse(rawText);
  return DOMPurify.sanitize(rawHtml, { ADD_ATTR: ['class'] });
});

// 监听答案变化
watch(() => props.fullAnswer, async () => {
  if (!answerBoxRef.value) return;
  
  // 只有当用户没有主动向上翻看时，才自动滚动到底部
  if (!userScrolledUp.value) {
    await nextTick();
    scrollToBottom();
  }
});

// 监听滚动事件，判断用户是否主动往上滑了
const handleScroll = (e) => {
  const { scrollTop, scrollHeight, clientHeight } = e.target;
  // 距离底部超过 50px 视为用户主动向上看，暂停自动滚动
  if (scrollHeight - scrollTop - clientHeight > 50) {
    userScrolledUp.value = true;
  } else {
    // 滚回底部了，恢复自动滚动
    userScrolledUp.value = false; 
  }
};

const scrollToBottom = () => {
  if (answerBoxRef.value) {
    answerBoxRef.value.scrollTop = answerBoxRef.value.scrollHeight;
  }
};

// 复制功能
const copyAnswer = async () => {
  try {
    // 去除 HTML 标签，只复制纯文本
    const tempDiv = document.createElement('div');
    tempDiv.innerHTML = parsedAnswer.value;
    const plainText = tempDiv.innerText || tempDiv.textContent;
    
    await navigator.clipboard.writeText(plainText);
    alert('解析内容已复制到剪贴板'); // 可替换为你的 Message 组件
  } catch (err) {
    console.error('复制失败', err);
  }
};

// 截图功能
const saveScreenshot = async () => {
  // 1. 获取要截图的原始元素 (即包含标题、牌阵、文字的 modal-content)
  const originalElement = document.querySelector('.answer-modal');
  if (!originalElement) return;

  try {
    // 2. 创建一个克隆节点，用于在“幕后”进行样式重置
    const clone = originalElement.cloneNode(true);
    
    // 3. 应用“长图专用样式”
    Object.assign(clone.style, {
      position: 'absolute',
      top: '-9999px',
      left: '0',
      width: originalElement.offsetWidth + 'px', // 保持宽度一致
      height: 'auto',       // 高度根据内容自适应
      maxHeight: 'none',    // 去除高度限制
      overflow: 'visible',  // 允许溢出内容可见
    });

    // 4. 特别处理克隆节点中的滚动区域
    const clonedAnswerBox = clone.querySelector('.answer-box');
    if (clonedAnswerBox) {
      Object.assign(clonedAnswerBox.style, {
        height: 'auto',      // 展开全部文字
        maxHeight: 'none',   // 去除限制
        overflow: 'visible', // 允许溢出
        border: 'none'       // 可选：截图时去掉文字框内边框，更像一张卡片
      });
    }

    // 5. 隐藏克隆节点中的按钮（用户通常不希望截图中包含“关闭”按钮）
    const clonedFooter = clone.querySelector('.modal-footer');
    if (clonedFooter) clonedFooter.style.display = 'none';

    // 6. 将克隆节点挂载到 body，否则无法渲染
    document.body.appendChild(clone);

    // 7. 执行截图
    const canvas = await html2canvas(clone, {
      backgroundColor: '#151518', // 保持背景色
      scale: 2,                  // 高清倍率
      useCORS: true,             // 允许跨域图片
      logging: false,
      windowHeight: clone.scrollHeight // 关键：告诉渲染引擎高度是全长的
    });

    // 8. 移除克隆节点
    document.body.removeChild(clone);

    // 9. 下载图片
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
/* 0. 基础全局修正：确保所有元素的宽度计算包含 padding 和 border */
* {
  box-sizing: border-box;
}

/* 背景蒙版 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.85); /* 稍微加深 */
  backdrop-filter: blur(10px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 2000;
  padding: 20px; /* 移动端边距 */
}

/* 弹窗主体 */
.modal-content.answer-modal {
  background: #151518;
  padding: 30px;
  border-radius: 20px;
  width: 900px;
  max-width: 100%; /* 确保不超出屏幕 */
  height: 90vh;
  max-height: 900px;
  border: 1px solid rgba(212, 175, 55, 0.4);
  box-shadow: 0 10px 50px rgba(0,0,0,0.9), 0 0 20px rgba(212, 175, 55, 0.1);
  display: flex;
  flex-direction: column;
  overflow: hidden; /* 关键：防止整体内容溢出 */
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

/* Modal Body 负责包含牌阵和文字 */
.modal-body {
  display: flex;
  flex-direction: column;
  flex: 1; 
  overflow: hidden; /* 极其重要：确保内部滚动条生效 */
  gap: 20px;
  min-height: 0; /* 修复 Flexbox 内部溢出的经典 Bug */
}

/* 牌阵容器 */
.stage-container {
  flex-shrink: 0;
  width: 100%;
  height: 240px;
  background: #0a0a0c;
  border-radius: 12px;
  border: 1px solid rgba(212, 175, 55, 0.3);
  position: relative;
  overflow: hidden; /* 防止牌阵溢出 */
  display: flex;
  justify-content: center;
  align-items: center;
}

/* 解析文本区域 */
.answer-box {
  flex: 1;
  background: #0a0a0c;
  color: #e0e0e0;
  border: 1px solid rgba(212, 175, 55, 0.3);
  border-radius: 12px;
  padding: 24px;
  font-size: 16px;
  line-height: 1.8;
  overflow-y: auto;
  word-break: break-word; /* 强制长单词换行，防止溢出 */
}

/* === 底部按钮组：修复按钮短的问题 === */
.modal-footer {
  flex-shrink: 0;
  display: flex;
  width: 100%; /* 确保容器撑满 */
  gap: 12px;   /* 按钮间距 */
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

/* 统一按钮样式 */
.modal-footer button {
  flex: 1;               /* 核心：让所有按钮平分宽度 */
  display: flex;         /* 开启 flex 让图标和文字居中 */
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px 0;       /* 增加上下高度 */
  border-radius: 10px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;   /* 防止按钮内文字换行 */
}

/* 次要按钮样式 */
.btn-secondary {
  background: rgba(255, 255, 255, 0.05);
  color: #D4AF37;
  border: 1px solid rgba(212, 175, 55, 0.5);
}

/* 主要按钮样式 */
.btn-primary {
  background: linear-gradient(135deg, #D4AF37, #AA7700);
  color: #000;
  border: none;
}

/* 按钮悬停与禁用 */
button:not(:disabled):hover {
  transform: translateY(-2px);
  filter: brightness(1.2);
}

button:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  filter: grayscale(1);
}

/* 滚动条美化 */
.answer-box::-webkit-scrollbar { width: 6px; }
.answer-box::-webkit-scrollbar-thumb { background: rgba(212, 175, 55, 0.3); border-radius: 3px; }

/* 动画 */
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>