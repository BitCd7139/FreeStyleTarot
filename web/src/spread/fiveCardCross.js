export default {
    name: "圣十字牌阵",
    cardCount: 5,
    match: (gridCards) => {
      if (gridCards.length !== 5) return null;
  
      // 1. 寻找中心牌：几何距离所有牌最近的那张
      const avgX = gridCards.reduce((sum, c) => sum + c.gridX, 0) / 5;
      const avgY = gridCards.reduce((sum, c) => sum + c.gridY, 0) / 5;
  
      const sortedByDistance = [...gridCards].sort((a, b) => {
        const distA = Math.pow(a.gridX - avgX, 2) + Math.pow(a.gridY - avgY, 2);
        const distB = Math.pow(b.gridX - avgX, 2) + Math.pow(b.gridY - avgY, 2);
        return distA - distB;
      });
  
      const centerCard = sortedByDistance[0];
      const outerCards = sortedByDistance.slice(1);
  
      // 2. 将外围四张牌按方向分类
      // 允许一定的偏差（例如 gridX 相等时通过 gridY 辅助判断）
      const topCard = outerCards.find(c => c.gridY < centerCard.gridY && Math.abs(c.gridX - centerCard.gridX) <= 2);
      const bottomCard = outerCards.find(c => c.gridY > centerCard.gridY && Math.abs(c.gridX - centerCard.gridX) <= 2);
      const leftCard = outerCards.find(c => c.gridX < centerCard.gridX && Math.abs(c.gridY - centerCard.gridY) <= 2);
      const rightCard = outerCards.find(c => c.gridX > centerCard.gridX && Math.abs(c.gridY - centerCard.gridY) <= 2);
  
      // 严谨性检查：必须上下左右四个方位各有一张牌
      if (!topCard || !bottomCard || !leftCard || !rightCard) return null;
  
      // 3. 校验纵向和横向跨度（确保不是挤在一起的）
      if (bottomCard.gridY - topCard.gridY < 3) return null;
      if (rightCard.gridX - leftCard.gridX < 3) return null;
  
      return [
        { id: centerCard.id, meaning: "现状：目前的核心" },
        { id: leftCard.id, meaning: "过去/阻碍" },
        { id: rightCard.id, meaning: "未来趋势" },
        { id: topCard.id, meaning: "目标：心愿或最佳可能" },
        { id: bottomCard.id, meaning: "基础：潜意识或原因" }
      ];
    }
  };