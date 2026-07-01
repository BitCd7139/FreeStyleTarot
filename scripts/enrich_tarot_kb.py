#!/usr/bin/env python3
"""
补充 tarot_knowledge_base.json 中的 facing、patternTags、impact 字段。

字段说明：
- facing: []string 支持多重含义
  - "front": 人物正面朝前
  - "left": 人物明确面朝左侧
  - "right": 人物明确面朝右侧
  - "surrounding": 环绕/发散式意象
  - "downward": 向下流动
  - "multi-direction": 多方向意象

- patternTags: []string 模式触发标签
  - choice: 天然带选择意（如2号牌）
  - guide: 有支配/引导意象
  - block: 破坏性牌
  - pivot: 转折/枢纽牌（如命运之轮）
  - revelation: 启示/觉醒牌（如审判）
  - duality: 二元对立牌
  - threshold: 门槛/过渡牌
  - completion: 完结/成就牌

- impact: 1-3 冲击力等级
  - 3: 极高冲击力（塔、太阳）
  - 2: 高冲击力（大阿卡纳）
  - 1: 普通冲击力（其他）
"""

import json
import os

# ============================================================
# Facing 定义 - 支持多重含义
# ============================================================

FACING = {
    # Major Arcana 大阿卡纳
    "thefool": ["left"],
    "themagician": ["front", "surrounding"],
    "thehighpriestess": ["duality"],
    "theempress": ["right"],
    "theemperor": ["front", "left"],
    "thehierophant": ["duality"],
    "thelovers": ["front", "duality"],
    "thechariot": ["front", "duality"],
    "thestrength": ["left"],
    "thehermit": ["left", "downward"],
    "wheeloffortune": ["front", "surrounding"],
    "justice": ["front", "right", "duality"],
    "thehangedman": ["front", "threshold"],
    "death": ["right", "threshold", "block"],
    "temperance": ["front", "duality"],
    "thedevil": ["front", "duality"],
    "thetower": ["duality", "block", "surrounding"],
    "thestar": ["front", "downward"],
    "themoon": ["front", "surrounding"],
    "thesun": ["front", "surrounding"],
    "judgement": ["front", "downward"],
    "theworld": ["front", "surrounding", "left"],

    # Wands 权杖
    "aceofwands": ["front"],
    "twoofwands": ["right", "choice"],
    "threeofwands": ["front"],
    "fourofwands": ["front"],
    "fiveofwands": ["front"],
    "sixofwands": ["right"],
    "sevenofwands": ["front"],
    "eightofwands": ["front"],
    "nineofwands": ["left"],
    "tenofwands": ["right"],
    "pageofwands": ["right"],
    "knightofwands": ["right"],
    "queenofwands": ["front"],
    "kingofwands": ["front"],

    # Cups 圣杯
    "aceofcups": ["front"],
    "twoofcups": ["front", "duality", "choice"],
    "threeofcups": ["right"],
    "fourofcups": ["left"],
    "fiveofcups": ["front"],
    "sixofcups": ["right"],
    "sevenofcups": ["front"],
    "eightofcups": ["left"],
    "nineofcups": ["front"],
    "tenofcups": ["front"],
    "pageofcups": ["front"],
    "knightofcups": ["right"],
    "queenofcups": ["front"],
    "kingofcups": ["front"],

    # Swords 宝剑
    "aceofswords": ["front"],
    "twoofswords": ["front", "choice"],
    "threeofswords": ["front"],
    "fourofswords": ["front"],
    "fiveofswords": ["right"],
    "sixofswords": ["right"],
    "sevenofswords": ["left"],
    "eightofswords": ["front"],
    "nineofswords": ["right"],
    "tenofswords": ["front"],
    "pageofswords": ["front"],
    "knightofswords": ["right"],
    "queenofswords": ["front"],
    "kingofswords": ["right"],

    # Pentacles 星币
    "aceofpentacles": ["front"],
    "twoofpentacles": ["front", "choice"],
    "threeofpentacles": ["front"],
    "fourofpentacles": ["front"],
    "fiveofpentacles": ["front"],
    "sixofpentacles": ["front"],
    "sevenofpentacles": ["front"],
    "eightofpentacles": ["front"],
    "nineofpentacles": ["front"],
    "tenofpentacles": ["front"],
    "pageofpentacles": ["front"],
    "knightofpentacles": ["right"],
    "queenofpentacles": ["front"],
    "kingofpentacles": ["front"],
}

# ============================================================
# Pattern Tags 定义
# ============================================================

PATTERN_TAGS = {
    # === Major Arcana 大阿卡纳 ===
    # 仅保留与视觉流动模式直接相关的标签
    # 标签含义: choice=选择型, guide=指导型, block=堵塞型, pivot=枢纽型, duality=二元对立型
    # 无视觉流动相关标签的卡牌设为空列表

    "thefool": [],
    "themagician": [],
    "thehighpriestess": [],
    "theempress": [],
    "theemperor": [],
    "thehierophant": [],
    "thelovers": ["choice", "duality"],   # 牌面两人对立 + 天使
    "thechariot": ["guide"],              # 战士驾驭战车，支配/引导意象
    "thestrength": [],
    "thehermit": [],
    "wheeloffortune": ["pivot"],          # 轮盘旋转，枢纽意象
    "justice": [],
    "thehangedman": [],
    "death": [],
    "temperance": [],
    "thedevil": ["guide"],                # 恶魔支配两人
    "thetower": ["block"],                # 闪电击塔，破坏/阻断意象
    "thestar": [],
    "themoon": [],
    "thesun": [],
    "judgement": [],
    "theworld": [],

    # === Minor Arcana 小阿卡纳 ===
    # 2号牌天然带选择意
    "twoofwands": ["choice"],
    "twoofcups": ["choice"],
    "twoofswords": ["choice"],
    "twoofpentacles": ["choice"],
}

