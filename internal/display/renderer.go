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
		sb.WriteString(fmt.Sprintf("  %s %s %s",
			PromptStyle.Render(fmt.Sprintf("[%d]", i+1)),
			rarityStyle.Render(item.Name),
			StatLabel.Render(item.Description)))
		sb.WriteString("\n")
	}
	return sb.String()
}

// RenderRoom 渲染房间
func (r *Renderer) RenderRoom(room *game.Room, dungeon *game.Dungeon) string {
	var sb strings.Builder

	// 房间类型图标
	icon := r.roomIcon(room.Type)
	sb.WriteString(r.LogInfo(fmt.Sprintf("%s 进入房间 #%d (%s)", icon, room.Index, string(room.Type))))
	sb.WriteString("\n")

	if room.Description != "" {
		sb.WriteString(r.LogInfo("  " + room.Description))
		sb.WriteString("\n")
	}

	// 显示可前往的房间
	connected := dungeon.GetConnectedRooms()
	if len(connected) > 0 {
		sb.WriteString("\n")
		sb.WriteString(StatLabel.Render("可前往:"))
		sb.WriteString("\n")
		for _, next := range connected {
			visited := ""
			if next.Visited {
				visited = StatLabel.Render(" (已探索)")
			}
			sb.WriteString(fmt.Sprintf("  %s 房间#%d [%s]%s",
				PromptStyle.Render(fmt.Sprintf("[%d]", next.Index)),
				next.Index,
				string(next.Type),
				visited))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (r *Renderer) roomIcon(t game.RoomType) string {
	switch t {
	case game.RoomEntrance:
		return "🚪"
	case game.RoomCombat:
		return "⚔️"
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

// RenderDungeonMap 渲染简易地牢地图
func (r *Renderer) RenderDungeonMap(dungeon *game.Dungeon) string {
	var sb strings.Builder
	sb.WriteString(r.LogInfo(fmt.Sprintf("=== 地牢地图 (第 %d 层) ===", dungeon.Depth)))
	sb.WriteString("\n\n")

	for i, room := range dungeon.Rooms {
		marker := " "
		if i == dungeon.Current {
			marker = "►"
		}
		if !room.Visited {
			sb.WriteString(fmt.Sprintf("  %s [???]\n", marker))
		} else {
			icon := r.roomIcon(room.Type)
			cleared := ""
			if room.Cleared {
				cleared = " ✓"
			}
			sb.WriteString(fmt.Sprintf("  %s %s %s%s\n", marker, icon, string(room.Type), cleared))
		}
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

	// 装备
	sb.WriteString("\n")
	sb.WriteString(StatLabel.Render("装备:"))
	sb.WriteString("\n")
	for slot, item := range p.Equipment {
		if item != nil {
			rarityStyle, ok := RarityStyle[string(item.Rarity)]
			if !ok {
				rarityStyle = lipgloss.NewStyle().Foreground(ColorWhite)
			}
			sb.WriteString(fmt.Sprintf("  %s: %s\n", string(slot), rarityStyle.Render(item.Name)))
		} else {
			sb.WriteString(fmt.Sprintf("  %s: -\n", string(slot)))
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

	return sb.String()
}

// RenderShop 渲染商店
func (r *Renderer) RenderShop(items []*game.Item, playerGold int) string {
	var sb strings.Builder
	sb.WriteString(r.LogInfo("=== 商店 ==="))
	sb.WriteString("\n")
	sb.WriteString(GoldStyle.Render(fmt.Sprintf("你的金币: %d", playerGold)))
	sb.WriteString("\n\n")

	for i, item := range items {
		rarityStyle, ok := RarityStyle[string(item.Rarity)]
		if !ok {
			rarityStyle = lipgloss.NewStyle().Foreground(ColorWhite)
		}
		price := 10 + len(item.Effects)*5
		sb.WriteString(fmt.Sprintf("  %s %s %s %s",
			PromptStyle.Render(fmt.Sprintf("[%d]", i+1)),
			rarityStyle.Render(item.Name),
			StatLabel.Render(item.Description),
			GoldStyle.Render(fmt.Sprintf("%dg", price))))
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("\n  %s 返回", PromptStyle.Render("[0]")))
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
		rarityStyle, ok := RarityStyle[string(item.Rarity)]
		if !ok {
			rarityStyle = lipgloss.NewStyle().Foreground(ColorWhite)
		}
		sb.WriteString(fmt.Sprintf("  %s %s\n",
			ItemStyle.Render("🎁"),
			rarityStyle.Render(item.Name)))
	}
	return sb.String()
}
