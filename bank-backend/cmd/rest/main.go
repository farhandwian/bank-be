package main

import (
	"bank-backend/internal/config"
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	rootCmd := &cobra.Command{}
	cmd := []*cobra.Command{
		{
			Use:   "serve-http",
			Short: "Run HTTP server",
			Run: func(cmd *cobra.Command, _ []string) {
				config.StartHTTPServer((ctx))
			},
		},
	}

	rootCmd.AddCommand(cmd...)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