# ============================================================
# Impact 定义
# ============================================================

IMPACT = {
    # Major Arcana
    "thefool": 2,
    "themagician": 2,
    "thehighpriestess": 2,
    "theempress": 2,
    "theemperor": 2,
    "thehierophant": 2,
    "thelovers": 2,
    "thechariot": 2,
    "thestrength": 2,
    "thehermit": 2,
    "wheeloffortune": 2,
    "justice": 2,
    "thehangedman": 2,
    "death": 2,
    "temperance": 2,
    "thedevil": 2,
    "thetower": 3,
    "thestar": 2,
    "themoon": 2,
    "thesun": 3,
    "judgement": 2,
    "theworld": 2,
    # Minor Arcana
    "aceofwands": 1, "twoofwands": 1, "threeofwands": 1, "fourofwands": 1,
    "fiveofwands": 1, "sixofwands": 1, "sevenofwands": 1, "eightofwands": 1,
    "nineofwands": 1, "tenofwands": 1, "pageofwands": 1, "knightofwands": 1,
    "queenofwands": 1, "kingofwands": 1,
    "aceofcups": 1, "twoofcups": 1, "threeofcups": 1, "fourofcups": 1,
    "fiveofcups": 1, "sixofcups": 1, "sevenofcups": 1, "eightofcups": 1,
    "nineofcups": 1, "tenofcups": 1, "pageofcups": 1, "knightofcups": 1,
    "queenofcups": 1, "kingofcups": 1,
    "aceofswords": 1, "twoofswords": 1, "threeofswords": 1, "fourofswords": 1,
    "fiveofswords": 1, "sixofswords": 1, "sevenofswords": 1, "eightofswords": 1,
    "nineofswords": 1, "tenofswords": 1, "pageofswords": 1, "knightofswords": 1,
    "queenofswords": 1, "kingofswords": 1,
    "aceofpentacles": 1, "twoofpentacles": 1, "threeofpentacles": 1, "fourofpentacles": 1,
    "fiveofpentacles": 1, "sixofpentacles": 1, "sevenofpentacles": 1, "eightofpentacles": 1,
    "nineofpentacles": 1, "tenofpentacles": 1, "pageofpentacles": 1, "knightofpentacles": 1,
    "queenofpentacles": 1, "kingofpentacles": 1,
}


def main():
    # 获取项目根目录
    base_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    kb_path = os.path.join(base_dir, "storage", "tarot_knowledge_base.json")

    with open(kb_path, "r", encoding="utf-8") as f:
        data = json.load(f)

    # 统计
    updated = 0
    missing_facing = []
    missing_impact = []
    missing_tags = []

    for card_key, card_data in data.items():
        changed = False

        # facing: 先删除再重新插入（始终使用最新定义）
        del card_data["facing"]
        if card_key in FACING:
            card_data["facing"] = FACING[card_key]
        else:
            card_data["facing"] = ["front"]
            missing_facing.append(card_key)
        changed = True

        # patternTags: 先删除再重新插入（始终使用最新定义）
        del card_data["patternTags"]
        if card_key in PATTERN_TAGS:
            card_data["patternTags"] = PATTERN_TAGS[card_key]
        else:
            card_data["patternTags"] = []
        changed = True

        # impact: 先删除再重新插入（始终使用最新定义）
        del card_data["impact"]
        if card_key in IMPACT:
            card_data["impact"] = IMPACT[card_key]
        else:
            card_data["impact"] = 1
            missing_impact.append(card_key)
        changed = True

        if changed:
            updated += 1

    # 检查未定义的标签
    for card_key in PATTERN_TAGS:
        if card_key not in data:
            missing_tags.append(card_key)

    # 写回文件
    with open(kb_path, "w", encoding="utf-8") as f:
        json.dump(data, f, ensure_ascii=False, indent=2)

    # 输出统计
    print("=" * 60)
    print("Knowledge Base Enrichment Complete")
    print("=" * 60)
    print(f"Updated cards: {updated}")
    print(f"Total cards: {len(data)}")
    print()
    if missing_facing:
        print(f"[!] Missing facing definitions ({len(missing_facing)}):")
        for k in missing_facing[:10]:
            print(f"   - {k}")
        if len(missing_facing) > 10:
            print(f"   ... and {len(missing_facing) - 10} more")
        print()
    else:
        print("[OK] All facing fields defined")
    if missing_impact:
        print(f"[!] Missing impact definitions ({len(missing_impact)}):")
        for k in missing_impact[:10]:
            print(f"   - {k}")
        if len(missing_impact) > 10:
            print(f"   ... and {len(missing_impact) - 10} more")
        print()
    else:
        print("[OK] All impact fields defined")
    if missing_tags:
        print(f"[!] patternTags reference cards not in KB ({len(missing_tags)}):")
        for k in missing_tags:
            print(f"   - {k}")
    else:
        print("[OK] patternTags references valid")
    print()
    print(f"Output file: {kb_path}")

    # 打印大阿卡纳的 facing 和 patternTags 摘要
    print()
    print("=" * 60)
    print("Major Arcana Summary (Please Verify)")
    print("=" * 60)
    major_arcana = [k for k in sorted(data.keys()) if data[k].get("arcana") == "Major"]
    for card_key in major_arcana:
        card = data[card_key]
        facing = card.get("facing", [])
        tags = card.get("patternTags", [])
        impact = card.get("impact", "N/A")
        facing_str = "|".join(facing) if facing else "N/A"
        tags_str = ",".join(tags) if tags else "-"
        print(f"{card_key:22} facing=[{facing_str:25}] tags=[{tags_str:30}] impact={impact}")


if __name__ == "__main__":
    main()
