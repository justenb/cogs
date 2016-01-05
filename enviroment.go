package main

import (
	"runtime"
)

type Enviroment struct {
	hosts  []*SshConfig
	filter Filter
	pubkey string
	shell  bool
}

type Filter struct {
	on  bool
	key string
}

func (env *Enviroment) SetEnv(csvFile, hosts, filter, pubkey, userName, port string, maxprocs int, shell bool) {
	if maxprocs > runtime.NumCPU() {
		runtime.GOMAXPROCS(0)
	}

	if !Empty(csvFile) {
		env.hosts = ImportCsvFile(csvFile)

	} else if !Empty(hosts) {
		env.hosts = ImportPubKeyConfig(hosts, pubkey, userName, port)
	}

	if !Empty(filter) {
		env.filter.key, env.filter.on = filter, true
	}

	if !shell {
		env.shell = false

	} else {
		env.shell = true
	}
}
