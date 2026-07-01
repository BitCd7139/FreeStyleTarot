<template>
  <Transition name="slide-right">
    <div v-if="showModal" class="modal-overlay" @mousedown="onOverlayMouseDown" @mouseup="onOverlayMouseUp">
      <div class="modal-content" @mousedown.stop>
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
              <!-- 匹配中 / 询问态 / 预览态 -->
              <template v-if="detectorState === 'asking' || detectorState === 'previewing'">
                <div class="detector-content">
                  <span class="icon">✨</span>
                  <p v-if="detectorState === 'asking'">
                    这是 <strong class="highlight-text">{{ currentMatchTemplate?.name }}</strong> 吗？
                  </p>
                  <p v-if="detectorState === 'previewing'">
                    正在预览：<strong class="highlight-text">{{ currentMatchTemplate?.name }}</strong>，确认使用吗？
                  </p>
                </div>
                <div class="detector-actions">
                  <!-- 询问态：显示预览按钮 -->
                  <button v-if="detectorState === 'asking'" class="btn-mini btn-yes" @click="previewSpread">预览</button>
                  <!-- 预览态：显示确认按钮 -->
                  <button v-if="detectorState === 'previewing'" class="btn-mini btn-yes" @click="acceptSpread">确认</button>
                  
                  <button class="btn-mini btn-no" @click="rejectSpread">换一个</button>
                </div>
              </template>
              
              <!-- 匹配成功态 -->
              <template v-else-if="detectorState === 'accepted'">
                <div class="detector-content">
                  <span class="icon">🔮</span>
                  <p>已应用牌阵模板：<strong class="highlight-text">{{ acceptedSpreadName }}</strong></p>
                </div>
              </template>

              <!-- 未匹配 / 自定义态 -->
              <template v-else-if="detectorState === 'custom'">
                <div class="detector-content">
                  <span class="icon">🌌</span>
                  <p class="custom-text">这可能是你的自定义牌阵，请选择：</p>
                </div>
                <div class="detector-actions">
                  <button class="btn-mini btn-yes" @click="enterManualMode">我来介绍牌阵</button>
                  <button class="btn-mini btn-no" @click="enterFreestyleMode">Freestyle Mode</button>
                </div>
              </template>

              <!-- Freestyle Mode 态 -->
              <template v-else-if="detectorState === 'freestyle'">
                <div class="detector-content">
                  <span class="icon">✨</span>
                  <p>已开启 <strong class="highlight-text">Freestyle Mode</strong>，将通过视觉流动分析解读牌阵</p>
                </div>
                <div class="detector-actions">
                  <button class="btn-mini btn-no" @click="exitFreestyleMode">返回</button>
                </div>
              </template>
            </div>
          </div>

          <!-- 第二部分：自定义牌阵含义（Freestyle Mode 时隐藏） -->
          <div class="form-group spread-definition" v-if="detectorState !== 'freestyle'">
            <div class="card-list">
              <div v-for="(card, index) in drawnCards" :key="card.id" class="card-item">
                <div class="card-preview">
                  <span class="card-order">{{ card.order }}</span>
                  <img v-if="card.name && card.name !== 'back'" :src="getCardUrl(card.name)" :class="{ 'is-reversed': card.isReversed }" alt="tarot" />
                  <CardBack v-else :class="{ 'is-reversed': card.isReversed }" />
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

        <!-- 第三部分：选择角色 -->
        <RoleSelector 
          v-model="selectedModel" 
          :options="roleOptions" 
        />

        <!-- 第四部分：操作按钮 -->
        <div class="modal-footer">
          <div class="modal-btns">
            <button class="cancel-btn" @click="showModal = false" :disabled="isAnalyzingIntent">返回调整</button>
            <button 
              class="confirm-btn" 
              @click="HandleSubmit" 
              :disabled="isSubmitDisabled || cdCount > 0 || isPending || isAnalyzingIntent">
              <template v-if="isAnalyzingIntent">正在理解问题…</template>
              <template v-else-if="isPending">提交中...</template>
              <template v-else-if="cdCount > 0">已提交 ({{ cd }}s)</template>
              <template v-else>确认提交</template>
            </button>
          </div>
        </div>

        <div v-if="isAnalyzingIntent" class="intent-loading-overlay" aria-live="polite">
          <span class="intent-spinner" aria-hidden="true"></span>
          <p>正在理解你的问题…</p>
        </div>
      </div>
    </div>
  </Transition>

  <IntentClarifyModal
    v-model:showModal="showClarifyModal"
    :questions="clarifyQuestions"
    :intent-summary="clarifyIntentSummary"
    @confirm="onClarifyConfirm"
    @skip="onClarifyConfirm"
    @cancel="onClarifyCancel"
  />
