export default {
  id: 'goldenTriangle',
  name: "金三角牌阵",
  description: "三张牌构成三角形，分别代表目标、基础与行动，简洁有力。",
  useCase: "对一件事做最精简的三维分析（目标、基础、行动）时使用。",
  cardCount: 3,
  slots: [
    { gridX: 5, gridY: 2, meaning: "目标：理想方向" },
    { gridX: 2, gridY: 8, meaning: "基础：已有条件" },
    { gridX: 8, gridY: 8, meaning: "行动：建议路径" }
  ],
  // 三角形三边
  connections: [[0, 1], [1, 2], [2, 0]],
  match: (gridCards) => {
    if (gridCards.length !== 3) return null;
    const sorted = [...gridCards].sort((a, b) => a.gridY - b.gridY);
    const top = sorted[0];
    const bottomLeft = sorted[1].gridX <= sorted[2].gridX ? sorted[1] : sorted[2];
    const bottomRight = sorted[1].gridX > sorted[2].gridX ? sorted[1] : sorted[2];
    if (bottomRight.gridX - bottomLeft.gridX < 3) return null;
    if (top.gridY >= bottomLeft.gridY) return null;
    return [
      { id: top.id, meaning: "目标：理想方向" },
      { id: bottomLeft.id, meaning: "基础：已有条件" },
      { id: bottomRight.id, meaning: "行动：建议路径" }
    ];
  }
};
