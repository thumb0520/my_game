package display

import (
	"fmt"
	"strings"
	"time"

	"dungeonlog/internal/game"

	"github.com/charmbracelet/lipgloss"
)

// Renderer 终端渲染器
type Renderer struct {
	stealthMode bool
}

// NewRenderer 创建渲染器
func NewRenderer() *Renderer {
	return &Renderer{}
}

// ToggleStealth 切换伪装模式
func (r *Renderer) ToggleStealth() {
	r.stealthMode = !r.stealthMode
}

// IsStealth 是否处于伪装模式
func (r *Renderer) IsStealth() bool {
	return r.stealthMode
}

// LogInfo 输出 INFO 级别日志
func (r *Renderer) LogInfo(msg string) string {
	return r.formatLog("INFO", ColorInfo, msg)
}

// LogWarn 输出 WARN 级别日志
func (r *Renderer) LogWarn(msg string) string {
	return r.formatLog("WARN", ColorWarn, msg)
}

// LogError 输出 ERROR 级别日志
func (r *Renderer) LogError(msg string) string {
	return r.formatLog("ERROR", ColorError, msg)
}

// LogFatal 输出 FATAL 级别日志
func (r *Renderer) LogFatal(msg string) string {
	return r.formatLog("FATAL", ColorFatal, msg)
}

func (r *Renderer) formatLog(level string, color lipgloss.Color, msg string) string {
	ts := time.Now().Format("2006-01-02 15:04:05")
	timestamp := LogTimestamp.Render(fmt.Sprintf("[%s]", ts))
	levelStr := lipgloss.NewStyle().Foreground(color).Bold(true).Render(fmt.Sprintf("%-5s", level))
	return fmt.Sprintf("%s %s %s", timestamp, levelStr, msg)
}

// RenderStatusBar 渲染顶部状态栏
func (r *Renderer) RenderStatusBar(p *game.Character) string {
	hpBar := r.renderHPBar(p.HP, p.MaxHP, 20)
	hpText := fmt.Sprintf("HP:%d/%d", p.HP, p.MaxHP)

	parts := []string{
		TitleStyle.Render(fmt.Sprintf("⚔ %s", p.Name)),
		StatLabel.Render(string(p.Class)),
		StatLabel.Render(fmt.Sprintf("Lv.%d", p.Level)),
		hpBar,
		HPStyle.Render(hpText),
		GoldStyle.Render(fmt.Sprintf("💰%d", p.Gold)),
	}

	return strings.Join(parts, " │ ")
}

// RenderHPBar 渲染 HP 条
func (r *Renderer) renderHPBar(current, max, width int) string {
	filled := int(float64(current) / float64(max) * float64(width))
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}
	empty := width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)

	percent := float64(current) / float64(max) * 100
	var color lipgloss.Color
	switch {
	case percent > 60:
		color = ColorPrimary
	case percent > 30:
		color = ColorWarning
	default:
		color = ColorDanger
	}

	return lipgloss.NewStyle().Foreground(color).Render(bar)
}

// RenderMonsterHP 渲染怪物 HP 条
func (r *Renderer) RenderMonsterHP(m *game.Monster) string {
	hpBar := r.renderHPBar(m.HP, m.MaxHP, 15)
	hpText := fmt.Sprintf("HP:%d/%d", m.HP, m.MaxHP)

	parts := []string{
		CombatTitle.Render(m.Name),
		hpBar,
		HPStyle.Render(hpText),
	}

	if m.IsBoss {
		parts = append([]string{BossStyle.Render("★BOSS★")}, parts...)
	}

	return strings.Join(parts, " ")
}

