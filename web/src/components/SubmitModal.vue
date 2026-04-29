<template>
    <Transition name="slide-right">
      <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
        <div class="modal-content">
          <div class="modal-header">
            <h3>冥想与牌阵确认</h3>
            <p class="modal-hint">请看着你的牌阵，定义它们的含义</p>
          </div>
          
          <!-- 核心滚动区域 -->
          <div class="modal-body">
            <!-- 新增：迷你牌阵预览组件 -->
            <MiniTarot 
              :drawnCards="drawnCards"
              :cardWidth="cardWidth"
              :cardHeight="cardHeight"
              :containerWidth="370"  
              :containerHeight="180" 
            />
  
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
  import { defineModel, defineProps, defineEmits } from 'vue';
  import MiniTarot from './MiniTarot.vue'; // 引入刚才抽离的组件
  
  const showModal = defineModel('showModal');
  const question = defineModel('question');
  const drawnCards = defineModel('drawnCards');
  
  // 接收外部传来的状态，以及为了渲染迷你牌阵需要的宽高等比参数
  const props = defineProps({
    isSubmitDisabled: Boolean,
    cardWidth: { type: Number, default: 120 },
    cardHeight: { type: Number, default: 210 }
  });
  
  const emit = defineEmits(['submitToBackend']);
  
  const submitToBackend = () => {
    emit('submitToBackend');
  };
  
  const getCardUrl = (name) => {
    if (name === 'back') return new URL(`../assets/tarots/back.jpeg`, import.meta.url).href;
    const ext = 'jpeg';
    return new URL(`../assets/tarots/${name}.${ext}`, import.meta.url).href;
  };
  </script>
  
  <style scoped>
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
    padding: 30px 40px;
    border-radius: 24px;
    width: 450px;
    /* 限制最大高度，配合内部滚动 */
    max-height: 90vh; 
    display: flex;
    flex-direction: column;
    border: 1px solid #333;
    text-align: center;
  }
  
  .modal-header { margin-bottom: 10px; }
  .modal-hint { color: #666; font-size: 14px; margin-bottom: 15px; }
  
  /* 中间增加内部滚动条，避免内容过多时弹窗撑破屏幕 */
  .modal-body {
    flex: 1;
    overflow-y: auto;
    padding-right: 5px; /* 留出滚动条空间 */
  }
  
  /* 隐藏自带的丑陋滚动条（针对webkit浏览器） */
  .modal-body::-webkit-scrollbar { width: 6px; }
  .modal-body::-webkit-scrollbar-thumb { background: #333; border-radius: 4px; }
  
  textarea {
    width: 100%;
    height: 100px;
    background: #0a0a0c;
    color: #fff;
    border: 1px solid #333;
    border-radius: 12px;
    padding: 15px;
    resize: none;
    font-size: 15px;
    margin-bottom: 20px;
  }
  
  .card-item {
    display: flex;
    align-items: center;
    background: #1a1a1f;
    padding: 10px;
    border-radius: 12px;
    margin-bottom: 10px;
    gap: 15px;
  }
  .card-preview { position: relative; width: 40px; height: 70px; }
  .card-preview img { width: 100%; height: 100%; border-radius: 4px; object-fit: cover;}
  .card-preview img.is-reversed { transform: rotate(180deg); }
  .card-order { position: absolute; bottom: -8px; left: 50%; transform: translateX(-50%); background: #fff; color: #000; border-radius: 50%; width: 16px; height: 16px; font-size: 10px; font-weight: bold;}
  .card-input { flex: 1; display: flex; align-items: center; gap: 10px;}
  .card-input input { flex: 1; padding: 10px; border-radius: 8px; border: 1px solid #333; background: #0a0a0c; color: #fff;}
  .orientation-tag { font-size: 12px; color: #888; background: #222; padding: 4px 8px; border-radius: 6px; }
  
  .modal-footer {
    margin-top: 20px;
    padding-top: 20px;
    border-top: 1px solid #222;
  }
  .modal-btns { display: flex; gap: 15px; }
  .modal-btns button { flex: 1; padding: 12px; border-radius: 12px; cursor: pointer; font-weight: bold; border: none; transition: 0.2s;}
  .cancel-btn { background: #222; color: #999; }
  .cancel-btn:hover { background: #333; color: #fff; }
  .confirm-btn { background: #fff; color: #000; }
  .confirm-btn:disabled { opacity: 0.3; cursor: not-allowed; }
  </style>