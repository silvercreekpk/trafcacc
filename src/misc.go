// misc.go 辅助的 helper functions

package trafcacc

import (
	"net"
	"runtime"
	"strings"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
)

const maxopenfile = 3267600

func dialTimeout(network, address string, timeout time.Duration) (conn net.Conn, err error) {
	m := int(timeout / time.Second)
	for i := 0; i < m; i++ {
		conn, err = net.DialTimeout(network, address, timeout)
		if err == nil || !strings.Contains(err.Error(), "can't assign requested address") {
			break
		}
		time.Sleep(time.Second)
	}
	return
}

func increaseMaxopenfile() {

	var lim syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &lim); err != nil {
		log.Infoln("failed to get NOFILE rlimit: ", err)
	}

	if lim.Cur < maxopenfile || lim.Max < maxopenfile {
		if lim.Cur < maxopenfile {
			lim.Cur = maxopenfile
		}
		if lim.Max < maxopenfile {
			lim.Max = maxopenfile
		}

		if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &lim); err != nil {
			log.Infoln("failed to set NOFILE rlimit: ", err)
		}
	}
}

func increaseGomaxprocs() {
	cpu := runtime.NumCPU()
	if cpu > runtime.GOMAXPROCS(-1) {
		runtime.GOMAXPROCS(cpu)
	}
}
