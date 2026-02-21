package main

import (
	"log/slog"
	"os"
)

func main() {
	api := api{
		config: config{
			addr: ":8000",
			db: dbConfig{
				dsn: "",
			},
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := api.run(api.mount()); err != nil {
		slog.Error("server has failed to start", "error", err)
		os.Exit(1)
	}
}
