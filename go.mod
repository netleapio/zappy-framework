module github.com/netleapio/zappy-framework

go 1.19

require tinygo.org/x/drivers v0.22.0

require github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect

replace (
	github.com/netleapio/zappy-framework => ../zappy-framework
	tinygo.org/x/drivers => ../tinygo-drivers/main
	tinygo.org/x/tinyfont => ../tinyfont
)
