package engine

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"dungeonlog/internal/display"
	"dungeonlog/internal/game"
	"dungeonlog/internal/llm"
)

// Game 游戏引擎
type Game struct {
	state    *GameState
	renderer *display.Renderer
	llm      llm.LLMProvider
	ctx      context.Context
	output   []string // 输出缓冲
}

// NewGame 创建游戏引擎
func NewGame(provider llm.LLMProvider) *Game {
	return &Game{
		state: &GameState{
			Phase: PhaseTitle,
		},
		renderer: display.NewRenderer(),
		llm:      provider,
		ctx:      context.Background(),
		output:   make([]string, 0),
	}
}

// GetOutput 获取并清空输出缓冲
func (g *Game) GetOutput() []string {
	out := g.output
	g.output = make([]string, 0)
	return out
}

// addOutput 添加输出
func (g *Game) addOutput(msg string) {
	g.output = append(g.output, msg)
}

// GetRenderer 获取渲染器
func (g *Game) GetRenderer() *display.Renderer {
	return g.renderer
}

// GetState 获取游戏状态
func (g *Game) GetState() *GameState {
	return g.state
}

// ProcessInput 处理用户输入
func (g *Game) ProcessInput(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	parts := strings.Fields(input)
	cmd := strings.ToLower(parts[0])
	args := parts[1:]

	switch g.state.Phase {
	case PhaseTitle:
		g.processTitleInput(cmd, args)
	case PhaseCharacter:
		g.processCharacterInput(cmd, args)
	case PhaseExplore:
		g.processExploreInput(cmd, args)
	case PhaseCombat:
		g.processCombatInput(cmd, args)
	case PhaseShop:
		g.processShopInput(cmd, args)
	case PhaseEvent:
		g.processEventInput(cmd, args)
	case PhaseGameOver, PhaseVictory:
		g.processEndInput(cmd, args)
	}
}

// === 标题界面 ===
func (g *Game) processTitleInput(cmd string, args []string) {
	switch cmd {
	case "start":
		if len(args) < 1 {
			g.addOutput(g.renderer.LogError("用法: start <名字> [职业]"))
			return
		}

		// 检查是否同时提供了职业
		if len(args) >= 2 {
			classMap := map[string]game.Class{
				"战士": game.Warrior, "warrior": game.Warrior,
				"游侠": game.Ranger, "ranger": game.Ranger,
				"法师": game.Mage, "mage": game.Mage,
				"盗贼": game.Rogue, "rogue": game.Rogue,
			}
			if class, ok := classMap[args[1]]; ok {
				g.startGame(args[0], class)
				return
			}
		}

		g.state.Phase = PhaseCharacter
		g.state.Player = &game.Character{Name: args[0]}
		g.addOutput(g.renderer.LogInfo(fmt.Sprintf("欢迎，%s！请选择你的职业:", args[0])))
		g.addOutput("")
		for class, info := range game.ClassData {
			g.addOutput(fmt.Sprintf("  %s - %s",
				display.TitleStyle.Render(string(class)),
				info.Description))
		}
		g.addOutput("")
		g.addOutput(g.renderer.LogInfo("输入职业名称选择: 战士 / 游侠 / 法师 / 盗贼"))

	case "load":
		state, err := LoadGame("data/saves/save.json")
		if err != nil {
			g.addOutput(g.renderer.LogError(fmt.Sprintf("加载失败: %v", err)))
			return
		}
		g.state = state
		g.addOutput(g.renderer.LogInfo("游戏加载成功！"))
		g.addOutput(g.renderer.RenderRoom(g.state.Dungeon.GetCurrentRoom(), g.state.Dungeon))

	case "help":
		g.addOutput(g.renderer.RenderHelp())

	case "quit", "exit":
		g.addOutput(g.renderer.LogInfo("再见，冒险者！"))
		// The main loop will handle exit

	default:
		g.addOutput(g.renderer.LogError("未知命令。输入 'help' 查看帮助，'start <名字>' 开始游戏。"))
	}
}

