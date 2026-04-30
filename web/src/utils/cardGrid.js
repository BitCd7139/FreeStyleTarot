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