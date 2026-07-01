import { ref, computed } from 'vue';
import {
  fetchAuthConfig,
  sendCode,
  verifyCode,
  completeCodeSignup,
  loginWithPassword,
  registerWithPassword,
  resetPassword as resetPasswordApi,
  getMe,
  patchMe,
} from '../utils/authApi.js';

const TOKEN_KEY = 'fst_auth_token';
const USER_KEY = 'fst_auth_user';

const token = ref(localStorage.getItem(TOKEN_KEY) || '');
const user = ref(null);
const canPredict = ref(true);
const nextPredictAt = ref(null);
const loading = ref(false);
const forceLogin = ref(true);
const skipVerify = ref(true);
const configLoaded = ref(false);

function loadCachedUser() {
  if (!token.value) {
    user.value = null;
    return;
  }
  try {
    const raw = localStorage.getItem(USER_KEY);
    user.value = raw ? JSON.parse(raw) : null;
  } catch {
    localStorage.removeItem(USER_KEY);
    user.value = null;
  }
}

function persistUser(nextUser) {
  user.value = nextUser;
  if (nextUser) {
    localStorage.setItem(USER_KEY, JSON.stringify(nextUser));
  } else {
    localStorage.removeItem(USER_KEY);
  }
}

loadCachedUser();

export function useAuth() {
  const isAuthenticated = computed(() => !!token.value && !!user.value);

  async function loadAuthConfig() {
    const cfg = await fetchAuthConfig();
    forceLogin.value = cfg.force_login;
    skipVerify.value = cfg.skip_verify;
    configLoaded.value = true;
    return cfg;
  }

  function setToken(newToken) {
    token.value = newToken;
    if (newToken) {
      localStorage.setItem(TOKEN_KEY, newToken);
    } else {
      localStorage.removeItem(TOKEN_KEY);
      persistUser(null);
    }
  }

  async function fetchMe() {
    if (!token.value) {
      persistUser(null);
      return null;
    }
    loading.value = true;
    try {
      const data = await getMe(token.value);
      persistUser(data.user);
      canPredict.value = data.can_predict;
      nextPredictAt.value = data.next_predict_at || null;
      return data;
    } catch (err) {
      if (err.status === 401) {
        logout();
      }
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function restoreSession() {
    const savedToken = localStorage.getItem(TOKEN_KEY);
    if (!savedToken) {
      logout();
      return false;
    }

    token.value = savedToken;
    loadCachedUser();

    try {
      await fetchMe();
      return isAuthenticated.value;
    } catch {
      return isAuthenticated.value;
    }
  }

  async function initAuth() {
    await loadAuthConfig();
    return restoreSession();
  }

  async function verifyEmailCode(email, code) {
    return verifyCode(email, code);
  }

  async function finishCodeSignup(email, nickname, password, confirmPassword) {
    const data = await completeCodeSignup(email, nickname, password, confirmPassword);
    setToken(data.token);
    persistUser(data.user);
    await fetchMe();
    return data;
  }

  function applyAuthResult(data) {
    setToken(data.token);
    persistUser(data.user);
    return fetchMe();
  }

  async function loginWithCode(email, code, nickname = '') {
    const data = await verifyCode(email, code);
    if (data.needs_setup) {
      const err = new Error('请先完成资料填写');
      err.needsSetup = true;
      throw err;
    }
    await applyAuthResult(data);
    return data;
  }

  async function login(email, password) {
    const data = await loginWithPassword(email, password);
    setToken(data.token);
    persistUser(data.user);
    await fetchMe();
    return data;
  }

  async function register(email, password, nickname = '', code = '') {
    const data = await registerWithPassword(email, password, nickname, code);
    setToken(data.token);
    persistUser(data.user);
    await fetchMe();
    return data;
  }

  async function resetPassword(email, code, password, confirmPassword) {
    return resetPasswordApi(email, code, password, confirmPassword);
  }

  async function requestCode(email) {
    return sendCode(email);
  }

  async function updateNickname(nickname) {
    const data = await patchMe(token.value, { nickname });
    persistUser(data.user);
    return data;
  }

  async function updateExperienceMode(mode) {
    const data = await patchMe(token.value, { experienceMode: mode });
    persistUser(data.user);
    return data;
  }

  function logout() {
    setToken('');
    canPredict.value = true;
    nextPredictAt.value = null;
  }

  function getToken() {
    return token.value;
  }

  return {
    token,
    user,
    canPredict,
    nextPredictAt,
    loading,
    forceLogin,
    skipVerify,
    configLoaded,
    isAuthenticated,
    loadAuthConfig,
    initAuth,
    restoreSession,
    fetchMe,
    login,
    loginWithCode,
    verifyEmailCode,
    finishCodeSignup,
    applyAuthResult,
    register,
    resetPassword,
    requestCode,
    updateNickname,
    updateExperienceMode,
    logout,
    getToken,
  };
}
