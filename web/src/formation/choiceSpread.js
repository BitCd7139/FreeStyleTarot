export default {
    name: "二选一牌阵",
    cardCount: 5,
    match: (gridCards) => {
      if (gridCards.length !== 5) return null;
  
      // 1. 按 X 轴排序，最左边的一张是现状
      const sortedByX = [...gridCards].sort((a, b) => a.gridX - b.gridX);
      const currentStatus = sortedByX[0];
      const others = sortedByX.slice(1);
  
      // 2. 将剩下的4张按 Y 轴切分为“上方路径”和“下方路径”
      const sortedByY = [...others].sort((a, b) => a.gridY - b.gridY);
      const topPath = sortedByY.slice(0, 2).sort((a, b) => a.gridX - b.gridX);
      const bottomPath = sortedByY.slice(2, 4).sort((a, b) => a.gridX - b.gridX);
  
      return [
        { id: currentStatus.id, meaning: "现状/询问者心境" },
        { id: topPath[0].id, meaning: "选项一：近期发展" },
        { id: topPath[1].id, meaning: "选项一：最终结果" },
        { id: bottomPath[0].id, meaning: "选项二：近期发展" },
        { id: bottomPath[1].id, meaning: "选项二：最终结果" }
      ];
    }
};