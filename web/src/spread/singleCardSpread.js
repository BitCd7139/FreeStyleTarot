export default {
    name: "单张牌启示",
    cardCount: 1,
    match: (gridCards) => {
      if (gridCards.length !== 1) return null;
  
      // 逻辑极其简单，只要是一张牌即匹配
      return [
        { id: gridCards[0].id, meaning: "核心启示：答案或现状" }
      ];
    }
  };