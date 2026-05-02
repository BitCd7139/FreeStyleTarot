/**
 * 改进后的离散化函数
 * @param {Array} cards - 卡牌数组
 * @param {Object} options - 配置项
 * @param {number} options.cardWidth - 单张卡牌的显示宽度
 * @param {number} options.cardHeight - 单张卡牌的显示高度
 * @param {number} options.paddingFactor - 边缘留白比例 (默认0.5，表示四周留出半张牌的宽度)
 */
export function discretizeCards(cards, cardWidth, cardHeight, { paddingFactor = 0.5 } = {}) {
  if (!cards || cards.length === 0) return [];

  // 1. 获取所有卡牌中心点的边界
  const xs = cards.map(c => c.x);
  const ys = cards.map(c => c.y);
  
  let minX = Math.min(...xs);
  let maxX = Math.max(...xs);
  let minY = Math.min(...ys);
  let maxY = Math.max(...ys);

  // 2. 核心改进：外扩边界 (Buffer)
  // 不再以点为边界，而是以“视觉包围盒”为边界，并增加 padding
  const paddingX = cardWidth * paddingFactor;
  const paddingY = cardHeight * paddingFactor;
  
  minX -= paddingX;
  maxX += paddingX;
  minY -= paddingY;
  maxY += paddingY;

  const width = maxX - minX;
  const height = maxY - minY;
  
  const gridSize = 11; // 0-11 索引

  return cards.map(c => {
    let gridX, gridY;

    // 3. 处理单张卡牌或范围极小的情况
    if (width < 1) {
      gridX = Math.floor(gridSize / 2);
    } else {
      // 归一化计算：(当前值 - 最小值) / 总宽度
      const ratioX = (c.x - minX) / width;
      // 限制范围在 0-1 之间，防止 padding 导致的溢出
      gridX = Math.round(Math.max(0, Math.min(1, ratioX)) * gridSize);
    }

    if (height < 1) {
      gridY = Math.floor(gridSize / 2);
    } else {
      const ratioY = (c.y - minY) / height;
      gridY = Math.round(Math.max(0, Math.min(1, ratioY)) * gridSize);
    }

    return {
      ...c,
      gridX,
      gridY
    };
  });
}


/**
 * 还原离散化后的卡牌坐标
 * @param {Array} cards - 卡牌数组 (需包含 x, y, gridX, gridY, 以及可选的 order)
 * @param {number} cardWidth - 单张卡牌的显示宽度
 * @param {number} cardHeight - 单张卡牌的显示高度
 * @param {Object} options - 配置项
 * @param {number} options.paddingFactor - 边缘留白比例 (必须与离散化时一致)
 * @param {number} options.offsetFactor - 重叠时的偏移比例 (默认0.1)
 */
export function recoverFromGrids(cards, cardWidth, cardHeight, { paddingFactor = 0.5, offsetFactor = 0.1 } = {}) {
  if (!cards || cards.length === 0) return [];

  // 1. 重新计算离散化时使用的边界 (保持与 discretizeCards 逻辑完全一致)
  const xs = cards.map(c => c.x);
  const ys = cards.map(c => c.y);
  
  let minX = Math.min(...xs);
  let maxX = Math.max(...xs);
  let minY = Math.min(...ys);
  let maxY = Math.max(...ys);

  const paddingX = cardWidth * paddingFactor;
  const paddingY = cardHeight * paddingFactor;
  
  // 这里的 minX, maxX 是考虑了 padding 后的视觉包围盒边界
  const minPaddedX = minX - paddingX;
  const maxPaddedX = maxX + paddingX;
  const minPaddedY = minY - paddingY;
  const maxPaddedY = maxY + paddingY;

  const totalWidth = maxPaddedX - minPaddedX;
  const totalHeight = maxPaddedY - minPaddedY;
  
  const gridSize = 11; // 必须与离散化一致

  // 2. 准备处理重叠逻辑
  // 按照 order 排序，确保 order 大的在处理时能感知到前面的牌
  const sortedCards = [...cards].sort((a, b) => (a.order || 0) - (b.order || 0));
  
  // 用于记录每个网格点已经放置了多少张牌
  const gridOccupancy = {}; 
  const offsetX = cardWidth * offsetFactor;
  const offsetY = cardHeight * offsetFactor;

  return sortedCards.map(c => {
    let fixedX, fixedY;

    // 3. 计算基础网格坐标
    if (totalWidth < 1) {
      fixedX = minX; // 如果没有宽度，回到原始中心
    } else {
      // 逆运算: ratio = gridX / gridSize
      fixedX = minPaddedX + (c.gridX / gridSize) * totalWidth;
    }

    if (totalHeight < 1) {
      fixedY = minY;
    } else {
      fixedY = minPaddedY + (c.gridY / gridSize) * totalHeight;
    }

    // 4. 处理重叠偏移
    const gridKey = `${c.gridX}-${c.gridY}`;
    const count = gridOccupancy[gridKey] || 0;
    
    // 如果该位置已有牌，根据 count 叠加偏移
    const finalX = fixedX + count * offsetX;
    const finalY = fixedY + count * offsetY;

    // 更新计数器
    gridOccupancy[gridKey] = count + 1;

    return {
      ...c,
      fixedX: finalX,
      fixedY: finalY,
      // 增加 zIndex 建议，确保 order 大的在视觉上也处于上方
      zIndex: (c.order || 0) + count 
    };
  });
}