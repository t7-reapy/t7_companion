package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <doc_id>",
		Short: "Print a document's full body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := openStore()
			if err != nil {
				return err
			}
			defer st.Close()

			doc, err := st.Get(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if doc == nil {
				return fmt.Errorf("no such doc_id: %s", args[0])
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "# %s\n\n", doc.Title)
			fmt.Fprintf(out, "doc_id: %s\nsource: %s  ·  reliability: %.2f\n",
				doc.DocID, doc.Source, doc.Reliability)
			if doc.URL != "" {
				fmt.Fprintf(out, "url: %s\n", doc.URL)
			}
			fmt.Fprintf(out, "\n%s\n", doc.Body)
			return nil
		},
	}
}
