export async function* predictStream(payload) {
    const response = await fetch('/api/predict', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });
  
    if (!response.ok) throw new Error("Network response was not ok");
  
    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    let buffer = "";
  
    while (true) {
      const { value, done } = await reader.read();
      if (done) break;
  
      buffer += decoder.decode(value, { stream: true });
      const lines = buffer.split("\n");
      buffer = lines.pop();
  
      for (const line of lines) {
        const trimmed = line.trim();
        if (trimmed.startsWith("data: ")) {
          const data = trimmed.substring(6);
          if (data === "[DONE]") return;
          try {
            const parsed = JSON.parse(data);
            yield parsed.content; 
          } catch (e) {}
        }
      }
    }
  }