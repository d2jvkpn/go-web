package wrap

import (
	"testing"

	. "github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	type Github struct {
		Name string
		Link string
	}

	type Gitlab struct {
		User    string `mapstructure:"user"`
		Address string `mapstructure:"address"`
	}

	github, gitlab := new(Github), new(Gitlab)

	var err error

	err = LoadConfig("testConfig", "config.demo.yaml", map[string]any{
		"github": github,
		"gitlab": gitlab,
	})
	NoError(t, err)

	Equal(t, github.Name, "d2jvkpn")
	Equal(t, github.Link, "https://github.com/d2jvkpn")

	Equal(t, gitlab.User, "d2jvkpn")
	Equal(t, gitlab.Address, "https://gitlab.com/d2jvkpn")
}
