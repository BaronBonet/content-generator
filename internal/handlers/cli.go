package handlers

import (
	"context"
	"fmt"
	"sort"

	"github.com/BaronBonet/content-generator/internal/core/ports"
	"github.com/BaronBonet/go-logger/logger"
	"github.com/urfave/cli/v2"
)

type CliHandler struct {
	app *cli.App
}

func NewCLIHandler(ctx context.Context, service ports.Service, logger logger.Logger) *CliHandler {
	app := &cli.App{
		Name:                 "Generate News Content",
		EnableBashCompletion: true,
		Usage:                "Cli tool for generating news content",
		Commands: []*cli.Command{
			{
				Name:  "generateNewsContent",
				Usage: "Run the news content generation process",
				Action: func(c *cli.Context) error {
					err := service.GenerateNewsContent(ctx)
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  "createPrompt",
				Usage: "Create a prompt",
				Action: func(c *cli.Context) error {
					prompt := c.Args().Get(0)
					if prompt == "" {
						return nil
					}

					prompt, err := service.CreatePrompt(ctx, prompt)
					if err != nil {
						return err
					}
					logger.Info("Created prompt", "prompt", prompt)
					return nil
				},
			},
			{
				Name:  "generateImage",
				Usage: "Create an image and get the url to the image, requires entering a prompt",
				Action: func(c *cli.Context) error {
					prompt := c.Args().Get(0)
					if prompt == "" {
						return nil
					}
					url, err := service.GenerateImage(ctx, prompt)
					if err != nil {
						return err
					}
					fmt.Println("Image generated at:")
					fmt.Println(url)
					fmt.Println("")
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
