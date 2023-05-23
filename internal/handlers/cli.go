package handlers

import (
	"context"
	"github.com/BaronBonet/content-generator/internal/core/ports"
	"github.com/urfave/cli/v2"
	"sort"
)

type CliHandler struct {
	app *cli.App
}

func NewCLIHandler(ctx context.Context, service ports.Service) *CliHandler {
	app := &cli.App{
		Name:                 "Generate News Content",
		EnableBashCompletion: true,
		Usage:                "Run the full content generation process",
		Commands: []*cli.Command{
			{
				Name:  "generateNewsContent",
				Usage: "List all raw maps",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name: "json",
					},
				},
				Action: func(c *cli.Context) error {
					err := service.GenerateNewsContent(ctx)
					if err != nil {
						return err
					}
					return nil
				},
			},
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	return &CliHandler{app: app}
}

func (h *CliHandler) Run(args []string) error {
	return h.app.Run(args)
}
