export default {
    id: 'choiceSpread',
    name: "二选一牌阵",
    description: "面对两难选择时使用：一张现状牌，加上两条路径的近期发展与最终结果。",
    useCase: "在两个选项之间犹豫不决时使用，例如「该选 A 还是 B」。",
    cardCount: 5,
    slots: [
      { gridX: 2, gridY: 5, meaning: "现状/询问者心境" },
      { gridX: 5, gridY: 2, meaning: "选项一：近期发展" },
      { gridX: 8, gridY: 2, meaning: "选项一：最终结果" },
      { gridX: 5, gridY: 8, meaning: "选项二：近期发展" },
      { gridX: 8, gridY: 8, meaning: "选项二：最终结果" }
    ],
    // 两条分叉路径：现状→选项一(近期→结果)，现状→选项二(近期→结果)
    connections: [[0, 1], [1, 2], [0, 3], [3, 4]],
    match: (gridCards) => {
      if (gridCards.length !== 5) return null;

    // 1. 识别现状牌：按 X 轴排序，最左边的一张
    const sortedByX = [...gridCards].sort((a, b) => a.gridX - b.gridX);
    const currentStatus = sortedByX[0];
    const others = gridCards.filter(c => c.id !== currentStatus.id);

    // 【验证1】：现状牌必须明显在其他牌的左侧
    const minXOfOthers = Math.min(...others.map(c => c.gridX));
    if (currentStatus.gridX >= minXOfOthers) return null;

    // 2. 将剩下的 4 张按 Y 轴切分为“上方路径”和“下方路径”
    const sortedByY = [...others].sort((a, b) => a.gridY - b.gridY);
    const topHalf = sortedByY.slice(0, 2);
    const bottomHalf = sortedByY.slice(2, 4);

    // 【验证2】：上方路径整体必须在下方路径的上面
    const maxYOfTop = Math.max(...topHalf.map(c => c.gridY));
    const minYOfBottom = Math.min(...bottomHalf.map(c => c.gridY));
    if (maxYOfTop >= minYOfBottom) return null; // 纵向重叠或位置不对

    // 【验证3】：现状牌的纵向位置应该处于上下两条路径之间
    // (允许一定误差，但不能比最顶上的还高，或比最底下的还低)
    if (currentStatus.gridY < topHalf[0].gridY || currentStatus.gridY > bottomHalf[1].gridY) {
      // 如果需要极度严谨可以开启此项，普通场景可略过
    }

    // 3. 在每条路径内部，按 X 轴排序区分“近期”和“结果”
    const topPath = topHalf.sort((a, b) => a.gridX - b.gridX);
    const bottomPath = bottomHalf.sort((a, b) => a.gridX - b.gridX);

    // 【验证4】：每条路径的横向必须有明显的递进（不是垂直堆叠）
    if (topPath[1].gridX - topPath[0].gridX < 0.5) return null;
    if (bottomPath[1].gridX - bottomPath[0].gridX < 0.5) return null;
  
      return [
        { id: currentStatus.id, meaning: "现状/询问者心境" },
        { id: topPath[0].id, meaning: "选项一：近期发展" },
        { id: topPath[1].id, meaning: "选项一：最终结果" },
        { id: bottomPath[0].id, meaning: "选项二：近期发展" },
        { id: bottomPath[1].id, meaning: "选项二：最终结果" }
      ];
    }
};