import { authHeaders } from './authApi.js';

async function parseResponse(response) {
  const data = await response.json().catch(() => ({}));
  if (!response.ok) {
    const err = new Error(data.error || '请求失败');
    err.status = response.status;
    err.data = data;
    throw err;
  }
  return data;
}

export async function fetchPredictHistory(token) {
  const response = await fetch('/api/auth/predict-history', {
    headers: authHeaders(token),
  });
  return parseResponse(response);
}
