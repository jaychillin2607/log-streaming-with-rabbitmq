package types

type LogLevel int

const (
	INFO LogLevel = iota
	DEBUG
	WARNING
	ERROR
)

type LogMessage struct {
	Level LogLevel `json:"level"`
	Msg   string   `json:"msg"`
}
