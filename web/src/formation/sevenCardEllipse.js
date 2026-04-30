export default {
    name: "七星阵",
    cardCount: 7,
    match: (gridCards) => {
      if (gridCards.length !== 7) return null;
  
      // 1. 找到中心牌
      const avgX = gridCards.reduce((sum, c) => sum + c.gridX, 0) / 7;
      const avgY = gridCards.reduce((sum, c) => sum + c.gridY, 0) / 7;
      const sortedByDist = [...gridCards].sort((a, b) => {
        const distA = Math.pow(a.gridX - avgX, 2) + Math.pow(a.gridY - avgY, 2);
        const distB = Math.pow(b.gridX - avgX, 2) + Math.pow(b.gridY - avgY, 2);
        return distA - distB;
      });
      const centerCard = sortedByDist[0];
      const orbitCards = sortedByDist.slice(1);
  
      // 2. 严谨性：验证中心牌是否真的在“中间”
      // 周围牌的平均距离应该远大于 0，且中心牌到均值的距离应较小
      const distToCenter = Math.sqrt(Math.pow(centerCard.gridX - avgX, 2) + Math.pow(centerCard.gridY - avgY, 2));
      if (distToCenter > 1.5) return null; 
  
      // 3. 将周围 6 张牌按角度排序（从正上方 12 点钟方向顺时针）
      const withAngles = orbitCards.map(c => ({
        ...c,
        angle: Math.atan2(c.gridY - centerCard.gridY, c.gridX - centerCard.gridX)
      })).sort((a, b) => a.angle - b.angle);
  
      return [
        { id: withAngles[0].id, meaning: "性格/当前处境" },
        { id: withAngles[1].id, meaning: "财运/物质" },
        { id: withAngles[2].id, meaning: "兄弟/沟通" },
        { id: withAngles[3].id, meaning: "家庭/根基" },
        { id: withAngles[4].id, meaning: "恋爱/创造力" },
        { id: withAngles[5].id, meaning: "健康/工作" },
        { id: centerCard.id, meaning: "核心建议" }
      ];
    }
};