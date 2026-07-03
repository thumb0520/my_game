package game

import "fmt"

// MetaProgress 局外养成数据
type MetaProgress struct {
	TotalGold       int            `json:"total_gold"`       // 累计获得金币
	TotalRuns       int            `json:"total_runs"`       // 累计通关次数
	BestFloor       int            `json:"best_floor"`       // 最高到达层数
	Upgrades        map[string]int `json:"upgrades"`         // 已购买的升级
	UnlockedClasses []string       `json:"unlocked_classes"` // 已解锁的职业
}

// Upgrade 升级选项
type Upgrade struct {
	ID          string
	Name        string
	Description string
	MaxLevel    int
	BaseCost    int
	CostScale   float64 // 每级成本倍率
	Effect      string  // 效果描述
}

// 可用升级列表
var AvailableUpgrades = []Upgrade{
	{
		ID: "start_hp", Name: "生命强化",
		Description: "增加初始生命值",
		MaxLevel:    10, BaseCost: 50, CostScale: 1.5,
		Effect: "每级 +10 初始 HP",
	},
	{
		ID: "start_atk", Name: "攻击强化",
		Description: "增加初始攻击力",
		MaxLevel:    10, BaseCost: 60, CostScale: 1.5,
		Effect: "每级 +2 初始 ATK",
	},
	{
		ID: "start_def", Name: "防御强化",
		Description: "增加初始防御力",
		MaxLevel:    10, BaseCost: 55, CostScale: 1.5,
		Effect: "每级 +2 初始 DEF",
	},
	{
		ID: "start_gold", Name: "财富祝福",
		Description: "增加初始金币",
		MaxLevel:    10, BaseCost: 40, CostScale: 1.4,
		Effect: "每级 +10 初始金币",
	},
	{
		ID: "start_crit", Name: "致命精准",
		Description: "增加初始暴击率",
		MaxLevel:    5, BaseCost: 100, CostScale: 2.0,
		Effect: "每级 +2% 初始暴击",
	},
	{
		ID: "potion_heal", Name: "药水强化",
		Description: "增加药水恢复量",
		MaxLevel:    5, BaseCost: 80, CostScale: 1.6,
		Effect: "每级药水多恢复 10 HP",
	},
	{
		ID: "loot_bonus", Name: "寻宝直觉",
		Description: "增加掉落率",
		MaxLevel:    5, BaseCost: 120, CostScale: 1.8,
		Effect: "每级 +5% 掉落率",
	},
}

// NewMetaProgress 创建新的养成数据
func NewMetaProgress() *MetaProgress {
	return &MetaProgress{
		Upgrades:        make(map[string]int),
		UnlockedClasses: []string{"战士"}, // 默认解锁战士
	}
}

// GetUpgradeLevel 获取升级等级
func (mp *MetaProgress) GetUpgradeLevel(upgradeID string) int {
	return mp.Upgrades[upgradeID]
}

// GetUpgradeCost 获取升级成本
func (mp *MetaProgress) GetUpgradeCost(upgrade Upgrade) int {
	level := mp.GetUpgradeLevel(upgrade.ID)
	if level >= upgrade.MaxLevel {
		return -1 // 已满级
	}
	cost := float64(upgrade.BaseCost)
	for i := 0; i < level; i++ {
		cost *= upgrade.CostScale
	}
	return int(cost)
}

// CanUpgrade 是否可以升级
func (mp *MetaProgress) CanUpgrade(upgrade Upgrade) bool {
	cost := mp.GetUpgradeCost(upgrade)
	return cost > 0 && mp.TotalGold >= cost
}

// PurchaseUpgrade 购买升级
func (mp *MetaProgress) PurchaseUpgrade(upgrade Upgrade) bool {
	cost := mp.GetUpgradeCost(upgrade)
	if cost <= 0 || mp.TotalGold < cost {
		return false
	}
	mp.TotalGold -= cost
	mp.Upgrades[upgrade.ID]++
	return true
}

// ApplyUpgrades 应用升级到角色
func (mp *MetaProgress) ApplyUpgrades(p *Character) {
	// 生命强化
	if level := mp.GetUpgradeLevel("start_hp"); level > 0 {
		p.MaxHP += level * 10
		p.HP = p.MaxHP
	}

	// 攻击强化
	if level := mp.GetUpgradeLevel("start_atk"); level > 0 {
		p.ATK += level * 2
	}

	// 防御强化
	if level := mp.GetUpgradeLevel("start_def"); level > 0 {
		p.DEF += level * 2
	}

	// 财富祝福
	if level := mp.GetUpgradeLevel("start_gold"); level > 0 {
		p.Gold += level * 10
	}

	// 致命精准
	if level := mp.GetUpgradeLevel("start_crit"); level > 0 {
		p.CRIT += level * 2
	}
}

// GetPotionBonus 获取药水加成
func (mp *MetaProgress) GetPotionBonus() int {
	return mp.GetUpgradeLevel("potion_heal") * 10
}

// GetLootBonus 获取掉落加成
func (mp *MetaProgress) GetLootBonus() int {
	return mp.GetUpgradeLevel("loot_bonus") * 5
}

// IsClassUnlocked 职业是否已解锁
func (mp *MetaProgress) IsClassUnlocked(class Class) bool {
	for _, c := range mp.UnlockedClasses {
		if c == string(class) {
			return true
		}
	}
	return false
}

// UnlockClass 解锁职业
func (mp *MetaProgress) UnlockClass(class Class, cost int) bool {
	if mp.TotalGold < cost {
		return false
	}
	if mp.IsClassUnlocked(class) {
		return false
	}
	mp.TotalGold -= cost
	mp.UnlockedClasses = append(mp.UnlockedClasses, string(class))
	return true
}

// RecordRun 记录一次通关
func (mp *MetaProgress) RecordRun(goldEarned int, floorReached int) {
	mp.TotalRuns++
	mp.TotalGold += goldEarned
	if floorReached > mp.BestFloor {
		mp.BestFloor = floorReached
	}
}

// GetClassName 获取职业中文名
func GetClassName(class Class) string {
	switch class {
	case Warrior:
		return "战士"
	case Ranger:
		return "游侠"
	case Mage:
		return "法师"
	case Rogue:
		return "盗贼"
	default:
		return string(class)
	}
}

// GetClassUnlockCost 获取职业解锁成本
func GetClassUnlockCost(class Class) int {
	switch class {
	case Ranger:
		return 200
	case Mage:
		return 300
	case Rogue:
		return 500
	default:
		return 0
	}
}

// FormatMetaStats 格式化养成统计
func FormatMetaStats(mp *MetaProgress) string {
	return fmt.Sprintf("通关次数: %d | 最高层数: %d | 累计金币: %d",
		mp.TotalRuns, mp.BestFloor, mp.TotalGold)
}
