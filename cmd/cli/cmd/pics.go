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

type getPicCommand struct {
	shared.ConnectionInfo
	animal string
}

func NewPicsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pics",
		Short: "Find pics and things",
	}

	cmd.AddCommand(newCatPicsCommand())
	cmd.AddCommand(newDogPicsCommand())

	return cmd
}

func newCatPicsCommand() *cobra.Command {
	var (
		picCommand getPicCommand
	)

	picCommand.animal = "cat"

	cmd := &cobra.Command{
		Use:   "cats",
		Short: "pics of cats",
		RunE:  picCommand.Run,
	}

	return cmd
}

func newDogPicsCommand() *cobra.Command {
	var (
		picCommand getPicCommand
	)

	picCommand.animal = "dog"

	cmd := &cobra.Command{
		Use:   "dogs",
		Short: "pics of dogs",
		RunE:  picCommand.Run,
	}

	picCommand.registerFlags(cmd)

	return cmd
}

func (pc *getPicCommand) registerFlags(cmd *cobra.Command) {
	shared.RegisterEnvironmentFlags(&pc.ConnectionInfo, cmd)
}

func (pc *getPicCommand) Run(cmd *cobra.Command, _ []string) error {
	env, _ := shared.GetEnvConfig(pc.Environment)

	var url string
	if pc.animal == "cat" {
		url = fmt.Sprintf("%s/%s", env.BaseUrl, "pics/cats")
	} else {
		url = fmt.Sprintf("%s/%s", env.BaseUrl, "pics/dogs")
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	pic := models.Pic{}
	err = json.NewDecoder(resp.Body).Decode(&pic)
	if err != nil {
		return err
	}

	fmt.Printf("%s", pic.Data)

	return nil
}
