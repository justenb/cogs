package main

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io/ioutil"
	"net"
	"os"
	"os/user"
)

type SshConfig struct {
	User     string
	Name     string
	Key      string
	Port     string
	Password string
}

func getKeyFile(keypath string) (ssh.Signer, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	file := usr.HomeDir + keypath
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	pubkey, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}

	return pubkey, nil
}

func (conf *SshConfig) Session() (*ssh.Session, error) {
	client, err := conf.Client()
	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (conf *SshConfig) Client() (*ssh.Client, error) {

	auths := []ssh.AuthMethod{}

	if !Empty(conf.Password) {
		auths = append(auths, ssh.Password(conf.Password))
	}

	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers))
		defer sshAgent.Close()
	}

	if pubkey, err := getKeyFile(conf.Key); err == nil {
		auths = append(auths, ssh.PublicKeys(pubkey))
	}

	config := &ssh.ClientConfig{
		User: conf.User,
		Auth: auths,
	}

	client, err := ssh.Dial("tcp", conf.Name+":"+conf.Port, config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (conf *SshConfig) Run(command string) (string, error) {
	session, err := conf.Session()

	if err != nil {
		return "", err
	}

	defer session.Close()

	var out bytes.Buffer
	session.Stdout = &out

	err = session.Run(command)
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
