package game

import (
	"fmt"
	"math/rand"
)

// RoomType 房间类型
type RoomType string

const (
	RoomEntrance RoomType = "entrance"
	RoomCombat   RoomType = "combat"
	RoomElite    RoomType = "elite"
	RoomTreasure RoomType = "treasure"
	RoomShop     RoomType = "shop"
	RoomRest     RoomType = "rest"
	RoomEvent    RoomType = "event"
	RoomBoss     RoomType = "boss"
)

// MapNode 地图节点
type MapNode struct {
	Row     int
	Col     int
	Type    RoomType
	Visited bool
	Cleared bool
}

// MapConnection 地图连接
type MapConnection struct {
	FromRow int
	FromCol int
	ToRow   int
	ToCol   int
}

// FloorMap 一层地图（杀戮尖塔风格）
type FloorMap struct {
	Depth       int
	RowCount    int
	Nodes       [][]*MapNode // [row][col]
	Connections []MapConnection
	CurrentRow  int
	CurrentCol  int
}

// NodeKey 生成节点的唯一键
func NodeKey(row, col int) string {
	return fmt.Sprintf("%d,%d", row, col)
}

// GenerateFloorMap 生成杀戮尖塔风格的网状地图
func GenerateFloorMap(depth int) *FloorMap {
	rowCount := 12 + rand.Intn(4) // 12-15 行
	fm := &FloorMap{
		Depth:    depth,
		RowCount: rowCount,
		Nodes:    make([][]*MapNode, rowCount),
	}

	// 第 0 行：入口（1 个节点）
	fm.Nodes[0] = []*MapNode{{Row: 0, Col: 0, Type: RoomEntrance, Visited: true, Cleared: true}}

	// 最后一行：Boss（1 个节点）
	fm.Nodes[rowCount-1] = []*MapNode{{Row: rowCount - 1, Col: 0, Type: RoomBoss}}

	// 中间行：随机节点
	for row := 1; row < rowCount-1; row++ {
		nodeCount := 3 + rand.Intn(5) // 3-7 个节点
		fm.Nodes[row] = make([]*MapNode, nodeCount)
		for col := 0; col < nodeCount; col++ {
			fm.Nodes[row][col] = &MapNode{
				Row:  row,
				Col:  col,
				Type: randomNodeType(row, rowCount, depth),
			}
		}
	}

	// 确保关键节点存在
	ensureKeyNodes(fm)

	// 生成连接
	fm.Connections = generateConnections(fm)

	// 设置当前位置为入口
	fm.CurrentRow = 0
	fm.CurrentCol = 0

	return fm
}

// randomNodeType 根据行位置随机生成节点类型
func randomNodeType(row, totalRows, depth int) RoomType {
	// 精英怪只出现在中后段（约 40%-70% 位置）
	isEliteZone := row > totalRows*2/5 && row < totalRows*7/10

	roll := rand.Intn(100)

	// 精英区域有更高概率出精英
	if isEliteZone && roll < 25 {
		return RoomElite
	}

	switch {
	case roll < 40:
		return RoomCombat
	case roll < 55:
		return RoomEvent
	case roll < 65:
		return RoomRest
	case roll < 75:
		return RoomShop
	case roll < 85:
		return RoomTreasure
	default:
		if isEliteZone {
			return RoomElite
		}
		return RoomCombat
	}
}

// ensureKeyNodes 确保地图中有必要的节点类型
func ensureKeyNodes(fm *FloorMap) {
	hasShop := false
	hasRest := false
	hasElite := false

	for row := 1; row < fm.RowCount-1; row++ {
		for _, node := range fm.Nodes[row] {
			switch node.Type {
			case RoomShop:
				hasShop = true
			case RoomRest:
				hasRest = true
			case RoomElite:
				hasElite = true
			}
		}
	}

	// 确保至少有 1 个商店
	if !hasShop {
		row := 2 + rand.Intn(fm.RowCount/3)
		if row < fm.RowCount-1 && len(fm.Nodes[row]) > 0 {
			col := rand.Intn(len(fm.Nodes[row]))
			fm.Nodes[row][col].Type = RoomShop
		}
	}

	// 确保至少有 1 个休息点
	if !hasRest {
		row := fm.RowCount/2 + rand.Intn(fm.RowCount/4)
		if row < fm.RowCount-1 && len(fm.Nodes[row]) > 0 {
			col := rand.Intn(len(fm.Nodes[row]))
			for fm.Nodes[row][col].Type == RoomShop {
				col = rand.Intn(len(fm.Nodes[row]))
			}
			fm.Nodes[row][col].Type = RoomRest
		}
	}

	// 确保至少有 1 个精英怪
	if !hasElite {
		row := fm.RowCount*2/5 + rand.Intn(fm.RowCount/5)
		if row < fm.RowCount-1 && len(fm.Nodes[row]) > 0 {
			col := rand.Intn(len(fm.Nodes[row]))
			fm.Nodes[row][col].Type = RoomElite
		}
	}
}