</template>
  
<script setup>
import { defineModel, defineProps, defineEmits, ref, watch, computed } from 'vue';
import MiniTarot from './MiniTarot.vue';
import CardBack from './CardBack.vue';

// 引入解耦后的算法
import { discretizeCards } from '../utils/cardGrid.js';
import { SPREAD_TEMPLATES } from '../spread/index.js';
import { rateLimiter, isVipUser } from '../utils/rateLimiter.js';
import { getName } from '../utils/cardInfo.js';
import { buttonCooldown } from '../utils/buttonCooldown.js';
import RoleSelector from './RoleSelector.vue';
import IntentClarifyModal from './IntentClarifyModal.vue';
import { useAuth } from '../composables/useAuth.js';
import { useCustomApi } from '../composables/useCustomApi.js';
import { useTarotFlow } from '../composables/useTarotFlow.js';
import { fetchClarify } from '../utils/clarifyApi.js';

const { user, getToken } = useAuth();
const { getConfig: getCustomApiConfig, isEnabled: customApiEnabled } = useCustomApi();
const { guidedMode, selectedSpread } = useTarotFlow();

// 定义角色卡字典
const roleOptions = {
  "辉夜姬": "neet",
  "白月光": "mate",
  "女仆猫娘": "neko",
  "雌小鬼": "zako",
  "馆长": "keeper",
};

// 选中的模型值
const selectedModel = ref(roleOptions["馆长"]);

const { count: cdCount, isPending, start: startCD, stop: resetCD } = buttonCooldown(3);
const errors = ref({
  question: false,
  meanings: []
});

const showModal = defineModel('showModal');
const question = defineModel('question');
const drawnCards = defineModel('drawnCards'); 

const props = defineProps({
  isSubmitDisabled: Boolean,
  cardWidth: { type: Number, default: 120 },
  cardHeight: { type: Number, default: 210 }
});

const emit = defineEmits(['submit']);

const isAnalyzingIntent = ref(false);
const showClarifyModal = ref(false);
const clarifyQuestions = ref([]);
const clarifyIntentSummary = ref('');
const pendingSubmitInfo = ref('Result');
const pendingIntentSummary = ref('');

// 状态机管理：新增 'previewing' 和 'freestyle' 状态
const detectorState = ref('custom'); 
const submitInfo = ref('Result'); // 使用 ref 保证响应式
const possibleMatches = ref([]);
const currentMatchIndex = ref(0);
const backupMeanings = ref(new Map());

// Freestyle Mode 状态
const freestyleMode = ref(false);
 
const currentMatchTemplate = computed(() => possibleMatches.value[currentMatchIndex.value] || null);

// 引导模式下直接使用新手选定的牌阵名，跳过识别
const acceptedSpreadName = computed(() => {
  if (guidedMode.value && selectedSpread.value) return selectedSpread.value.name;
  return currentMatchTemplate.value?.name || '';
});

