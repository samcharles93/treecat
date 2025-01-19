package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/samcharles93/treecat/tree"
	"github.com/urfave/cli/v3"
)

var (
	Version   = "0.1.1" // Added depth control and safety features
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	cmd := &cli.Command{
		Name:    "treecat",
		Usage:   "Display a directory tree with file contents",
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "e",
				Usage:   "Pattern to exclude (glob syntax)",
				Aliases: []string{"exclude"},
			},
			&cli.StringFlag{
				Name:    "i",
				Usage:   "Pattern to include (glob syntax)",
				Aliases: []string{"include"},
			},
			&cli.StringFlag{
				Name:    "o",
				Usage:   "Output file path",
				Aliases: []string{"output", "out"},
			},
			&cli.IntFlag{
				Name:    "d",
				Usage:   "Maximum depth to traverse (default: 1)",
				Aliases: []string{"depth"},
				Value:   1,
			},
			&cli.BoolFlag{
				Name:    "f",
				Usage:   "Force processing of large directories",
				Aliases: []string{"force"},
				Value:   false,
			},
		},
		Action: func(_ context.Context, c *cli.Command) error {
			startDir := "."
			if c.NArg() > 0 {
				startDir = c.Args().First()
			}

			absPath, err := tree.ResolveAbsolutePath(startDir)
			if err != nil {
				return fmt.Errorf("error getting absolute path: %w", err)
			}

			root, err := tree.BuildTree(absPath, c.String("e"), c.String("i"), startDir, int(c.Int("d")), c.Bool("f"))
			if err != nil {
				return fmt.Errorf("error building tree: %w", err)
			}

			outputPath := c.String("o")
			if outputPath != "" {
				file, err := os.Create(outputPath)
				if err != nil {
					return fmt.Errorf("error creating output file: %w", err)
				}
				defer file.Close()

				fmt.Fprintln(file, absPath)
				tree.PrintTreeWithOutput(root, "", true, file, absPath)
			} else {
				fmt.Fprintln(os.Stdout, absPath)
				tree.PrintTreeWithOutput(root, "", true, os.Stdout, absPath)
			}

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
