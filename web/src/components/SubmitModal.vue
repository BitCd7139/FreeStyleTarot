<template>
  <Transition name="slide-right">
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3>冥想与牌阵确认</h3>
          <p class="modal-hint">请看着你的牌阵，定义它们的含义</p>
        </div>
        
        <div class="modal-body">
          <!-- 迷你牌阵预览组件 -->
          <MiniTarot 
            :drawnCards="drawnCards"
            :cardWidth="cardWidth"
            :cardHeight="cardHeight"
            :containerWidth="370"  
            :containerHeight="180" 
          />

          <!-- 第一部分：问题输入 -->
          <div class="form-group">
            <textarea 
              v-model="question" 
              :class="{ 'error-border': errors.question }"
              placeholder="请在此输入你的困惑，例如：我近期的事业运势如何？（5-500字）"
              minlength="5"
              maxlength="500"
            ></textarea>
            <span v-if="errors.question" class="error-text">问题长度需在 5 到 500 个字符之间</span>
          </div>

          <!-- 新增：牌阵自动识别区域 -->
          <div class="form-group spread-detector" v-if="drawnCards && drawnCards.length > 0">
            <div class="detector-box">
            <!-- 匹配中 / 询问态 -->
            <template v-if="detectorState === 'asking'">
              <div class="detector-content">
                <span class="icon">✨</span>
                <p>这是 <strong class="highlight-text">{{ currentMatchTemplate?.name }}</strong> 吗？</p>
              </div>
              <div class="detector-actions">
                <!-- 应用了统一的按钮样式 -->
                <button class="btn-mini btn-yes" @click="acceptSpread">确认</button>
                <button class="btn-mini btn-no" @click="rejectSpread">换一个</button>
              </div>
            </template>
            
            <!-- 匹配成功态 -->
            <template v-else-if="detectorState === 'accepted'">
              <div class="detector-content">
                <span class="icon">🔮</span>
                <p>已应用牌阵模板：<strong class="highlight-text">{{ currentMatchTemplate?.name }}</strong></p>
              </div>
            </template>

            <!-- 未匹配 / 自定义态 -->
            <template v-else-if="detectorState === 'custom'">
              <div class="detector-content">
                <span class="icon">🌌</span>
                <p class="custom-text">这可能是你的自定义牌阵，可以介绍一下吗？</p>
              </div>
            </template>
          </div>
        </div>

          <!-- 第二部分：自定义牌阵含义 -->
          <div class="form-group spread-definition">
            <div class="card-list">
              <div v-for="(card, index) in drawnCards" :key="card.id" class="card-item">
                <div class="card-preview">
                  <span class="card-order">{{ card.order }}</span>
                  <img :src="getCardUrl(card.name)" :class="{ 'is-reversed': card.isReversed }" alt="tarot" />
                </div>
                <div class="card-input">
                  <div class="card-info-header">
                    <span class="orientation-tag">{{ getName(card.name) + " "}}</span>
                    <span class="orientation-tag" :class="{ 'rev-text': card.isReversed }">
                      {{ card.isReversed ? '逆位' : '正位' }}
                    </span>
                  </div>
                  <input 
                    type="text" 
                    v-model="card.meaning" 
                    maxlength="10"
                    :class="{ 'error-border': errors.meanings[index] }"
                    :placeholder="`定义第 ${card.order} 张牌 (如: 过去)`"
                  />
                  <span v-if="errors.meanings[index]" class="error-text">含义不能为空且最多10字</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 第三部分：操作按钮 -->
        <div class="modal-footer">
          <div class="modal-btns">
            <button class="cancel-btn" @click="showModal = false">返回调整</button>
            <button class="confirm-btn" @click="submitToBackend" :disabled="isSubmitDisabled">
              确认提交
            </button>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>
  
  <script setup>
import { defineModel, defineProps, defineEmits, ref, watch, computed } from 'vue';
import MiniTarot from './MiniTarot.vue';

// 引入解耦后的算法
import { discretizeCards } from '../utils/cardGrid.js';
import { SPREAD_TEMPLATES } from '../formation/index.js';
import { rateLimiter } from '../utils/rateLimiter.js';  
import { getName } from '../utils/cardDict.js';

const errors = ref({
  question: false,
  meanings: [] // 存储每张牌是否有错的布尔值数组
});

const showModal = defineModel('showModal');
const question = defineModel('question');
const drawnCards = defineModel('drawnCards'); 

const props = defineProps({
  isSubmitDisabled: Boolean,
  cardWidth: { type: Number, default: 120 },
  cardHeight: { type: Number, default: 210 }
});

const emit = defineEmits(['submitToBackend']);

// 状态机管理
const detectorState = ref('custom'); 
const possibleMatches = ref([]);
const currentMatchIndex = ref(0);

const currentMatchTemplate = computed(() => possibleMatches.value[currentMatchIndex.value] || null);

