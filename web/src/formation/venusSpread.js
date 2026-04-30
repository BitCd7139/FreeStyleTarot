export default {
  name: "维纳斯爱情牌阵",
  cardCount: 7,
  match: (gridCards) => {
    if (gridCards.length !== 7) return null;

    // 1. 识别底部的顶点牌（结果）
    const sortedByY = [...gridCards].sort((a, b) => b.gridY - a.gridY);
    const resultCard = sortedByY[0]; // Y值最大的（最下方）

    // 2. 验证层级：除去顶点，剩下的牌必须在 Y 轴上明显高于顶点
    const others = gridCards.filter(c => c.id !== resultCard.id);
    const maxYOfOthers = Math.max(...others.map(c => c.gridY));
    if (resultCard.gridY <= maxYOfOthers) return null; // 顶点必须在最下面

    // 3. 左右分流：按 X 轴相对于顶点牌的位置划分为左翼和右翼
    const leftSide = others.filter(c => c.gridX < resultCard.gridX).sort((a, b) => a.gridY - b.gridY);
    const rightSide = others.filter(c => c.gridX > resultCard.gridX).sort((a, b) => a.gridY - b.gridY);

    // 严谨性检查：必须是左右对称的 3-3 结构
    if (leftSide.length !== 3 || rightSide.length !== 3) return null;

    // 4. 验证 V 字形的纵向梯度：侧翼的牌在 Y 轴上应该有至少两个高度
    const leftYs = new Set(leftSide.map(c => Math.round(c.gridY)));
    if (leftYs.size < 2) return null; // 如果左边三张水平排成一线，则不是维纳斯

    return [
      { id: leftSide[0].id, meaning: "询问者的真心" },
      { id: leftSide[1].id, meaning: "询问者受到的影响" },
      { id: leftSide[2].id, meaning: "询问者的外在表现" },
      { id: rightSide[0].id, meaning: "对方的真心" },
      { id: rightSide[1].id, meaning: "对方受到的影响" },
      { id: rightSide[2].id, meaning: "对方的外在表现" },
      { id: resultCard.id, meaning: "两人的最终关系结果" }
    ];
  }
};