// generateConnections 生成网状连接
func generateConnections(fm *FloorMap) []MapConnection {
	var conns []MapConnection

	for row := 0; row < fm.RowCount-1; row++ {
		currentNodes := fm.Nodes[row]
		nextNodes := fm.Nodes[row+1]

		if len(currentNodes) == 0 || len(nextNodes) == 0 {
			continue
		}

		// 为每个当前节点分配 1-2 个下一行的连接
		for _, node := range currentNodes {
			// 找到最近的下一行节点
			connected := make(map[int]bool)

			// 第一个连接：最近的节点
			bestCol := findClosest(node.Col, len(nextNodes))
			connected[bestCol] = true
			conns = append(conns, MapConnection{
				FromRow: row, FromCol: node.Col,
				ToRow: row + 1, ToCol: bestCol,
			})

			// 40% 概率添加第二个连接
			if rand.Intn(100) < 40 && len(nextNodes) > 1 {
				secondCol := bestCol
				attempts := 0
				for connected[secondCol] && attempts < 10 {
					secondCol = rand.Intn(len(nextNodes))
					attempts++
				}
				if !connected[secondCol] {
					connected[secondCol] = true
					conns = append(conns, MapConnection{
						FromRow: row, FromCol: node.Col,
						ToRow: row + 1, ToCol: secondCol,
					})
				}
			}
		}

		// 确保下一行每个节点至少有一个入边
		hasIncoming := make(map[int]bool)
		for _, c := range conns {
			if c.ToRow == row+1 {
				hasIncoming[c.ToCol] = true
			}
		}
		for col := 0; col < len(nextNodes); col++ {
			if !hasIncoming[col] {
				// 连接到上一行最近的节点
				bestFrom := findClosest(col, len(currentNodes))
				conns = append(conns, MapConnection{
					FromRow: row, FromCol: bestFrom,
					ToRow: row + 1, ToCol: col,
				})
			}
		}
	}

	return conns
}

// findClosest 找到最接近目标列的节点列
func findClosest(targetCol, maxCol int) int {
	if targetCol < 0 {
		return 0
	}
	if targetCol >= maxCol {
		return maxCol - 1
	}
	return targetCol
}

// GetCurrentNode 获取当前节点
func (fm *FloorMap) GetCurrentNode() *MapNode {
	if fm.CurrentRow >= 0 && fm.CurrentRow < fm.RowCount &&
		fm.CurrentCol >= 0 && fm.CurrentCol < len(fm.Nodes[fm.CurrentRow]) {
		return fm.Nodes[fm.CurrentRow][fm.CurrentCol]
	}
	return nil
}

// GetReachableNodes 获取可前往的节点
func (fm *FloorMap) GetReachableNodes() []*MapNode {
	var nodes []*MapNode
	seen := make(map[string]bool)

	for _, conn := range fm.Connections {
		if conn.FromRow == fm.CurrentRow && conn.FromCol == fm.CurrentCol {
			key := NodeKey(conn.ToRow, conn.ToCol)
			if !seen[key] {
				seen[key] = true
				if conn.ToRow < fm.RowCount && conn.ToCol < len(fm.Nodes[conn.ToRow]) {
					nodes = append(nodes, fm.Nodes[conn.ToRow][conn.ToCol])
				}
			}
		}
	}
	return nodes
}

// MoveTo 移动到指定节点
func (fm *FloorMap) MoveTo(row, col int) bool {
	// 检查是否是可达节点
	reachable := fm.GetReachableNodes()
	for _, node := range reachable {
		if node.Row == row && node.Col == col {
			fm.CurrentRow = row
			fm.CurrentCol = col
			node.Visited = true
			return true
		}
	}
	return false
}

// IsComplete 地牢是否通关
func (fm *FloorMap) IsComplete() bool {
	bossNode := fm.Nodes[fm.RowCount-1][0]
	return bossNode.Cleared
}

// GetNodeIcon 获取节点显示图标
func GetNodeIcon(n *RoomType) string {
	switch *n {
	case RoomEntrance:
		return "🚪"
	case RoomCombat:
		return "⚔"
	case RoomElite:
		return "🔥"
	case RoomTreasure:
		return "💎"
	case RoomShop:
		return "🏪"
	case RoomRest:
		return "🏕"
	case RoomEvent:
		return "❓"
	case RoomBoss:
		return "💀"
	default:
		return "·"
	}
}

// GetNodeLabel 获取节点显示标签
func GetNodeLabel(t RoomType) string {
	switch t {
	case RoomEntrance:
		return "入口"
	case RoomCombat:
		return "战斗"
	case RoomElite:
		return "精英"
	case RoomTreasure:
		return "宝箱"
	case RoomShop:
		return "商店"
	case RoomRest:
		return "休息"
	case RoomEvent:
		return "事件"
	case RoomBoss:
		return "BOSS"
	default:
		return "?"
	}
}

// contains 检查切片是否包含值
func contains(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
