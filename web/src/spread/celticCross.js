export default {
    name: "凯尔特十字牌阵",
    cardCount: 10,
    match: (gridCards) => {
      if (gridCards.length !== 10) return null;
  
      // 规则1：将牌按 X 轴排序，切分为左侧十字区(6张) 和 右侧权杖区(4张)
      const sortedByX = [...gridCards].sort((a, b) => a.gridX - b.gridX);
      const leftCrossGroup = sortedByX.slice(0, 6);
      const rightStaffGroup = sortedByX.slice(6, 10);
  
      // 规则2：右侧权杖区必须明显在左侧十字区的右边
      const maxLeftX = Math.max(...leftCrossGroup.map(c => c.gridX));
      const minRightX = Math.min(...rightStaffGroup.map(c => c.gridX));
      if (maxLeftX >= minRightX) return null;
  
      // --- 解析右侧权杖区 ---
      // 按 Y 轴降序排序（在屏幕上 Y 值越大越靠下，所以这代表 从下到上 的空间顺序）
      const staffBottomToTop = [...rightStaffGroup].sort((a, b) => b.gridY - a.gridY);
  
      // --- 解析左侧十字区 ---
      // 找到十字架的四个顶点（最左、最右、最上、最下）
      const leftPoint = [...leftCrossGroup].sort((a, b) => a.gridX - b.gridX)[0];
      const rightPoint = [...leftCrossGroup].sort((a, b) => b.gridX - a.gridX)[0];
      const topPoint = [...leftCrossGroup].sort((a, b) => a.gridY - b.gridY)[0];
      const bottomPoint = [...leftCrossGroup].sort((a, b) => b.gridY - a.gridY)[0];
  
      // 中心的两张牌：剔除四个顶点后剩下的两张
      const edgeIds = new Set([leftPoint.id, rightPoint.id, topPoint.id, bottomPoint.id]);
      const centerCards = leftCrossGroup.filter(c => !edgeIds.has(c.id));
  
      // 如果形状太乱导致顶点重合，无法析出2张中心牌，则匹配失败
      if (centerCards.length !== 2) return null;
  
      // 完全通过物理空间位置映射塔罗含义！
      return [
        // 中心重叠区
        { id: centerCards[0].id, meaning: "现状" },
        { id: centerCards[1].id, meaning: "障碍" },
        // 十字架四周
        { id: leftPoint.id, meaning: "过去" },
        { id: rightPoint.id, meaning: "未来" },
        { id: topPoint.id, meaning: "显意识 (目标)" },
        { id: bottomPoint.id, meaning: "潜意识 (基础)" },
        // 侧边柱状区（从下到上）
        { id: staffBottomToTop[0].id, meaning: "自我" },
        { id: staffBottomToTop[1].id, meaning: "环境" },
        { id: staffBottomToTop[2].id, meaning: "希望与恐惧" },
        { id: staffBottomToTop[3].id, meaning: "最终结果" }
      ];
    }
  };