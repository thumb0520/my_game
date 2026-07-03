package game

import "math/rand"

// Slot 装备槽位
type Slot string

const (
	SlotWeapon  Slot = "武器"
	SlotHelmet  Slot = "头盔"
	SlotArmor   Slot = "胸甲"
	SlotBoots   Slot = "靴子"
	SlotAcc1    Slot = "饰品1"
	SlotAcc2    Slot = "饰品2"
)

// ItemType 物品类型
type ItemType string

const (
	Weapon     ItemType = "weapon"
	Helmet     ItemType = "helmet"
	Armor      ItemType = "armor"
	Boots      ItemType = "boots"
	Accessory  ItemType = "accessory"
	Consumable ItemType = "consumable"
)

// Rarity 品质
type Rarity string

const (
	Common    Rarity = "普通"
	Uncommon  Rarity = "优秀"
	Rare      Rarity = "稀有"
	Epic      Rarity = "史诗"
	Legendary Rarity = "传说"
)

// RarityColors 品质对应的颜色名（用于显示）
var RarityColors = map[Rarity]string{
	Common:    "white",
	Uncommon:  "green",
	Rare:      "blue",
	Epic:      "purple",
	Legendary: "orange",
}

// Item 物品
type Item struct {
	Name        string
	Type        ItemType
	Rarity      Rarity
	Description string
	Effects     map[string]int // 属性加成：atk, def, hp, crit, dodge, heal 等
	Slot        Slot           // 装备槽位（仅装备有效）
}

// SlotForType 根据物品类型获取装备槽位
func SlotForType(t ItemType) Slot {
	switch t {
	case Weapon:
		return SlotWeapon
	case Helmet:
		return SlotHelmet
	case Armor:
		return SlotArmor
	case Boots:
		return SlotBoots
	case Accessory:
		return SlotAcc1
	default:
		return ""
	}
}

// 词缀系统
type Affix struct {
	Name   string
	Stat   string
	MinVal int
	MaxVal int
}

var PrefixAffixes = []Affix{
	{"锋利的", "atk", 2, 8},
	{"坚固的", "def", 2, 8},
	{"迅捷的", "dodge", 1, 5},
	{"致命的", "crit", 2, 8},
	{"健壮的", "hp", 10, 30},
	{"烈焰", "atk", 5, 12},
	{"寒冰", "def", 5, 12},
	{"暗影", "crit", 5, 10},
	{"神圣", "hp", 20, 50},
}

var SuffixAffixes = []Affix{
	{"之力", "atk", 1, 6},
	{"守护", "def", 1, 6},
	{"疾风", "dodge", 1, 4},
	{"破甲", "crit", 1, 5},
	{"生命", "hp", 5, 20},
}

// 基础装备模板
type ItemTemplate struct {
	Name   string
	Type   ItemType
	Slot   Slot
	Effects map[string]int
}

var WeaponTemplates = []ItemTemplate{
	{"铁剑", Weapon, SlotWeapon, map[string]int{"atk": 5}},
	{"钢刀", Weapon, SlotWeapon, map[string]int{"atk": 8}},
	{"长矛", Weapon, SlotWeapon, map[string]int{"atk": 6, "crit": 3}},
	{"法杖", Weapon, SlotWeapon, map[string]int{"atk": 7, "hp": 10}},
	{"匕首", Weapon, SlotWeapon, map[string]int{"atk": 4, "crit": 5, "dodge": 3}},
	{"战斧", Weapon, SlotWeapon, map[string]int{"atk": 10}},
	{"细剑", Weapon, SlotWeapon, map[string]int{"atk": 4, "crit": 8}},
	{"巨锤", Weapon, SlotWeapon, map[string]int{"atk": 12, "dodge": -3}},
}

var HelmetTemplates = []ItemTemplate{
	{"皮盔", Helmet, SlotHelmet, map[string]int{"def": 3}},
	{"铁盔", Helmet, SlotHelmet, map[string]int{"def": 5}},
	{"法师帽", Helmet, SlotHelmet, map[string]int{"def": 2, "hp": 15}},
	{"刺客兜帽", Helmet, SlotHelmet, map[string]int{"def": 2, "dodge": 4}},
}

var ArmorTemplates = []ItemTemplate{
	{"皮甲", Armor, SlotArmor, map[string]int{"def": 5}},
	{"锁子甲", Armor, SlotArmor, map[string]int{"def": 8}},
	{"板甲", Armor, SlotArmor, map[string]int{"def": 12, "dodge": -2}},
	{"法师袍", Armor, SlotArmor, map[string]int{"def": 3, "hp": 20}},
}