// 触发模式识别
const runSpreadDetection = () => {
  if (!drawnCards.value || drawnCards.value.length === 0) {
    detectorState.value = 'custom';
    return;
  }

  // 新手引导模式：直接使用已选定的牌阵，跳过识别与确认
  if (guidedMode.value && selectedSpread.value) {
    // freestyle 牌阵：不填充含义，直接进入自由流式模式
    if (selectedSpread.value.freestyle) {
      freestyleMode.value = true;
      detectorState.value = 'freestyle';
      return;
    }
    const slots = selectedSpread.value.slots || [];
    drawnCards.value.forEach(card => {
      if (card.slotIndex !== undefined && card.slotIndex !== null && slots[card.slotIndex]) {
        card.meaning = slots[card.slotIndex].meaning;
      }
    });
    detectorState.value = 'accepted';
    return;
  }

  const gridCards = discretizeCards(drawnCards.value, props.cardWidth, props.cardHeight, 12);
  
  const matches = [];
  for (const tpl of SPREAD_TEMPLATES) {
    if (tpl.cardCount !== gridCards.length) continue;
    
    const mapping = tpl.match(gridCards);
    if (mapping) {
      matches.push({
        name: tpl.name,
        mapping: mapping 
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

// 点击预览：将牌阵含义填入输入框供预览
const previewSpread = () => {
  const matchObj = currentMatchTemplate.value;
  if (matchObj && matchObj.mapping) {
    // 备份用户当前的输入，防止预览覆盖后无法找回
    backupMeanings.value.clear();
    drawnCards.value.forEach(card => {
      backupMeanings.value.set(card.id, card.meaning || '');
    });

    // 填入计算出的牌阵含义
    drawnCards.value.forEach(card => {
      const mappingItem = matchObj.mapping.find(m => m.id === card.id);
      if (mappingItem) {
        card.meaning = mappingItem.meaning;
      }
    });
    
    // 进入预览状态
    detectorState.value = 'previewing';
  }
};

// 同意应用牌阵
const acceptSpread = () => {
  // 因为 previewSpread 已经把含义赋值给 card 了，直接修改状态即可
  detectorState.value = 'accepted';
};

// 拒绝该牌阵，尝试下一个匹配或进入自定义
const rejectSpread = () => {
  if (detectorState.value === 'previewing') {
    drawnCards.value.forEach(card => {
      if (backupMeanings.value.has(card.id)) {
        card.meaning = backupMeanings.value.get(card.id);
      }
    });
  }

  if (currentMatchIndex.value < possibleMatches.value.length - 1) {
    currentMatchIndex.value++;
    detectorState.value = 'asking';
  } else {
    detectorState.value = 'custom';
  }
};

// 进入手动模式（我来介绍牌阵）
const enterManualMode = () => {
  freestyleMode.value = false;
};

// 进入 Freestyle Mode
const enterFreestyleMode = () => {
  freestyleMode.value = true;
  detectorState.value = 'freestyle';
};

// 退出 Freestyle Mode
const exitFreestyleMode = () => {
  freestyleMode.value = false;
  detectorState.value = 'custom';
};

// 提取卡牌的拓扑特征值（只有当卡牌的 id 或坐标发生变化时，才重新计算牌阵）
const cardTopologyData = computed(() => {
  if (!drawnCards.value) return '';
  return drawnCards.value.map(c => `${c.id}_${c.x}_${c.y}`).join('|');
});

// 防止输入 meaning 时导致牌阵识别被重置
watch(showModal, (newVal) => {
  if (newVal) {
    runSpreadDetection();
    errors.value = { question: false, meanings: [] };
  }
});

// 仅在弹窗开启且卡牌位置真正改变时，才重新识别牌阵
watch(cardTopologyData, (newVal, oldVal) => {
  if (showModal.value && newVal !== oldVal) {
    runSpreadDetection();
  }
});

const buildClarifyPayload = () => ({
  question: question.value.trim(),
  freestylemode: freestyleMode.value,
  cards: drawnCards.value.map((card) => ({
    order: card.order,
    name: card.name,
    x: card.x,
    y: card.y,
    orientation: card.isReversed ? 'reversed' : 'upright',
    meaning: card.meaning.trim(),
  })),
  custom_api: customApiEnabled.value ? getCustomApiConfig() : undefined,
});

const finalizeSubmit = (clarifications = []) => {
  submitInfo.value = pendingSubmitInfo.value;
  startCD();
  if (pendingSubmitInfo.value === 'Result') {
    rateLimiter.recordSubmission();
  }
  emit('submit', {
    clarifications,
    intentSummary: pendingIntentSummary.value,
    freestylemode: freestyleMode.value,
  });
};

const onClarifyConfirm = (clarifications) => {
  finalizeSubmit(clarifications);
};

const onClarifyCancel = () => {
  // 用户取消澄清，不消耗限流
};

const HandleSubmit = async () => {
  // 自定义 API 模式下跳过前端限流
  if (!customApiEnabled.value) {
    const limitStatus = rateLimiter.checkLimit({ isVip: isVipUser(user.value) });
    if (!limitStatus.allowed) {
      alert(limitStatus.message);
      return;
    }
    pendingSubmitInfo.value = limitStatus.submitInfo;
  } else {
    pendingSubmitInfo.value = 'Result';
  }

  const qLen = question.value ? question.value.trim().length : 0;
  errors.value.question = qLen < 5 || qLen > 500;
  
  let hasMeaningError = false;
  if (!freestyleMode.value) {
    errors.value.meanings = drawnCards.value.map((card) => {
      return !card.meaning || card.meaning.trim().length === 0;
    });
    hasMeaningError = errors.value.meanings.some((err) => err);
  }

  if (errors.value.question || hasMeaningError) {
    return;
  }

  isAnalyzingIntent.value = true;
  try {
    const result = await fetchClarify(buildClarifyPayload(), getToken());
    if (result.needs_clarification && result.questions?.length > 0) {
      pendingIntentSummary.value = result.intent_summary || '';
      clarifyQuestions.value = result.questions;
      clarifyIntentSummary.value = result.intent_summary || '';
      showClarifyModal.value = true;
      return;
    }
    pendingIntentSummary.value = result.intent_summary || '';
    finalizeSubmit([]);
  } catch (err) {
    console.error(err);
    const proceed = confirm(
      (err.message || '意图分析失败') + '\n\n是否仍使用当前问题直接开始占卜？'
    );
    if (proceed) {
      finalizeSubmit([]);
    }
  } finally {
    isAnalyzingIntent.value = false;
  }
};
defineExpose({
  unlockSubmit: resetCD,
  selectedModel,
  submitInfo
});


const getCardUrl = (name) => {
  return new URL(`../assets/tarots/${name}.jpeg`, import.meta.url).href;
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
  overflow-y: auto;      /* 内容多了就在内部垂直滚动 */
  overflow-x: hidden;    /* 彻底杜绝左右滚动 */
  padding: 20px;         /* 桌面端内边距 */

  /* 针对 Chrome, Safari, Edge (WebKit) */
  &::-webkit-scrollbar {
    display: none;
  }

  /* 针对 Firefox */
  scrollbar-width: none;

  /* 针对 IE 和 老版 Edge */
  -ms-overflow-style: none;
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

.intent-loading-overlay {
  position: absolute;
  inset: 0;
  background: rgba(10, 10, 12, 0.75);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 14px;
  z-index: 10;
  border-radius: inherit;
}

.intent-loading-overlay p {
  margin: 0;
  color: #e5d8b0;
  font-size: 14px;
}

.intent-spinner {
  width: 36px;
  height: 36px;
  border: 3px solid rgba(194, 163, 95, 0.25);
  border-top-color: #c2a35f;
  border-radius: 50%;
  animation: intentSpin 0.8s linear infinite;
}

@keyframes intentSpin {
  to { transform: rotate(360deg); }
}

.modal-content {
  position: relative;
}
</style>
