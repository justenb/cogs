package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/ryanuber/go-glob"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (conn *SshConfig) transport() (*sftp.Client, error) {
	c, err := conn.Client()
	if err != nil {
		return nil, err
	}

	sftp, err := sftp.NewClient(c)
	if err != nil {
		return nil, err
	}

	return sftp, nil
}

func (conn *SshConfig) Get(rPath string, env *Enviroment) error {
	sftp, err := conn.transport()
	if err != nil {
		return err
	}

	defer sftp.Close()

	// handler for *glob in rpath
	if GlobIn(PathChomp(rPath)) {
		err := conn.GlobHandler(rPath, sftp, env)
		if err != nil {
			return err
		}

		return nil
	}

	remoteFile, err := sftp.Open(rPath)
	if err != nil {
		return err
	}

	fStat, err := remoteFile.Stat()
	if err != nil {
		return err
	}

	if _, err := os.Stat(conn.Name); os.IsNotExist(err) {
		os.Mkdir(conn.Name, 0755)
	}

	switch {

	case (!IsRegularFile(fStat)):
		w := sftp.Walk(rPath)
		for w.Step() {
			if w.Err() != nil {
				continue
			}

			if IsRegularFile(w.Stat()) {
				err := conn.GetFile(w.Path(), PrependParent(conn.Name, w.Path()), sftp, env)
				if err != nil {
					return err
				}

			} else {
				_, err := os.Stat(PrependParent(conn.Name, w.Path()))

				if os.IsNotExist(err) {
					os.Mkdir(PrependParent(conn.Name, w.Path()), 0755)
				}
			}
		}

	default:
		conn.GetFile(rPath, filepath.Join(conn.Name, PathChomp(rPath)), sftp, env)
	}

	return nil
}

func (conn *SshConfig) GetFile(rPath, lPath string, sftp *sftp.Client, env *Enviroment) error {
	remoteFile, err := sftp.Open(rPath)
	if err != nil {
		return err
	}

	localFile, err := os.OpenFile(lPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer localFile.Close()

	if env.filter.on {
		// filter lines in the file that match the filter key
		fmt.Printf("downloading and applying filter to %s:%s\n", conn.Name, rPath)

		scanner := bufio.NewScanner(remoteFile)
		for scanner.Scan() {
			if LineMatch(scanner.Bytes(), env.filter.key) {

				if _, err = localFile.Write([]byte(scanner.Text() + "\n")); err != nil {
					return err
				}

			}
		}

		srcStat, err := localFile.Stat()
		if err != nil {
			return err
		}

		if (srcStat.Size()) == 0 {
			localFile.Write([]byte("No matches found."))
		}

		return nil
	}

	srcStat, err := remoteFile.Stat()
	if err != nil {
		return err
	}

	writer := io.Writer(localFile)

	io.Copy(writer, remoteFile)
	fmt.Printf("%s @filesize %dB <-- %s\n", lPath, srcStat.Size(), PathChomp(rPath))

	localFile.Sync()

	return nil
}

func (conn *SshConfig) Put(local, remote string) error {
	sftp, err := conn.transport()
	if err != nil {
		return err
	}

	defer sftp.Close()

	src, err := ioutil.ReadFile(local)
	if err != nil {
		return err
	}

	dst, err := sftp.Create(remote)
	if err != nil {
		return errors.New("Could not create the remote file.")
	}

	fmt.Printf("%s --> %s:%s\n", local, conn.Name, remote)
	if _, err := dst.Write(src); err != nil {
		return err
	}

	return nil
}

func (conn *SshConfig) GlobHandler(rPath string, sftp *sftp.Client, env *Enviroment) error {
	rPaths := make([]string, 0)

	files, err := sftp.ReadDir(PathCutFile(rPath))
	if err != nil {
		return err
	}

	for _, file := range files {
		if glob.Glob(PathChomp(rPath), file.Name()) {

			if IsRegularFile(file) {
				rPaths = append(rPaths, JoinPath(rPath, file.Name()))
			}
		}
	}

	if len(rPaths) > 0 {
		if _, err := os.Stat(conn.Name); os.IsNotExist(err) {
			os.Mkdir(conn.Name, 0755)
		}

		for _, rp := range rPaths {
			conn.GetFile(rp, filepath.Join(conn.Name, PathChomp(rp)), sftp, env)
		}
	}

	return nil
}
