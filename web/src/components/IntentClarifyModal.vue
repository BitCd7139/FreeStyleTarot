<template>
  <Transition name="fade">
    <div v-if="showModal" class="clarify-overlay" @click.self="handleCancel">
      <div class="clarify-panel">
        <div class="modal-header">
          <div class="clarify-header-row">
            <h3>理解你的问题</h3>
            <span class="clarify-progress">{{ currentIndex + 1 }} / {{ questions.length }}</span>
          </div>
          <p v-if="intentSummary" class="modal-hint">{{ intentSummary }}</p>
          <p v-else class="modal-hint">补充一点背景，解读会更贴合你的情况</p>
        </div>

        <div class="modal-body clarify-body">
          <p class="clarify-question">{{ currentQuestion?.text }}</p>
          <p v-if="currentQuestion?.multi_select" class="multi-select-hint">可多选</p>

          <div class="option-list">
            <button
              v-for="opt in currentQuestion?.options || []"
              :key="opt.id"
              type="button"
              class="option-chip"
              :class="{ selected: isOptionSelected(opt.id) }"
              @click="toggleOption(opt.id)"
            >
              {{ opt.label }}
            </button>
          </div>

          <div v-if="currentQuestion?.allow_custom" class="custom-input-wrap">
            <input
              v-model="customText"
              type="text"
              class="custom-input"
              maxlength="200"
              :placeholder="currentQuestion.custom_placeholder || '其他情况…'"
              @input="onCustomInput"
            />
          </div>
        </div>

        <div class="modal-footer clarify-footer">
          <button type="button" class="skip-btn" @click="handleSkip">直接开始占卜</button>
          <div class="clarify-nav">
            <button
              v-if="currentIndex > 0"
              type="button"
              class="cancel-btn"
              @click="goBack"
            >
              上一题
            </button>
            <button
              type="button"
              class="confirm-btn"
              :disabled="!canProceed"
              @click="handleNext"
            >
              {{ isLastQuestion ? '开始占卜' : '下一题' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { computed, ref, watch } from 'vue';

const showModal = defineModel('showModal');

const props = defineProps({
  questions: {
    type: Array,
    default: () => [],
  },
  intentSummary: {
    type: String,
    default: '',
  },
});

const emit = defineEmits(['confirm', 'skip', 'cancel']);

const currentIndex = ref(0);
const answers = ref([]);
const selectedOptionId = ref('');
const selectedOptionIds = ref([]);
const customText = ref('');

const currentQuestion = computed(() => props.questions[currentIndex.value] || null);
const isLastQuestion = computed(() => currentIndex.value >= props.questions.length - 1);
const canProceed = computed(() => {
  if (customText.value.trim()) return true;
  if (currentQuestion.value?.multi_select) {
    return selectedOptionIds.value.length > 0;
  }
  return Boolean(selectedOptionId.value);
});

const resetForQuestion = (index) => {
  const saved = answers.value[index];
  if (saved) {
    selectedOptionId.value = saved.optionId || '';
    selectedOptionIds.value = saved.optionIds ? [...saved.optionIds] : [];
    customText.value = saved.customText || '';
    return;
  }
  selectedOptionId.value = '';
  selectedOptionIds.value = [];
  customText.value = '';
};

watch(showModal, (open) => {
  if (open) {
    currentIndex.value = 0;
    answers.value = [];
    resetForQuestion(0);
  }
});

watch(currentIndex, (idx) => {
  resetForQuestion(idx);
});

const isOptionSelected = (optionId) => {
  if (currentQuestion.value?.multi_select) {
    return selectedOptionIds.value.includes(optionId) && !customText.value.trim();
  }
  return selectedOptionId.value === optionId && !customText.value.trim();
};

const toggleOption = (optionId) => {
  if (currentQuestion.value?.multi_select) {
    const idx = selectedOptionIds.value.indexOf(optionId);
    if (idx >= 0) {
      selectedOptionIds.value.splice(idx, 1);
    } else {
      selectedOptionIds.value.push(optionId);
    }
    customText.value = '';
  } else {
    selectedOptionId.value = optionId;
    customText.value = '';
  }
};

const onCustomInput = () => {
  if (customText.value.trim()) {
    selectedOptionId.value = '';
    selectedOptionIds.value = [];
  }
};

const buildCurrentAnswer = () => {
  const q = currentQuestion.value;
  if (!q) return null;

  const trimmedCustom = customText.value.trim();
  if (trimmedCustom) {
    return {
      questionId: q.id,
      question: q.text,
      answer: trimmedCustom,
    };
  }

  if (q.multi_select) {
    if (selectedOptionIds.value.length === 0) return null;
    const labels = selectedOptionIds.value
      .map((id) => q.options.find((o) => o.id === id))
      .filter(Boolean)
      .map((o) => o.label);
    if (labels.length === 0) return null;
    return {
      questionId: q.id,
      question: q.text,
      optionIds: [...selectedOptionIds.value],
      answer: labels.join('、'),
    };
  }

  const opt = q.options.find((o) => o.id === selectedOptionId.value);
  if (!opt) return null;

  return {
    questionId: q.id,
    question: q.text,
    optionId: opt.id,
    answer: opt.label,
  };
};

const saveCurrentAnswer = () => {
  const entry = buildCurrentAnswer();
  if (!entry) return;
  answers.value[currentIndex.value] = {
    ...entry,
    customText: customText.value.trim(),
    optionId: entry.optionId || '',
    optionIds: entry.optionIds || [],
  };
};

const buildClarifications = () =>
  answers.value
    .filter(Boolean)
    .map(({ questionId, question, optionId, optionIds, answer }) => ({
      questionId,
      question,
      optionId: optionId || undefined,
      optionIds: optionIds && optionIds.length > 0 ? optionIds : undefined,
      answer,
    }));

const handleNext = () => {
  saveCurrentAnswer();
  if (isLastQuestion.value) {
    emit('confirm', buildClarifications());
    showModal.value = false;
    return;
  }
  currentIndex.value += 1;
};

const goBack = () => {
  saveCurrentAnswer();
  if (currentIndex.value > 0) {
    currentIndex.value -= 1;
  }
};

const handleSkip = () => {
  saveCurrentAnswer();
  emit('skip', buildClarifications());
  showModal.value = false;
};

const handleCancel = () => {
  emit('cancel');
  showModal.value = false;
};
</script>

<style scoped>
.clarify-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.65);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1100;
  padding: 16px;
}

