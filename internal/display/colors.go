package display

import "github.com/charmbracelet/lipgloss"

// 颜色主题 - 暗色系，适合长时间游玩
var (
	// 基础颜色
	ColorPrimary   = lipgloss.Color("#00FF88") // 亮绿 - 日志风格的主色
	ColorSecondary = lipgloss.Color("#00AAFF") // 亮蓝
	ColorAccent    = lipgloss.Color("#FFD700") // 金色 - 重要信息
	ColorDanger    = lipgloss.Color("#FF4444") // 红色 - 危险/伤害
	ColorWarning   = lipgloss.Color("#FF8800") // 橙色 - 警告
	ColorMuted     = lipgloss.Color("#666666") // 灰色 - 次要信息
	ColorWhite     = lipgloss.Color("#CCCCCC") // 浅灰 - 普通文本

	// 日志级别颜色
	ColorInfo  = lipgloss.Color("#00CC66") // INFO - 绿色
	ColorWarn  = lipgloss.Color("#FFAA00") // WARN - 橙色
	ColorError = lipgloss.Color("#FF3333") // ERROR - 红色
	ColorFatal = lipgloss.Color("#FF0000") // FATAL - 亮红
)

// 品质颜色
var RarityStyle = map[string]lipgloss.Style{
	"普通": lipgloss.NewStyle().Foreground(ColorWhite),
	"优秀": lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")),
	"稀有": lipgloss.NewStyle().Foreground(lipgloss.Color("#0088FF")),
	"史诗": lipgloss.NewStyle().Foreground(lipgloss.Color("#AA44FF")),
	"传说": lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8800")),
}

// 样式定义
var (
	// 日志格式样式
	LogTimestamp = lipgloss.NewStyle().Foreground(ColorMuted)
	LogInfo      = lipgloss.NewStyle().Foreground(ColorInfo).Bold(true)
	LogWarn      = lipgloss.NewStyle().Foreground(ColorWarning).Bold(true)
	LogError     = lipgloss.NewStyle().Foreground(ColorDanger).Bold(true)
	LogFatal     = lipgloss.NewStyle().Foreground(ColorFatal).Bold(true)
	LogMessage   = lipgloss.NewStyle().Foreground(ColorWhite)

	// 游戏 UI 样式
	TitleStyle    = lipgloss.NewStyle().Foreground(ColorAccent).Bold(true)
	HPStyle       = lipgloss.NewStyle().Foreground(ColorDanger)
	HPBarStyle    = lipgloss.NewStyle().Foreground(ColorPrimary)
	MPStyle       = lipgloss.NewStyle().Foreground(ColorSecondary)
	StatLabel     = lipgloss.NewStyle().Foreground(ColorMuted)
	StatValue     = lipgloss.NewStyle().Foreground(ColorWhite).Bold(true)
	CombatTitle   = lipgloss.NewStyle().Foreground(ColorDanger).Bold(true)
	SkillStyle    = lipgloss.NewStyle().Foreground(ColorSecondary)
	ItemStyle     = lipgloss.NewStyle().Foreground(ColorAccent)
	GoldStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	BossStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true)
	CritStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).Bold(true)
	DodgeStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00AAFF")).Italic(true)
	BuffStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF88"))
	DebuffStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444"))
	PromptStyle   = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	ErrorStyle    = lipgloss.NewStyle().Foreground(ColorDanger)
	SuccessStyle  = lipgloss.NewStyle().Foreground(ColorPrimary)
	StealthHeader = lipgloss.NewStyle().Foreground(ColorMuted)
	AccentStyle   = lipgloss.NewStyle().Foreground(ColorAccent)
)
