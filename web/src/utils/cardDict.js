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