.clarify-panel {
  background: rgba(30, 27, 36, 0.96);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(194, 163, 95, 0.4);
  border-radius: 16px;
  width: min(480px, 100%);
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.6);
}

.clarify-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.clarify-progress {
  font-size: 13px;
  color: #c2a35f;
  white-space: nowrap;
}

.clarify-body {
  padding: 20px 25px;
}

.clarify-question {
  margin: 0 0 16px;
  color: #e5d8b0;
  font-size: 16px;
  line-height: 1.5;
}

.multi-select-hint {
  margin: -8px 0 12px;
  color: #c2a35f;
  font-size: 12px;
  font-weight: 500;
}

.option-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.option-chip {
  width: 100%;
  text-align: left;
  padding: 12px 14px;
  border-radius: 10px;
  border: 1px solid rgba(194, 163, 95, 0.35);
  background: rgba(194, 163, 95, 0.06);
  color: #e5d8b0;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.option-chip:hover {
  background: rgba(194, 163, 95, 0.12);
  border-color: rgba(194, 163, 95, 0.55);
}

.option-chip.selected {
  background: rgba(194, 163, 95, 0.22);
  border-color: #c2a35f;
  box-shadow: 0 0 12px rgba(194, 163, 95, 0.25);
}

.custom-input-wrap {
  margin-top: 14px;
}

.custom-input {
  width: 100%;
  box-sizing: border-box;
  padding: 10px 12px;
  border-radius: 8px;
  border: 1px solid rgba(194, 163, 95, 0.35);
  background: rgba(0, 0, 0, 0.35);
  color: #e5d8b0;
  font-size: 14px;
}

.custom-input:focus {
  outline: none;
  border-color: #c2a35f;
}

.clarify-footer {
  flex-direction: column;
  gap: 12px;
  align-items: stretch;
}

.skip-btn {
  width: 100%;
  padding: 10px 0;
  border: none;
  background: transparent;
  color: #a89f91;
  font-size: 13px;
  cursor: pointer;
  text-decoration: underline;
  text-underline-offset: 3px;
}

.skip-btn:hover {
  color: #c2a35f;
}

.clarify-nav {
  display: flex;
  gap: 12px;
  width: 100%;
}

.clarify-nav .cancel-btn,
.clarify-nav .confirm-btn {
  flex: 1;
  margin-right: 0;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
