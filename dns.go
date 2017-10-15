package main

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/johanliu/essos"
	"github.com/johanliu/essos/components"
	"github.com/johanliu/essos/interfaces"
	"github.com/johanliu/vidar"
)

const (
	version    = "1.0"
	root       = "/coredns/"
	defaultTTL = 300
)

var cli client.Client

type norvos struct {
	ops map[string]essos.Operation
}

type Message struct {
	result string
}

func (n *norvos) Discover() map[string]essos.Operation {
	return n.ops
}

func (n *norvos) Start(params interface{}) error {
	var err error
	p, ok := params.(interfaces.DNS)
	if !ok {
		return errors.New("Params format is mismatched with pre-defined one")
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	cfg := client.Config{
		Endpoints:               []string{p.Etcd},
		Transport:               transport,
		HeaderTimeoutPerRequest: time.Second,
	}

	cli, err = client.New(cfg)
	if err != nil {
		return err
	}
	return nil
}

func (n *norvos) Stop() error {
	return nil
}

func init() {
	components.Add("dns",
		&norvos{
			ops: map[string]essos.Operation{
				"create": create("Create"),
				"read":   read("Read"),
				"update": update("Update"),
				"delete": delete("Delete"),
			},
		})
}

// For test only
func main() {}

func reverse(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

type create string

func (create) Description() string {
	return "Create operation for component of dns"
}

func (create) Do(ctx context.Context, args []string) (context.Context, error) {
	type payload struct {
		Host        string `json:"host"`
		Port        int    `json:"port"`
		Priority    int    `json:"priority"`
		Weight      int    `json:"weight"`
		TTL         int    `json:"ttl"`
		TargetStrip int    `json:"targetstrip"`
	}

	type request struct {
		Name string
		payload
	}

	type response struct {
		Message *client.Response `json:"message"`
	}

	if ctx == nil {
		return ctx, errors.New("context is nil")
	}

	input := ctx.Value("input").([]byte)
	req := new(request)
	if err := json.Unmarshal(input, req); err != nil {
		return nil, err
	}

	if req.TTL == 0 {
		req.TTL = defaultTTL
	}

	api := client.NewKeysAPI(cli)

	name := strings.Split(req.Name, ".")
	reverse(name)
	req.Name = strings.Join(name, "/")

	key := root + req.Name
	value, err := json.Marshal(req.payload)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := api.Create(ctx, key, string(value))
	if err != nil {
		ce, ok := err.(client.Error)
		if ok {
			if ce.Code == client.ErrorCodeNodeExist {
				return nil, vidar.NewHTTPError(http.StatusAlreadyReported, ce.Message)
			} else {
				return nil, vidar.NewHTTPError(http.StatusBadGateway, ce.Message)
			}
		} else if err == context.DeadlineExceeded {
			return nil, vidar.GatewayTimeoutError
		} else {
			return nil, err
		}
	}

	result := response{
		Message: resp,
	}

	ctx = context.WithValue(ctx, "result", result)
	return ctx, nil
}

type delete string

func (delete) Description() string {
	return "Delete operation for component of dns"
}

func (delete) Do(ctx context.Context, args []string) (context.Context, error) {
	type request struct {
		Name string
	}

	type response struct {
		Message *client.Response `json:"message"`
	}

	if ctx == nil {
		return ctx, errors.New("context is nil")
	}

	input := ctx.Value("input").([]byte)
	req := new(request)
	if err := json.Unmarshal(input, req); err != nil {
		return nil, err
	}

	api := client.NewKeysAPI(cli)

	name := strings.Split(req.Name, ".")
	reverse(name)
	req.Name = strings.Join(name, "/")

	key := root + req.Name

	resp, err := api.Delete(context.Background(), key, &client.DeleteOptions{Recursive: false})
	if err != nil {
		ce, ok := err.(client.Error)
		if ok {
			if ce.Code == client.ErrorCodeKeyNotFound {
				return nil, vidar.NewHTTPError(http.StatusNoContent, ce.Message)
			} else {
				return nil, vidar.NewHTTPError(http.StatusBadGateway, ce.Message)
			}
		} else if err == context.DeadlineExceeded {
			return nil, vidar.GatewayTimeoutError
		} else {
			return nil, err
		}
	}

	result := response{
		Message: resp,
	}

	ctx = context.WithValue(ctx, "result", result)
	return ctx, nil
}

type update string

func (update) Description() string {
	return "Update operation for component of dns. Sorry, no Upsert"
}

func (update) Do(ctx context.Context, args []string) (context.Context, error) {
	type payload struct {
		Host        string `json:"host"`
		Port        int    `json:"port"`
		Priority    int    `json:"priority"`
		Weight      int    `json:"weight"`
		TTL         int    `json:"ttl"`
		TargetStrip int    `json:"targetstrip"`
	}

	type request struct {
		Name string
		payload
	}

	type response struct {
		Message *client.Response `json:"message"`
	}

	if ctx == nil {
		return ctx, errors.New("context is nil")
	}

	input := ctx.Value("input").([]byte)
	req := new(request)
	if err := json.Unmarshal(input, req); err != nil {
		return nil, err
	}

	if req.TTL == 0 {
		req.TTL = defaultTTL
	}

	api := client.NewKeysAPI(cli)

	name := strings.Split(req.Name, ".")
	reverse(name)
	req.Name = strings.Join(name, "/")

	key := root + req.Name
	value, err := json.Marshal(req.payload)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := api.Update(ctx, key, string(value))
	if err != nil {
		ce, ok := err.(client.Error)
		if ok {
			if ce.Code == client.ErrorCodeKeyNotFound {
				return nil, vidar.NewHTTPError(http.StatusNoContent, ce.Message)
			} else {
				return nil, vidar.NewHTTPError(http.StatusBadGateway, ce.Message)
			}
		} else if err == context.DeadlineExceeded {
			return nil, vidar.GatewayTimeoutError
		} else {
			return nil, err
		}
	}

	result := response{
		Message: resp,
	}

	ctx = context.WithValue(ctx, "result", result)
	return ctx, nil
}

type read string

func (read) Description() string {
	return "Read operation for component of dns"
}

func (read) Do(ctx context.Context, args []string) (context.Context, error) {
	type request struct {
		Name string
	}

	type response struct {
		Message *client.Response `json:"message"`
	}

	if ctx == nil {
		return ctx, errors.New("context is nil")
	}

	input := ctx.Value("input").([]byte)
	req := new(request)
	if err := json.Unmarshal(input, req); err != nil {
		return nil, err
	}

	api := client.NewKeysAPI(cli)

	name := strings.Split(req.Name, ".")
	reverse(name)
	req.Name = strings.Join(name, "/")

	key := root + req.Name

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := api.Get(ctx, key, &client.GetOptions{Quorum: true})
	if err != nil {
		ce, ok := err.(client.Error)
		if ok {
			if ce.Code == client.ErrorCodeKeyNotFound {
				return nil, vidar.NewHTTPError(http.StatusNoContent, ce.Message)
			} else {
				return nil, vidar.NewHTTPError(http.StatusBadGateway, ce.Message)
			}
		} else if err == context.DeadlineExceeded {
			return nil, vidar.GatewayTimeoutError
		} else {
			return nil, err
		}
	}

	result := response{
		Message: resp,
	}

	ctx = context.WithValue(ctx, "result", result)
	return ctx, nil
}
