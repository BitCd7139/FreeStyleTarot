export default {
  id: 'universalFiveCard',
  name: "万能五张牌阵",
  description: "五张牌一字排开，以自由流式解读整体走向，适合开放式提问。",
  useCase: "无需固定含义，希望从整体流式角度自由解读一件事的脉络时使用。",
  cardCount: 5,
  freestyle: true,
  slots: [
    { gridX: 1, gridY: 5, meaning: "万能牌1" },
    { gridX: 3, gridY: 5, meaning: "万能牌2" },
    { gridX: 5, gridY: 5, meaning: "万能牌3" },
    { gridX: 7, gridY: 5, meaning: "万能牌4" },
    { gridX: 9, gridY: 5, meaning: "万能牌5" }
  ],
  // 五张一排，顺次相连
  connections: [[0, 1], [1, 2], [2, 3], [3, 4]],
  match: (gridCards) => {
    if (gridCards.length !== 5) return null;
    // 验证大致在同一水平线
    const ys = gridCards.map(c => c.gridY);
    if (Math.max(...ys) - Math.min(...ys) > 2) return null;
    const sorted = [...gridCards].sort((a, b) => a.gridX - b.gridX);
    const meanings = ["1", "2", "3", "4", "5"];
    return sorted.map((c, i) => ({ id: c.id, meaning: meanings[i] }));
  }
};
