package logger

import (
	"flag"
	"log"
	"os"
)

var (
	Log *log.Logger
)

func init() {

	var logPath = "./logger/info.log"
	// set location of log file
	// if Config.Logpath == ""{
	// 	logPath = Config.LogPath
	// }

	flag.Parse()
	var file, err1 = os.Create(logPath)

	if err1 != nil {
		panic(err1)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
	Log.Println("LogFile : " + logPath)
}