// startGame 开始游戏（创建角色并进入探索）
func (g *Game) startGame(name string, class game.Class) {
	player := game.NewCharacter(name, class)
	g.state.Player = player
	g.state.Depth = 1
	g.state.Phase = PhaseExplore
	g.state.Dungeon = game.GenerateDungeon(g.state.Depth)

	g.addOutput(g.renderer.LogInfo(fmt.Sprintf("你选择了 %s！", string(class))))
	g.addOutput(g.renderer.LogInfo(fmt.Sprintf("%s 的冒险即将开始……", player.Name)))
	g.addOutput("")
	g.addOutput(g.renderer.RenderCharacterSheet(player))
	g.addOutput("")
	g.addOutput(g.renderer.RenderRoom(g.state.Dungeon.GetCurrentRoom(), g.state.Dungeon))
}

// === 角色创建 ===
func (g *Game) processCharacterInput(cmd string, args []string) {
	classMap := map[string]game.Class{
		"战士": game.Warrior, "warrior": game.Warrior,
		"游侠": game.Ranger, "ranger": game.Ranger,
		"法师": game.Mage, "mage": game.Mage,
		"盗贼": game.Rogue, "rogue": game.Rogue,
	}

	class, ok := classMap[cmd]
	if !ok {
		g.addOutput(g.renderer.LogError("无效的职业。请输入: 战士 / 游侠 / 法师 / 盗贼"))
		return
	}

	g.startGame(g.state.Player.Name, class)
}

// === 探索阶段 ===
func (g *Game) processExploreInput(cmd string, args []string) {
	switch cmd {
	case "go":
		if len(args) < 1 {
			g.addOutput(g.renderer.LogError("用法: go <房间编号>"))
			return
		}
		idx, err := strconv.Atoi(args[0])
		if err != nil {
			g.addOutput(g.renderer.LogError("无效的房间编号"))
			return
		}
		g.moveToRoom(idx)

	case "look", "interact":
		g.interactWithRoom()

	case "info", "status":
		g.addOutput(g.renderer.RenderCharacterSheet(g.state.Player))

	case "bag", "inventory":
		g.addOutput(g.renderer.RenderInventory(g.state.Player.Inventory))

	case "skills":
		g.addOutput(g.renderer.RenderSkills(g.state.Player.Skills))

	case "map":
		g.addOutput(g.renderer.RenderDungeonMap(g.state.Dungeon))

	case "equip":
		if len(args) < 1 {
			g.addOutput(g.renderer.LogError("用法: equip <物品编号>"))
			return
		}
		g.equipItem(args[0])

	case "save":
		err := SaveGame(g.state, "data/saves/save.json")
		if err != nil {
			g.addOutput(g.renderer.LogError(fmt.Sprintf("保存失败: %v", err)))
		} else {
			g.addOutput(g.renderer.LogInfo("游戏已保存。"))
		}

	case "load":
		state, err := LoadGame("data/saves/save.json")
		if err != nil {
			g.addOutput(g.renderer.LogError(fmt.Sprintf("加载失败: %v", err)))
			return
		}
		g.state = state
		g.addOutput(g.renderer.LogInfo("游戏加载成功！"))
		g.addOutput(g.renderer.RenderRoom(g.state.Dungeon.GetCurrentRoom(), g.state.Dungeon))

	case "stealth":
		g.renderer.ToggleStealth()
		if g.renderer.IsStealth() {
			g.addOutput(g.renderer.RenderStealthScreen(g.state.Player))
		} else {
			g.addOutput(g.renderer.LogInfo("伪装模式已关闭"))
		}

	case "config":
		g.showConfig()

	case "help":
		g.addOutput(g.renderer.RenderHelp())

	case "quit", "exit":
		g.addOutput(g.renderer.LogInfo("再见，冒险者！"))

	default:
		g.addOutput(g.renderer.LogError("未知命令。输入 'help' 查看帮助。"))
	}
}

