<template>
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
  </template>
  
  <script setup>
  import { defineModel, defineProps, defineEmits } from 'vue';
  
  // 使用 defineModel 接收父组件传来的变量，这允许我们在不改变量名的前提下直接保持双向绑定
  const showModal = defineModel('showModal');
  const question = defineModel('question');
  const drawnCards = defineModel('drawnCards');
  
  // 接收原代码模板中用到的 disabled 状态
  const props = defineProps(['isSubmitDisabled']);
  
  // 定义向父组件发送提交请求的事件
  const emit = defineEmits(['submitToBackend']);
  
  // 保持变量名不变，当点击确认时，通知父组件去执行真正的 submitToBackend
  const submitToBackend = () => {
    emit('submitToBackend');
  };
  
  // 弹窗内部需要预览图片，所以原封不动保留这个工具函数
  const getCardUrl = (name) => {
    if (name === 'back') return new URL(`../assets/tarots/back.jpeg`, import.meta.url).href;
    const ext = name === 'TheLovers' ? 'jpg' : 'jpeg';
    return new URL(`../assets/tarots/${name}.${ext}`, import.meta.url).href;
  };
  </script>
  
  <style scoped>
  /* 弹窗专属样式提取 */
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