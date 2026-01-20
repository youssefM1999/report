package retry

import (
	"fmt"
	"math"
	"time"
)

func Retry(fn func() error, maxRetries int) error {
	for i := range maxRetries {
		if err := fn(); err != nil {
			time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
			continue
		}
		return nil
	}
	return fmt.Errorf("failed to execute function after %d attempts", maxRetries)
}
