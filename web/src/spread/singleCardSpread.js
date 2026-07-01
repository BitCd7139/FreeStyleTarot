export default {
    id: 'singleCard',
    name: "单张牌启示",
    description: "只抽一张牌，快速获得对当下问题的核心启示。",
    useCase: "适合每日指引、简单是非题或快速获取一个核心提示。",
    cardCount: 1,
    slots: [
      { gridX: 5, gridY: 5, meaning: "核心启示：答案或现状" }
    ],
    // 单张牌无连接线
    connections: [],
    match: (gridCards) => {
      if (gridCards.length !== 1) return null;
  
      // 逻辑极其简单，只要是一张牌即匹配
      return [
        { id: gridCards[0].id, meaning: "核心启示：答案或现状" }
      ];
    }
  };