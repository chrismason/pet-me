package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/chrismason/pet-me/cmd/cli/shared"
	"github.com/chrismason/pet-me/internal/models"

	"github.com/spf13/cobra"
)

type catFactCommand struct {
	shared.ConnectionInfo
	factCount uint32
}

func NewFactsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "facts",
		Short: "Things you can find out",
	}

	cmd.AddCommand(newCatFactsCommand())

	return cmd
}

func newCatFactsCommand() *cobra.Command {
	var (
		factCmd catFactCommand
	)

	cmd := &cobra.Command{
		Use:   "cats",
		Short: "Things you can learn about cats",
		RunE:  factCmd.Run,
	}

	factCmd.registerFlags(cmd)

	return cmd
}

func (cfc *catFactCommand) registerFlags(cmd *cobra.Command) {
	shared.RegisterEnvironmentFlags(&cfc.ConnectionInfo, cmd)
	cmd.Flags().Uint32VarP(&cfc.factCount, "count", "c", 1, "How many facts you would like between 1 and 10")
}

func (cfc *catFactCommand) Run(cmd *cobra.Command, _ []string) error {
	if cfc.factCount > 10 {
		return fmt.Errorf("you must specify between 1 and 10 facts but you provided %d", cfc.factCount)
	}

	env, _ := shared.GetEnvConfig(cfc.Environment)

	url := fmt.Sprintf("%s/%s?factCount=%d", env.BaseUrl, "facts/cats", cfc.factCount)

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	facts := models.FactResponse{}
	err = json.NewDecoder(resp.Body).Decode(&facts)
	if err != nil {
		return err
	}

	for _, fact := range facts.Facts {
		fmt.Printf("%s\n", fact)
	}
	return nil
}
