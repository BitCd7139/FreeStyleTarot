export default {
    name: "圣三角牌阵",
    cardCount: 3,
    match: (gridCards) => {
      // 规则1：必须是3张牌
      if (gridCards.length !== 3) return null;
  
      // 规则2：Y轴高度差不能太大（认定为近似在同一水平线上）
      const ys = gridCards.map(c => c.gridY);
      if (Math.max(...ys) - Math.min(...ys) > 4) return null;
  
      // 重点：无视原本顺序，只按 X 坐标（从左到右）在网格上排序
      const sortedLeftToRight = [...gridCards].sort((a, b) => a.gridX - b.gridX);
  
      // 根据空间位置（左，中，右）赋予含义
      return [
        { id: sortedLeftToRight[0].id, meaning: "过去" },
        { id: sortedLeftToRight[1].id, meaning: "现在" },
        { id: sortedLeftToRight[2].id, meaning: "未来" }
      ];
    }
  };