// moveToRoom 移动到指定房间
func (g *Game) moveToRoom(idx int) {
	if !g.state.Dungeon.MoveTo(idx) {
		g.addOutput(g.renderer.LogError("无法前往该房间！请检查连接路径。"))
		return
	}

	room := g.state.Dungeon.GetCurrentRoom()
	g.addOutput(g.renderer.RenderRoom(room, g.state.Dungeon))

	// 根据房间类型触发事件
	switch room.Type {
	case game.RoomCombat:
		g.startCombat(false)
	case game.RoomBoss:
		g.startCombat(true)
	case game.RoomShop:
		g.startShop()
	case game.RoomEvent:
		g.startEvent()
	case game.RoomTreasure:
		g.openTreasure()
	case game.RoomRest:
		g.restAtCamp()
	case game.RoomEntrance:
		room.Cleared = true
	}
}

// interactWithRoom 与当前房间交互
func (g *Game) interactWithRoom() {
	room := g.state.Dungeon.GetCurrentRoom()
	if room.Cleared {
		g.addOutput(g.renderer.LogInfo("这个房间已经被清理过了。"))
		return
	}

	switch room.Type {
	case game.RoomCombat:
		g.startCombat(false)
	case game.RoomBoss:
		g.startCombat(true)
	case game.RoomShop:
		g.startShop()
	case game.RoomEvent:
		g.startEvent()
	case game.RoomTreasure:
		g.openTreasure()
	case game.RoomRest:
		g.restAtCamp()
	default:
		g.addOutput(g.renderer.LogInfo("这里没有什么特别的。"))
	}
}

// === 战斗系统 ===
func (g *Game) startCombat(isBoss bool) {
	// 生成怪物
	monsterResp, err := g.llm.GenerateMonster(g.ctx, llm.MonsterRequest{
		Depth:    g.state.Depth,
		RoomType: "combat",
		IsBoss:   isBoss,
	})
	if err != nil {
		g.addOutput(g.renderer.LogError(fmt.Sprintf("生成怪物失败: %v", err)))
		return
	}

	monster := game.NewMonsterFromData(
		monsterResp.Name,
		monsterResp.Description,
		monsterResp.Traits,
		monsterResp.Dialogue,
		monsterResp.Strategy,
		g.state.Depth,
		isBoss,
	)

	g.state.CombatState = game.NewCombatState(g.state.Player, monster)
	g.state.Phase = PhaseCombat

	g.addOutput(g.renderer.LogWarn(fmt.Sprintf("⚔ 遭遇敌人：%s！", monster.Name)))
	if monster.Dialogue != "" {
		g.addOutput(g.renderer.LogInfo(fmt.Sprintf("  \"%s\"", monster.Dialogue)))
	}
	g.addOutput("")
	g.addOutput(g.renderer.RenderCombatView(g.state.CombatState))
	g.addOutput("")
	g.addOutput(g.renderer.RenderActions([]string{
		"普通攻击",
		"技能",
		"使用物品",
		"逃跑",
	}))
}

