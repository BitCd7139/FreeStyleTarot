const API_BASE = '/api';

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

export async function fetchAnnouncement() {
  const response = await fetch(`${API_BASE}/announcement`);
  return parseResponse(response);
}
