package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	hosts    string
	cmd      string
	get      string
	put      string
	rpath    string
	filter   string
	key      string
	csvFile  string
	username string
	port     string
	maxprocs int
	shell    bool

	usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
)

func init() {
	/*
	 *       cli args
	 */

	// operations that can be performed on remote host(s).
	flag.StringVar(&cmd, "cmd", "", "execute a command on remote a host(s).")
	flag.StringVar(&cmd, "c", "", "shorthand for --cmd")
	flag.StringVar(&get, "get", "", "use sftp to retrieve a remote file from 1 or more hosts.")
	flag.StringVar(&get, "g", "", "shorthand for --get")
	flag.StringVar(&put, "put", "", "use sftp to upload a local file to a remote location on 1 or more hosts.")
	flag.StringVar(&rpath, "rpath", "", "the target remote path for use with the --put option.")

	// options for importing remote host(s) configuration.
	flag.StringVar(&csvFile, "csv", "", "import a list of servers contained within a csv file. acceptable format: server,user,password,port.")
	flag.StringVar(&hosts, "hosts", "", "used to import hosts declared with the cli itself contrary to importing from a file. Use this option with a public key.")
	flag.StringVar(&hosts, "h", "", "shorthand for --hosts")
	flag.StringVar(&username, "user", "", "user associated with the server(s) declared using the --hosts option.")
	flag.StringVar(&username, "u", "", "shorthand for --user")
	flag.StringVar(&port, "port", "22", "port associated with the server(s) declared using the --hosts option.")
	flag.StringVar(&port, "p", "", "shorthand for --port")

	// import public key config
	flag.StringVar(&key, "key", "", "private key --path")
	flag.StringVar(&key, "k", "", "shorthand for --key")

	// misc options
	flag.StringVar(&filter, "filter", "", "Grep like option for filtering command output and downloaded text based files. Think log filtering.")
	flag.StringVar(&filter, "f", "", "shorthand for --filter")
	flag.IntVar(&maxprocs, "maxprocs", 0, "Max num of OS threads used to execute the operation.")
	flag.IntVar(&maxprocs, "max", 1, "shorthand for --maxprocs")
	flag.BoolVar(&shell, "shell", true, "Executes Bash shell in login mode")

	flag.Parse()
}

func main() {
	env := &Enviroment{}
	env.SetEnv(csvFile, hosts, filter, key, username, port, maxprocs, shell)

	if hosts != "" && csvFile != "" {
		usage()
	}

	switch {

	case cmd != "":
		opsExecCmd(cmd, env)

	case put != "":
		opsPut(put, rpath, env)

	case get != "":
		opsGet(get, env)

	default:
		usage()
	}
}
