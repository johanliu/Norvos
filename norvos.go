package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/johanliu/essos"
	"github.com/johanliu/essos/components"
)

const (
	VERSION = "1.0"

	ADD = iota
	DELETE
	UPDATE
	READ

	API_CONTENT_HEADER  = "application/json;charset=utf-8"
	ETCD_CONTENT_HEADER = "application/x-www-form-urlencoded"

	ADD_METHOD     = "PUT"
	DELETE_METHOD  = "DELETE"
	CONTENT_HEADER = "Content-Type"

	NORVOS_PATH             = "/norvos"
	DEFAULT_ETCD_NORVOS_LOC = "/v2/keys/norvos/node/purge/"
	DEFAULT_ETCD_PORT       = "2379"
	DEFAULT_ETCD_TIMEOUT    = 5
	DEFAULT_SKYDNS_LOC      = "/v2/keys/skydns/"
	DEFAULT_ARPA_LOC        = "/v2/keys/skydns/arpa/in-addr/"
	DEFAULT_SCAN_ARGS       = "/?recursive=true"
	DEFAULT_ETCD_CAS        = "?prevValue="
	DEFAULT_TTL             = 60

	DEFAULT_TRIM_KEY        = "/skydns/"
	DEFAULT_TRIM_ARPA_KEY   = DEFAULT_TRIM_KEY + "arpa/in-addr/"
	DEFAULT_TRIM_SERVER_KEY = "/norvos/node/purge/"

	DEFAULT_TIMEOUT = 5
)

var Entry norvos

type norvos struct {
	ops map[string]essos.Operation
}

func (n *norvos) Discover() map[string]essos.Operation {
	return n.ops
}

type create string

func (create) Description() string {
	return "This is Create"
}

func (create) Do(ctx context.Context, args []string) (context.Context, error) {
	if ctx == nil {
		return ctx, errors.New("context is nil")
	}

	fmt.Print("Call successful")

	return ctx, nil
}

type delete string

func (delete) Description() string {
	return "This is Delete"
}

func (delete) Do(ctx context.Context, args []string) (context.Context, error) {
	if ctx == nil {
		return ctx, errors.New("context is nil")
	}

	fmt.Print("Call successful")

	return ctx, nil
}

type update string

func (update) Description() string {
	return "This is Update"
}

func (update) Do(ctx context.Context, args []string) (context.Context, error) {
	if ctx == nil {
		return ctx, errors.New("context is nil")
	}

	fmt.Print("Call successful")

	return ctx, nil
}

type read string

func (read) Description() string {
	return "This is Read"
}

func (read) Do(ctx context.Context, args []string) (context.Context, error) {
	if ctx == nil {
		return ctx, errors.New("context is nil")
	}

	fmt.Print("Call successful")

	return ctx, nil
}

func do(args []string, data interface{}) error {
	return nil
}

func init() {
	components.Add("dns", func() essos.Component {
		return &norvos{
			ops: map[string]essos.Operation{
				"create": create("Create"),
				"read":   read("Read"),
				"delete": delete("Delete"),
				"update": update("Update"),
			},
		}
	})
}