func (g *Game) processCombatInput(cmd string, args []string) {
	cs := g.state.CombatState
	if cs == nil {
		g.state.Phase = PhaseExplore
		return
	}

	var result *game.CombatResult

	switch cmd {
	case "1", "attack":
		result = cs.ExecutePlayerAttack()

	case "2", "skill":
		if len(args) < 1 {
			g.addOutput(g.renderer.RenderSkills(g.state.Player.Skills))
			g.addOutput(g.renderer.LogInfo("输入 'skill <编号>' 使用技能"))
			return
		}
		idx, err := strconv.Atoi(args[0])
		if err != nil || idx < 1 || idx > len(g.state.Player.Skills) {
			g.addOutput(g.renderer.LogError("无效的技能编号"))
			return
		}
		result = cs.ExecutePlayerSkill(idx - 1)

	case "3", "use":
		if len(args) < 1 {
			g.addOutput(g.renderer.RenderInventory(g.state.Player.Inventory))
			g.addOutput(g.renderer.LogInfo("输入 'use <编号>' 使用物品"))
			return
		}
		idx, err := strconv.Atoi(args[0])
		if err != nil || idx < 1 || idx > len(g.state.Player.Inventory) {
			g.addOutput(g.renderer.LogError("无效的物品编号"))
			return
		}
		result = cs.ExecutePlayerUseItem(idx - 1)

	case "4", "flee":
		result = cs.ExecutePlayerFlee()
		if cs.Fled {
			g.state.Phase = PhaseExplore
			g.state.Dungeon.GetCurrentRoom().Cleared = true
			for _, msg := range result.Messages {
				g.addOutput(g.renderer.LogInfo(msg))
			}
			return
		}

	case "info", "status":
		g.addOutput(g.renderer.RenderCharacterSheet(g.state.Player))
		return

	case "skills":
		g.addOutput(g.renderer.RenderSkills(g.state.Player.Skills))
		return

	case "bag", "inventory":
		g.addOutput(g.renderer.RenderInventory(g.state.Player.Inventory))
		return

	default:
		g.addOutput(g.renderer.LogError("无效的行动。输入 1-4 选择行动。"))
		return
	}

	// 显示玩家行动结果
	for _, msg := range result.Messages {
		if result.PlayerCrit {
			g.addOutput(display.CritStyle.Render("★ " + msg))
		} else {
			g.addOutput(g.renderer.LogInfo(msg))
		}
	}

	// 检查战斗是否结束
	if result.IsOver {
		if result.PlayerWon {
			g.combatVictory()
		} else {
			g.gameOver(false)
		}
		return
	}

	// 怪物回合 - LLM 生成怪物行动描述
	monsterAction := ""
	if cs.Monster.IsAlive() {
		actionDesc, err := g.llm.GenerateCombatAction(g.ctx, llm.CombatActionRequest{
			MonsterName:  cs.Monster.Name,
			Traits:       cs.Monster.Traits,
			Strategy:     cs.Monster.Strategy,
			MonsterHP:    cs.Monster.HP,
			MonsterMaxHP: cs.Monster.MaxHP,
			PlayerHP:     g.state.Player.HP,
			PlayerMaxHP:  g.state.Player.MaxHP,
			Round:        cs.Round,
		})
		if err == nil {
			monsterAction = actionDesc
		}
	}

	monsterResult := cs.ExecuteMonsterTurn(monsterAction)
	for _, msg := range monsterResult.Messages {
		if monsterResult.MonsterCrit {
			g.addOutput(display.CritStyle.Render("★ " + msg))
		} else if monsterResult.PlayerDodge {
			g.addOutput(display.DodgeStyle.Render("↗ " + msg))
		} else {
			g.addOutput(g.renderer.LogWarn(msg))
		}
	}

	// 回合结束处理
	turnEndMsgs := cs.ProcessTurnEnd()
	for _, msg := range turnEndMsgs {
		g.addOutput(g.renderer.LogInfo(msg))
	}

	// 再次检查战斗状态
	if !g.state.Player.IsAlive() {
		g.gameOver(false)
		return
	}
	if !cs.Monster.IsAlive() {
		g.combatVictory()
		return
	}

	// 显示当前战斗状态
	g.addOutput("")
	g.addOutput(g.renderer.RenderCombatView(cs))
	g.addOutput("")
	g.addOutput(g.renderer.RenderActions([]string{
		"普通攻击",
		"技能",
		"使用物品",
		"逃跑",
	}))
}

func (g *Game) combatVictory() {
	cs := g.state.CombatState
	room := g.state.Dungeon.GetCurrentRoom()
	room.Cleared = true

	// 获得经验和金币
	leveledUp := g.state.Player.GainExp(cs.Monster.ExpReward)
	g.state.Player.Gold += cs.Monster.GoldReward

	g.addOutput(g.renderer.LogInfo(fmt.Sprintf("🎉 击败了 %s！", cs.Monster.Name)))
	g.addOutput(g.renderer.LogInfo(fmt.Sprintf("  获得 %d 经验，%d 金币", cs.Monster.ExpReward, cs.Monster.GoldReward)))

	if leveledUp {
		g.addOutput(display.AccentStyle.Render(fmt.Sprintf("  ★ 升级了！当前等级: %d", g.state.Player.Level)))
	}

	// 掉落物品
	lootChance := 40 + g.state.Depth*5
	if cs.Monster.IsBoss {
		lootChance = 100
	}
	if rand.Intn(100) < lootChance {
		var loot *game.Item
		if rand.Intn(100) < 30 {
			loot = game.GenerateConsumableLoot(g.state.Depth)
		} else {
			loot = game.GenerateLoot(g.state.Depth)
		}
		g.state.Player.Inventory = append(g.state.Player.Inventory, loot)
		g.addOutput(g.renderer.RenderLoot([]*game.Item{loot}, 0, 0))
	}

	// Boss 击败后进入下一层
	if cs.Monster.IsBoss {
		g.addOutput(g.renderer.LogInfo("═══════════════════════════"))
		g.addOutput(g.renderer.LogInfo(fmt.Sprintf("  地牢第 %d 层通关！", g.state.Depth)))
		g.addOutput(g.renderer.LogInfo("═══════════════════════════"))
		g.state.Depth++
		g.state.Dungeon = game.GenerateDungeon(g.state.Depth)
		g.addOutput(g.renderer.LogInfo(fmt.Sprintf("  进入第 %d 层……", g.state.Depth)))
		g.addOutput(g.renderer.RenderRoom(g.state.Dungeon.GetCurrentRoom(), g.state.Dungeon))
	}

	g.state.CombatState = nil
	g.state.Phase = PhaseExplore
}

