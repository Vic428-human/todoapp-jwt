package utils

import (
	"context"
	"fmt"
	"time"
)

func PerformOperation(ctx context.Context) {
	select {
	case <-time.After(2 * time.Second):
		fmt.Println("Operation completed")
	case <-ctx.Done():
		fmt.Println("Operation timed out or cancelled")
	}
}
