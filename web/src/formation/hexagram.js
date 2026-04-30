export default {
  name: "六芒星牌阵",
  cardCount: 7,
  match: (gridCards) => {
    if (gridCards.length !== 7) return null;

    // 1. 验证垂直层级：按 Y 轴排序
    const sortedByY = [...gridCards].sort((a, b) => a.gridY - b.gridY);
    const topPoint = sortedByY[0];        // 最高的牌
    const bottomPoint = sortedByY[6];     // 最低的牌
    
    // 严谨性检查：最高和最低之间必须有足够的纵向空间（假设网格间距）
    if (bottomPoint.gridY - topPoint.gridY < 4) return null;

    // 2. 识别中心牌：计算所有牌的几何中心，离中心最近的那张
    const avgX = gridCards.reduce((sum, c) => sum + c.gridX, 0) / 7;
    const avgY = gridCards.reduce((sum, c) => sum + c.gridY, 0) / 7;
    const centerCard = [...gridCards].sort((a, b) => {
      const distA = Math.pow(a.gridX - avgX, 2) + Math.pow(a.gridY - avgY, 2);
      const distB = Math.pow(b.gridX - avgX, 2) + Math.pow(b.gridY - avgY, 2);
      return distA - distB;
    })[0];

    // 中心牌不能是最高或最低的顶点
    if (centerCard.id === topPoint.id || centerCard.id === bottomPoint.id) return null;

    // 3. 剩下的 4 张牌作为“侧翼”，应该分布在中心牌的左右两侧
    const wings = gridCards.filter(c => 
      c.id !== topPoint.id && c.id !== bottomPoint.id && c.id !== centerCard.id
    );
    const leftWings = wings.filter(c => c.gridX < centerCard.gridX).sort((a, b) => a.gridY - b.gridY);
    const rightWings = wings.filter(c => c.gridX > centerCard.gridX).sort((a, b) => a.gridY - b.gridY);

    // 严谨性检查：左右两侧必须各有两个点
    if (leftWings.length !== 2 || rightWings.length !== 2) return null;

    return [
      { id: topPoint.id, meaning: "目标：事情的理想" },
      { id: bottomPoint.id, meaning: "基础：事情的现实" },
      { id: leftWings[0].id, meaning: "过去：已经发生的影响" },
      { id: leftWings[1].id, meaning: "潜意识：内在的动向" },
      { id: rightWings[0].id, meaning: "未来：即将发生的情况" },
      { id: rightWings[1].id, meaning: "显意识：外在的表现" },
      { id: centerCard.id, meaning: "结论：最终的结果" }
    ];
  }
};