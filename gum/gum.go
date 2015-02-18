package gum

import (
	"errors"
	"net/http"
	"os"
	"syscall"

	"hsproom/config"
	"hsproom/controllers"
	"hsproom/models"
	"hsproom/plugins"
	"hsproom/templates"
	"hsproom/utils"

	"github.com/gorilla/context"
)

func init() {

	utils.LogFile = config.LogFile
	utils.DisplayLog = config.DisplayLog
	utils.LogLevel = config.LogLevel

}

func Del() {
	models.Del()
	controllers.Del()
	templates.Del()
	plugins.Del()
}

func Start() {
	for route := range controllers.Router.Iterator() {
		http.HandleFunc(route.Path, route.Function)
		utils.PromulgateDebugStr(os.Stdout, route.Path+"に関数を割当")
	}

	http.Handle("/"+config.StaticPath, http.StripPrefix("/"+config.StaticPath, http.FileServer(http.Dir(config.StaticPath))))
	utils.PromulgateDebugStr(os.Stdout, "/"+config.StaticPath+"に静的コンテンツを割当")

	utils.PromulgateInfoStr(os.Stdout, "ポート"+config.ServerPort+"でサーバを開始...")
	http.ListenAndServe(":"+config.ServerPort, context.ClearHandler(http.DefaultServeMux))
}

func Daemonize() error {
	var ret uintptr
	var err syscall.Errno

	ret, _, err = syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if err != 0 {
		return errors.New("デーモン化に失敗")
	}
	switch ret {
	case 0:
		break
	default:
		os.Exit(0)
	}

	pid, _ := syscall.Setsid()
	if pid == -1 {
		return errors.New("デーモン化に失敗")
	}

	os.Chdir("/")

	syscall.Umask(0)

	f, e := os.OpenFile("/dev/null", os.O_RDWR, 0)
	if e == nil {
		fd := int(f.Fd())
		syscall.Dup2(fd, int(os.Stdin.Fd()))
		syscall.Dup2(fd, int(os.Stdout.Fd()))
		syscall.Dup2(fd, int(os.Stderr.Fd()))
	}

	os.Chdir(config.RootPath)

	return nil
}