// === 商店系统 ===
func (g *Game) startShop() {
	// 生成商店物品
	items := make([]*game.Item, 5)
	for i := range items {
		if rand.Intn(100) < 30 {
			items[i] = game.GenerateConsumableLoot(g.state.Depth)
		} else {
			items[i] = game.GenerateLoot(g.state.Depth)
		}
	}
	g.state.ShopItems = items
	g.state.Phase = PhaseShop

	// LLM 生成店主对话
	dialogue, _ := g.llm.GenerateDialogue(g.ctx, llm.DialogueRequest{
		NPCType:    "shop",
		Context:    fmt.Sprintf("地牢第%d层", g.state.Depth),
		PlayerInfo: g.state.Player.String(),
		Topic:      "欢迎",
	})

	g.addOutput(g.renderer.LogInfo("🏪 你走进了一家地牢商店"))
	if dialogue != "" {
		g.addOutput(g.renderer.LogInfo(fmt.Sprintf("  店主说: \"%s\"", dialogue)))
	}
	g.addOutput("")
	g.addOutput(g.renderer.RenderShop(items, g.state.Player.Gold))
}

func (g *Game) processShopInput(cmd string, args []string) {
	if cmd == "0" || cmd == "back" || cmd == "leave" {
		g.state.Phase = PhaseExplore
		g.state.Dungeon.GetCurrentRoom().Cleared = true
		g.addOutput(g.renderer.LogInfo("你离开了商店。"))
		g.addOutput(g.renderer.RenderRoom(g.state.Dungeon.GetCurrentRoom(), g.state.Dungeon))
		return
	}

	idx, err := strconv.Atoi(cmd)
	if err != nil || idx < 1 || idx > len(g.state.ShopItems) {
		g.addOutput(g.renderer.LogError("无效的选择。输入 0 返回。"))
		return
	}

	item := g.state.ShopItems[idx-1]
	price := 10 + len(item.Effects)*5

	if g.state.Player.Gold < price {
		g.addOutput(g.renderer.LogError("金币不足！"))
		return
	}

	g.state.Player.Gold -= price
	g.state.Player.Inventory = append(g.state.Player.Inventory, item)
	g.addOutput(g.renderer.LogInfo(fmt.Sprintf("购买了 %s！", item.Name)))

	// 从商店移除
	g.state.ShopItems = append(g.state.ShopItems[:idx-1], g.state.ShopItems[idx:]...)
	g.addOutput(g.renderer.RenderShop(g.state.ShopItems, g.state.Player.Gold))
}

// === 事件系统 ===
func (g *Game) startEvent() {
	eventResp, err := g.llm.GenerateEvent(g.ctx, llm.EventRequest{
		Depth:      g.state.Depth,
		RoomType:   "event",
		PlayerInfo: g.state.Player.String(),
	})
	if err != nil {
		g.addOutput(g.renderer.LogError("事件生成失败"))
		g.state.Dungeon.GetCurrentRoom().Cleared = true
		g.state.Phase = PhaseExplore
		return
	}

	g.state.EventData = &EventState{
		Title:       eventResp.Title,
		Description: eventResp.Description,
	}
	for _, c := range eventResp.Choices {
		g.state.EventData.Choices = append(g.state.EventData.Choices, EventChoiceState{
			Text:        c.Text,
			Description: c.Description,
			Outcome:     c.Outcome,
		})
	}
	g.state.Phase = PhaseEvent

	g.addOutput(g.renderer.RenderEvent(&game.Event{
		Title:       eventResp.Title,
		Description: eventResp.Description,
	}))
}

