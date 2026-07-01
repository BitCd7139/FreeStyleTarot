/**
 * @typedef {{ key: string, phase: string, label: string, loading: boolean, text: string }} StreamBlock
 */

/**
 * @param {import('vue').Ref<StreamBlock[]>} blocksRef
 * @param {{ type: string, phase?: string, label?: string, text?: string }} event
 */
export function applyStreamEvent(blocksRef, event) {
  const blocks = blocksRef.value;

  switch (event.type) {
    case 'intro':
    case 'outro': {
      const existing = blocks.find((b) => b.phase === event.type);
      if (existing) {
        existing.text += event.text || '';
        existing.loading = false;
      } else {
        blocks.push({
          key: event.type,
          phase: event.type,
          label: '',
          loading: false,
          text: event.text || '',
        });
      }
      break;
    }
    case 'phase_start': {
      blocks.push({
        key: event.phase,
        phase: event.phase,
        label: event.label || '',
        loading: true,
        text: '',
      });
      break;
    }
    case 'delta': {
      const block = blocks.find((b) => b.phase === event.phase);
      if (block) {
        block.text += event.text || '';
        block.loading = false;
      }
      break;
    }
    case 'phase_end': {
      const block = blocks.find((b) => b.phase === event.phase);
      if (block) {
        block.loading = false;
      }
      break;
    }
    default:
      break;
  }

  blocksRef.value = [...blocks];
}

/** @param {StreamBlock[]} blocks */
export function blocksToAnswer(blocks) {
  return blocks.map((b) => b.text).join('');
}

/** @returns {StreamBlock[]} */
export function createEmptyStreamBlocks() {
  return [];
}
