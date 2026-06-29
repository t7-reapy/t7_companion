// Package cli builds the t7kb command tree (cobra).
package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/t7-reapy/t7_companion/internal/store"
)

var dbFlag string

// Execute runs the root command, exiting non-zero on error.
func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "t7kb",
		Short: "Hybrid search over a local Black Ops 3 modding knowledge base",
		Long: "t7kb queries a local t7kb.db with hybrid retrieval (BM25 + vector).\n\n" +
			"It is built to be driven by an AI agent over MCP (`t7kb mcp`), with a\n" +
			"small CLI for direct use.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.PersistentFlags().StringVar(&dbFlag, "db", "",
		"path to t7kb.db (default: $T7KB_DB, then beside the binary, then ./t7kb.db)")
	root.AddCommand(newSearchCmd(), newGetCmd(), newMCPCmd(), newEmbedCmd())
	return root
}

// resolveDB picks the t7kb.db path: --db > $T7KB_DB > beside-binary > cwd.
func resolveDB() string {
	if dbFlag != "" {
		return dbFlag
	}
	if env := os.Getenv("T7KB_DB"); env != "" {
		return env
	}
	if exe, err := os.Executable(); err == nil {
		beside := filepath.Join(filepath.Dir(exe), "t7kb.db")
		if _, err := os.Stat(beside); err == nil {
			return beside
		}
	}
	return "t7kb.db"
}

func openStore() (*store.Store, error) {
	return store.Open(resolveDB())
}
