package command

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/funkygao/gafka/ctx"
	"github.com/funkygao/gocli"
	"github.com/go-ozzo/ozzo-dbx"
	_ "github.com/go-sql-driver/mysql"
	"github.com/olekukonko/tablewriter"
)

type WhoisAppInfo struct {
	AppId            string `db:"AppId"`
	ApplicationName  string `db:"ApplicationName"`
	ApplicationIntro string `db:"ApplicationIntro"`
	Cluster          string `db:"Cluster"`
	CreateBy         string `db:"CreateBy"`
	CreateTime       string `db:"CreateTime"`
	Status           string `db:"Status"`
}

type WhoisTopicInfo struct {
	AppId      string `db:"AppId"`
	AppName    string
	TopicName  string `db:"TopicName"`
	TopicIntro string `db:"TopicIntro"`
	CreateBy   string `db:"CreateBy"`
	CreateTime string `db:"CreateTime"`
	Status     string `db:"Status"`
}

type WhoisGroupInfo struct {
	AppId      string `db:"AppId"`
	AppName    string
	GroupName  string `db:"GroupName"`
	GroupIntro string `db:"GroupIntro"`
	CreateBy   string `db:"CreateBy"`
	CreateTime string `db:"CreateTime"`
	Status     string `db:"Status"`
}

type Whois struct {
	Ui  cli.Ui
	Cmd string

	zone     string
	app      string
	topic    string
	group    string
	likeMode bool

	appInfos   []WhoisAppInfo
	topicInfos []WhoisTopicInfo
	groupInfos []WhoisGroupInfo
}

func (this *Whois) Run(args []string) (exitCode int) {
	cmdFlags := flag.NewFlagSet("whois", flag.ContinueOnError)
	cmdFlags.Usage = func() { this.Ui.Output(this.Help()) }
	cmdFlags.StringVar(&this.zone, "z", ctx.ZkDefaultZone(), "")
	cmdFlags.StringVar(&this.app, "app", "", "")
	cmdFlags.StringVar(&this.group, "g", "", "")
	cmdFlags.StringVar(&this.topic, "t", "", "")
	cmdFlags.BoolVar(&this.likeMode, "l", false, "")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	ensureZoneValid(this.zone)

	mysqlDsns := map[string]string{
		"prod": "user_pubsub:p0nI7mEL6OLW@tcp(m3342.wdds.mysqldb.com:3342)/pubsub?charset=utf8&timeout=10s",
		"sit":  "pubsub:pubsub@tcp(10.209.44.12:10043)/pubsub?charset=utf8&timeout=10s",
		"test": "pubsub:pubsub@tcp(10.209.44.14:10044)/pubsub?charset=utf8&timeout=10s",
	}
	this.loadFromManager(mysqlDsns[this.zone])

	table := tablewriter.NewWriter(os.Stdout)
	switch {
	case this.topic+this.group == "":
		// list apps
		table.SetHeader([]string{"Id", "Name", "Cluster", "Ctime"})
		for _, ai := range this.appInfos {
			table.Append([]string{ai.AppId, ai.ApplicationName, ai.Cluster, ai.CreateTime})
		}

	case this.topic != "":
		table.SetHeader([]string{"topic", "desc", "app", "owner", "ctime"})
		for _, ti := range this.topicInfos {
			table.Append([]string{ti.TopicName, ti.TopicIntro, ti.AppName, ti.CreateBy, ti.CreateTime})
		}

	case this.group != "":
		table.SetHeader([]string{"group", "desc", "app", "owner", "ctime"})
		for _, gi := range this.groupInfos {
			table.Append([]string{gi.GroupName, gi.GroupIntro, gi.AppName, gi.CreateBy, gi.CreateTime})
		}

	case this.topic != "" && this.group != "":
		this.Ui.Error("-t and -g cannot be set at the same time！")
		return 1
	}

	table.Render()

	return
}

func (this *Whois) loadFromManager(dsn string) {
	db, err := dbx.Open("mysql", dsn)
	swallow(err)

	// TODO fetch from topics_version
	sql := "SELECT AppId,ApplicationName,ApplicationIntro,Cluster,CreateBy,CreateTime,Status FROM application"
	if this.app != "" {
		sql += " WHERE AppId IN (" + this.app + ")"
	}
	sql += " ORDER BY AppId"
	q := db.NewQuery(sql)

	swallow(q.All(&this.appInfos))

	op := "="
	if this.likeMode {
		op = "LIKE"
	}
	if this.topic != "" {
		if this.likeMode {
			this.topic = "%" + this.topic + "%"
		}
		sql = fmt.Sprintf("SELECT AppId,TopicName,TopicIntro,CreateBy,CreateTime,Status FROM topics WHERE TopicName %s '%s'",
			op, this.topic)
		q = db.NewQuery(sql)
		swallow(q.All(&this.topicInfos))

		for i, ti := range this.topicInfos {
			this.topicInfos[i].AppName = this.appName(ti.AppId)
		}
	}

	if this.group != "" {
		if this.likeMode {
			this.group = "%" + this.group + "%"
		}
		sql = fmt.Sprintf("SELECT AppId,GroupName,GroupIntro,CreateBy,CreateTime,Status FROM application_group WHERE GroupName %s '%s'",
			op, this.group)
		q = db.NewQuery(sql)
		swallow(q.All(&this.groupInfos))
		for i, gi := range this.groupInfos {
			this.groupInfos[i].AppName = this.appName(gi.AppId)
		}
	}
}

func (this *Whois) appName(appId string) string {
	for _, ai := range this.appInfos {
		if ai.AppId == appId {
			return ai.ApplicationName
		}
	}

	return "NotFound"
}

func (*Whois) Synopsis() string {
	return "Lookup PubSub App Information"
}

func (this *Whois) Help() string {
	help := fmt.Sprintf(`
Usage: %s whois [options]

    Lookup PubSub App Information

Options:

    -z zone

    -app comma seperated appId

    -g group

    -t topic

    -l
      Like mode. 
      Pattern wildcard match of group or topic name.

`, this.Cmd)
	return strings.TrimSpace(help)
}