var BootsTemplates = []ItemTemplate{
	{"皮靴", Boots, SlotBoots, map[string]int{"def": 2, "dodge": 2}},
	{"铁靴", Boots, SlotBoots, map[string]int{"def": 4}},
	{"轻便靴", Boots, SlotBoots, map[string]int{"dodge": 5}},
	{"重甲靴", Boots, SlotBoots, map[string]int{"def": 6, "dodge": -1}},
}

var AccessoryTemplates = []ItemTemplate{
	{"力量戒指", Accessory, SlotAcc1, map[string]int{"atk": 3}},
	{"守护护符", Accessory, SlotAcc1, map[string]int{"def": 3, "hp": 10}},
	{"敏捷之靴", Accessory, SlotAcc1, map[string]int{"dodge": 4, "crit": 2}},
	{"生命宝石", Accessory, SlotAcc1, map[string]int{"hp": 25}},
}

// GenerateLoot 根据地牢深度生成随机装备
func GenerateLoot(depth int) *Item {
	// 随机选择装备类型
	templates := [][]ItemTemplate{WeaponTemplates, HelmetTemplates, ArmorTemplates, BootsTemplates, AccessoryTemplates}
	pool := templates[rand.Intn(len(templates))]
	tmpl := pool[rand.Intn(len(pool))]

	item := &Item{
		Name:        tmpl.Name,
		Type:        tmpl.Type,
		Slot:        tmpl.Slot,
		Description: "",
		Effects:     make(map[string]int),
	}

	// 复制基础属性
	for k, v := range tmpl.Effects {
		item.Effects[k] = v + depth // 深度越深，基础属性越高
	}

	// 决定品质
	rarityRoll := rand.Intn(100)
	switch {
	case rarityRoll < 5:
		item.Rarity = Legendary
	case rarityRoll < 15:
		item.Rarity = Epic
	case rarityRoll < 35:
		item.Rarity = Rare
	case rarityRoll < 60:
		item.Rarity = Uncommon
	default:
		item.Rarity = Common
	}

	// 根据品质添加词缀
	affixCount := 0
	switch item.Rarity {
	case Legendary:
		affixCount = 3
	case Epic:
		affixCount = 2
	case Rare:
		affixCount = 1
	}

	prefixUsed := ""
	suffixUsed := ""
	for i := 0; i < affixCount; i++ {
		if i%2 == 0 && len(PrefixAffixes) > 0 {
			affix := PrefixAffixes[rand.Intn(len(PrefixAffixes))]
			val := affix.MinVal + rand.Intn(affix.MaxVal-affix.MinVal+1) + depth/2
			item.Effects[affix.Stat] += val
			prefixUsed = affix.Name
		} else if len(SuffixAffixes) > 0 {
			affix := SuffixAffixes[rand.Intn(len(SuffixAffixes))]
			val := affix.MinVal + rand.Intn(affix.MaxVal-affix.MinVal+1) + depth/2
			item.Effects[affix.Stat] += val
			suffixUsed = affix.Name
		}
	}

	// 生成名字
	if prefixUsed != "" {
		item.Name = prefixUsed + item.Name
	}
	if suffixUsed != "" {
		item.Name = item.Name + "·" + suffixUsed
	}

	return item
}

// GenerateConsumableLoot 生成消耗品
func GenerateConsumableLoot(depth int) *Item {
	items := []*Item{
		{Name: "小型药水", Type: Consumable, Rarity: Common, Description: "恢复 30 点生命值", Effects: map[string]int{"heal": 30}},
		{Name: "中型药水", Type: Consumable, Rarity: Uncommon, Description: "恢复 60 点生命值", Effects: map[string]int{"heal": 60}},
		{Name: "大型药水", Type: Consumable, Rarity: Rare, Description: "恢复 100 点生命值", Effects: map[string]int{"heal": 100}},
		{Name: "解毒草", Type: Consumable, Rarity: Common, Description: "清除中毒状态", Effects: map[string]int{"cure_poison": 1}},
		{Name: "力量药剂", Type: Consumable, Rarity: Uncommon, Description: "临时增加 5 点攻击力", Effects: map[string]int{"temp_atk": 5}},
		{Name: "铁皮药剂", Type: Consumable, Rarity: Uncommon, Description: "临时增加 5 点防御力", Effects: map[string]int{"temp_def": 5}},
	}

	// 深度越深，好物品概率越高
	idx := rand.Intn(len(items))
	if depth > 3 && idx < 2 {
		idx = 2 + rand.Intn(len(items)-2)
	}
	if idx >= len(items) {
		idx = len(items) - 1
	}
	return items[idx]
}
