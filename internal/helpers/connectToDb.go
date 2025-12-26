package helpers

import (
	"log/slog"
	"os"

	"git.sr.ht/~jamesponddotco/bunnystorage-go"
)

var BunnyClient *bunnystorage.Client

func ConnectToBunny() {
	readOnlyKey := os.Getenv("BUNNYNET_READ_API_KEY")
	readWriteKey := os.Getenv("BUNNYNET_WRITE_API_KEY")

	if readOnlyKey == "" || readWriteKey == "" {
		slog.Error("missing bunny env vars")
		os.Exit(1)
	}

	client, err := bunnystorage.NewClient(&bunnystorage.Config{
		StorageZone: "fillipsportfolio",
		Key:         readWriteKey,
		ReadOnlyKey: readOnlyKey,
		Endpoint:    bunnystorage.EndpointFalkenstein,
	})

	if err != nil {
		slog.Error("Error connecting to Image Server", "error", err)
		os.Exit(1)
	}

	BunnyClient = client
}
