<template>
    <div 
      class="tarot-container" 
      :style="backgroundStyle"
      @pointermove="handlepointerMove" 
      @pointerup="handlepointerUp"
      @pointerleave="handlepointerUp"
      @pointercancel="handlepointerUp"
    >
      <!-- 垃圾桶组件 -->
      <TrashCan :isHovered="isNearTrash"
                @clear-all="clearAllCards"  />

      <!-- 牌阵舞台 -->
      <div class="stage" ref="stage">
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
          <img :src="backCardUrl" alt="Card Back" class="pile-back" draggable="false" />
          <div class="pile-shadow"></div>
        </div>
        <span class="pile-label">长按拖拽抽取</span>
      </div>
  
      <!-- 右下角完成按钮 -->
      <button class="finish-btn" @click="showModal = true" :disabled="drawnCards.length === 0">
        完成抽取 ({{ drawnCards.length }})
      </button>
  
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
        :drawnCards="drawnCards"
        :cardWidth="baseWidth"
        :cardHeight="cardHeight"
       />
    </div>
  </template>
  
  <script setup>
  import { ref, computed } from 'vue';
  import SubmitModal from './SubmitModal.vue';
  import AnswerModal from './AnswerModal.vue';
  import Card from './Card.vue'; 
  import TrashCan from './TrashCan.vue';
  import { allCardNames, ASPECT_RATIO } from '../utils/cardInfo.js';
  import { predictStream } from '../utils/predictStream.js';

  // 响应式状态
  const stage = ref(null);
  const drawnCards = ref([]);
  const availableCards = ref([...allCardNames]);
  const activeCard = ref(null);
  const showModal = ref(false);
  const question = ref('');
  const submitModalRef = ref(null);

  const showAnswerModal = ref(false);
  const backendAnswer = ref('');
  const isStreaming = ref(false);
  
  const isSubmitDisabled = computed(() => false); 
  
  // 尺寸控制 (全局状态)
  const STORAGE_KEY = 'tarot_card_base_width';
  const savedWidth = localStorage.getItem(STORAGE_KEY);
  const baseWidth = ref(savedWidth ? parseFloat(savedWidth) : 120); 

  const cardHeight = computed(() => baseWidth.value * ASPECT_RATIO);

  // 交互状态
  const isResizing = ref(false);
  const isNearTrash = ref(false);
  let dragOffset = { x: 0, y: 0 };
  let resizeStartData = { x: 0, width: 0 };
  
  // 3. 背景与图片资源逻辑
  const backgroundStyle = computed(() => ({
    backgroundImage: `url(${new URL('../assets/background.webp', import.meta.url).href})`,
    backgroundSize: 'cover',
    backgroundPosition: 'center'
  }));
  
  // 牌堆背面的图片URL
  const backCardUrl = new URL(`../assets/tarots/back.jpeg`, import.meta.url).href;

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
    if (showModal.value || showAnswerModal.value) return; 
    isResizing.value = true;
    resizeStartData = {
      x: e.clientX,
      width: baseWidth.value
    };
  };
  
  // 5. 交互逻辑：抽取与接收拖拽事件
  const onpointerDownPile = (e) => {
    if (showModal.value || showAnswerModal.value) return; 

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
    if (showModal.value || showAnswerModal.value) return; 

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

    const trashX = window.innerWidth - 80; 
    const trashY = 80;
    const distance = Math.sqrt(
      Math.pow(e.clientX - trashX, 2) + Math.pow(e.clientY - trashY, 2)
    );
    isNearTrash.value = distance < 100;
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
        activeCard.value = null;
      }
    }
  };
  
  // 6. 提交
  const submitToBackend = async () => {
    try{
    const payload = {
      question: question.value,
      model: submitModalRef.value.selectedModel,
      cardSize: { width: baseWidth.value, height: cardHeight.value },
      cards: drawnCards.value.map(card => ({
        order: card.order,
        name: card.name,
        x: Math.round(card.x),
        y: Math.round(card.y),
        orientation: card.isReversed ? 'reversed' : 'upright',
        meaning: card.meaning.trim() 
      }))
    };
  
    backendAnswer.value = "";
    showAnswerModal.value = true;
    isStreaming.value = true;
    
    console.log(payload);
    console.log(submitModalRef.value.submitInfo.value);

    try {
      const stream = predictStream(payload, submitModalRef.value.submitInfo);
      for await (const chunk of stream) {
        backendAnswer.value += chunk;
      }
    } catch (err) {
      console.error(err);
    } finally {
      isStreaming.value = false;
    }
  } catch (err){
    console.error("提交失败", err);
    submitModalRef.value?.unlockSubmit();
  }
  };
</script>
  
<style scoped>
  /* 布局级与组件级样式 */
  .tarot-container {
    width: 100vw;
    height: 100vh;
    background-color: #050508;
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
  </style>