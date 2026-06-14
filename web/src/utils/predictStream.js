export async function* predictStream(payload, submitInfo, token) {
  let postdir = '/api/predict';
  if (submitInfo === 'Prompt') {
    postdir = '/api/prompt';
  } else if (submitInfo === 'Result') {
    postdir = '/api/predict';
  }

  const headers = { 'Content-Type': 'application/json' };
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(postdir, {
    method: 'POST',
    headers,
    body: JSON.stringify(payload),
  });

  if (response.status === 401) {
    const err = new Error('登录已失效，请重新登录');
    err.status = 401;
    throw err;
  }

  if (response.status === 429) {
    const data = await response.json().catch(() => ({}));
    const mins = data.wait_minutes || '?';
    const err = new Error(`提问次数已达上限，请 ${mins} 分钟后再试`);
    err.status = 429;
    err.data = data;
    throw err;
  }

  if (!response.ok) {
    const data = await response.json().catch(() => ({}));
    throw new Error(data.error || 'Network response was not ok');
  }

  const reader = response.body.getReader();
  const decoder = new TextDecoder();
  let buffer = '';

  while (true) {
    const { value, done } = await reader.read();
    if (done) break;

    buffer += decoder.decode(value, { stream: true });
    const lines = buffer.split('\n');
    buffer = lines.pop();

    for (const line of lines) {
      const trimmed = line.trim();
      if (trimmed.startsWith('data: ')) {
        const data = trimmed.substring(6);
        if (data === '[DONE]') return;
        try {
          const parsed = JSON.parse(data);
          if (parsed.type === 'status') {
            yield {
              type: 'status',
              phase: parsed.phase,
              label: parsed.label,
              streamPublic: !!parsed.stream_public,
            };
          } else if (typeof parsed.content === 'string') {
            yield { type: 'content', content: parsed.content };
          }
        } catch {
          // ignore malformed chunks
        }
      }
    }
  }
}
