<template>
  <button
    type="button"
    class="picker-container"
    :class="{ 'is-open': isOpen }"
    @click="$emit('open')"
  >
    <div class="icon-wrapper">
      <svg class="picker-icon" viewBox="0 0 24 24" width="32" height="32" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
        <line x1="12" y1="8" x2="12" y2="16"></line>
        <line x1="8" y1="12" x2="16" y2="12"></line>
      </svg>
    </div>
    <span class="picker-label">{{ isOpen ? '选择中...' : '自选卡牌' }}</span>
  </button>
</template>

<script setup>
defineProps({
  isOpen: Boolean,
});
defineEmits(['open']);
</script>

<style scoped>
.picker-container {
  position: absolute;
  top: 30px;
  left: 30px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  color: rgba(255, 255, 255, 0.4);
  transition: all 0.3s ease;
  z-index: 1000;
  padding: 20px;
  cursor: pointer;
  user-select: none;
  pointer-events: auto;
  touch-action: none;
  -webkit-tap-highlight-color: transparent;
  background: none;
  border: none;
  outline: none;
}

.icon-wrapper {
  position: relative;
  background: rgba(255, 255, 255, 0.05);
  padding: 12px;
  border-radius: 50%;
  border: 1px solid rgba(255, 255, 255, 0.1);
  transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  display: flex;
  align-items: center;
  justify-content: center;
}

.picker-container:hover {
  color: #c2a35f;
}

.picker-container:hover .icon-wrapper {
  background: rgba(194, 163, 95, 0.15);
  border-color: #c2a35f;
  transform: scale(1.1);
  box-shadow: 0 0 20px rgba(194, 163, 95, 0.25);
}

.picker-container.is-open {
  color: #c2a35f;
}

.picker-container.is-open .icon-wrapper {
  background: rgba(194, 163, 95, 0.2);
  border-color: #c2a35f;
  transform: scale(1.15);
  box-shadow: 0 0 25px rgba(194, 163, 95, 0.3);
}

.picker-label {
  font-size: 12px;
  opacity: 0;
  transition: opacity 0.3s;
}

.picker-container:hover .picker-label,
.picker-container.is-open .picker-label {
  opacity: 1;
}
</style>