func (g *Game) processEventInput(cmd string, args []string) {
	idx, err := strconv.Atoi(cmd)
	if err != nil || idx < 1 || idx > len(g.state.EventData.Choices) {
		g.addOutput(g.renderer.LogError("无效的选择"))
		return
	}

	choice := g.state.EventData.Choices[idx-1]

	// 根据结果类型处理
	switch choice.Outcome {
	case "good":
		// 好结果：恢复 HP 或获得物品
		if rand.Intn(100) < 50 {
			heal := 20 + g.state.Depth*5
			g.state.Player.Heal(heal)
			g.addOutput(g.renderer.LogInfo(fmt.Sprintf("恢复了 %d 点生命值！", heal)))
		} else {
			item := game.GenerateConsumableLoot(g.state.Depth)
			g.state.Player.Inventory = append(g.state.Player.Inventory, item)
			g.addOutput(g.renderer.LogInfo(fmt.Sprintf("获得了 %s！", item.Name)))
		}
	case "bad":
		// 坏结果：受到伤害
		damage := 10 + g.state.Depth*5
		g.state.Player.HP -= damage
		if g.state.Player.HP < 0 {
			g.state.Player.HP = 1
		}
		g.addOutput(g.renderer.LogError(fmt.Sprintf("受到了 %d 点伤害！", damage)))
	}

	g.addOutput(g.renderer.RenderEventChoiceResult(game.EventChoice{
		Text:        choice.Text,
		Description: choice.Description,
		Outcome:     choice.Outcome,
	}))

	g.state.Dungeon.GetCurrentRoom().Cleared = true
	g.state.EventData = nil
	g.state.Phase = PhaseExplore

	if !g.state.Player.IsAlive() {
		g.gameOver(false)
		return
	}

	g.addOutput(g.renderer.RenderRoom(g.state.Dungeon.GetCurrentRoom(), g.state.Dungeon))
}

// === 宝箱 ===
func (g *Game) openTreasure() {
	room := g.state.Dungeon.GetCurrentRoom()

	// LLM 生成宝箱描述
	desc, _ := g.llm.GenerateNarrative(g.ctx, llm.NarrativeRequest{
		RoomType:  "treasure",
		Depth:     g.state.Depth,
		RoomIndex: room.Index,
	})
	if desc != "" {
		g.addOutput(g.renderer.LogInfo(desc))
	}

	// 掉落
	var items []*game.Item
	gold := 15 + g.state.Depth*8 + rand.Intn(20)

	// 1-2 件物品
	itemCount := 1 + rand.Intn(2)
	for i := 0; i < itemCount; i++ {
		if rand.Intn(100) < 40 {
			items = append(items, game.GenerateConsumableLoot(g.state.Depth))
		} else {
			items = append(items, game.GenerateLoot(g.state.Depth))
		}
	}

	g.state.Player.Gold += gold
	for _, item := range items {
		g.state.Player.Inventory = append(g.state.Player.Inventory, item)
	}

	g.addOutput(g.renderer.RenderLoot(items, gold, 0))
	room.Cleared = true
}

// === 休息点 ===
func (g *Game) restAtCamp() {
	room := g.state.Dungeon.GetCurrentRoom()

	desc, _ := g.llm.GenerateNarrative(g.ctx, llm.NarrativeRequest{
		RoomType:  "rest",
		Depth:     g.state.Depth,
		RoomIndex: room.Index,
	})
	if desc != "" {
		g.addOutput(g.renderer.LogInfo(desc))
	}

	healAmount := g.state.Player.MaxHP / 3
	g.state.Player.Heal(healAmount)
	g.addOutput(g.renderer.LogInfo(fmt.Sprintf("🏕 你休息了一会，恢复了 %d 点生命值。", healAmount)))
	g.addOutput(g.renderer.LogInfo(fmt.Sprintf("  当前 HP: %d/%d", g.state.Player.HP, g.state.Player.MaxHP)))

	// 休息时有几率恢复技能冷却
	for _, skill := range g.state.Player.Skills {
		if skill.CurCD > 0 {
			skill.CurCD--
		}
	}

	room.Cleared = true
}

