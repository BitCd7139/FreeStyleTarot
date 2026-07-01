import { ref } from 'vue';
import { fetchPredictHistory } from '../utils/predictHistoryApi.js';
import { useAuth } from './useAuth.js';

const items = ref([]);
const limit = ref(3);
const loading = ref(false);
const error = ref('');
const loaded = ref(false);

export function usePredictHistory() {
  const { getToken } = useAuth();

  async function loadHistory() {
    const token = getToken();
    if (!token) {
      items.value = [];
      error.value = '';
      loaded.value = false;
      return;
    }
    loading.value = true;
    error.value = '';
    try {
      const data = await fetchPredictHistory(token);
      items.value = data.items || [];
      limit.value = data.limit ?? 3;
      loaded.value = true;
    } catch (err) {
      error.value = err.message || '加载提问历史失败';
      items.value = [];
    } finally {
      loading.value = false;
    }
  }

  async function refreshHistory() {
    if (!loaded.value) return;
    await loadHistory();
  }

  function resetHistory() {
    items.value = [];
    limit.value = 3;
    error.value = '';
    loaded.value = false;
  }

  return {
    items,
    limit,
    loading,
    error,
    loaded,
    loadHistory,
    refreshHistory,
    resetHistory,
  };
}
