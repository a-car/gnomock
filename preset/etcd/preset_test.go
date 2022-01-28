package etcd_test

import (
	"testing"

    etcd_client "go.etcd.io/etcd/client/v3"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/redis"
	"github.com/stretchr/testify/require"
)

func TestPreset(t *testing.T) {
	t.Parallel()

	for _, version := range []string{"3", "latest"} {
		t.Run(version, testPreset(version))
	}
}

func testPreset(version string) func(t *testing.T) {
	return func(t *testing.T) {
		vs := make(map[string]interface{})

		vs["test_1"] = "value_1"
		vs["test_2"] = "value_2"

		p := etcd.Preset(
			etcd.WithValues(vs),
			etcd.WithVersion(version),
		)
		container, err := gnomock.Start(p)

		defer func() { require.NoError(t, gnomock.Stop(container)) }()

		require.NoError(t, err)

		addr := container.DefaultAddress()
        client, err = etcd_client.New(etcd_client.Config{
            Endpoints: []string{addr},
        })

        for k, v := range vs {
            resp, err := client.Get(cx,k)
            if resp != v {
                t.Error()
            }
        }
	}
}
