// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package utils

import (
	"errors"
	"fmt"
	"io"
	"net"

	"strings"
	"unicode/utf8"

	"crb/models"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SftpClient holds sftp connection and ssh connection to the copy repo.
type SftpClient struct {
	sftpCon *sftp.Client
	sshCon  *ssh.Client
}

// These errors are defined in the net package but are not public,
// exposing here for use in error checking
const MissingPort = "missing port in address"
const TooManyColons = "too many colons in address"

// RemotePath defines the path being used on remote copy repo file server
const RemotePath = "/var/lib/crb/"

// OpenSftpConnection opens both ssh connection and sftp connection to the copy repo.
func OpenSftpConnection(repoInfo *models.RepositoryInfo) (*SftpClient, error) {
	const defaultPort = "22"

	var sftpcon SftpClient
	var err error

	config := ssh.ClientConfig{
		User: *(repoInfo.User),
		Auth: []ssh.AuthMethod{
			ssh.Password(repoInfo.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	host, port, err := net.SplitHostPort(*repoInfo.Addr)
	if err != nil {
		// a host with no port will return an error "address <host>: missing port in address"
		if strings.Contains(err.Error(), MissingPort) {
			host = *repoInfo.Addr
			port = defaultPort
		} else if strings.Contains(err.Error(), TooManyColons) {
			tmp := net.ParseIP(*repoInfo.Addr) //if too many colon error but ParseIP returns a value it's
			if tmp != nil {                    //an ipv6 address with no port.  SplitHostPort works on [::1]:8080 format
				host = *repoInfo.Addr
				port = defaultPort
			}
		} else {
			return nil, err
		}
	}

	addr := net.JoinHostPort(host, port)
	sftpcon.sshCon, err = ssh.Dial("tcp", addr, &config)
	if err != nil {
		return nil, err
	}

	sftpcon.sftpCon, err = sftp.NewClient(sftpcon.sshCon)
	return &sftpcon, err
}

// CloseSftpConnection closes both ssh connection and sftp connection to the copy repo.
// It should be called for clean up from the same place where OpenSftpConnection is called.
func CloseSftpConnection(sftpClient *SftpClient) error {
	if sftpClient == nil {
		return nil
	}

	if sftpClient.sftpCon != nil {
		err := sftpClient.sftpCon.Close()
		if err != nil {
			return err
		}
	}

	if sftpClient.sshCon != nil {
		err := sftpClient.sshCon.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// ValidateCopyID check is copyID is a valid linux file name
// Ref: http://www.linfo.org/file_name.html
func ValidateCopyID(copyID string) error {
	if idLen := utf8.RuneCountInString(copyID); idLen == 0 || idLen > 255 {
		return errors.New("copyID does not meet linux file name length")
	}

	if strings.Contains(copyID, "/") || strings.Index(copyID, "-") == 0 {
		return errors.New("copyID does not meet linux file naming rules")
	}
	return nil
}

// StoreCopy stores copy data on the copy repo at given copyID
// Returns number of bytes that have been copied i.e. size of the file
func (sftpClient *SftpClient) StoreCopy(copyID string, copyData io.ReadCloser) (int64, error) {
	defer copyData.Close()
	cmd := fmt.Sprintf("mkdir -p %s", RemotePath)
	session, err := sftpClient.sshCon.NewSession()
	if err != nil {
		return 0, err
	}
	defer session.Close()

	err = session.Run(cmd)
	if err != nil {
		return 0, err
	}
	f, err := sftpClient.sftpCon.Create(RemotePath + copyID)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	size, err := io.Copy(f, copyData)

	if err != nil {
		return 0, err
	}
	return size, nil
}

// RetrieveCopy gets copy file handler for the given copyID from the copy repo
// The returned file handler needs to be closed in the caller of this function
func (sftpClient *SftpClient) RetrieveCopy(copyID string) (io.ReadCloser, error) {
	file, error := sftpClient.sftpCon.Open(RemotePath + copyID)
	if error != nil {
		return nil, error
	}

	return file, error
}

// DeleteCopy deletes the copy on the copy repo
func (sftpClient *SftpClient) DeleteCopy(copyID string) error {
	err := sftpClient.sftpCon.Remove(RemotePath + copyID)
	if err != nil {
		return err
	}

	return nil
}
