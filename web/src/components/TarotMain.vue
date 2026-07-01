<template>
    <div 
      class="tarot-container" 
      :style="backgroundStyle"
      @pointermove="handlepointerMove" 
      @pointerup="handlepointerUp"
      @pointerleave="handlepointerUp"
      @pointercancel="handlepointerUp"
    >
      <!-- 左侧工具栏：自选卡牌 + 垃圾桶 -->
      <div class="left-toolbar">
        <CardPickerButton :isOpen="showCardPicker" @open="showCardPicker = true" />
        <TrashCan
          ref="trashCanRef"
          :isHovered="isNearTrash"
          @clear-all="clearAllCards"
        />
      </div>

      <!-- 右上角 Profile（游客或未登录均可点击） -->
      <button
        type="button"
        class="profile-btn"
        @click="handleProfileClick"
        :title="isAuthenticated ? '个人资料' : '登录 / 注册'"
      >
        <div class="profile-icon-wrap">
          <svg class="profile-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
        </div>
        <span v-if="isAuthenticated" class="profile-label">{{ profileLabel }}</span>
      </button>

      <AuthDrawer v-model:show="showAuth" />
      <ProfileDrawer v-model:show="showProfile" />

      <CardPickerModal
        v-model:show="showCardPicker"
        @confirm="onCardPickerConfirm"
      />

      <div v-if="guidedMode" class="spread-banner">
        <span class="spread-banner-name">{{ selectedSpread?.name }}</span>
        <button type="button" class="respread-btn" @click="respread">返回</button>
      </div>

      <!-- 牌阵舞台 -->
      <div class="stage" ref="stage">
        <div
          v-for="slot in stageSlots"
          v-show="!occupiedSlotIndices.has(slot.slotIndex)"
          :key="'slot-'+slot.slotIndex"
          class="slot-frame"
          :style="{
            left: slot.x + 'px',
            top: slot.y + 'px',
            width: slot.width + 'px',
            height: slot.height + 'px'
          }"
        >
          <span class="slot-meaning">{{ slot.meaning }}</span>
        </div>
        <Card
          v-for="card in drawnCards" 
          :key="card.id"
          :card="card"
          :width="baseWidth"
          :height="cardHeight"
          :isActive="activeCard?.id === card.id"
          :isGlobalResizing="isResizing"
          @drag-start="onpointerDownExisting"
          @resize-start="onpointerDownResize"
        />
      </div>
  
      <!-- 左下角牌堆 -->
      <div class="card-pile" @pointerdown.stop="onpointerDownPile">
        <div class="pile-wrapper" :style="{ width: baseWidth + 'px', height: cardHeight + 'px' }">
          <CardBack class="pile-back" />
          <div class="pile-shadow"></div>
        </div>
        <span class="pile-label">长按拖拽抽取</span>
      </div>
  
      <!-- 右下角完成按钮 -->
      <button class="finish-btn" @click="onFinish" :disabled="drawnCards.length === 0">
        完成抽取 ({{ drawnCards.length }})
      </button>

      <!-- 引导模式空缺提醒 -->
      <Transition name="modal-fade">
        <div v-if="showMissingToast" class="missing-toast">
          还有 {{ missingCount }} 张牌未摆放，请继续抽牌
        </div>
      </Transition>
  
      <!-- 弹窗组件 -->
      <SubmitModal 
        v-model:showModal="showModal"
        v-model:question="question"
        v-model:drawnCards="drawnCards"
        :isSubmitDisabled="isSubmitDisabled"
        :cardWidth="baseWidth"     
        :cardHeight="cardHeight"  
        ref="submitModalRef"
        @submit="submitToBackend"
      />

      <AnswerModal 
        v-model:showModal="showAnswerModal"
        :fullAnswer="backendAnswer"
        :streamBlocks="streamBlocks"
        :drawnCards="drawnCards"
        :cardWidth="baseWidth"
        :cardHeight="cardHeight"
        :isStreaming="isStreaming"
       />
    </div>
  </template>
  
  <script setup>
  import { ref, computed, onMounted, onUnmounted } from 'vue';
  import SubmitModal from './SubmitModal.vue';
  import AnswerModal from './AnswerModal.vue';
  import Card from './Card.vue'; 
