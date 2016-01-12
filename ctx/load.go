package ctx

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	jsconf "github.com/funkygao/jsconf"
)

func LoadConfig(fn string) {
	cf, err := jsconf.Load(fn)
	if err != nil {
		panic(err)
	}

	conf = new(config)
	conf.hostname, _ = os.Hostname()
	conf.kafkaHome = cf.String("kafka_home", "")
	conf.logLevel = cf.String("loglevel", "info")
	conf.influxdbHost = cf.String("influxdb_host", "")
	conf.zones = make(map[string]string)
	conf.consulBootstrap = cf.String("consul_bootstrap", "")
	conf.zkDefaultZone = cf.String("zk_default_zone", "")
	conf.tunnels = make(map[string]string)
	conf.aliases = make(map[string]string)
	for i := 0; i < len(cf.List("aliases", nil)); i++ {
		section, err := cf.Section(fmt.Sprintf("aliases[%d]", i))
		if err != nil {
			panic(err)
		}

		conf.aliases[section.String("cmd", "")] = section.String("alias", "")
	}

	for i := 0; i < len(cf.List("zones", nil)); i++ {
		section, err := cf.Section(fmt.Sprintf("zones[%d]", i))
		if err != nil {
			panic(err)
		}

		z := new(zone)
		z.loadConfig(section)
		conf.zones[z.name] = z.zk
		conf.tunnels[z.name] = z.tunnel
	}

	conf.reverseDns = make(map[string][]string)
	for _, entry := range cf.StringList("reverse_dns", nil) {
		if entry != "" {
			// entry e,g. k11000b.sit.wdds.kfk.com:10.213.33.149
			parts := strings.SplitN(entry, ":", 2)
			if len(parts) != 2 {
				panic("invalid reverse_dns record")
			}

			ip, host := strings.TrimSpace(parts[1]), strings.TrimSpace(parts[0])
			if _, present := conf.reverseDns[ip]; !present {
				conf.reverseDns[ip] = make([]string, 0)
			}

			conf.reverseDns[ip] = append(conf.reverseDns[ip], host)
		}
	}

}

