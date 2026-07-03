package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"dungeonlog/internal/game"
)

// GameState 游戏状态
type GameState struct {
	Player      *game.Character
	FloorMap    *game.FloorMap
	Depth       int
	Phase       GamePhase
	CombatState *game.CombatState
	ShopItems   []*game.Item
	EventData   *EventState
	Meta        *game.MetaProgress // 局外养成

	// 战斗暂停控制
	Waiting              string // "" = 正常输入, "continue" = 等待按回车继续
	PendingMonsterAction string // 待执行的怪物行动描述
}

// GamePhase 游戏阶段
type GamePhase string

const (
	PhaseTitle     GamePhase = "title"
	PhaseTown      GamePhase = "town" // 城镇（局外养成）
	PhaseExplore   GamePhase = "explore"
	PhaseCombat    GamePhase = "combat"
	PhaseShop      GamePhase = "shop"
	PhaseEvent     GamePhase = "event"
	PhaseGameOver  GamePhase = "gameover"
	PhaseVictory   GamePhase = "victory"
	PhaseCharacter GamePhase = "character" // 角色创建
)

// EventState 事件状态
type EventState struct {
	Title       string
	Description string
	Choices     []EventChoiceState
}

// EventChoiceState 事件选项状态
type EventChoiceState struct {
	Text        string
	Description string
	Outcome     string
}

// SaveData 存档数据
type SaveData struct {
	PlayerName  string              `json:"player_name"`
	PlayerClass string              `json:"player_class"`
	Level       int                 `json:"level"`
	Exp         int                 `json:"exp"`
	HP          int                 `json:"hp"`
	MaxHP       int                 `json:"max_hp"`
	Gold        int                 `json:"gold"`
	Depth       int                 `json:"depth"`
	STR         int                 `json:"str"`
	DEX         int                 `json:"dex"`
	INT         int                 `json:"int"`
	VIT         int                 `json:"vit"`
	LUK         int                 `json:"luk"`
	Inventory   []SaveItem          `json:"inventory"`
	Equipment   map[string]SaveItem `json:"equipment"`
	Meta        *game.MetaProgress  `json:"meta"` // 局外养成
}

// SaveItem 存档物品
type SaveItem struct {
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Rarity      string         `json:"rarity"`
	Description string         `json:"description"`
	Effects     map[string]int `json:"effects"`
	Slot        string         `json:"slot,omitempty"`
}

// SaveGame 保存游戏
func SaveGame(state *GameState, path string) error {
	if state.Player == nil {
		return fmt.Errorf("no player to save")
	}

	saveData := SaveData{
		PlayerName:  state.Player.Name,
		PlayerClass: string(state.Player.Class),
		Level:       state.Player.Level,
		Exp:         state.Player.Exp,
		HP:          state.Player.HP,
		MaxHP:       state.Player.MaxHP,
		Gold:        state.Player.Gold,
		Depth:       state.Depth,
		STR:         state.Player.STR,
		DEX:         state.Player.DEX,
		INT:         state.Player.INT,
		VIT:         state.Player.VIT,
		LUK:         state.Player.LUK,
		Inventory:   make([]SaveItem, 0),
		Equipment:   make(map[string]SaveItem),
		Meta:        state.Meta,
	}

	for _, item := range state.Player.Inventory {
		saveData.Inventory = append(saveData.Inventory, SaveItem{
			Name:        item.Name,
			Type:        string(item.Type),
			Rarity:      string(item.Rarity),
			Description: item.Description,
			Effects:     item.Effects,
			Slot:        string(item.Slot),
		})
	}

	for slot, item := range state.Player.Equipment {
		if item != nil {
			saveData.Equipment[string(slot)] = SaveItem{
				Name:        item.Name,
				Type:        string(item.Type),
				Rarity:      string(item.Rarity),
				Description: item.Description,
				Effects:     item.Effects,
				Slot:        string(item.Slot),
			}
		}
	}

	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// LoadGame 加载游戏
func LoadGame(path string) (*GameState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var saveData SaveData
	if err := json.Unmarshal(data, &saveData); err != nil {
		return nil, err
	}

	class := game.Class(saveData.PlayerClass)
	player := game.NewCharacter(saveData.PlayerName, class)
	player.Level = saveData.Level
	player.Exp = saveData.Exp
	player.HP = saveData.HP
	player.MaxHP = saveData.MaxHP
	player.Gold = saveData.Gold
	player.STR = saveData.STR
	player.DEX = saveData.DEX
	player.INT = saveData.INT
	player.VIT = saveData.VIT
	player.LUK = saveData.LUK

	// 加载物品
	for _, si := range saveData.Inventory {
		player.Inventory = append(player.Inventory, &game.Item{
			Name:        si.Name,
			Type:        game.ItemType(si.Type),
			Rarity:      game.Rarity(si.Rarity),
			Description: si.Description,
			Effects:     si.Effects,
			Slot:        game.Slot(si.Slot),
		})
	}

	// 加载装备
	for slot, si := range saveData.Equipment {
		player.Equipment[game.Slot(slot)] = &game.Item{
			Name:        si.Name,
			Type:        game.ItemType(si.Type),
			Rarity:      game.Rarity(si.Rarity),
			Description: si.Description,
			Effects:     si.Effects,
			Slot:        game.Slot(si.Slot),
		}
	}

	player.RecalcStats()

	state := &GameState{
		Player: player,
		Depth:  saveData.Depth,
		Phase:  PhaseExplore,
	}

	// 加载养成数据
	if saveData.Meta != nil {
		state.Meta = saveData.Meta
	} else {
		state.Meta = game.NewMetaProgress()
	}

	// 生成新地图
	state.FloorMap = game.GenerateFloorMap(state.Depth)

	return state, nil
}