import CardBack from './CardBack.vue';
import TrashCan from './TrashCan.vue';
  import CardPickerButton from './CardPickerButton.vue';
  import CardPickerModal from './CardPickerModal.vue';
  import ProfileDrawer from './ProfileDrawer.vue';
  import AuthDrawer from './AuthDrawer.vue';
  import { allCardNames, ASPECT_RATIO } from '../utils/cardInfo.js';
  import { predictStream } from '../utils/predictStream.js';
  import { applyStreamEvent, blocksToAnswer, createEmptyStreamBlocks } from '../utils/streamBlocks.js';
  import { useAuth } from '../composables/useAuth.js';
  import { usePredictHistory } from '../composables/usePredictHistory.js';
  import { useCustomApi } from '../composables/useCustomApi.js';
  import { useTarotFlow } from '../composables/useTarotFlow.js';
  import { layoutSlotsInStage, findNearestSlot } from '../utils/spreadLayout.js';

  const { getToken, logout, fetchMe, isAuthenticated, forceLogin, user } = useAuth();
  const { refreshHistory } = usePredictHistory();
  const { getConfig: getCustomApiConfig, isEnabled: customApiEnabled } = useCustomApi();

  const trashCanRef = ref(null);
  const showAuth = ref(false);

  const profileLabel = computed(() => user.value?.nickname || '用户');

  function handleProfileClick() {
    if (isAuthenticated.value) {
      showProfile.value = true;
    } else {
      showAuth.value = true;
    }
  }
  const stage = ref(null);
  const drawnCards = ref([]);
  const availableCards = ref([...allCardNames]);
  const activeCard = ref(null);
  const showModal = ref(false);
  const showCardPicker = ref(false);
  const showProfile = ref(false);
  const question = ref('');
  const submitModalRef = ref(null);

  const showAnswerModal = ref(false);
  const backendAnswer = ref('');
  const streamBlocks = ref(createEmptyStreamBlocks());
  const isStreaming = ref(false);
  
  const isSubmitDisabled = computed(() => false);

  // 引导模式空缺提醒
  const showMissingToast = ref(false);
  const missingCount = computed(() => {
    if (!guidedMode.value || !selectedSpread.value) return 0;
    return selectedSpread.value.cardCount - occupiedSlotIndices.value.size;
  });
  let missingToastTimer = null;
  function onFinish() {
    if (guidedMode.value && missingCount.value > 0) {
      showMissingToast.value = true;
      clearTimeout(missingToastTimer);
      missingToastTimer = setTimeout(() => { showMissingToast.value = false; }, 2500);
      return;
    }
    showModal.value = true;
  }
  
  // 尺寸控制 (全局状态)
  const STORAGE_KEY = 'tarot_card_base_width';
  const savedWidth = localStorage.getItem(STORAGE_KEY);
  const baseWidth = ref(savedWidth ? parseFloat(savedWidth) : 120); 

  const cardHeight = computed(() => baseWidth.value * ASPECT_RATIO);

  // 引导模式（新手体验）：透明阵位框 + 吸附 + 抽牌上限
  const { guidedMode, selectedSpread, clearSpread } = useTarotFlow();
  const stageWidth = ref(window.innerWidth);
  const stageHeight = ref(window.innerHeight);

  const stageSlots = computed(() => {
    if (!guidedMode.value || !selectedSpread.value) return [];
    return layoutSlotsInStage(
      selectedSpread.value.slots,
      stageWidth.value,
      stageHeight.value,
      baseWidth.value,
      cardHeight.value
    );
  });

  const cardLimit = computed(() =>
    guidedMode.value && selectedSpread.value ? selectedSpread.value.cardCount : 15
  );

  const occupiedSlotIndices = computed(() => {
    const set = new Set();
    drawnCards.value.forEach((c) => {
      if (c.slotIndex !== undefined && c.slotIndex !== null) set.add(c.slotIndex);
    });
    return set;
  });

  function updateStageSize() {
    if (stage.value) {
      const r = stage.value.getBoundingClientRect();
      stageWidth.value = r.width;
      stageHeight.value = r.height;
    }
  }

  function respread() {
    clearSpread();
    drawnCards.value = [];
    activeCard.value = null;
  }

  // 交互状态
  const isResizing = ref(false);
  const isNearTrash = ref(false);
  let dragOffset = { x: 0, y: 0 };
  let resizeStartData = { x: 0, width: 0 };
  
  // 3. 背景与图片资源逻辑
  const backgroundStyle = computed(() => ({
    backgroundColor: '#433843',
  }));

  const removeCardByName = (cardName, returnToPool = true) => {
    const index = drawnCards.value.findIndex(c => c.name === cardName);
    if (index === -1) return;
    if (returnToPool) {
      availableCards.value.push(cardName);
    }
    drawnCards.value.splice(index, 1);
    drawnCards.value.forEach((card, idx) => {
      card.order = idx + 1;
    });
  };

  const onCardPickerConfirm = ({ name, isReversed }) => {
    if (drawnCards.value.length >= cardLimit.value) return;

    removeCardByName(name, false);

    const poolIndex = availableCards.value.indexOf(name);
    if (poolIndex !== -1) {
      availableCards.value.splice(poolIndex, 1);
    }

    const stageRect = stage.value?.getBoundingClientRect();
    const centerX = stageRect ? stageRect.width / 2 - baseWidth.value / 2 : 200;
    const centerY = stageRect ? stageRect.height / 2 - cardHeight.value / 2 : 200;

    const offset = drawnCards.value.length * 12;
    drawnCards.value.push({
      id: Date.now(),
      name,
      x: centerX + offset,
      y: centerY + offset,
      isRevealed: true,
      isReversed,
      order: drawnCards.value.length + 1,
      meaning: '',
    });
  };

  const removeCard = (cardId) => {
    const index = drawnCards.value.findIndex(c => c.id === cardId);
    if (index !== -1) {
      const card = drawnCards.value[index];
      
      // 如果卡牌已经有名字了，把它还给牌堆池
      if (card.name) {
        availableCards.value.push(card.name);
      }
      
      // 从已抽取的牌中移除
      drawnCards.value.splice(index, 1);
      
      // 关键：重新排序
      drawnCards.value.forEach((card, idx) => {
        card.order = idx + 1;
      });
    }
  };

  const clearAllCards = () => {
  if (drawnCards.value.length === 0) return;
    drawnCards.value.forEach(card => {
      if (card.name) {
        availableCards.value.push(card.name);
      }
    });
    
    drawnCards.value = [];
    activeCard.value = null;
    
    console.log("所有卡牌已清空");
};

  const getRelativeCoords = (e) => {
    const rect = stage.value.getBoundingClientRect();
    return { x: e.clientX - rect.left, y: e.clientY - rect.top };
  };
  
  // 4. 交互逻辑：接收缩放事件
  const onpointerDownResize = (e) => {
    if (showModal.value || showAnswerModal.value || showCardPicker.value) return; 
    isResizing.value = true;
    resizeStartData = {
      x: e.clientX,
      width: baseWidth.value
    };
  };
  
  // 5. 交互逻辑：抽取与接收拖拽事件
  const onpointerDownPile = (e) => {
    if (showModal.value || showAnswerModal.value || showCardPicker.value) return; 

    if (availableCards.value.length === 0 || drawnCards.value.length >= 15) return;
    
    const coords = getRelativeCoords(e);
    const newCard = {
      id: Date.now(),
      name: '',
      x: coords.x - baseWidth.value / 2,
      y: coords.y - cardHeight.value / 2,
      isRevealed: false,
      isReversed: false,
      order: drawnCards.value.length + 1,
      meaning: '' 
    };
  
    drawnCards.value.push(newCard);
    activeCard.value = newCard;
    dragOffset = { x: baseWidth.value / 2, y: cardHeight.value / 2 };
  };
  
  const onpointerDownExisting = (card, e) => {
    if (showModal.value || showAnswerModal.value || showCardPicker.value) return; 

    activeCard.value = card;
    const coords = getRelativeCoords(e);
    dragOffset = { x: coords.x - card.x, y: coords.y - card.y };
  };
  
  // 鼠标全局移动监听
  const handlepointerMove = (e) => {
    if (isResizing.value) {
      const deltaX = e.clientX - resizeStartData.x;
      const newWidth = resizeStartData.width + deltaX;
      if (newWidth >= 40 && newWidth <= 400) {
        baseWidth.value = newWidth;
      }
      return;
    }
  
    if (!activeCard.value) return;
    const coords = getRelativeCoords(e);
    activeCard.value.x = coords.x - dragOffset.x;
    activeCard.value.y = coords.y - dragOffset.y;

    const trashEl = trashCanRef.value?.rootEl;
    if (trashEl) {
      const rect = trashEl.getBoundingClientRect();
      const trashX = rect.left + rect.width / 2;
      const trashY = rect.top + rect.height / 2;
      const distance = Math.sqrt(
        Math.pow(e.clientX - trashX, 2) + Math.pow(e.clientY - trashY, 2)
      );
      isNearTrash.value = distance < 80;
    }
  };
  
  // 鼠标松开监听
  const handlepointerUp = () => {
    if (isResizing.value) {
      isResizing.value = false;
      localStorage.setItem(STORAGE_KEY, baseWidth.value.toString());
      return;
    }
  
    if (isNearTrash.value) {
      removeCard(activeCard.value.id);
    } else {
      if (activeCard.value) {
        if (!activeCard.value.isRevealed) {
          const randomIndex = Math.floor(Math.random() * availableCards.value.length);
          activeCard.value.name = availableCards.value.splice(randomIndex, 1)[0];
          activeCard.value.isReversed = Math.random() > 0.5;
          activeCard.value.isRevealed = true;
        }
        if (guidedMode.value) {
          const cx = activeCard.value.x + baseWidth.value / 2;
          const cy = activeCard.value.y + cardHeight.value / 2;
          const slot = findNearestSlot(
            stageSlots.value,
            cx,
            cy,
            occupiedSlotIndices.value,
            baseWidth.value * 1.5
          );
          if (slot) {
            // 卡牌中心对齐 slot 中心：解决不同尺寸下卡牌位置不统一
            activeCard.value.x = slot.x + (slot.width - baseWidth.value) / 2;
            activeCard.value.y = slot.y + (slot.height - cardHeight.value) / 2;
            activeCard.value.slotIndex = slot.slotIndex;
          } else {
            activeCard.value.slotIndex = undefined;
          }
        }
        activeCard.value = null;
      }
    }
  };
  
  // 6. 提交
  const submitToBackend = async ({ clarifications = [], intentSummary = '', freestylemode = false } = {}) => {
    try{
    const customApi = customApiEnabled.value ? getCustomApiConfig() : null;
    const payload = {
      question: question.value,
      intent_summary: intentSummary || undefined,
      clarifications,
      freestylemode,
      model: submitModalRef.value.selectedModel,
      cardSize: { width: baseWidth.value, height: cardHeight.value },
      cards: drawnCards.value.map(card => ({
        order: card.order,
        name: card.name,
        x: Math.round(card.x),
        y: Math.round(card.y),
        orientation: card.isReversed ? 'reversed' : 'upright',
        meaning: card.meaning.trim() 
      })),
      custom_api: customApi || undefined,
    };
  
    backendAnswer.value = "";
    streamBlocks.value = createEmptyStreamBlocks();
    showAnswerModal.value = true;
    isStreaming.value = true;

    try {
      const submitInfo = submitModalRef.value.submitInfo?.value ?? submitModalRef.value.submitInfo ?? 'Result';
      const stream = predictStream(payload, submitInfo, getToken());
      for await (const event of stream) {
        applyStreamEvent(streamBlocks, event);
        backendAnswer.value = blocksToAnswer(streamBlocks.value);
      }
      await fetchMe();
      await refreshHistory();
    } catch (err) {
      console.error(err);
      if (err.status === 401 && forceLogin.value) {
        logout();
        alert('登录已失效，请重新登录');
      } else if (err.status === 429) {
        alert(err.message);
        showAnswerModal.value = false;
      } else {
        alert(err.message || '提交失败');
        showAnswerModal.value = false;
      }
    } finally {
      isStreaming.value = false;
    }
  } catch (err){
    console.error("提交失败", err);
    submitModalRef.value?.unlockSubmit();
  }
  };

  onMounted(() => {
    updateStageSize();
    window.addEventListener('resize', updateStageSize);
  });

  onUnmounted(() => {
    window.removeEventListener('resize', updateStageSize);
  });
