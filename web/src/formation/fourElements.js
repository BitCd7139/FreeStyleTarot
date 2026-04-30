export default {
    name: "四元素牌阵",
    cardCount: 4,
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