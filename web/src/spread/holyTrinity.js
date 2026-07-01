export default {
    id: 'holyTrinity',
    name: "圣三角牌阵",
    description: "三张牌横向排列，分别代表过去、现在与未来，是最经典的入门牌阵。",
    useCase: "适合时间线类问题，例如「这件事会如何发展」「过去现在未来的趋势」。",
    cardCount: 3,
    slots: [
      { gridX: 2, gridY: 5, meaning: "过去" },
      { gridX: 5, gridY: 5, meaning: "现在" },
      { gridX: 8, gridY: 5, meaning: "未来" }
    ],
    // 卡牌间关系连接（slot 索引对）：过去→现在→未来 时间线
    connections: [[0, 1], [1, 2]],
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