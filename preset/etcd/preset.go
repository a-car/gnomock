// Package redis includes Redis implementation of Gnomock Preset interface.
// This Preset can be passed to gnomock.Start() function to create a configured
// Redis container to use in tests.
package etcd

import (
	"context"
	"fmt"
	"time"

	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/internal/registry"
	etcd_client "go.etcd.io/etcd/client/v3"
)

const defaultVersion = "latest"

func init() {
	registry.Register("etcd", func() gnomock.Preset { return &P{} })
}

// Preset creates a new Gmomock Etcd preset. This preset includes a Etcd
// specific healthcheck function, default Redis image and port, and allows to
// optionally set up initial state.
func Preset(opts ...Option) gnomock.Preset {
	p := &P{}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// P is a Gnomock Preset implementation for etcd storage.
type P struct {
	Values  map[string]interface{} `json:"values"`
	Version string                 `json:"version"`
}

// Image returns an image that should be pulled to create this container.
func (p *P) Image() string {
    return fmt.Sprintf("bitnami/etcd:%s", p.Version)
}

// Ports returns ports that should be used to access this container.
func (p *P) Ports() gnomock.NamedPorts {
	return gnomock.DefaultTCP(2379)
}

// Options returns a list of options to configure this container.
func (p *P) Options() []gnomock.Option {
	p.setDefaults()

	opts := []gnomock.Option{
		gnomock.WithHealthCheck(healthcheck),
	}

	if p.Values != nil {
		initf := func(ctx context.Context, c *gnomock.Container) error {
			addr := c.Address(gnomock.DefaultPort)
            client, _ := etcd_client.New(etcd_client.Config{
                Endpoints: []string{addr},
            })

			for k, v := range p.Values {
                cx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
				_, err := client.Put(cx, k, v)
                cancel()
				if err != nil {
					return fmt.Errorf("can't set '%s'='%v': %w", k, v, err)
				}
			}

			return nil
		}

		opts = append(opts, gnomock.WithInit(initf))
	}

	return opts
}

func (p *P) setDefaults() {
	if p.Version == "" {
		p.Version = defaultVersion
	}
}

func healthcheck(ctx context.Context, c *gnomock.Container) error {
	addr := c.Address(gnomock.DefaultPort)
    client, err := etcd_client.New(etcd_client.Config{
        Endpoints: []string{addr},
    })
    cx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
    k := "test_key"
    v := "test_value"
    _, err = client.Put(cx, k, v)
    cancel()
    if err != nil {
        return fmt.Errorf("can't set '%s'='%v': %w", k, v, err)
    }

	return err
}
