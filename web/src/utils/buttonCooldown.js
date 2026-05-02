import { ref, onUnmounted } from 'vue';

export function buttonCooldown(initialSeconds = 30) {
  const count = ref(0);
  const isPending = ref(false); 
  let timer = null;

  const start = () => {
    count.value = initialSeconds;
    isPending.value = true;
    
    // 清除可能存在的旧定时器
    if (timer) clearInterval(timer);

    timer = setInterval(() => {
      count.value--;
      if (count.value <= 0) {
        stop();
      }
    }, 1000);
  };

  const stop = () => {
    if (timer) {
      clearInterval(timer);
      timer = null;
    }
    count.value = 0;
    isPending.value = false;
  };

  // 组件销毁时自动清理
  onUnmounted(stop);

  return {
    count,
    isPending,
    start,
    stop
  };
}