func LoadFromHome() {
	const defaultConfig = `
{
    zones: [
        {
            name: "sit"
            zk: "10.77.144.87:10181,10.77.144.88:10181,10.77.144.89:10181"
            tunnel: "gaopeng27@10.77.130.14"
        }
        {
            name: "test"
            zk: "10.77.144.101:10181,10.77.144.132:10181,10.77.144.182:10181"
            tunnel: "gaopeng27@10.77.130.14"
        }      
        {
            name: "pre"
            zk: "10.213.33.154:2181,10.213.42.48:2181,10.213.42.49:2181"
            tunnel: "gaopeng27@10.209.11.11"
        }  
        {
            name: "prod"
            zk: "10.209.33.69:2181,10.209.37.19:2181,10.209.37.68:2181"
            tunnel: "gaopeng27@10.209.11.11"
        }
    ]

    aliases: [
    	{
    		cmd: "zone"
    		alias: "zone"
    	}
    ]

	reverse_dns: [
		// pre zk
		"z2181a.sit.wdds.zk.com:10.213.33.154"
		"z2181b.sit.wdds.zk.com:10.213.42.48"
		"z2181c.sit.wdds.zk.com:10.213.42.49"

		// pre kafka brokers
		"k10101a.sit.wdds.kfk.com:10.213.33.148"
		"k10101b.sit.wdds.kfk.com:10.213.33.149"
		"k10102a.sit.wdds.kfk.com:10.213.33.148"
		"k10102b.sit.wdds.kfk.com:10.213.33.149"
		"k10103a.sit.wdds.kfk.com:10.213.33.148"
		"k10103b.sit.wdds.kfk.com:10.213.33.149"
		"k10104a.sit.wdds.kfk.com:10.213.33.148"
		"k10104b.sit.wdds.kfk.com:10.213.33.149"
		"k10105a.sit.wdds.kfk.com:10.213.33.148"
		"k10105b.sit.wdds.kfk.com:10.213.33.149"
		"k10106a.sit.wdds.kfk.com:10.213.33.148"
		"k10106b.sit.wdds.kfk.com:10.213.33.149"
		"k10107a.sit.wdds.kfk.com:10.213.33.148"
		"k10107b.sit.wdds.kfk.com:10.213.33.149"
		"k10108a.sit.wdds.kfk.com:10.213.33.148"
		"k10108b.sit.wdds.kfk.com:10.213.33.149"
		"k10109a.sit.wdds.kfk.com:10.213.33.148"
		"k10109b.sit.wdds.kfk.com:10.213.33.149"
		"k10110a.sit.wdds.kfk.com:10.213.33.148"
		"k10110b.sit.wdds.kfk.com:10.213.33.149"
		"k10111a.sit.wdds.kfk.com:10.213.33.148"
		"k10111b.sit.wdds.kfk.com:10.213.33.149"
		"k10112a.sit.wdds.kfk.com:10.213.33.148"
		"k10112b.sit.wdds.kfk.com:10.213.33.149"
		"k10113a.sit.wdds.kfk.com:10.213.33.148"
		"k10113b.sit.wdds.kfk.com:10.213.33.149"
		"k10114a.sit.wdds.kfk.com:10.213.33.148"
		"k10114b.sit.wdds.kfk.com:10.213.33.149"
		"k10115a.sit.wdds.kfk.com:10.213.33.148"
		"k10115b.sit.wdds.kfk.com:10.213.33.149"
		"k10116a.sit.wdds.kfk.com:10.213.33.148"
		"k10116b.sit.wdds.kfk.com:10.213.33.149"
		"k10117c.sit.wdds.kfk.com:10.213.33.150"
		"k10117d.sit.wdds.kfk.com:10.213.33.151"
		"k10118c.sit.wdds.kfk.com:10.213.33.150"
		"k10118d.sit.wdds.kfk.com:10.213.33.151"
		"k11000a.sit.wdds.kfk.com:10.213.33.148"
		"k11000b.sit.wdds.kfk.com:10.213.33.149"
		"k11001a.sit.wdds.kfk.com:10.213.33.148"
		"k11001b.sit.wdds.kfk.com:10.213.33.149"

		// prod zk
		"zk2181a.wdds.zk.com:10.209.33.69"
		"zk2181b.wdds.zk.com:10.209.37.19"
		"zk2181c.wdds.zk.com:10.209.37.68"

		// prod kafka brokers
		"k10101a.wdds.kfk.com:10.209.37.39"
		"k10101b.wdds.kfk.com:10.209.33.20"
		"k10102a.wdds.kfk.com:10.209.37.39"
		"k10102b.wdds.kfk.com:10.209.33.20"
		"k10103a.wdds.kfk.com:10.209.37.39"
		"k10103b.wdds.kfk.com:10.209.33.20"
		"k10104a.wdds.kfk.com:10.209.37.39"
		"k10104b.wdds.kfk.com:10.209.33.20"
		"k10105a.wdds.kfk.com:10.209.37.39"
		"k10105b.wdds.kfk.com:10.209.33.20"
		"k10106a.wdds.kfk.com:10.209.37.39"
		"k10106b.wdds.kfk.com:10.209.33.20"
		"k10107a.wdds.kfk.com:10.209.37.39"
		"k10107b.wdds.kfk.com:10.209.33.20"
		"k10108a.wdds.kfk.com:10.209.37.39"
		"k10108b.wdds.kfk.com:10.209.33.20"
		"k10109a.wdds.kfk.com:10.209.37.39"
		"k10109b.wdds.kfk.com:10.209.33.20"
		"k10110a.wdds.kfk.com:10.209.37.39"
		"k10110b.wdds.kfk.com:10.209.33.20"
		"k10111a.wdds.kfk.com:10.209.37.39"
		"k10111b.wdds.kfk.com:10.209.33.20"
		"k10112a.wdds.kfk.com:10.209.37.39"
		"k10112b.wdds.kfk.com:10.209.33.20"
		"k10113a.wdds.kfk.com:10.209.37.69"
		"k10113b.wdds.kfk.com:10.209.33.40"
		"k10114a.wdds.kfk.com:10.209.37.69"
		"k10114b.wdds.kfk.com:10.209.33.40"
		"k10115a.wdds.kfk.com:10.209.37.69"
		"k10115b.wdds.kfk.com:10.209.33.40"
		"k10116a.wdds.kfk.com:10.209.37.69"
		"k10116b.wdds.kfk.com:10.209.33.40"
		"k10117a.wdds.kfk.com:10.209.37.69"
		"k10117b.wdds.kfk.com:10.209.33.40"
		"k10118a.wdds.kfk.com:10.209.37.69"
		"k10118b.wdds.kfk.com:10.209.33.40"
		"k10119a.wdds.kfk.com:10.209.37.69"
		"k10119b.wdds.kfk.com:10.209.33.40"
		"k10120a.wdds.kfk.com:10.209.37.69"
		"k10120b.wdds.kfk.com:10.209.33.40"
		"k10121a.wdds.kfk.com:10.209.37.69"
		"k10121b.wdds.kfk.com:10.209.33.40"
		"k10122a.wdds.kfk.com:10.209.37.69"
		"k10122b.wdds.kfk.com:10.209.33.40"
		"k11000a.wdds.kfk.com:10.209.37.39"
		"k11000b.wdds.kfk.com:10.209.33.20"
		"k11001a.wdds.kfk.com:10.209.37.69"
		"k11001b.wdds.kfk.com:10.209.33.40"
		"k10120a.wdds.kfk.com:10.209.10.161"
		"k10120b.wdds.kfk.com:10.209.10.141"
		"k10121a.wdds.kfk.com:10.209.10.161"
		"k10121b.wdds.kfk.com:10.209.10.141"
		"k10118a.wdds.kfk.com:10.209.11.166"
		"k10118b.wdds.kfk.com:10.209.11.195"
	]

    zk_default_zone: "prod"
    consul_bootstrap: "10.209.33.69:8500"
    kafka_home: "/opt/kafka_2.10-0.8.1.1"
    loglevel: "info"
}
`
	var configFile string
	if usr, err := user.Current(); err == nil {
		configFile = filepath.Join(usr.HomeDir, ".gafka.cf")
	} else {
		panic(err)
	}

	_, err := os.Stat(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			// create the config file on the fly
			if e := ioutil.WriteFile(configFile,
				[]byte(strings.TrimSpace(defaultConfig)), 0644); e != nil {
				panic(e)
			}
		} else {
			panic(err)
		}
	}

	LoadConfig(configFile)
}
