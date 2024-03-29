package main

import (
	"fmt"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"www.velocidex.com/golang/timetracker"
)

var (
	app = kingpin.New("timetracker",
		"A tool for tracking time - based on wakatime.")

	filename_flag = app.Flag("file", "File being edited").String()
	plugin_flag   = app.Flag("plugin", "Caller plugin name").String()
	time_flag     = app.Flag("time", "Timestamp").Float64()
	write_flag    = app.Flag("write", "Set when we write the file").Bool()
	key_flag      = app.Flag("key", "API Key (ignored)").String()

	dump_cmd_arg = app.Flag("parse", "The file to parse").String()
	debug_flag   = app.Flag("debug", "Pring Debugging").Bool()
)

func doParse(filename string) {
	fd, err := os.Open(filename)
	kingpin.FatalIfError(err, "ProcessFile")

	context := timetracker.NewContext(nil)
	context.Debug = *debug_flag

	err = timetracker.ProcessFile(fd, context)
	kingpin.FatalIfError(err, "ProcessFile")

	fmt.Println(context.Dump())
}

func main() {
	app.HelpFlag.Short('h')
	app.UsageTemplate(kingpin.CompactUsageTemplate)
	kingpin.MustParse(app.Parse(os.Args[1:]))
	if *dump_cmd_arg != "" {
		doParse(*dump_cmd_arg)
		return
	}

	err := timetracker.WriteDataPoint(*filename_flag, *time_flag)
	kingpin.FatalIfError(err, "WriteDataPoint")
}
