export default {
    id: 'fourElements',
    name: "四元素牌阵",
    description: "四张牌摆成菱形，对应火、土、风、水四元素，从精神到情感全面看待现状。",
    useCase: "想从行动、情感、思想、物质四个维度全面了解当前处境时使用。",
    cardCount: 4,
    slots: [
      { gridX: 5, gridY: 2, meaning: "火：精神/行动力" },
      { gridX: 5, gridY: 8, meaning: "土：物质/现实" },
      { gridX: 2, gridY: 5, meaning: "风：逻辑/沟通" },
      { gridX: 8, gridY: 5, meaning: "水：情感/潜意识" }
    ],
    // 菱形连接：火→风→土→水→火
    connections: [[0, 2], [2, 1], [1, 3], [3, 0]],
    match: (gridCards) => {
      if (gridCards.length !== 4) return null;
  
      // 分别找出最上、最下、最左、最右
      const top = [...gridCards].sort((a, b) => a.gridY - b.gridY)[0];
      const bottom = [...gridCards].sort((a, b) => b.gridY - a.gridY)[0];
      const left = [...gridCards].sort((a, b) => a.gridX - b.gridX)[0];
      const right = [...gridCards].sort((a, b) => b.gridX - a.gridX)[0];
  
      // 容错检查：如果四个顶点有重复，说明不是清晰的正方形/菱形
      const ids = new Set([top.id, bottom.id, left.id, right.id]);
      if (ids.size !== 4) return null;
  
      return [
        { id: top.id, meaning: "火：精神/行动力" },
        { id: bottom.id, meaning: "土：物质/现实" },
        { id: left.id, meaning: "风：逻辑/沟通" },
        { id: right.id, meaning: "水：情感/潜意识" }
      ];
    }
};