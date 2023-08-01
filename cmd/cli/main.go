package cli

import (
	"context"
	"fmt"
	"os"
)

func Execute() {
	ctx := context.Background()

	rootCmd := NewRootCmd()

	// TODO: make embedded config on build time
	// if err := config.New().Load(); err != nil {
	// log.Fatalf("unable to load config: %v", err)
	//}

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