</script>
  
<style scoped>
  /* 布局级与组件级样式 */
  .tarot-container {
    width: 100vw;
    height: 100vh;
    background-color: #433843;
    position: relative;
    user-select: none;
    overflow: hidden;
    touch-action: none;
  }
  
  .stage {
    width: 100%;
    height: 100%;
    position: absolute;
  }
  
  /* 牌堆样式 */
  .card-pile {
    position: absolute;
    bottom: 40px;
    left: 40px;
    cursor: pointer;
    z-index: 500;
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  .pile-wrapper {
    position: relative;
    transition: transform 0.3s;
  }
  
  .card-pile:hover .pile-wrapper {
    transform: translateY(-5px);
  }
  
  .pile-back {
    width: 100%;
    height: 100%;
    border-radius: 8px;
    box-shadow: 0 5px 15px rgba(0,0,0,0.5);
  }
  
  .pile-shadow {
    position: absolute;
    top: 6px;
    left: 6px;
    width: 100%;
    height: 100%;
    background: rgba(0,0,0,0.4);
    border-radius: 8px;
    z-index: -1;
  }
  
  .pile-label {
    color: #aaa;
    margin-top: 15px;
    font-size: 13px;
    text-shadow: 0 2px 4px rgba(0,0,0,0.5);
  }
  
  /* 完成按钮样式 */
  .finish-btn {
    position: absolute;
    bottom: 50px;
    right: 50px;
    padding: 15px 40px;
    background: rgba(255, 255, 255, 0.05);
    backdrop-filter: blur(15px);
    color: #fff;
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 30px;
    cursor: pointer;
    font-size: 16px;
    transition: all 0.3s;
    z-index: 500;
    
  }
  
  .finish-btn:hover:not(:disabled) {
    background: rgba(255, 255, 255, 0.2);
    transform: scale(1.05);
  }
  
  .finish-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .left-toolbar {
    position: absolute;
    top: 30px;
    left: 30px;
    z-index: 1000;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;
  }

  .profile-btn {
    position: absolute;
    top: 30px;
    right: 30px;
    z-index: 1000;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    padding: 0;
    background: none;
    border: none;
    color: rgba(255, 255, 255, 0.4);
    cursor: pointer;
    transition: all 0.3s ease;
  }

  .profile-icon-wrap {
    width: 56px;
    height: 56px;
    box-sizing: border-box;
    background: rgba(255, 255, 255, 0.05);
    padding: 12px;
    border-radius: 50%;
    border: 1px solid rgba(255, 255, 255, 0.1);
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  }

  .profile-icon {
    width: 32px;
    height: 32px;
    display: block;
  }

  .profile-label {
    font-size: 12px;
    opacity: 0;
    transition: opacity 0.3s;
    white-space: nowrap;
  }

  .profile-btn:hover {
    color: #c2a35f;
  }

  .profile-btn:hover .profile-icon-wrap {
    background: rgba(194, 163, 95, 0.15);
    border-color: #c2a35f;
    transform: scale(1.1);
    box-shadow: 0 0 20px rgba(194, 163, 95, 0.25);
  }

  .profile-btn:hover .profile-label {
    opacity: 1;
  }

  .spread-banner {
    position: absolute;
    top: 30px;
    left: 50%;
    transform: translateX(-50%);
    z-index: 500;
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 10px 18px;
    background: rgba(0, 0, 0, 0.45);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(194, 163, 95, 0.4);
    border-radius: 24px;
  }

  .spread-banner-name {
    color: #e5d8b0;
    font-size: 14px;
  }

  .respread-btn {
    padding: 4px 12px;
    background: rgba(194, 163, 95, 0.15);
    border: 1px solid rgba(194, 163, 95, 0.5);
    color: #c2a35f;
    border-radius: 12px;
    cursor: pointer;
    font-size: 12px;
  }

  .respread-btn:hover {
    background: rgba(194, 163, 95, 0.3);
  }

  .slot-frame {
    position: absolute;
    border: 2px dashed rgba(194, 163, 95, 0.5);
    border-radius: 8px;
    box-sizing: border-box;
    pointer-events: none;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: border-color 0.2s, background 0.2s;
  }

  .slot-meaning {
    position: absolute;
    bottom: -22px;
    left: 50%;
    transform: translateX(-50%);
    color: #c2a35f;
    font-size: 12px;
    white-space: nowrap;
    text-shadow: 0 1px 3px rgba(0, 0, 0, 0.8);
    text-align: center;
    padding: 2px 6px;
    border-radius: 4px;
    background: rgba(0, 0, 0, 0.3);
  }

  .missing-toast {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 800;
    padding: 14px 22px;
    background: rgba(0, 0, 0, 0.75);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 140, 80, 0.5);
    border-radius: 12px;
    color: #ffd9b3;
    font-size: 14px;
    text-align: center;
    box-shadow: 0 6px 24px rgba(0, 0, 0, 0.6);
    pointer-events: none;
  }

  .modal-fade-enter-active, .modal-fade-leave-active {
    transition: opacity 0.25s ease;
  }
  .modal-fade-enter-from, .modal-fade-leave-to {
    opacity: 0;
  }
  </style>