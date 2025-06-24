package opt

import (
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"os"
)

type Global struct {
	Url            *string
	Query          *string
	Timeout        *int
	ConnectTimeout *int
	Limit          *int
	ConvertJson    *bool
}

type CommandsOption struct {
	GlobalCommand *Global
}

func NewCommands() *CommandsOption {
	return &CommandsOption{
		GlobalCommand: &Global{},
	}
}

func ParseOptions(Args []string) (*CommandsOption, error) {
	app := kingpin.New("qnoracle", "a versatile command-line tool").Version("1.0.0")
	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')

	command := NewCommands()
	command.GlobalCommand.Url = app.Flag("url", "oracle connect url string").Short('u').Required().String()
	command.GlobalCommand.Query = app.Flag("query", "sql query").Short('q').Default("select 1 from dual").String()
	command.GlobalCommand.Timeout = app.Flag("timeout", "query sql timeout").Short('t').Default("32").Int()
	command.GlobalCommand.ConnectTimeout = app.Flag("connect-timeout", "connect to oracle timeout").Short('c').Default("6").Int()
	command.GlobalCommand.Limit = app.Flag("limit", "limit rows").Short('l').Default("1000").Int()
	command.GlobalCommand.ConvertJson = app.Flag("json", "convert data to json format").Short('j').Bool()
	_, err := app.Parse(Args[1:])
	if err != nil {
		return nil, err
	}

	return command, nil
}

func IsFile(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("error checking path %s: %w", path, err)
	}
	return !fileInfo.IsDir(), nil
}
