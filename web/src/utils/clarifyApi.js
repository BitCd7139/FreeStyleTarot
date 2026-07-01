import { authHeaders } from './authApi.js';

const API_BASE = '/api';

export async function fetchClarify(payload, token) {
  const response = await fetch(`${API_BASE}/clarify`, {
    method: 'POST',
    headers: authHeaders(token),
    body: JSON.stringify(payload),
  });

  const data = await response.json().catch(() => ({}));
  if (!response.ok) {
    const err = new Error(data.error || '意图分析失败');
    err.status = response.status;
    err.data = data;
    throw err;
  }
  return data;
}
