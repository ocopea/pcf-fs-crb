// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"crb/models"
	"crb/utils"
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func Test_validate(t *testing.T) {
	type args struct {
		repoinfo *models.RepositoryInfo
	}
	addr := "1.2.3.4"
	user := "user"
	ipv6Addr := "::1"
	dnsAddr := "localhost"
	ipv6AddrPort := "[::1]:8080"
	addrPort := "1.2.3.4:8080"
	dnsAddrPort := "localhost:8080"
	invalidAddrPort := "blah:8080"
	addrInvalidPort := "1.2.3.4:65536"
	ipv6AddrInvalidPort := "[::1]:65536"
	dnsAddrInvalidPort := "localhost:blah"
	invalidIpv6AddrFormat := "2001:0db8:85a3:0000:0000:8a2e:0370:7334:8080" // missing [] around ip literal before the port 8080
	validIpv6AddrFormat := "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080"
	emptyAddr := ""
	invalidAddr := "blah"
	emptyUser := ""
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "testValidIPv4Addr",
			args:    args{&models.RepositoryInfo{Addr: &addr, User: &user, Password: "Pass"}},
			wantErr: false,
		},
		{
			name:    "testValidIPv6Addr",
			args:    args{&models.RepositoryInfo{Addr: &ipv6Addr, User: &user, Password: "Pass"}},
			wantErr: false,
		},
		{
			name:    "testDnsAddr",
			args:    args{&models.RepositoryInfo{Addr: &dnsAddr, User: &user, Password: "Pass"}},
			wantErr: false,
		},
		{
			name:    "testDnsAddrPort",
			args:    args{&models.RepositoryInfo{Addr: &dnsAddrPort, User: &user, Password: "Pass"}},
			wantErr: false,
		},
		{
			name:    "testIPV6Port",
			args:    args{&models.RepositoryInfo{Addr: &ipv6AddrPort, User: &user, Password: "Pass"}},
			wantErr: false,
		},
		{
			name:    "testIPV4Port",
			args:    args{&models.RepositoryInfo{Addr: &addrPort, User: &user, Password: "Pass"}},
			wantErr: false,
		},
		{
			name:    "testEmptyAddr",
			args:    args{&models.RepositoryInfo{Addr: &emptyAddr, User: &user, Password: "Pass"}},
			wantErr: true,
		},
		{
			name:    "testInvalidAddr",
			args:    args{&models.RepositoryInfo{Addr: &invalidAddr, User: &user, Password: "Pass"}},
			wantErr: true,
		},
		{
			name:    "testEmptyUser",
			args:    args{&models.RepositoryInfo{Addr: &addr, User: &emptyUser, Password: "Pass"}},
			wantErr: true,
		},
		{
			name:    "testEmptyPassword",
			args:    args{&models.RepositoryInfo{Addr: &addr, User: &user, Password: ""}},
			wantErr: false,
		},

		{
			name:    "testInvalidAddrValidPort",
			args:    args{&models.RepositoryInfo{Addr: &invalidAddrPort, User: &user, Password: "Pass"}},
			wantErr: true,
		},
		{
			name:    "testIPV4InvalidPort",
			args:    args{&models.RepositoryInfo{Addr: &addrInvalidPort, User: &user, Password: "Pass"}},
			wantErr: true,
		},
		{
			name:    "testIPV6InvalidPort",
			args:    args{&models.RepositoryInfo{Addr: &ipv6AddrInvalidPort, User: &user, Password: "Pass"}},
			wantErr: true,
		},
		{
			name:    "testDnsAddrInvalidPort",
			args:    args{&models.RepositoryInfo{Addr: &dnsAddrInvalidPort, User: &user, Password: "Pass"}},
			wantErr: true,
		},
		{
			name:    "testInvalidIPV6Format",
			args:    args{&models.RepositoryInfo{Addr: &invalidIpv6AddrFormat, User: &user, Password: "Pass"}},
			wantErr: true,
		},
		{
			name:    "testValidInvalidIPV6Format",
			args:    args{&models.RepositoryInfo{Addr: &validIpv6AddrFormat, User: &user, Password: "Pass"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validate(tt.args.repoinfo); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostrepohandler(t *testing.T) {
	var mockDB *utils.Mockdb
	addr := "1.2.3.4"
	user := "user"

	type args struct {
		repoinfo *models.RepositoryInfo
	}

	tests := []struct {
		name         string
		args         args
		wantErr      bool
		expectations utils.Wants
	}{
		{name: "testPostSuccessful",
			args:         args{repoinfo: &models.RepositoryInfo{Addr: &addr, User: &user, Password: "Pass"}},
			wantErr:      false,
			expectations: utils.Wants{RepoTableCreate: nil, AddRepo: nil},
		},
		{name: "testPostAddRepoFails",
			args:         args{repoinfo: &models.RepositoryInfo{Addr: &addr, User: &user, Password: "Pass"}},
			wantErr:      true,
			expectations: utils.Wants{RepoTableCreate: nil, AddRepo: errors.New("Error while adding repo info")},
		},
		{name: "testPostTableCreateFails",
			args:         args{repoinfo: &models.RepositoryInfo{Addr: &addr, User: &user, Password: "Pass"}},
			wantErr:      true,
			expectations: utils.Wants{RepoTableCreate: errors.New("Error while creating table"), AddRepo: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB = utils.SetMockDbExpectations(&tt.expectations)
			if err := postrepohandler(tt.args.repoinfo, mockDB); (err != nil) != tt.wantErr {
				t.Errorf("postrepohandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getRepoInstance(t *testing.T) {
	tests := []struct {
		name        string
		postRequest string
		want        *models.RepositoryInstance
		expectFail  bool
	}{
		{
			name:        "TestPostRequestSuccess",
			postRequest: "http://127.0.0.1:8080/crb/repositories",
			want:        &models.RepositoryInstance{CopyRepoURL: "http://127.0.0.1:8080/crb/repositories"},
			expectFail:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", tt.postRequest, nil)
			if got := getRepoInstance(req); reflect.DeepEqual(got, tt.want) {
				if tt.expectFail {
					t.Errorf("getRepoInstance() got = %v, want %v", got, tt.want)
				}
			} else {
				if !tt.expectFail {
					t.Errorf("getRepoInstance() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
