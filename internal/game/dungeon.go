package game

import "math/rand"

// RoomType 房间类型
type RoomType string

const (
	RoomEntrance RoomType = "entrance"
	RoomCombat   RoomType = "combat"
	RoomTreasure RoomType = "treasure"
	RoomShop     RoomType = "shop"
	RoomRest     RoomType = "rest"
	RoomEvent    RoomType = "event"
	RoomBoss     RoomType = "boss"
)

// Room 房间
type Room struct {
	Type        RoomType
	Index       int
	Description string // LLM 生成的描述
	Visited     bool
	Cleared     bool
	Connections []int  // 连接的房间索引
}

// Dungeon 地牢
type Dungeon struct {
	Depth  int
	Rooms  []*Room
	Current int // 当前房间索引
	Floors int  // 总层数（每层一个地牢）
}

// GenerateDungeon 生成地牢
func GenerateDungeon(depth int) *Dungeon {
	roomCount := 8 + rand.Intn(5) // 8-12 个房间
	rooms := make([]*Room, roomCount)

	// 第一个房间是入口
	rooms[0] = &Room{
		Type:    RoomEntrance,
		Index:   0,
		Visited: true,
	}

	// 最后一个房间是 Boss
	rooms[roomCount-1] = &Room{
		Type:    RoomBoss,
		Index:   roomCount - 1,
	}

	// 中间房间随机分配类型
	for i := 1; i < roomCount-1; i++ {
		rooms[i] = &Room{
			Type:  randomRoomType(depth),
			Index: i,
		}
	}

	// 确保至少有 1 个商店和 1 个休息点
	hasShop := false
	hasRest := false
	for _, r := range rooms {
		if r.Type == RoomShop {
			hasShop = true
		}
		if r.Type == RoomRest {
			hasRest = true
		}
	}
	if !hasShop {
		idx := 1 + rand.Intn(roomCount-2)
		rooms[idx].Type = RoomShop
	}
	if !hasRest {
		idx := 1 + rand.Intn(roomCount-2)
		for rooms[idx].Type == RoomShop {
			idx = 1 + rand.Intn(roomCount-2)
		}
		rooms[idx].Type = RoomRest
	}

	// 生成连接（简单线性 + 一些分支）
	for i := 0; i < roomCount-1; i++ {
		rooms[i].Connections = append(rooms[i].Connections, i+1)
		rooms[i+1].Connections = append(rooms[i+1].Connections, i)
	}

	// 添加一些随机分支连接
	extraConnections := rand.Intn(3) + 1
	for i := 0; i < extraConnections; i++ {
		from := rand.Intn(roomCount - 2)
		to := from + 2 + rand.Intn(min(3, roomCount-from-2))
		if to < roomCount && !contains(rooms[from].Connections, to) {
			rooms[from].Connections = append(rooms[from].Connections, to)
			rooms[to].Connections = append(rooms[to].Connections, from)
		}
	}

	return &Dungeon{
		Depth:   depth,
		Rooms:   rooms,
		Current: 0,
		Floors:  depth,
	}
}

func randomRoomType(depth int) RoomType {
	roll := rand.Intn(100)
	switch {
	case roll < 40:
		return RoomCombat
	case roll < 55:
		return RoomTreasure
	case roll < 65:
		return RoomEvent
	case roll < 75:
		return RoomShop
	case roll < 85:
		return RoomRest
	default:
		return RoomCombat
	}
}

func contains(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// GetCurrentRoom 获取当前房间
func (d *Dungeon) GetCurrentRoom() *Room {
	return d.Rooms[d.Current]
}

// MoveTo 移动到指定房间
func (d *Dungeon) MoveTo(index int) bool {
	if index < 0 || index >= len(d.Rooms) {
		return false
	}
	// 检查是否可以从当前房间移动到目标房间
	for _, conn := range d.Rooms[d.Current].Connections {
		if conn == index {
			d.Current = index
			d.Rooms[index].Visited = true
			return true
		}
	}
	return false
}

// GetConnectedRooms 获取可前往的房间
func (d *Dungeon) GetConnectedRooms() []*Room {
	var rooms []*Room
	for _, idx := range d.Rooms[d.Current].Connections {
		rooms = append(rooms, d.Rooms[idx])
	}
	return rooms
}

// IsComplete 地牢是否通关
func (d *Dungeon) IsComplete() bool {
	return d.Rooms[len(d.Rooms)-1].Cleared
}
