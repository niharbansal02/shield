package cmd_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/odpf/shield/cmd"
	"github.com/stretchr/testify/assert"
)

func TestClientAction(t *testing.T) {
	t.Run("without config file", func(t *testing.T) {
		tests := []struct {
			name        string
			cliConfig   *cmd.Config
			subCommands []string
			want        string
			err         error
		}{
			{
				name:        "`action` list only should throw error host not found",
				want:        "",
				subCommands: []string{"list"},
				err:         cmd.ErrClientConfigHostNotFound,
			},
			{
				name:        "`action` list with host flag should pass",
				want:        "",
				subCommands: []string{"list", "-h", "test"},
				err:         context.DeadlineExceeded,
			},
			{
				name:        "`action` create only should throw error host not found",
				want:        "",
				subCommands: []string{"create"},
				err:         cmd.ErrClientConfigHostNotFound,
			},
			{
				name:        "`action` create with host flag should throw error missing required flag",
				want:        "",
				subCommands: []string{"create", "-h", "test"},
				err:         errors.New("required flag(s) \"file\", \"header\" not set"),
			},
			{
				name:        "`action` edit without host should throw error host not found",
				want:        "",
				subCommands: []string{"edit", "123"},
				err:         cmd.ErrClientConfigHostNotFound,
			},
			{
				name:        "`action` edit with host flag should throw error missing required flag",
				want:        "",
				subCommands: []string{"edit", "123", "-h", "test"},
				err:         errors.New("required flag(s) \"file\" not set"),
			},
			{
				name:        "`action` view without host should throw error host not found",
				want:        "",
				subCommands: []string{"view", "123"},
				err:         cmd.ErrClientConfigHostNotFound,
			},
			{
				name:        "`action` view with host flag should pass",
				want:        "",
				subCommands: []string{"view", "123", "-h", "test"},
				err:         context.DeadlineExceeded,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.cliConfig = &cmd.Config{}
				cli := cmd.New(tt.cliConfig)

				buf := new(bytes.Buffer)
				cli.SetOutput(buf)
				args := append([]string{"action"}, tt.subCommands...)
				cli.SetArgs(args)

				err := cli.Execute()
				got := buf.String()

				assert.Equal(t, tt.err, err)
				assert.Equal(t, tt.want, got)
			})
		}
	})
}
