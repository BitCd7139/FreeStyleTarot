import holyTrinity from './holyTrinity';
import celticCross from './celticCross';
import choiceSpread from './choiceSpread';
import fourElements from './fourElements';
import hexagram from './hexagram';
import venusSpread from './venusSpread';
import sevenCardEllipse from './sevenCardEllipse';
import threeChoiceSpread from './threeChoiceSpread';
import singleCard from './singleCardSpread';
import fiveCardCross from './fiveCardCross';
import universalFiveCard from './universalFiveCard';
import elementCross from './elementCross';
import goldenTriangle from './goldenTriangle';
import flowerMystery from './flowerMystery';

// 导出所有支持的牌阵模板
export const SPREAD_TEMPLATES = [
  holyTrinity,
  celticCross,
  choiceSpread,
  fourElements,
  hexagram,
  venusSpread,
  sevenCardEllipse,
  threeChoiceSpread,
  singleCard,
  fiveCardCross,
  universalFiveCard,
  elementCross,
  goldenTriangle,
  flowerMystery,
  // 以后在这里无缝扩展
];

/**
 * 返回适合新手页展示的牌阵目录（UI 安全，不含 match 函数）
 * 仅包含同时具备 id / description / slots 且 slots 数量与 cardCount 一致的模板
 */
export function getSpreadCatalog() {
  return SPREAD_TEMPLATES
    .filter(t => t.id && t.description && Array.isArray(t.slots) && t.slots.length === t.cardCount)
    .map(t => ({
      id: t.id,
      name: t.name,
      description: t.description,
      useCase: t.useCase || '',
      cardCount: t.cardCount,
      freestyle: t.freestyle || false,
      slots: t.slots.map(s => ({ gridX: s.gridX, gridY: s.gridY, meaning: s.meaning })),
      connections: Array.isArray(t.connections) ? t.connections.map(pair => [pair[0], pair[1]]) : [],
    }));
}