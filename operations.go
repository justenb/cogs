package main

import (
	"errors"
	"fmt"
)

type result struct {
	server  string
	output  []string
	err     error
	success bool
}

func opsExecCmd(cmd string, env *Enviroment) {
	// SSH command execution
	results := make(chan result, len(env.hosts))

	for _, host := range env.hosts {
		go func(host *SshConfig) {

			if !env.shell {
				/* If the -shell cli argument is set to false then cogs
				   will not wrap the command as such -> bash -l -c "command"' */
				data, err := host.Run(WithBash(cmd))
				if err != nil {
					results <- result{server: host.Name, output: []string{""}, err: err}
				}

				results <- result{server: host.Name, output: ProcessOutput(data, host, env), err: nil}

			} else {
				data, err := host.Run(cmd)
				if err != nil {
					results <- result{server: host.Name, output: []string{""}, err: err}
				}

				results <- result{server: host.Name, output: ProcessOutput(data, host, env), err: nil}
			}

		}(host)
	}

	errorCollect := make([]string, 0)

	for i := 0; i < len(env.hosts); i++ {
		select {

		case data := <-results:
			if data.err != nil {
				errorCollect = append(errorCollect, fmt.Sprintf("%s: %q", data.server, data.err))
			}

			for _, processed := range data.output {
				fmt.Println(processed)
			}
		}
	}

	if len(errorCollect) > 0 {
		for _, e := range errorCollect {
			fmt.Println(e)
		}
	}
}

func opsGet(rPath string, env *Enviroment) {
	// SFTP Get Operation
	results := make(chan result, len(env.hosts))

	for _, host := range env.hosts {
		go func(host *SshConfig) {
			err := host.Get(rPath, env)

			if err != nil {
				results <- result{success: false, err: err, server: host.Name}

			} else {
				results <- result{success: true}
			}

		}(host)
	}

	errorCollect := make([]string, 0)

	for i := 0; i < len(env.hosts); i++ {
		select {
		case done := <-results:
			if !done.success {
				errorCollect = append(errorCollect, fmt.Sprintf("%s: %q", done.server, done.err))
			}
		}
	}

	if len(errorCollect) > 0 {
		for _, e := range errorCollect {
			fmt.Println(e)
		}
	}
}

func opsPut(lPath, rPath string, env *Enviroment) {
	// SFTP Put Operation
	results := make(chan result, len(env.hosts))
	errorCollect := make([]string, 0)

	if rPath == "" {
		errorCollect = append(errorCollect, fmt.Sprintf("%s: %q", "Put error", errors.New("remote path not specified or not valid.")))
	}

	for _, host := range env.hosts {
		go func(host *SshConfig) {
			err := host.Put(lPath, rPath)

			if err != nil {
				results <- result{success: false, err: err, server: host.Name}

			} else {
				results <- result{success: true}
			}

		}(host)
	}

	for i := 0; i < len(env.hosts); i++ {
		select {

		case done := <-results:
			if !done.success {
				errorCollect = append(errorCollect, fmt.Sprintf("%s: %q", done.server, done.err))
			}
		}
	}

	if len(errorCollect) > 0 {
		for _, e := range errorCollect {
			fmt.Println(e)
		}
	}

}
