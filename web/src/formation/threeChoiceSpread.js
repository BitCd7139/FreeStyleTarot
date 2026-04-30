export default {
    name: "三选一牌阵",
    cardCount: 7,
    match: (gridCards) => {
      if (gridCards.length !== 7) return null;
  
      // 1. 识别最左侧的“求测者”牌
      const sortedByX = [...gridCards].sort((a, b) => a.gridX - b.gridX);
      const requesterCard = sortedByX[0];
      const rightGroup = sortedByX.slice(1);
  
      // 严谨性检查：求测者牌的 X 坐标必须明显小于右侧所有牌
      const minRightX = Math.min(...rightGroup.map(c => c.gridX));
      if (requesterCard.gridX >= minRightX) return null;
  
      // 2. 解析右侧的三个选项（3行 x 2列）
      // 首先按 Y 轴（从上到下）排序，以便切分行
      const sortedByY = [...rightGroup].sort((a, b) => a.gridY - b.gridY);
  
      // 严谨性检查：验证是否至少存在三个不同的高度层级
      const yCoords = sortedByY.map(c => c.gridY);
      const rowGap = Math.max(...yCoords) - Math.min(...yCoords);
      if (rowGap < 4) return null; // 如果高度差太小，可能只是排成了一横排
  
      // 将 6 张牌切分为三组（顶行、中行、底行）
      const rowTopRaw = sortedByY.slice(0, 2);
      const rowMidRaw = sortedByY.slice(2, 4);
      const rowBottomRaw = sortedByY.slice(4, 6);
  
      // 对每一行内的两张牌按 X 轴排序，区分“过程”和“结果”
      const option1 = rowTopRaw.sort((a, b) => a.gridX - b.gridX);
      const option2 = rowMidRaw.sort((a, b) => a.gridX - b.gridX);
      const option3 = rowBottomRaw.sort((a, b) => a.gridX - b.gridX);
  
      // 3. 返回含义映射
      return [
        { id: requesterCard.id, meaning: "求测者当前状态" },
        
        // 选项一 (Top Row)
        { id: option1[0].id, meaning: "选项一：发展过程" },
        { id: option1[1].id, meaning: "选项一：最终结果" },
        
        // 选项二 (Middle Row)
        { id: option2[0].id, meaning: "选项二：发展过程" },
        { id: option2[1].id, meaning: "选项二：最终结果" },
        
        // 选项三 (Bottom Row)
        { id: option3[0].id, meaning: "选项三：发展过程" },
        { id: option3[1].id, meaning: "选项三：最终结果" }
      ];
    }
  };