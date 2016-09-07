package main

import (
	"runtime"
)

type Enviroment struct {
	hosts  []*SshConfig
	pubkey string
	shell  bool
}

func (env *Enviroment) SetEnv( csvFile string,  hosts string,  pubkey string, userName string, port string, maxprocs int, shell bool){
	if maxprocs > runtime.NumCPU() {
		runtime.GOMAXPROCS(0)
	}

    // Import host(s) credentials from a csvfile or from inline arguments.
	if !Empty(csvFile) {
		env.hosts = ImportCsvFile(csvFile)
	} else if !Empty(hosts) {
		env.hosts = ImportPubKeyConfig(hosts, pubkey, userName, port)
	}

	if !shell {
		env.shell = false
	} else {
		env.shell = true
	}
}
