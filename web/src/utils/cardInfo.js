export const ASPECT_RATIO = 1.75;
export const allCardNames = ["aceofcups", "aceofpentacles", "aceofswords", "aceofwands", "death", "eightofcups", "eightofpentacles", "eightofswords", "eightofwands", "fiveofcups", "fiveofpentacles", "fiveofswords", "fiveofwands", "fourofcups", "fourofpentacles", "fourofswords", "fourofwands", "judgement", "justice", "kingofcups", "kingofpentacles", "kingofswords", "kingofwands", "knightofcups", "knightofpentacles", "knightofswords", "knightofwands", "nineofcups", "nineofpentacles", "nineofswords", "nineofwands", "pageofcups", "pageofpentacles", "pageofswords", "pageofwands", "queenofcups", "queenofpentacles", "queenofswords", "queenofwands", "sevenofcups", "sevenofpentacles", "sevenofswords", "sevenofwands", "sixofcups", "sixofpentacles", "sixofswords", "sixofwands", "temperance", "tenofcups", "tenofpentacles", "tenofswords", "tenofwands", "thechariot", "thedevil", "theemperor", "theempress", "thefool", "thehangedman", "thehermit", "thehierophant", "thehighpriestess", "thelovers", "themagician", "themoon", "thestar", "thestrength", "thesun", "thetower", "theworld", "threeofcups", "threeofpentacles", "threeofswords", "threeofwands", "twoofcups", "twoofpentacles", "twoofswords", "twoofwands", "wheeloffortune"];

const TAROT_DICT = {
    "thefool": "愚者", "themagician": "魔术师", "thehighpriestess": "女祭司", 
    "theempress": "女皇", "theemperor": "皇帝", "thehierophant": "教皇", 
    "thelovers": "恋人", "thechariot": "战车", "thestrength": "力量", 
    "thehermit": "隐士", "wheeloffortune": "命运之轮", "justice": "正义", 
    "thehangedman": "倒吊人", "death": "死神", "temperance": "节制", 
    "thedevil": "恶魔", "thetower": "塔", "thestar": "星星", 
    "themoon": "月亮", "thesun": "太阳", "judgement": "审判", "theworld": "世界",
    "aceofcups": "圣杯首牌", "twoofcups": "圣杯二", "threeofcups": "圣杯三", 
    "fourofcups": "圣杯四", "fiveofcups": "圣杯五", "sixofcups": "圣杯六", 
    "sevenofcups": "圣杯七", "eightofcups": "圣杯八", "nineofcups": "圣杯九", 
    "tenofcups": "圣杯十", "pageofcups": "圣杯侍从", "knightofcups": "圣杯骑士", 
    "queenofcups": "圣杯女王", "kingofcups": "圣杯国王",
    "aceofpentacles": "星币首牌", "twoofpentacles": "星币二", "threeofpentacles": "星币三", 
    "fourofpentacles": "星币四", "fiveofpentacles": "星币五", "sixofpentacles": "星币六", 
    "sevenofpentacles": "星币七", "eightofpentacles": "星币八", "nineofpentacles": "星币九", 
    "tenofpentacles": "星币十", "pageofpentacles": "星币侍从", "knightofpentacles": "星币骑士", 
    "queenofpentacles": "星币女王", "kingofpentacles": "星币国王",
    "aceofswords": "宝剑首牌", "twoofswords": "宝剑二", "threeofswords": "宝剑三", 
    "fourofswords": "宝剑四", "fiveofswords": "宝剑五", "sixofswords": "宝剑六", 
    "sevenofswords": "宝剑七", "eightofswords": "宝剑八", "nineofswords": "宝剑九", 
    "tenofswords": "宝剑十", "pageofswords": "宝剑侍从", "knightofswords": "宝剑骑士", 
    "queenofswords": "宝剑女王", "kingofswords": "宝剑国王",
    "aceofwands": "权杖首牌", "twoofwands": "权杖二", "threeofwands": "权杖三", 
    "fourofwands": "权杖四", "fiveofwands": "权杖五", "sixofwands": "权杖六", 
    "sevenofwands": "权杖七", "eightofwands": "权杖八", "nineofwands": "权杖九", 
    "tenofwands": "权杖十", "pageofwands": "权杖侍从", "knightofwands": "权杖骑士", 
    "queenofwands": "权杖女王", "kingofwands": "权杖国王"
  };

