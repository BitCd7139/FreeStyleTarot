<!-- RoleSelector.vue -->
<template>
    <div class="role-selector-container">
      <div class="model-select-wrapper">
        <span class="model-label">选择占卜角色：</span>
        <div class="select-container">
          <select 
            :value="modelValue" 
            @change="$emit('update:modelValue', $event.target.value)" 
            class="custom-select"
          >
            <option v-for="(val, key) in options" :key="val" :value="val">
              {{ key }}
            </option>
          </select>
        </div>
      </div>
    </div>
  </template>
  
  <script setup>
  defineProps({
    modelValue: String,
    options: {
      type: Object,
      default: () => ({})
    }
  });
  
  defineEmits(['update:modelValue']);
  </script>
  
  <style scoped>
  .role-selector-container {
    margin: 20px 0;
    padding: 12px;
    background: rgba(194, 163, 95, 0.05); /* 淡淡的金棕色背景，增加层次感 */
    border-radius: 8px;
  }
  
  .model-select-wrapper {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  
  .model-label {
    color: #c2a35f;
    font-size: 14px;
    font-weight: bold;
    white-space: nowrap;
  }
  
  .select-container {
    flex: 1;
    margin-left: 15px;
    position: relative;
  }
  
  .select-container::after {
    content: '▼';
    position: absolute;
    right: 12px;
    top: 50%;
    transform: translateY(-50%);
    color: #c2a35f;
    font-size: 10px;
    pointer-events: none;
  }
  
  .custom-select {
    width: 100%;
    padding: 10px 15px;
    background: #1a1a1a; 
    color: #c2a35f;
    border: 1px solid #c2a35f;
    border-radius: 6px;
    font-size: 14px;
    outline: none;
    cursor: pointer;
    appearance: none;
    -webkit-appearance: none;
    transition: all 0.3s ease;
  }
  
  .custom-select:focus {
    border-color: #e6d2a5;
    box-shadow: 0 0 10px rgba(194, 163, 95, 0.2);
  }
  
  .custom-select option {
    background-color: #1a1a1a;
    color: #c2a35f;
  }
  </style>