/**
 * 基础邮箱格式校验（与后端 ValidateEmail 规则大致对齐）
 */
export function validateEmail(email) {
  const trimmed = email.trim().toLowerCase();
  if (!trimmed) {
    return { ok: false, message: '请输入邮箱' };
  }
  if (trimmed.length > 255) {
    return { ok: false, message: '邮箱格式无效' };
  }
  const emailRe = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!emailRe.test(trimmed)) {
    return { ok: false, message: '邮箱格式无效' };
  }
  return { ok: true, value: trimmed };
}

/**
 * 注册密码：6-20 位，数字 / 字母 / 符号至少包含两种
 */
export function validateRegisterPassword(password) {
  const length = [...password].length;
  if (length < 6) {
    return { ok: false, message: '密码长度至少 6 位' };
  }
  if (length > 20) {
    return { ok: false, message: '密码长度不能超过 20 位' };
  }

  let hasDigit = false;
  let hasLetter = false;
  let hasSymbol = false;

  for (const ch of password) {
    if (/\d/.test(ch)) {
      hasDigit = true;
    } else if (/[a-zA-Z]/.test(ch)) {
      hasLetter = true;
    } else {
      hasSymbol = true;
    }
  }

  const categories = [hasDigit, hasLetter, hasSymbol].filter(Boolean).length;
  if (categories < 2) {
    return { ok: false, message: '密码需包含数字、字母、符号中的至少两种' };
  }

  return { ok: true };
}

export function validateNickname(nickname) {
  const trimmed = nickname.trim();
  if (!trimmed) return { ok: true };
  if (trimmed.length < 2 || trimmed.length > 32) {
    return { ok: false, message: '昵称长度需在 2-32 字符之间' };
  }
  return { ok: true, value: trimmed };
}

export function validateRequiredNickname(nickname) {
  const trimmed = nickname.trim();
  if (!trimmed) {
    return { ok: false, message: '请输入昵称' };
  }
  if (trimmed.length < 2 || trimmed.length > 32) {
    return { ok: false, message: '昵称长度需在 2-32 字符之间' };
  }
  return { ok: true, value: trimmed };
}
