export default {
  id: 'elementCross',
  name: "元素十字牌阵",
  description: "以四元素（火风水土）加中心整合构成十字，从精神、思想、情感、物质四维度分析。",
  useCase: "想从四元素维度（精神、思想、情感、物质）全面理解一件事时使用。",
  cardCount: 6,
  slots: [
    { gridX: 5, gridY: 2, meaning: "空气，神秘和自我认同" },
    { gridX: 2, gridY: 5, meaning: "水，困扰和情绪" },
    { gridX: 8, gridY: 5, meaning: "火，意志与目标" },
    { gridX: 5, gridY: 8, meaning: "土，现实与基础" },
    { gridX: 5, gridY: 5, meaning: "基本问题" },
    { gridX: 5, gridY: 11, meaning: "精神，解决问题的方案" }
  ],
  // 十字：四臂汇聚于中心
  connections: [[0, 4], [1, 4], [2, 4], [3, 4], [5, 4]],
  match: (gridCards) => {
    if (gridCards.length !== 6) return null;

    const avgX = gridCards.reduce((s, c) => s + c.gridX, 0) / 6;
    const avgY = gridCards.reduce((s, c) => s + c.gridY, 0) / 6;
    const sortedByDist = [...gridCards].sort((a, b) =>
      (Math.pow(a.gridX - avgX, 2) + Math.pow(a.gridY - avgY, 2)) -
      (Math.pow(b.gridX - avgX, 2) + Math.pow(b.gridY - avgY, 2))
    );
    const centerCard = sortedByDist[0];
    const others = gridCards.filter(c => c.id !== centerCard.id);

    const topCard = others.find(c => c.gridY < centerCard.gridY && Math.abs(c.gridX - centerCard.gridX) <= 2);
    const bottomCard = others.find(c => c.gridY > centerCard.gridY && Math.abs(c.gridX - centerCard.gridX) <= 2);
    const leftCard = others.find(c => c.gridX < centerCard.gridX && Math.abs(c.gridY - centerCard.gridY) <= 2);
    const rightCard = others.find(c => c.gridX > centerCard.gridX && Math.abs(c.gridY - centerCard.gridY) <= 2);

    if (!topCard || !bottomCard || !leftCard || !rightCard) return null;

    return [
      { id: topCard.id, meaning: "火 / 精神：意志与目标" },
      { id: leftCard.id, meaning: "风 / 思想：理性与沟通" },
      { id: rightCard.id, meaning: "水 / 情感：感受与关系" },
      { id: bottomCard.id, meaning: "土 / 物质：现实与基础" },
      { id: centerCard.id, meaning: "核心 / 整合：综合结论" }
    ];
  }
};
