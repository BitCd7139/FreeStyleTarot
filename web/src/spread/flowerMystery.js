export default {
  id: 'flower-mystery',
  name: "花朵之谜牌阵",
  description: "四张牌从左到右，如四朵神秘之花依次绽放：风信子带来消息，罂粟藏起秘密，忍冬散发甜蜜，莲花照见真相。",
  useCase: "想看清一段关系中传递的信息、隐藏的秘密、表面的甜蜜与内在的真相时使用。",
  cardCount: 4,
  slots: [
    { gridX: 1, gridY: 1, meaning: "消息的含义（鸢尾花）" },
    { gridX: 3, gridY: 1, meaning: "潜在的秘密（康乃馨）" },
    { gridX: 5, gridY: 1, meaning: "外在的甜蜜（百合花）" },
    { gridX: 7, gridY: 1, meaning: "关系的真相（玫瑰花）" }
  ],
  // 花朵藤蔓般的连线：消息→秘密→甜蜜→真相，以及消息直指真相
  connections: [[0, 1], [1, 2], [2, 3], [0, 3]],
  match: (gridCards) => {
    if (gridCards.length !== 4) return null;

    // 严格按从左到右排列：gridX 必须递增
    const sorted = [...gridCards].sort((a, b) => a.gridX - b.gridX);
    for (let i = 1; i < sorted.length; i++) {
      if (sorted[i].gridX <= sorted[i - 1].gridX) return null;
    }

    return [
      { id: sorted[0].id, meaning: "消息的含义（风信子）" },
      { id: sorted[1].id, meaning: "潜在的秘密（罂粟）" },
      { id: sorted[2].id, meaning: "外在的甜蜜（忍冬）" },
      { id: sorted[3].id, meaning: "关系的真相（莲花）" }
    ];
  }
};