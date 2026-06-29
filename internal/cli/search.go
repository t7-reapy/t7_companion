package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/t7-reapy/t7_companion/internal/embed"
)

func newSearchCmd() *cobra.Command {
	var (
		num  int
		bm25 bool
	)
	cmd := &cobra.Command{
		Use:   "search <query>...",
		Short: "Hybrid (BM25 + vector) search",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := strings.Join(args, " ")

			st, err := openStore()
			if err != nil {
				return err
			}
			defer st.Close()

			var qvec []float32
			if !bm25 {
				emb, err := embed.New()
				if err != nil {
					return err
				}
				if qvec, err = emb.Embed(query); err != nil {
					return err
				}
			}

			hits, err := st.SearchHybrid(cmd.Context(), query, qvec, num)
			if err != nil {
				return err
			}
			if len(hits) == 0 {
				fmt.Fprintln(cmd.ErrOrStderr(), "no results")
				return nil
			}
			out := cmd.OutOrStdout()
			for i, h := range hits {
				fmt.Fprintf(out, "%2d. %s  (%s)\n", i+1, h.DocID, h.Source)
				fmt.Fprintf(out, "    %s\n", h.Title)
				if h.Snippet != "" {
					fmt.Fprintf(out, "    %s\n", h.Snippet)
				}
			}
			return nil
		},
	}
	cmd.Flags().IntVarP(&num, "num", "n", 10, "number of results")
	cmd.Flags().BoolVar(&bm25, "bm25", false, "keyword-only (skip vector embedding)")
	return cmd
}
