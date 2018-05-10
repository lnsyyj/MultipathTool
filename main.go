package main

import (
	"flag"
	"os"
	"log"
	"fmt"
	"os/exec"
)

var MyLog = &log.Logger{}
var VdbenchLog = &log.Logger{}

func getMultipathStatus() {
	cmd := exec.Command("multipath", "-l")
	out, err := cmd.Output()
	if err != nil {
		MyLog.Printf("\n%s", err)
	}
	MyLog.Printf("\n%s", string(out))
}

func configureMultipath() {
	file, err := os.Create("/etc/multipath.conf")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	fmt.Fprintf(file,
		`defaults{
        user_friendly_names yes
        polling_interval 10
        checker_timeout 120
        queue_without_daemon no
}

blacklist {
        devnode "^(ram|raw|loop|fd|md|dm-|sr|scd|st)[0-9]*"
        devnode "^hd[a-z]"
}

devices {
        device{
                path_grouping_policy failover
        }
}`)
}

func execVdbench() {
	cmd := exec.Command("/home/yujiang/vdbench/vdbench", "-f /home/yujiang/vdbench/vtest.file")
	out, err := cmd.Output()
	if err != nil {
		MyLog.Printf("\n%s", err)
	}
	MyLog.Printf("\n%s", string(out))
}

func main() {

	var CutOffMultipath, LogPath,VdbenchLogPath string

	flag.StringVar(&LogPath, "LogPath", "./main.log", "program execution log")
	flag.StringVar(&VdbenchLogPath, "VdbenchLogPath", "./vdbench.log", "Collect vdbench output")
	flag.StringVar(&CutOffMultipath, "CutOffMultipath", "active", "active/standby mode, which link is cut off, Cut off standby mode: enabled")
	flag.Parse()
	//plog.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	file_main, err_main := os.OpenFile(LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err_main != nil {
		log.Fatalln("Failed to open log file : ", err_main)
	}
	defer file_main.Close()
	MyLog = log.New(file_main, "",	log.Ldate|log.Ltime|log.Lshortfile)

	file_vdbench, err_vdbench := os.OpenFile(LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err_vdbench != nil {
		log.Fatalln("Failed to open log file : ", err_vdbench)
	}
	defer file_vdbench.Close()
	VdbenchLog = log.New(file_vdbench, "",	log.Ldate|log.Ltime|log.Lshortfile)

	getMultipathStatus()
	execVdbench()

	//executionlog.Info.Println("Special Information")
	//executionlog.Warning.Println("There is something you need to know about")
	//executionlog.Error.Println("Something has failed")

}