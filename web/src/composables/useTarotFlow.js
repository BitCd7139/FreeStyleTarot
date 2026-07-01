import { ref, computed } from 'vue';
import { useAuth } from './useAuth.js';

// 跨视图共享：选阵状态在 TarotBeginner 选定后由 TarotMain 读取
const selectedSpread = ref(null);

export function useTarotFlow() {
  const { user } = useAuth();

  // 登录用户的体验偏好为新手时才进入引导流程；游客与高级用户始终走高级模式
  const isBeginner = computed(() => !!user.value && user.value.experience_mode === 'beginner');
  const guidedMode = computed(() => isBeginner.value && !!selectedSpread.value);

  function selectSpread(spread) {
    selectedSpread.value = spread;
  }

  function clearSpread() {
    selectedSpread.value = null;
  }

  return {
    selectedSpread,
    isBeginner,
    guidedMode,
    selectSpread,
    clearSpread,
  };
}
