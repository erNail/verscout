package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRootCmdPrintsHelpWithoutError(t *testing.T) {
	t.Parallel()

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"-h"})

	err := cmd.Execute()

	require.NoError(t, err)
}

func TestRootCmdPrintsVersionWithoutError(t *testing.T) {
	t.Parallel()

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--version"})

	err := cmd.Execute()

	require.NoError(t, err)
}

func TestRootCmdCallsLatestSubcommand(t *testing.T) {
	t.Parallel()

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"latest", "-h"})

	err := cmd.Execute()

	require.NoError(t, err)
}

func TestRootCmdCallsNextSubcommand(t *testing.T) {
	t.Parallel()

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"next", "-h"})

	err := cmd.Execute()

	require.NoError(t, err)
}