// 触发模式识别
const runSpreadDetection = () => {
  if (!drawnCards.value || drawnCards.value.length === 0) {
    detectorState.value = 'custom';
    return;
  }

  // 1. 转为 12x12 网格矩阵
  //console.log('原始卡牌数据:', drawnCards.value);
  const gridCards = discretizeCards(drawnCards.value, props.cardWidth, props.cardHeight, 12);
  //console.log('网格化卡牌数据:', gridCards);
  
  // 2. 遍历所有牌阵模板，只要 match() 返回非 null 的映射表，就认为是潜在匹配
  const matches = [];
  for (const tpl of SPREAD_TEMPLATES) {
    // 基础过滤：卡牌数量必须一致
    if (tpl.cardCount !== gridCards.length) continue;
    
    // 执行高阶拓扑匹配算法
    const mapping = tpl.match(gridCards);
    if (mapping) {
      matches.push({
        name: tpl.name,
        mapping: mapping // 将返回的 {id, meaning} 映射表存下来
      });
    }
  }

  if (matches.length > 0) {
    possibleMatches.value = matches;
    currentMatchIndex.value = 0;
    detectorState.value = 'asking';
  } else {
    possibleMatches.value = [];
    detectorState.value = 'custom';
  }
};

// 同意应用牌阵
const acceptSpread = () => {
  const matchObj = currentMatchTemplate.value;
  if (matchObj && matchObj.mapping) {
    // 根据算法计算出的拓扑位置映射表，更新页面上卡牌的含义，与抽卡顺序无关
    drawnCards.value.forEach(card => {
      const mappingItem = matchObj.mapping.find(m => m.id === card.id);
      if (mappingItem) {
        card.meaning = mappingItem.meaning;
      }
    });
  }
  detectorState.value = 'accepted';
};

// 拒绝该牌阵，尝试下一个匹配或进入自定义
const rejectSpread = () => {
  if (currentMatchIndex.value < possibleMatches.value.length - 1) {
    currentMatchIndex.value++;
  } else {
    detectorState.value = 'custom';
  }
};

// 监听弹窗与卡牌坐标的变动
watch([showModal, () => drawnCards.value], ([newShowModal]) => {
  if (newShowModal) {
    runSpreadDetection();
  }
}, { deep: true });

const submitToBackend = () => {
  const qLen = question.value ? question.value.trim().length : 0;
  errors.value.question = qLen < 5 || qLen > 500;
  errors.value.meanings = drawnCards.value.map(card => {
    return !card.meaning || card.meaning.trim().length === 0;
  });

  const hasMeaningError = errors.value.meanings.some(err => err);
  if (errors.value.question || hasMeaningError) {
    return;
  }

  //const limitStatus = rateLimiter.checkLimit();
  //console.log('提交检查结果:', limitStatus);
  // if (!limitStatus.allowed) {
  //   alert(limitStatus.message);
  //   return;
  // }
  emit('submitToBackend');
  //rateLimiter.recordSubmission();
};

watch(showModal, (newVal) => {
  if (newVal) {
    errors.value = { question: false, meanings: [] };
  }
});

const getCardUrl = (name) => {
  if (name === 'back') return new URL(`../assets/tarots/back.jpeg`, import.meta.url).href;
  return new URL(`../assets/tarots/${name}.jpeg`, import.meta.url).href;
};
</script>

<style scoped>
/* 错误边框 */
.error-border {
  border: 1.5px solid #ff4d4f;
  background-color: #000000;
  color: #ffffff
}

/* 错误文字提示 */
.error-text {
  color: #ff4d4f;
  font-size: 12px;
  margin-top: 4px;
  display: block;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-5px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 调整 textarea 的高度以适应错误提示 */
.form-group textarea {
  width: 100%;            
  box-sizing: border-box;  
  min-height: 100px;       
  max-width: 100%;         
  resize: vertical;       
  
  padding: 12px;
  border-radius: 8px;
  font-size: 14px;
  line-height: 1.5;
  margin-bottom: 4px;
}

/* 调整输入框容器 */
.card-input {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.modal-body {
  flex: 1;
  overflow-y: hidden;      /* 内容多了就在内部垂直滚动 */
  overflow-x: hidden;    /* 彻底杜绝左右滚动 */
  padding: 20px;         /* 桌面端内边距 */
}

.modal-btns {
  display: flex;
  gap: 12px;       /* 两个按钮之间的间距 */
  width: 100%;
}

/* 让两个按钮都具备 flex: 1，它们会自动平分剩余空间 */
.cancel-btn, .confirm-btn {
  flex: 1; 
  padding: 12px 0; /* 增加上下内边距 */
  border-radius: 8px;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
}

/* 5. 移动端适配：减少内边距，让内容更宽 */
@media (max-width: 600px) {
  .modal-overlay {
    padding: 0;         /* 移动端去掉最外层 padding，改为由 content 控制 */
  }
  
  .modal-content {
    width: 100%;
    height: 100%;       /* 移动端全屏或半全屏 */
    max-height: 100vh;
    border-radius: 0;   /* 手机上通常不需要圆角，或者只给顶部圆角 */
  }

  .modal-body {
    padding: 15px;      /* 手机端内边距缩小，留出更多空间给内容 */
  }
}
</style>
