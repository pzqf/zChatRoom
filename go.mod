module zChatRoom

go 1.18

require (
	github.com/huichen/sego v0.0.0-20210824061530-c87651ea5c76
	github.com/jroimartin/gocui v0.5.0
	github.com/pkg/profile v1.6.0
	github.com/pzqf/zEngine v0.0.1
	github.com/pzqf/zUtil v0.0.1
	github.com/satori/go.uuid v1.2.0
	go.uber.org/zap v1.21.0
)

require (
	github.com/adamzy/cedar-go v0.0.0-20170805034717-80a9c64b256d // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/nsf/termbox-go v1.1.1 // indirect
	github.com/panjf2000/ants v1.3.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace github.com/pzqf/zEngine => ../zEngine

replace github.com/pzqf/zUtil => ../zUtil
