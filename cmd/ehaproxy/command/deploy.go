package command

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"github.com/funkygao/gocli"
)

type Deploy struct {
	Ui  cli.Ui
	Cmd string

	root    string
	rsyslog bool
}

func (this *Deploy) Run(args []string) (exitCode int) {
	cmdFlags := flag.NewFlagSet("deploy", flag.ContinueOnError)
	cmdFlags.Usage = func() { this.Ui.Output(this.Help()) }
	cmdFlags.StringVar(&this.root, "p", defaultPrefix, "")
	cmdFlags.BoolVar(&this.rsyslog, "rsyslog", true, "")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	// must useradd haproxy before deploy
	if _, err := user.Lookup("haproxy"); err != nil {
		panic(err)
	}

	err := os.MkdirAll(this.root, 0755)
	swalllow(err)
	err = os.MkdirAll(fmt.Sprintf("%s/bin", this.root), 0755)
	swalllow(err)
	err = os.MkdirAll(fmt.Sprintf("%s/sbin", this.root), 0755)
	swalllow(err)
	err = os.MkdirAll(fmt.Sprintf("%s/logs", this.root), 0755)
	swalllow(err)
	err = os.MkdirAll(fmt.Sprintf("%s/src", this.root), 0755)
	swalllow(err)

	// install files
	b, _ := Asset("templates/haproxy-1.6.3.tar.gz")
	srcPath := fmt.Sprintf("%s/src/haproxy-1.6.3.tar.gz", this.root)
	err = ioutil.WriteFile(srcPath, b, 0644)
	swalllow(err)
	b, _ = Asset("templates/hatop-0.7.7.tar.gz")
	hatop := fmt.Sprintf("%s/src/hatop-0.7.7.tar.gz", this.root)
	err = ioutil.WriteFile(hatop, b, 0644)
	swalllow(err)
	b, _ = Asset("templates/init.ehaproxy")
	initPath := fmt.Sprintf("%s/src/init.ehaproxy", this.root)
	err = ioutil.WriteFile(initPath, b, 0755)
	swalllow(err)

	this.Ui.Info("will read zones from $HOME/.gafka.cf")
	this.Ui.Info("check /etc/security/limits.conf")
	this.Ui.Info(fmt.Sprintf("compile haproxy to %s/sbin: make TARGET=xxx USE_ZLIB=yes", this.root))
	this.Ui.Info(fmt.Sprintf("cp %s to /etc/init.d/ehaproxy", initPath))
	this.Ui.Info(fmt.Sprintf("chkconfig --add ehaproxy"))

	this.configKernal()

	if this.rsyslog {
		this.configRsyslog()
	}

	return
}

func (this *Deploy) configKernal() {
	this.Ui.Warn("net.core.somaxconn = 16384")
	this.Ui.Warn("net.core.netdev_max_backlog = 2500")
}

func (this *Deploy) configRsyslog() {
	this.Ui.Output("install and setup rsyslog for haproxy")
	this.Ui.Output(fmt.Sprintf(`
vim /etc/security/limits.conf
*          soft    nofile          409600
*          hard    nofile          409600

*          soft    nproc          65535
*          hard    nproc          65535

vim /etc/rsyslog.conf		
$ModLoad imudp
$UDPServerAddress 127.0.0.1
$UDPServerRun 514

vim  /etc/rsyslog.d/haproxy.conf
local3.*     /var/log/haproxy.log

vim /etc/sysconfig/rsyslog
SYSLOGD_OPTIONS=”-c 2 -r -m 0″
#-c 2 使用兼容模式，默认是 -c 5
#-r 开启远程日志
#-m 0 标记时间戳。单位是分钟，为0时，表示禁用该功能	
		`))
}

func (this *Deploy) Synopsis() string {
	return fmt.Sprintf("Deploy %s system on localhost", this.Cmd)
}

func (this *Deploy) Help() string {
	help := fmt.Sprintf(`
Usage: %s deploy [options]

    Deploy %s system on localhost

Options:

    -p prefix dir
      Defaults %s

    -rsyslog
      Display rsyslog integration with haproxy

`, this.Cmd, this.Cmd, defaultPrefix)
	return strings.TrimSpace(help)
}