export const getName = (name) => {
  return TAROT_DICT[name] || name;
};

/** Rider-Waite 大阿卡纳编号 0–21 */
const MAJOR_ARCANA_ORDER = [
  'thefool', 'themagician', 'thehighpriestess', 'theempress', 'theemperor',
  'thehierophant', 'thelovers', 'thechariot', 'thestrength', 'thehermit',
  'wheeloffortune', 'justice', 'thehangedman', 'death', 'temperance',
  'thedevil', 'thetower', 'thestar', 'themoon', 'thesun', 'judgement', 'theworld',
];

const MAJOR_ARCANA = new Set(MAJOR_ARCANA_ORDER);

const MAJOR_ARCANA_KEYWORDS = ['大阿卡纳', '大阿尔卡那', '大阿卡罗纳', 'major', '大牌'];

const MAJOR_NUMBER_BY_ID = Object.fromEntries(
  MAJOR_ARCANA_ORDER.map((id, index) => [id, index])
);

const SUIT_META = {
  cups: { keywords: ['圣杯', '杯', 'cups'] },
  pentacles: { keywords: ['星币', '钱币', '金币', 'pentacles', 'pentacle'] },
  swords: { keywords: ['宝剑', '剑', 'swords', 'sword'] },
  wands: { keywords: ['权杖', '杖', 'wands', 'wand'] },
};

const RANK_META = {
  ace: { keywords: ['首牌', 'ace', 'a', '一'] },
  two: { keywords: ['二', '2', 'two', 'ii'] },
  three: { keywords: ['三', '3', 'three', 'iii'] },
  four: { keywords: ['四', '4', 'four', 'iv'] },
  five: { keywords: ['五', '5', 'five', 'v'] },
  six: { keywords: ['六', '6', 'six', 'vi'] },
  seven: { keywords: ['七', '7', 'seven', 'vii'] },
  eight: { keywords: ['八', '8', 'eight', 'viii'] },
  nine: { keywords: ['九', '9', 'nine', 'ix'] },
  ten: { keywords: ['十', '10', 'ten', 'x'] },
  page: { keywords: ['侍从', 'page', '侍', '11'] },
  knight: { keywords: ['骑士', 'knight', '12'] },
  queen: { keywords: ['女王', '皇后', 'queen', '13'] },
  king: { keywords: ['国王', 'king', '14'] },
};

function majorNumberKeywords(id) {
  const n = MAJOR_NUMBER_BY_ID[id];
  if (n === undefined) return [];
  const num = String(n);
  const padded = n < 10 ? `0${n}` : String(n);
  return [num, padded, `第${num}号`, `no${num}`, `no.${num}`];
}

function parseCardMeta(id) {
  if (MAJOR_ARCANA.has(id)) {
    return {
      suit: null,
      rank: null,
      isMajor: true,
      keywords: [...MAJOR_ARCANA_KEYWORDS, ...majorNumberKeywords(id)],
    };
  }
  for (const [suit, meta] of Object.entries(SUIT_META)) {
    if (id.includes(`of${suit}`)) {
      const prefix = id.replace(`of${suit}`, '');
      const rankMeta = RANK_META[prefix];
      return {
        suit,
        rank: prefix,
        isMajor: false,
        keywords: [...meta.keywords, ...(rankMeta?.keywords ?? [])],
      };
    }
  }
  return { suit: null, rank: null, isMajor: false, keywords: [] };
}

const CARD_SEARCH_INDEX = allCardNames.map((id) => {
  const meta = parseCardMeta(id);
  const cnName = TAROT_DICT[id] || id;
  return {
    id,
    cnName,
    searchText: [id, cnName, ...meta.keywords].join(' ').toLowerCase(),
    meta,
  };
});

export function searchCards(query) {
  const q = query.trim().toLowerCase();
  if (!q) return CARD_SEARCH_INDEX;
  const tokens = q.split(/\s+/).filter(Boolean);
  return CARD_SEARCH_INDEX.filter((card) =>
    tokens.every((token) => card.searchText.includes(token))
  );
}

export function getCardImageUrl(name) {
  return new URL(`../assets/tarots/${name}.jpeg`, import.meta.url).href;
}