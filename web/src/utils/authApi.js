const API_BASE = '/api';

export function authHeaders(token) {
  const headers = { 'Content-Type': 'application/json' };
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  return headers;
}

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

export async function fetchAuthConfig() {
  const response = await fetch(`${API_BASE}/auth/config`);
  return parseResponse(response);
}

export async function sendCode(email) {
  const response = await fetch(`${API_BASE}/auth/send-code`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify({ email }),
  });
  return parseResponse(response);
}

export async function verify(email, code, nickname = '') {
  const response = await fetch(`${API_BASE}/auth/verify`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify({ email, code, nickname }),
  });
  return parseResponse(response);
}

export async function verifyCode(email, code) {
  const response = await fetch(`${API_BASE}/auth/verify-code`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify({ email, code }),
  });
  return parseResponse(response);
}

export async function completeCodeSignup(email, nickname, password, confirmPassword) {
  const response = await fetch(`${API_BASE}/auth/complete-code-signup`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify({
      email,
      nickname,
      password,
      confirm_password: confirmPassword,
    }),
  });
  return parseResponse(response);
}

export async function loginWithPassword(email, password) {
  const response = await fetch(`${API_BASE}/auth/login`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify({ email, password }),
  });
  return parseResponse(response);
}

export async function registerWithPassword(email, password, nickname = '', code = '') {
  const response = await fetch(`${API_BASE}/auth/register`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify({ email, password, nickname, code }),
  });
  return parseResponse(response);
}

export async function resetPassword(email, code, password, confirmPassword) {
  const response = await fetch(`${API_BASE}/auth/reset-password`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify({
      email,
      code,
      password,
      confirm_password: confirmPassword,
    }),
  });
  return parseResponse(response);
}

export async function getMe(token) {
  const response = await fetch(`${API_BASE}/auth/me`, {
    headers: authHeaders(token),
  });
  return parseResponse(response);
}

export async function patchMe(token, nickname) {
  const response = await fetch(`${API_BASE}/auth/me`, {
    method: 'PATCH',
    headers: authHeaders(token),
    body: JSON.stringify({ nickname }),
  });
  return parseResponse(response);
}
