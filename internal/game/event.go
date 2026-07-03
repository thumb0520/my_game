package game

// Event 随机事件
type Event struct {
	Title       string
	Description string
	Choices     []EventChoice
}

// EventChoice 事件选项
type EventChoice struct {
	Text        string
	Description string
	Outcome     string // good / neutral / bad
}
