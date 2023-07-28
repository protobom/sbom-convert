package cli

import (
	"context"
	"log"
)

func Execute() {
	ctx := context.Background()

	rootCmd := NewRootCmd()

	// TODO: make logger configurable
	// TODO: make embedded config on build time
	// if err := config.New().Load(); err != nil {
	// log.Fatalf("unable to load config: %v", err)
	//}

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
}
