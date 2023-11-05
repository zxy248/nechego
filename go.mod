module nechego

go 1.21

require (
	github.com/antonmedv/expr v1.9.0
	golang.org/x/exp v0.0.0-20230626212559-97b1e661b5df
	golang.org/x/text v0.3.7
	gopkg.in/telebot.v3 v3.1.2
)

replace gopkg.in/telebot.v3 => ../telebot
