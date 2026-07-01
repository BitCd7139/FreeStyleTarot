/**
 * 牌阵阵位布局工具
 * slots 中的 gridX / gridY 为 0-11 的抽象网格坐标（与 cardGrid.discretizeCards 输出同坐标系）。
 * 通过固定参考画布把网格坐标映射为像素，再等比缩放居中到实际舞台。
 */

export const GRID_SIZE = 11;

/**
 * 构造 MiniTarot 可用的伪卡牌数组，用于新手页牌阵预览
 * @param {Array} slots - 牌阵 slots
 * @param {number} cardWidth
 * @param {number} cardHeight
 * @returns {Array} 伪 drawnCards（牌背 + order 徽章）
 */
export function slotsToPreviewCards(slots, cardWidth, cardHeight) {
  if (!Array.isArray(slots) || slots.length === 0) return [];
  const unitX = cardWidth;
  const unitY = cardHeight * 0.6;
  return slots.map((slot, index) => ({
    id: `preview-${index}`,
    name: 'back',
    x: slot.gridX * unitX,
    y: slot.gridY * unitY,
    isRevealed: true,
    isReversed: false,
    order: index + 1,
    meaning: slot.meaning,
  }));
}

/**
 * 在实际舞台中计算每个 slot 的像素坐标（居中缩放）
 * @param {Array} slots
 * @param {number} stageWidth - 舞台宽
 * @param {number} stageHeight - 舞台高
 * @param {number} cardWidth
 * @param {number} cardHeight
 * @returns {Array} [{ index, slotIndex, x, y, meaning, width, height }] 像素坐标（左上角）
 */
export function layoutSlotsInStage(slots, stageWidth, stageHeight, cardWidth, cardHeight) {
  if (!Array.isArray(slots) || slots.length === 0 || !stageWidth || !stageHeight) return [];
  const unitX = cardWidth;
  const unitY = cardHeight * 0.6;

  const xs = slots.map(s => s.gridX * unitX);
  const ys = slots.map(s => s.gridY * unitY);
  const minX = Math.min(...xs);
  const maxX = Math.max(...xs) + cardWidth;
  const minY = Math.min(...ys);
  const maxY = Math.max(...ys) + cardHeight;

  const refWidth = Math.max(maxX - minX, 1);
  const refHeight = Math.max(maxY - minY, 1);

  const padding = 0.1;
  const availW = stageWidth * (1 - padding * 2);
  const availH = stageHeight * (1 - padding * 2);
  const scale = Math.min(availW / refWidth, availH / refHeight, 1.5);

  const scaledW = refWidth * scale;
  const scaledH = refHeight * scale;
  const offsetX = (stageWidth - scaledW) / 2 - minX * scale;
  const offsetY = (stageHeight - scaledH) / 2 - minY * scale;

  return slots.map((slot, index) => ({
    index,
    slotIndex: index,
    x: slot.gridX * unitX * scale + offsetX,
    y: slot.gridY * unitY * scale + offsetY,
    meaning: slot.meaning,
    width: cardWidth * scale,
    height: cardHeight * scale,
  }));
}

/**
 * 计算牌阵在 fixed 模式下的 bbox（基于统一间距单位）
 * @param {Array} slots
 * @param {number} cardWidth
 * @param {number} cardHeight
 * @returns {{ bboxW: number, bboxH: number }}
 */
export function getSpreadBBox(slots, cardWidth, cardHeight) {
  if (!Array.isArray(slots) || slots.length === 0) return { bboxW: 0, bboxH: 0 };
  const unitX = cardWidth * 0.95;
  const unitY = cardHeight * 0.5;
  const xs = slots.map(s => s.gridX * unitX);
  const ys = slots.map(s => s.gridY * unitY);
  const bboxW = Math.max(...xs) + cardWidth - Math.min(...xs);
  const bboxH = Math.max(...ys) + cardHeight - Math.min(...ys);
  return { bboxW, bboxH };
}

/**
 * 构造 fixed 模式预览卡牌：在统一舞台尺寸内居中
 * 不同牌阵的卡牌显示大小完全一致，舞台尺寸由外部统一指定（强制一致）。
 * @param {Array} slots
 * @param {number} cardWidth - 固定卡牌宽
 * @param {number} cardHeight - 固定卡牌高
 * @param {number} stageWidth - 统一舞台宽
 * @param {number} stageHeight - 统一舞台高
 * @returns {Array} 伪 drawnCards（已居中的绝对坐标）
 */
export function slotsToFixedPreviewCards(slots, cardWidth, cardHeight, stageWidth, stageHeight) {
  if (!Array.isArray(slots) || slots.length === 0) return [];
  const unitX = cardWidth * 0.95;
  const unitY = cardHeight * 0.5;
  const positions = slots.map(s => ({
    x: s.gridX * unitX,
    y: s.gridY * unitY,
  }));
  const minX = Math.min(...positions.map(p => p.x));
  const maxX = Math.max(...positions.map(p => p.x)) + cardWidth;
  const minY = Math.min(...positions.map(p => p.y));
  const maxY = Math.max(...positions.map(p => p.y)) + cardHeight;
  const offsetX = (stageWidth - (maxX - minX)) / 2 - minX;
  const offsetY = (stageHeight - (maxY - minY)) / 2 - minY;
  return slots.map((slot, index) => ({
    id: `fixed-preview-${index}`,
    name: 'back',
    x: positions[index].x + offsetX,
    y: positions[index].y + offsetY,
    isRevealed: true,
    isReversed: false,
    order: index + 1,
    meaning: slot.meaning,
  }));
}

/**
 * 找到距离卡牌中心最近的未占用 slot
 * @param {Array} stageSlots - layoutSlotsInStage 的输出
 * @param {number} cardCenterX
 * @param {number} cardCenterY
 * @param {Set} occupiedSlotIndices - 已占用的 slot 索引
 * @param {number} threshold - 吸附距离阈值，默认 Infinity（始终吸附最近 slot）
 * @returns {Object|null} 最近的 slot 或 null（超出阈值）
 */
export function findNearestSlot(stageSlots, cardCenterX, cardCenterY, occupiedSlotIndices, threshold = Infinity) {
  if (!Array.isArray(stageSlots) || stageSlots.length === 0) return null;
  let best = null;
  let bestDist = Infinity;
  for (const slot of stageSlots) {
    if (occupiedSlotIndices && occupiedSlotIndices.has(slot.slotIndex)) continue;
    const slotCenterX = slot.x + slot.width / 2;
    const slotCenterY = slot.y + slot.height / 2;
    const dist = Math.hypot(cardCenterX - slotCenterX, cardCenterY - slotCenterY);
    if (dist < bestDist && dist <= threshold) {
      bestDist = dist;
      best = slot;
    }
  }
  return best;
}
