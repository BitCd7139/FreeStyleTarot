<template>
  <div 
    class="tarot-container" 
    :style="backgroundStyle"
    @mousemove="handleMouseMove" 
    @mouseup="handleMouseUp"
    @mouseleave="handleMouseUp"
  >
    <!-- 牌阵舞台 -->
    <div class="stage" ref="stage">
      <div 
        v-for="(card, index) in drawnCards" 
        :key="card.id"
        class="card-item"
        :class="{ 'dragging': activeCard?.id === card.id, 'resizing': isResizing }"
        :style="{ 
          left: card.x + 'px', 
          top: card.y + 'px',
          width: baseWidth + 'px',
          height: cardHeight + 'px',
          zIndex: activeCard?.id === card.id ? 1000 : index 
        }"
        @mousedown.stop="onMouseDownExisting(card, $event)"
      >
        <!-- 卡牌图片 -->
        <div class="card-inner">
          <img 
            :src="card.isRevealed ? getCardUrl(card.name) : getCardUrl('back')" 
            :class="{ 'reversed': card.isReversed && card.isRevealed }"
            class="tarot-img"
            draggable="false"
          />
          <div v-if="card.isRevealed" class="order-badge">{{ card.order }}</div>
        </div>

        <!-- 四角缩放手柄 -->
        <div class="resizer nw" @mousedown.stop="onMouseDownResize"></div>
        <div class="resizer ne" @mousedown.stop="onMouseDownResize"></div>
        <div class="resizer sw" @mousedown.stop="onMouseDownResize"></div>
        <div class="resizer se" @mousedown.stop="onMouseDownResize"></div>
      </div>
    </div>

    <!-- 左下角牌堆 -->
    <div class="card-pile" @mousedown.stop="onMouseDownPile">
      <div class="pile-wrapper" :style="{ width: baseWidth + 'px', height: cardHeight + 'px' }">
        <img :src="getCardUrl('back')" alt="Card Back" class="pile-back" draggable="false" />
        <div class="pile-shadow"></div>
      </div>
      <span class="pile-label">长按拖拽抽取</span>
    </div>

    <!-- 右下角完成按钮 -->
    <button class="finish-btn" @click="showModal = true" :disabled="drawnCards.length === 0">
      完成抽取 ({{ drawnCards.length }})
    </button>

    <!-- 问题输入及牌阵定义抽屉（右侧） -->
    <Transition name="slide-right">
      <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
        <div class="modal-content">
          <div class="modal-header">
            <h3>冥想与牌阵确认</h3>
            <p class="modal-hint">请看着你的牌阵，定义它们的含义</p>
          </div>
          
          <!-- 核心滚动区域 -->
          <div class="modal-body">
            <!-- 第一部分：问题输入 -->
            <div class="form-group">
              <textarea v-model="question" placeholder="请在此输入你的困惑，例如：我近期的事业运势如何？"></textarea>
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
                    <input 
                      type="text" 
                      v-model="card.meaning" 
                      :placeholder="`定义第 ${card.order} 张牌 (如: 过去)`"
                    />
                    <span class="orientation-tag">{{ card.isReversed ? '逆位' : '正位' }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 第三部分：固定在底部的操作按钮 -->
          <div class="modal-footer">
            <button class="cancel-btn" @click="showModal = false">返回调整</button>
            <button class="confirm-btn" @click="submitToBackend" :disabled="isSubmitDisabled">
              确认提交
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import axios from 'axios';

// 1. 基础配置
const ASPECT_RATIO = 1.618; // 黄金比例
const allCardNames = ["aceofcups", "aceofpentacles", "aceofswords", "aceofwands", "death", "eightofcups", "eightofpentacles", "eightofswords", "eightofwands", "fiveofcups", "fiveofpentacles", "fiveofswords", "fiveofwands", "fourofcups", "fourofpentacles", "fourofswords", "fourofwands", "judgement", "justice", "kingofcups", "kingofpentacles", "kingofswords", "kingofwands", "knightofcups", "knightofpentacles", "knightofswords", "knightofwands", "nineofcups", "nineofpentacles", "nineofswords", "nineofwands", "pageofcups", "pageofpentacles", "pageofswords", "pageofwands", "queenofcups", "queenofpentacles", "queenofswords", "queenofwands", "sevenofcups", "sevenofpentacles", "sevenofswords", "sevenofwands", "sixofcups", "sixofpentacles", "sixofswords", "sixofwands", "temperance", "tenofcups", "tenofpentacles", "tenofswords", "tenofwands", "thechariot", "thedevil", "theemperor", "theempress", "thefool", "thehangedman", "thehermit", "thehierophant", "thehighpriestess", "TheLovers", "themagician", "themoon", "thestar", "thestrength", "thesun", "thetower", "theworld", "threeofcups", "threeofpentacles", "threeofswords", "threeofwands", "twoofcups", "twoofpentacles", "twoofswords", "twoofwands", "wheeloffortune"];

// 2. 响应式状态
const stage = ref(null);
const drawnCards = ref([]);
const availableCards = ref([...allCardNames]);
const activeCard = ref(null);
const showModal = ref(false);
const question = ref('');

// 尺寸控制
const baseWidth = ref(120); // 初始宽度
const cardHeight = computed(() => baseWidth.value * ASPECT_RATIO);

// 交互状态
const isResizing = ref(false);
let dragOffset = { x: 0, y: 0 };
let resizeStartData = { x: 0, width: 0 };

// 3. 背景与图片资源逻辑
const backgroundStyle = computed(() => ({
  backgroundImage: `url(${new URL('../assets/background.png', import.meta.url).href})`,
  backgroundSize: 'cover',
  backgroundPosition: 'center'
}));

const getCardUrl = (name) => {
  if (name === 'back') return new URL(`../assets/tarots/back.jpeg`, import.meta.url).href;
  const ext = name === 'TheLovers' ? 'jpg' : 'jpeg';
  return new URL(`../assets/tarots/${name}.${ext}`, import.meta.url).href;
};

const getRelativeCoords = (e) => {
  const rect = stage.value.getBoundingClientRect();
  return { x: e.clientX - rect.left, y: e.clientY - rect.top };
};

// 4. 交互逻辑：缩放
const onMouseDownResize = (e) => {
  isResizing.value = true;
  resizeStartData = {
    x: e.clientX,
    width: baseWidth.value
  };
};

// 5. 交互逻辑：抽取与拖拽
const onMouseDownPile = (e) => {
  if (availableCards.value.length === 0) return;
  
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

const onMouseDownExisting = (card, e) => {
  activeCard.value = card;
  const coords = getRelativeCoords(e);
  dragOffset = { x: coords.x - card.x, y: coords.y - card.y };
};

const handleMouseMove = (e) => {
  // 处理缩放
  if (isResizing.value) {
    const deltaX = e.clientX - resizeStartData.x;
    const newWidth = resizeStartData.width + deltaX;
    // 限制缩放范围：80px 到 400px
    if (newWidth >= 80 && newWidth <= 400) {
      baseWidth.value = newWidth;
    }
    return;
  }

  // 处理拖拽
  if (!activeCard.value) return;
  const coords = getRelativeCoords(e);
  activeCard.value.x = coords.x - dragOffset.x;
  activeCard.value.y = coords.y - dragOffset.y;
};

const handleMouseUp = () => {
  if (isResizing.value) {
    isResizing.value = false;
    return;
  }

  if (activeCard.value) {
    // 如果是新抽取的牌，则翻面
    if (!activeCard.value.isRevealed) {
      const randomIndex = Math.floor(Math.random() * availableCards.value.length);
      activeCard.value.name = availableCards.value.splice(randomIndex, 1)[0];
      activeCard.value.isReversed = Math.random() > 0.5;
      activeCard.value.isRevealed = true;
    }
    activeCard.value = null;
  }
};

// 6. 提交
const submitToBackend = async () => {
  const payload = {
    question: question.value,
    cardSize: { width: baseWidth.value, height: cardHeight.value },
    cards: drawnCards.value.map(card => ({
      order: card.order,
      name: card.name,
      x: Math.round(card.x),
      y: Math.round(card.y),
      orientation: card.isReversed ? 'reversed' : 'upright',
      meaning: card.meaning.trim() // 【新增】将用户填写的含义发给后端
    }))
  };

  try {
    // 建议：测试阶段可以先打印看看 payload 长什么样
    console.log("Submitting payload:", payload);
    await axios.post('http://localhost:8080/api/tarot', payload);
    alert("契约已成，命运的齿轮开始转动...");
    showModal.value = false;
  } catch (error) {
    console.error(error);
    alert("无法连接至命运之门（后端错误）");
  }
};

</script>

<style scoped>
.tarot-container {
  width: 100vw;
  height: 100vh;
  background-color: #050508;
  position: relative;
  user-select: none;
  overflow: hidden;
}

.stage {
  width: 100%;
  height: 100%;
  position: absolute;
}

/* 卡牌基础样式 */
.card-item {
  position: absolute;
  cursor: grab;
  transition: transform 0.2s cubic-bezier(0.25, 0.46, 0.45, 0.94);
}

.card-item.dragging {
  cursor: grabbing;
  transform: scale(1.02);
  filter: brightness(1.1);
}

.card-inner {
  position: relative;
  width: 100%;
  height: 100%;
  border-radius: 8px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.8);
  overflow: visible;
}

.tarot-img {
  width: 100%;
  height: 100%;
  border-radius: 8px;
  object-fit: cover;
  pointer-events: none;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.reversed {
  transform: rotate(180deg);
}

/* 缩放手柄样式 */
.resizer {
  position: absolute;
  width: 16px;
  height: 16px;
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-radius: 50%;
  z-index: 10;
  opacity: 0;
  transition: opacity 0.3s, background 0.2s;
}

.card-item:hover .resizer {
  opacity: 1;
}

.resizer:hover {
  background: #fff;
  box-shadow: 0 0 10px #fff;
}

/* 手柄位置 */
.nw { top: -8px; left: -8px; cursor: nwse-resize; }
.ne { top: -8px; right: -8px; cursor: nesw-resize; }
.sw { bottom: -8px; left: -8px; cursor: nesw-resize; }
.se { bottom: -8px; right: -8px; cursor: nwse-resize; }

.order-badge {
  position: absolute;
  top: -12px;
  right: -12px;
  background: linear-gradient(135deg, #ffd700, #b8860b);
  color: #000;
  border-radius: 50%;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: bold;
  border: 2px solid #000;
}

/* 牌堆 */
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

/* 按钮 */
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

/* 弹窗 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.85);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 2000;
}

.modal-content {
  background: #151518;
  padding: 40px;
  border-radius: 24px;
  width: 450px;
  border: 1px solid #333;
  text-align: center;
}

.modal-hint {
  color: #666;
  font-size: 14px;
  margin-bottom: 20px;
}

textarea {
  width: 100%;
  height: 120px;
  background: #0a0a0c;
  color: #fff;
  border: 1px solid #333;
  border-radius: 12px;
  padding: 15px;
  resize: none;
  font-size: 16px;
  margin-bottom: 25px;
}

.modal-btns {
  display: flex;
  gap: 15px;
}

.modal-btns button {
  flex: 1;
  padding: 12px;
  border-radius: 12px;
  cursor: pointer;
  font-weight: bold;
  border: none;
}

.cancel-btn { background: #222; color: #999; }
.confirm-btn { background: #fff; color: #000; }
.confirm-btn:disabled { opacity: 0.3; }

/* 动画 */
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>