// === 装备系统 ===
func (g *Game) equipItem(args string) {
	idx, err := strconv.Atoi(args)
	if err != nil || idx < 1 || idx > len(g.state.Player.Inventory) {
		g.addOutput(g.renderer.LogError("无效的物品编号"))
		return
	}

	item := g.state.Player.Inventory[idx-1]
	if item.Type == game.Consumable {
		g.addOutput(g.renderer.LogError("消耗品无法装备"))
		return
	}

	// 装备物品
	slot := item.Slot
	if slot == "" {
		slot = game.SlotForType(item.Type)
	}

	// 如果已有装备，卸下
	if old, ok := g.state.Player.Equipment[slot]; ok && old != nil {
		g.state.Player.Inventory = append(g.state.Player.Inventory, old)
	}

	g.state.Player.Equipment[slot] = item
	// 从背包移除
	g.state.Player.Inventory = append(g.state.Player.Inventory[:idx-1], g.state.Player.Inventory[idx:]...)
	g.state.Player.RecalcStats()

	g.addOutput(g.renderer.LogInfo(fmt.Sprintf("装备了 %s！", item.Name)))
}

// === 游戏结束 ===
func (g *Game) gameOver(won bool) {
	if won {
		g.state.Phase = PhaseVictory
	} else {
		g.state.Phase = PhaseGameOver
	}
	g.addOutput(g.renderer.RenderGameOver(won, g.state.Player))
}

func (g *Game) processEndInput(cmd string, args []string) {
	switch cmd {
	case "restart":
		g.state = &GameState{Phase: PhaseTitle}
		g.addOutput(g.renderer.RenderWelcome())
	case "quit", "exit":
		g.addOutput(g.renderer.LogInfo("再见！"))
	default:
		g.addOutput(g.renderer.LogInfo("输入 'restart' 重新开始，或 'quit' 退出。"))
	}
}

// RenderCurrentView 渲染当前视图
func (g *Game) RenderCurrentView() string {
	var sb strings.Builder

	// 如果是伪装模式，显示伪装界面
	if g.renderer.IsStealth() {
		return g.renderer.RenderStealthScreen(g.state.Player)
	}

	// 状态栏
	if g.state.Player != nil && g.state.Phase != PhaseTitle {
		sb.WriteString(g.renderer.RenderStatusBar(g.state.Player))
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat("─", 60))
		sb.WriteString("\n")
	}

	// 输出缓冲
	for _, line := range g.GetOutput() {
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	return sb.String()
}

// showConfig 显示当前配置信息
func (g *Game) showConfig() {
	g.addOutput(g.renderer.LogInfo("=== LLM 配置 ==="))
	g.addOutput(fmt.Sprintf("  %s %s", display.StatLabel.Render("状态:"),
		display.SuccessStyle.Render("Mock 模式（内置）")))
	g.addOutput("")
	g.addOutput(g.renderer.LogInfo("如需启用 LLM，编辑 data/config.yaml:"))
	g.addOutput(fmt.Sprintf("  %s", display.StatLabel.Render("llm.enabled: true")))
	g.addOutput(fmt.Sprintf("  %s", display.StatLabel.Render("llm.api_key: \"your-api-key\"")))
	g.addOutput(fmt.Sprintf("  %s", display.StatLabel.Render("llm.base_url: \"https://api.openai.com/v1\"")))
	g.addOutput(fmt.Sprintf("  %s", display.StatLabel.Render("llm.model: \"gpt-4o-mini\"")))
	g.addOutput("")
	g.addOutput(g.renderer.LogInfo("支持的 API 格式: OpenAI 兼容"))
	g.addOutput(fmt.Sprintf("  %s", display.StatLabel.Render("OpenAI / DeepSeek / Moonshot / Ollama / 自定义")))
}