// RenderCombatView 渲染战斗界面
func (r *Renderer) RenderCombatView(cs *game.CombatState) string {
	var sb strings.Builder

	// 战斗标题
	sb.WriteString(r.LogWarn(fmt.Sprintf("=== 战斗：vs %s ===", cs.Monster.Name)))
	sb.WriteString("\n\n")

	// 怪物信息
	if cs.Monster.Desc != "" {
		sb.WriteString(r.LogInfo(cs.Monster.Desc))
		sb.WriteString("\n")
	}
	sb.WriteString(r.RenderMonsterHP(cs.Monster))
	sb.WriteString("\n")

	// 怪物状态
	if len(cs.Monster.Debuffs) > 0 {
		debuffs := make([]string, 0)
		for _, d := range cs.Monster.Debuffs {
			debuffs = append(debuffs, DebuffStyle.Render(fmt.Sprintf("%s(%d)", d.Name, d.Duration)))
		}
		sb.WriteString(StatLabel.Render("状态: ") + strings.Join(debuffs, " "))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")

	// 玩家状态
	sb.WriteString(r.RenderStatusBar(cs.Player))
	sb.WriteString("\n")

	// 玩家 Buff
	if len(cs.Player.Buffs) > 0 {
		buffs := make([]string, 0)
		for _, b := range cs.Player.Buffs {
			buffs = append(buffs, BuffStyle.Render(fmt.Sprintf("%s(%d)", b.Name, b.Duration)))
		}
		sb.WriteString(StatLabel.Render("增益: ") + strings.Join(buffs, " "))
		sb.WriteString("\n")
	}
	if len(cs.Player.Debuffs) > 0 {
		debuffs := make([]string, 0)
		for _, d := range cs.Player.Debuffs {
			debuffs = append(debuffs, DebuffStyle.Render(fmt.Sprintf("%s(%d)", d.Name, d.Duration)))
		}
		sb.WriteString(StatLabel.Render("减益: ") + strings.Join(debuffs, " "))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")

	// 战斗日志
	for _, msg := range cs.Log {
		sb.WriteString(r.LogInfo("  " + msg))
		sb.WriteString("\n")
	}

	return sb.String()
}

// RenderActions 渲染行动选项
func (r *Renderer) RenderActions(actions []string) string {
	var sb strings.Builder
	sb.WriteString(PromptStyle.Render("可用行动:"))
	sb.WriteString("\n")
	for i, action := range actions {
		sb.WriteString(fmt.Sprintf("  %s %s",
			PromptStyle.Render(fmt.Sprintf("[%d]", i+1)),
			action))
		sb.WriteString("\n")
	}
	return sb.String()
}

// RenderSkills 渲染技能列表
func (r *Renderer) RenderSkills(skills []*game.Skill) string {
	var sb strings.Builder
	sb.WriteString(PromptStyle.Render("技能列表:"))
	sb.WriteString("\n")
	for i, skill := range skills {
		cdText := ""
		if skill.CurCD > 0 {
			cdText = DebuffStyle.Render(fmt.Sprintf(" (冷却%d)", skill.CurCD))
		}
		sb.WriteString(fmt.Sprintf("  %s %s - %s%s",
			PromptStyle.Render(fmt.Sprintf("[%d]", i+1)),
			SkillStyle.Render(skill.Name),
			skill.Desc,
			cdText))
		sb.WriteString("\n")
	}
	return sb.String()
}

// getSlotIcon 获取装备槽位图标
func getSlotIcon(slot game.Slot) string {
	switch slot {
	case game.SlotWeapon:
		return "⚔"
	case game.SlotHelmet:
		return "🪖"
	case game.SlotArmor:
		return "🛡"
	case game.SlotBoots:
		return "👢"
	case game.SlotAcc1, game.SlotAcc2:
		return "💍"
	default:
		return "?"
	}
}

// RenderInventory 渲染物品栏
func (r *Renderer) RenderInventory(items []*game.Item) string {
	var sb strings.Builder
	sb.WriteString(PromptStyle.Render("物品栏:"))
	sb.WriteString("\n")
	if len(items) == 0 {
		sb.WriteString(StatLabel.Render("  (空)"))
		sb.WriteString("\n")
	}
	for i, item := range items {
		rarityStyle, ok := RarityStyle[string(item.Rarity)]
		if !ok {
			rarityStyle = lipgloss.NewStyle().Foreground(ColorWhite)
		}

		// 显示装备槽位
		slotInfo := ""
		if item.Type != game.Consumable && item.Slot != "" {
			slotInfo = fmt.Sprintf(" %s%s", getSlotIcon(item.Slot), string(item.Slot))
		} else if item.Type == game.Consumable {
			slotInfo = " 🧪消耗品"
		}

		sb.WriteString(fmt.Sprintf("  %s %s%s %s",
			PromptStyle.Render(fmt.Sprintf("[%d]", i+1)),
			rarityStyle.Render(item.Name),
			StatLabel.Render(slotInfo),
			StatLabel.Render(item.Description)))
		sb.WriteString("\n")
	}
	return sb.String()
}

// RenderShopItem 渲染商店物品详情（带属性对比）
func (r *Renderer) RenderShopItem(item *game.Item, price int, player *game.Character) string {
	var sb strings.Builder

	// 品质颜色
	rarityStyle, ok := RarityStyle[string(item.Rarity)]
	if !ok {
		rarityStyle = lipgloss.NewStyle().Foreground(ColorWhite)
	}

	// 物品名称和品质
	sb.WriteString(fmt.Sprintf("  %s %s %s",
		rarityStyle.Render(item.Name),
		StatLabel.Render(string(item.Rarity)),
		GoldStyle.Render(fmt.Sprintf("%dg", price))))
	sb.WriteString("\n")

	if item.Description != "" {
		sb.WriteString(fmt.Sprintf("    %s", StatLabel.Render(item.Description)))
		sb.WriteString("\n")
	}

	// 属性详情
	if item.Type != game.Consumable {
		sb.WriteString(fmt.Sprintf("    %s", StatLabel.Render("属性:")))
		sb.WriteString("\n")

		for stat, val := range item.Effects {
			statName := r.getStatName(stat)
			currentVal := r.getPlayerStat(player, stat)

			// 计算差值
			diff := val - currentVal
			diffStr := ""
			diffColor := ColorWhite
			if diff > 0 {
				diffStr = fmt.Sprintf("+%d", diff)
				diffColor = ColorPrimary
			} else if diff < 0 {
				diffStr = fmt.Sprintf("%d", diff)
				diffColor = ColorDanger
			} else {
				diffStr = "±0"
			}

			sb.WriteString(fmt.Sprintf("      %s: %s %s",
				statName,
				lipgloss.NewStyle().Foreground(ColorWhite).Render(fmt.Sprintf("%d", val)),
				lipgloss.NewStyle().Foreground(diffColor).Render(fmt.Sprintf("(%s)", diffStr))))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// RenderLootItem 渲染掉落物品详情
func (r *Renderer) RenderLootItem(item *game.Item) string {
	var sb strings.Builder

	rarityStyle, ok := RarityStyle[string(item.Rarity)]
	if !ok {
		rarityStyle = lipgloss.NewStyle().Foreground(ColorWhite)
	}

	sb.WriteString(fmt.Sprintf("  %s %s",
		ItemStyle.Render("🎁"),
		rarityStyle.Render(item.Name)))
	sb.WriteString("\n")

	if item.Description != "" {
		sb.WriteString(fmt.Sprintf("    %s", StatLabel.Render(item.Description)))
		sb.WriteString("\n")
	}

	// 属性详情
	if item.Type != game.Consumable {
		for stat, val := range item.Effects {
			statName := r.getStatName(stat)
			sb.WriteString(fmt.Sprintf("    %s +%d", statName, val))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// getStatName 获取属性中文名
func (r *Renderer) getStatName(stat string) string {
	switch stat {
	case "atk":
		return "攻击力"
	case "def":
		return "防御力"
	case "hp":
		return "生命值"
	case "crit":
		return "暴击率"
	case "dodge":
		return "闪避率"
	case "spd":
		return "速度"
	case "heal":
		return "恢复"
	case "cure_poison":
		return "解毒"
	case "temp_atk":
		return "临时攻击"
	case "temp_def":
		return "临时防御"
	default:
		return stat
	}
}

// getPlayerStat 获取玩家当前属性值
func (r *Renderer) getPlayerStat(player *game.Character, stat string) int {
	switch stat {
	case "atk":
		return player.ATK
	case "def":
		return player.DEF
	case "hp":
		return player.MaxHP
	case "crit":
		return player.CRIT
	case "dodge":
		return player.Dodge
	case "spd":
		return player.SPD
	default:
		return 0
	}
}

func (r *Renderer) roomIcon(t game.RoomType) string {
	switch t {
	case game.RoomEntrance:
		return "🚪"
	case game.RoomCombat:
		return "⚔️"
	case game.RoomElite:
		return "🔥"
	case game.RoomTreasure:
		return "💎"
	case game.RoomShop:
		return "🏪"
	case game.RoomRest:
		return "🏕️"
	case game.RoomEvent:
		return "❓"
	case game.RoomBoss:
		return "💀"
	default:
		return "📍"
	}
}

// RenderFloorMap 渲染杀戮尖塔风格的网状地图
func (r *Renderer) RenderFloorMap(fm *game.FloorMap) string {
	var sb strings.Builder

	sb.WriteString(r.LogInfo(fmt.Sprintf("══════════ 第 %d 层 ══════════", fm.Depth)))
	sb.WriteString("\n\n")

	// 获取当前节点和可达节点
	currentNode := fm.GetCurrentNode()
	reachable := fm.GetReachableNodes()
	reachableMap := make(map[string]bool)
	for _, n := range reachable {
		reachableMap[game.NodeKey(n.Row, n.Col)] = true
	}

	// 计算进度
	clearedCount := 0
	totalNodes := 0
	for _, row := range fm.Nodes {
		for range row {
			totalNodes++
		}
	}
	for _, row := range fm.Nodes {
		for _, node := range row {
			if node.Cleared {
				clearedCount++
			}
		}
	}
	sb.WriteString(StatLabel.Render(fmt.Sprintf("  进度: %d/%d 节点已探索", clearedCount, totalNodes)))
	sb.WriteString("\n\n")

	// 显示完整地图（从入口到 Boss，每行一个节点组）
	for row := 0; row < fm.RowCount; row++ {
		nodes := fm.Nodes[row]
		if len(nodes) == 0 {
			continue
		}

		// 行标签
		isCurrentRow := currentNode != nil && row == currentNode.Row
		if isCurrentRow {
			sb.WriteString(PromptStyle.Render("  ► "))
		} else {
			sb.WriteString(StatLabel.Render(fmt.Sprintf("  %2d ", row)))
		}

		// 绘制这一行的所有节点
		for i, node := range nodes {
			key := game.NodeKey(row, node.Col)
			icon := game.GetNodeIcon(&node.Type)
			isCurrent := isCurrentRow && node.Col == currentNode.Col
			isReach := reachableMap[key]
			isCleared := node.Cleared

			if i > 0 {
				sb.WriteString("  ")
			}

			// 节点显示
			if isCurrent {
				// 当前位置：高亮框
				sb.WriteString(PromptStyle.Render(fmt.Sprintf("[%s]", icon)))
			} else if isReach {
				// 可前往：高亮 + 下箭头标记
				sb.WriteString(lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render(fmt.Sprintf("▼%s", icon)))
			} else if isCleared {
				// 已清理：灰色+勾
				sb.WriteString(StatLabel.Render(fmt.Sprintf(" %s✓", icon)))
			} else {
				// 未探索
				sb.WriteString(fmt.Sprintf(" %s ", icon))
			}
		}
		sb.WriteString("\n")

		// 连接线（显示分支结构）
		if row < fm.RowCount-1 && len(fm.Nodes[row+1]) > 0 {
			// 收集当前行哪些节点有连接
			hasConn := make(map[int]bool) // col -> has connection
			for _, conn := range fm.Connections {
				if conn.FromRow == row {
					hasConn[conn.FromCol] = true
				}
			}

			// 绘制连接线
			connLine := "      "
			for i, node := range nodes {
				if i > 0 {
					connLine += "  "
				}
				if hasConn[node.Col] {
					connLine += " │ "
				} else {
					connLine += "   "
				}
			}
			sb.WriteString(StatLabel.Render(connLine))
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")

	// 显示当前可前往的节点列表
	if len(reachable) > 0 {
		sb.WriteString(PromptStyle.Render("  可前往:"))
		sb.WriteString("\n")
		for i, node := range reachable {
			icon := game.GetNodeIcon(&node.Type)
			label := game.GetNodeLabel(node.Type)

			var nodeColor lipgloss.Style
			switch node.Type {
			case game.RoomElite:
				nodeColor = lipgloss.NewStyle().Foreground(ColorDanger).Bold(true)
			case game.RoomBoss:
				nodeColor = BossStyle
			case game.RoomRest:
				nodeColor = lipgloss.NewStyle().Foreground(ColorPrimary)
			case game.RoomShop:
				nodeColor = GoldStyle
			case game.RoomTreasure:
				nodeColor = AccentStyle
			default:
				nodeColor = lipgloss.NewStyle().Foreground(ColorWhite)
			}

			sb.WriteString(fmt.Sprintf("    %s %s %s",
				PromptStyle.Render(fmt.Sprintf("[%d]", i+1)),
				icon,
				nodeColor.Render(label)))
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")

	// 图例
	sb.WriteString(StatLabel.Render("  图例: "))
	sb.WriteString(fmt.Sprintf("%s 战斗  ", "⚔"))
	sb.WriteString(fmt.Sprintf("%s 精英  ", "🔥"))
	sb.WriteString(fmt.Sprintf("%s 宝箱  ", "💎"))
	sb.WriteString(fmt.Sprintf("%s 商店  ", "🏪"))
	sb.WriteString(fmt.Sprintf("%s 休息  ", "🏕"))
	sb.WriteString(fmt.Sprintf("%s 事件  ", "❓"))
	sb.WriteString(fmt.Sprintf("%s Boss", "💀"))
	sb.WriteString("\n")

	return sb.String()
}

// isNodeReachable 检查节点是否可达
func (r *Renderer) isNodeReachable(fm *game.FloorMap, row, col int) bool {
	reachable := fm.GetReachableNodes()
	for _, n := range reachable {
		if n.Row == row && n.Col == col {
			return true
		}
	}
	return false
}

// RenderNodeSelection 渲染节点选择界面
func (r *Renderer) RenderNodeSelection(nodes []*game.MapNode) string {
	var sb strings.Builder
	sb.WriteString(PromptStyle.Render("选择目的地:"))
	sb.WriteString("\n")
	for i, node := range nodes {
		icon := game.GetNodeIcon(&node.Type)
		label := game.GetNodeLabel(node.Type)
		sb.WriteString(fmt.Sprintf("  %s %s %s (行%d 列%d)",
			PromptStyle.Render(fmt.Sprintf("[%d]", i+1)),
			icon,
			label,
			node.Row,
			node.Col))
		sb.WriteString("\n")
	}
	return sb.String()
}

// RenderCharacterSheet 渲染角色信息
func (r *Renderer) RenderCharacterSheet(p *game.Character) string {
	var sb strings.Builder
	sb.WriteString(r.LogInfo(fmt.Sprintf("=== %s ===", p.Name)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  %s %s\n", StatLabel.Render("职业:"), string(p.Class)))
	sb.WriteString(fmt.Sprintf("  %s %d  %s %d/%d\n",
		StatLabel.Render("等级:"), p.Level,
		StatLabel.Render("经验:"), p.Exp, p.ExpToNext))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  %s %d/%d\n", HPStyle.Render("HP:"), p.HP, p.MaxHP))
	sb.WriteString(fmt.Sprintf("  %s %d  %s %d  %s %d\n",
		StatLabel.Render("ATK:"), p.ATK,
		StatLabel.Render("DEF:"), p.DEF,
		StatLabel.Render("SPD:"), p.SPD))
	sb.WriteString(fmt.Sprintf("  %s %d%%  %s %d%%\n",
		StatLabel.Render("暴击:"), p.CRIT,
		StatLabel.Render("闪避:"), p.Dodge))
	sb.WriteString(fmt.Sprintf("  %s %d\n", GoldStyle.Render("金币:"), p.Gold))

	// 装备（显示所有槽位）
	sb.WriteString("\n")
	sb.WriteString(StatLabel.Render("装备:"))
	sb.WriteString("\n")

	// 按固定顺序显示槽位
	slots := []game.Slot{game.SlotWeapon, game.SlotHelmet, game.SlotArmor, game.SlotBoots, game.SlotAcc1, game.SlotAcc2}
	for _, slot := range slots {
		icon := getSlotIcon(slot)
		item, exists := p.Equipment[slot]
		if exists && item != nil {
			rarityStyle, ok := RarityStyle[string(item.Rarity)]
			if !ok {
				rarityStyle = lipgloss.NewStyle().Foreground(ColorWhite)
			}
			// 显示装备属性
			effects := ""
			for stat, val := range item.Effects {
				if effects != "" {
					effects += " "
				}
				effects += fmt.Sprintf("%s+%d", r.getStatName(stat), val)
			}
			sb.WriteString(fmt.Sprintf("  %s %s %s %s\n",
				StatLabel.Render(fmt.Sprintf("%s%s:", icon, string(slot))),
				rarityStyle.Render(item.Name),
				StatLabel.Render(string(item.Rarity)),
				StatLabel.Render(effects)))
		} else {
			sb.WriteString(fmt.Sprintf("  %s %s\n",
				StatLabel.Render(fmt.Sprintf("%s%s:", icon, string(slot))),
				StatLabel.Render("(空)")))
		}
	}

	return sb.String()
}

// RenderEvent 渲染随机事件
func (r *Renderer) RenderEvent(event *game.Event) string {
	var sb strings.Builder
	sb.WriteString(r.LogWarn(fmt.Sprintf("=== %s ===", event.Title)))
	sb.WriteString("\n")
	sb.WriteString(r.LogInfo(event.Description))
	sb.WriteString("\n\n")

	for i, choice := range event.Choices {
		sb.WriteString(fmt.Sprintf("  %s %s",
			PromptStyle.Render(fmt.Sprintf("[%d]", i+1)),
			choice.Text))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(StatLabel.Render("  输入编号做出选择"))

	return sb.String()
}

// RenderShop 渲染商店
func (r *Renderer) RenderShop(items []*game.Item, playerGold int, player *game.Character) string {
	var sb strings.Builder
	sb.WriteString(r.LogInfo("=== 商店 ==="))
	sb.WriteString("\n")
	sb.WriteString(GoldStyle.Render(fmt.Sprintf("你的金币: %d", playerGold)))
	sb.WriteString("\n\n")

	for i, item := range items {
		price := 10 + len(item.Effects)*5
		sb.WriteString(fmt.Sprintf("  %s", PromptStyle.Render(fmt.Sprintf("[%d]", i+1))))
		sb.WriteString(r.RenderShopItem(item, price, player))
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("  %s 返回", PromptStyle.Render("[0]")))
	sb.WriteString("\n")

	return sb.String()
}

// RenderGameOver 渲染游戏结束
func (r *Renderer) RenderGameOver(won bool, p *game.Character) string {
	var sb strings.Builder
	if won {
		sb.WriteString(r.LogInfo("═══════════════════════════"))
		sb.WriteString("\n")
		sb.WriteString(r.LogInfo("  🎉 地牢通关！恭喜你！"))
		sb.WriteString("\n")
		sb.WriteString(r.LogInfo("═══════════════════════════"))
	} else {
		sb.WriteString(r.LogFatal("═══════════════════════════"))
		sb.WriteString("\n")
		sb.WriteString(r.LogFatal("  💀 你倒下了……"))
		sb.WriteString("\n")
		sb.WriteString(r.LogFatal("═══════════════════════════"))
	}
	sb.WriteString("\n")
	sb.WriteString(r.RenderCharacterSheet(p))
	return sb.String()
}

// RenderStealthScreen 渲染伪装界面（类似 htop）
func (r *Renderer) RenderStealthScreen(p *game.Character) string {
	var sb strings.Builder

	// htop 风格头部
	sb.WriteString(StealthHeader.Render("  PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND"))
	sb.WriteString("\n")

	// 用游戏数据伪装成进程信息
	hpPercent := float64(p.HP) / float64(p.MaxHP) * 100
	sb.WriteString(fmt.Sprintf(" %4d production  20   0  128456  32768   8192 S  %4.1f  2.1   1:23.45 dungeonlog\n",
		1234+p.Level, hpPercent))
	sb.WriteString(fmt.Sprintf(" %4d production  20   0   65536  16384   4096 S  %4.1f  1.1   0:45.67 redis-server\n",
		5678, float64(p.ATK)))
	sb.WriteString(fmt.Sprintf(" %4d production  20   0  262144  65536  16384 S  %4.1f  4.2   3:21.09 postgres\n",
		9012, float64(p.DEF)))

	sb.WriteString("\n")
	sb.WriteString(StealthHeader.Render("  MiB Mem: 16384.0 total, 8192.0 free, 4096.0 used, 4096.0 buff/cache"))
	sb.WriteString(StealthHeader.Render("  MiB Swap: 4096.0 total, 4096.0 free, 0.0 used. 10240.0 avail Mem"))
	sb.WriteString("\n")
	sb.WriteString(StealthHeader.Render(fmt.Sprintf("  [F5]Refresh  [F9]Kill  [F10]Quit   DungeonLog v1.0 - %s", p.Name)))

	return sb.String()
}

// RenderWelcome 渲染欢迎界面
func (r *Renderer) RenderWelcome() string {
	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString(TitleStyle.Render("  ╔═══════════════════════════════════════╗"))
	sb.WriteString("\n")
	sb.WriteString(TitleStyle.Render("  ║         ⚔  DUNGEON LOG  ⚔           ║"))
	sb.WriteString("\n")
	sb.WriteString(TitleStyle.Render("  ║      终端地牢探险 RPG v1.0           ║"))
	sb.WriteString("\n")
	sb.WriteString(TitleStyle.Render("  ╚═══════════════════════════════════════╝"))
	sb.WriteString("\n\n")

	sb.WriteString(r.LogInfo("欢迎来到 DungeonLog，一款伪装成日志的终端地牢 RPG。"))
	sb.WriteString("\n")
	sb.WriteString(r.LogInfo("输入 'help' 查看可用命令，'start' 开始冒险。"))
	sb.WriteString("\n")

	return sb.String()
}

// RenderHelp 渲染帮助信息
func (r *Renderer) RenderHelp() string {
	var sb strings.Builder
	sb.WriteString(r.LogInfo("=== 可用命令 ==="))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  %s 开始新游戏\n", PromptStyle.Render("start <名字> <职业>")))
	sb.WriteString(fmt.Sprintf("  %s 查看角色信息\n", PromptStyle.Render("info / status")))
	sb.WriteString(fmt.Sprintf("  %s 查看物品栏\n", PromptStyle.Render("bag / inventory")))
	sb.WriteString(fmt.Sprintf("  %s 查看技能\n", PromptStyle.Render("skills")))
	sb.WriteString(fmt.Sprintf("  %s 查看地图\n", PromptStyle.Render("map")))
	sb.WriteString(fmt.Sprintf("  %s 前往指定房间\n", PromptStyle.Render("go <编号>")))
	sb.WriteString(fmt.Sprintf("  %s 与当前房间交互\n", PromptStyle.Render("look / interact")))
	sb.WriteString(fmt.Sprintf("  %s 保存游戏\n", PromptStyle.Render("save")))
	sb.WriteString(fmt.Sprintf("  %s 加载游戏\n", PromptStyle.Render("load")))
	sb.WriteString(fmt.Sprintf("  %s 切换伪装模式\n", PromptStyle.Render("stealth / Ctrl+H")))
	sb.WriteString(fmt.Sprintf("  %s 退出游戏\n", PromptStyle.Render("quit / exit")))
	sb.WriteString("\n")
	sb.WriteString(StatLabel.Render("职业选择:"))
	sb.WriteString("\n")
	for class, info := range game.ClassData {
		sb.WriteString(fmt.Sprintf("  %s - %s\n", TitleStyle.Render(string(class)), info.Description))
	}
	sb.WriteString("\n")
	sb.WriteString(StatLabel.Render("战斗中可用命令:"))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  %s 普通攻击\n", PromptStyle.Render("1 / attack")))
	sb.WriteString(fmt.Sprintf("  %s 使用技能\n", PromptStyle.Render("2 / skill <编号>")))
	sb.WriteString(fmt.Sprintf("  %s 使用物品\n", PromptStyle.Render("3 / use <编号>")))
	sb.WriteString(fmt.Sprintf("  %s 尝试逃跑\n", PromptStyle.Render("4 / flee")))

	return sb.String()
}

// RenderEventChoiceResult 渲染事件选择结果
func (r *Renderer) RenderEventChoiceResult(choice game.EventChoice) string {
	var sb strings.Builder
	switch choice.Outcome {
	case "good":
		sb.WriteString(r.LogInfo("✦ " + choice.Description))
	case "bad":
		sb.WriteString(r.LogError("✗ " + choice.Description))
	default:
		sb.WriteString(r.LogInfo("○ " + choice.Description))
	}
	return sb.String()
}

// RenderLoot 渲染战利品
func (r *Renderer) RenderLoot(items []*game.Item, gold int, exp int) string {
	var sb strings.Builder
	sb.WriteString(r.LogInfo("=== 战利品 ==="))
	sb.WriteString("\n")
	if gold > 0 {
		sb.WriteString(fmt.Sprintf("  %s\n", GoldStyle.Render(fmt.Sprintf("💰 %d 金币", gold))))
	}
	if exp > 0 {
		sb.WriteString(fmt.Sprintf("  %s\n", AccentStyle.Render(fmt.Sprintf("✦ %d 经验", exp))))
	}
	for _, item := range items {
		sb.WriteString(r.RenderLootItem(item))
	}
	return sb.String()
}

// RenderTownMenu 渲染城镇菜单
func (r *Renderer) RenderTownMenu(meta *game.MetaProgress) string {
	var sb strings.Builder

	sb.WriteString(GoldStyle.Render(fmt.Sprintf("  💰 可用金币: %d", meta.TotalGold)))
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf("  %s 升级商店\n", PromptStyle.Render("[1]")))
	sb.WriteString(fmt.Sprintf("  %s 解锁职业\n", PromptStyle.Render("[2]")))
	sb.WriteString(fmt.Sprintf("  %s 查看统计\n", PromptStyle.Render("[3]")))
	sb.WriteString(fmt.Sprintf("  %s 开始冒险\n", PromptStyle.Render("[4]")))
	sb.WriteString(fmt.Sprintf("  %s 返回标题\n", PromptStyle.Render("[5]")))

	return sb.String()
}

// RenderUpgradeShop 渲染升级商店
func (r *Renderer) RenderUpgradeShop(meta *game.MetaProgress) string {
	var sb strings.Builder

	sb.WriteString(r.LogInfo("=== 升级商店 ==="))
	sb.WriteString("\n")
	sb.WriteString(GoldStyle.Render(fmt.Sprintf("  💰 可用金币: %d", meta.TotalGold)))
	sb.WriteString("\n\n")

	for i, upgrade := range game.AvailableUpgrades {
		level := meta.GetUpgradeLevel(upgrade.ID)
		cost := meta.GetUpgradeCost(upgrade)

		// 等级显示
		levelStr := fmt.Sprintf("Lv.%d/%d", level, upgrade.MaxLevel)
		if level >= upgrade.MaxLevel {
			levelStr = SuccessStyle.Render("MAX")
		} else {
			levelStr = StatLabel.Render(levelStr)
		}

		// 费用显示
		costStr := ""
		if cost < 0 {
			costStr = StatLabel.Render("已满级")
		} else if meta.TotalGold >= cost {
			costStr = GoldStyle.Render(fmt.Sprintf("%dg", cost))
		} else {
			costStr = ErrorStyle.Render(fmt.Sprintf("%dg", cost))
		}

		sb.WriteString(fmt.Sprintf("  %s %s %s %s\n      %s\n",
			PromptStyle.Render(fmt.Sprintf("[%d]", i+1)),
			TitleStyle.Render(upgrade.Name),
			levelStr,
			costStr,
			StatLabel.Render(upgrade.Description)))
		sb.WriteString(fmt.Sprintf("      %s\n", StatLabel.Render(upgrade.Effect)))
	}

	sb.WriteString(fmt.Sprintf("\n  %s 返回", PromptStyle.Render("[0]")))
	sb.WriteString("\n")

	return sb.String()
}

// RenderClassShop 渲染职业商店
func (r *Renderer) RenderClassShop(meta *game.MetaProgress) string {
	var sb strings.Builder

	sb.WriteString(r.LogInfo("=== 职业解锁 ==="))
	sb.WriteString("\n")
	sb.WriteString(GoldStyle.Render(fmt.Sprintf("  💰 可用金币: %d", meta.TotalGold)))
	sb.WriteString("\n\n")

	classes := []struct {
		Class game.Class
		Name  string
		Desc  string
	}{
		{game.Warrior, "战士", "近战之王，高防高血，擅长持久战"},
		{game.Ranger, "游侠", "远程射手，高暴击高闪避"},
		{game.Mage, "法师", "元素法师，AOE伤害+控制"},
		{game.Rogue, "盗贼", "暗影刺客，高暴击+连击"},
	}

	for _, c := range classes {
		unlocked := meta.IsClassUnlocked(c.Class)
		cost := game.GetClassUnlockCost(c.Class)

		status := ""
		if unlocked {
			status = SuccessStyle.Render("✓ 已解锁")
		} else {
			if meta.TotalGold >= cost {
				status = GoldStyle.Render(fmt.Sprintf("%dg", cost))
			} else {
				status = ErrorStyle.Render(fmt.Sprintf("%dg", cost))
			}
		}

		sb.WriteString(fmt.Sprintf("  %s %s - %s %s\n",
			PromptStyle.Render(c.Name),
			TitleStyle.Render(c.Desc),
			status,
			StatLabel.Render(c.Desc)))
	}

	sb.WriteString(fmt.Sprintf("\n  %s 返回", PromptStyle.Render("[0]")))
	sb.WriteString("\n")

	return sb.String()
}
