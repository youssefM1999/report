package retry

import (
	"errors"
	"testing"
	"time"
)

func TestRetry_SuccessOnFirstAttempt(t *testing.T) {
	attempts := 0
	fn := func() error {
		attempts++
		return nil
	}

	err := Retry(fn, 3)
	if err != nil {
		t.Errorf("Retry() should not return error on success, got: %v", err)
	}
	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}

func TestRetry_SuccessAfterRetries(t *testing.T) {
	attempts := 0
	fn := func() error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary error")
		}
		return nil
	}

	err := Retry(fn, 5)
	if err != nil {
		t.Errorf("Retry() should not return error on success, got: %v", err)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestRetry_FailureAfterMaxRetries(t *testing.T) {
	attempts := 0
	expectedAttempts := 3
	fn := func() error {
		attempts++
		return errors.New("persistent error")
	}

	err := Retry(fn, expectedAttempts)
	if err == nil {
		t.Error("Retry() should return error after max retries")
	}
	if attempts != expectedAttempts {
		t.Errorf("Expected %d attempts, got %d", expectedAttempts, attempts)
	}
	expectedError := "failed to execute function after 3 attempts"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestRetry_BackoffDelay(t *testing.T) {
	attempts := 0
	attemptTimes := []time.Time{}

	fn := func() error {
		attempts++
		attemptTimes = append(attemptTimes, time.Now())
		if attempts < 3 {
			return errors.New("temporary error")
		}
		return nil
	}

	err := Retry(fn, 5)
	if err != nil {
		t.Errorf("Retry() should not return error on success, got: %v", err)
	}

	// Verify that delays increase exponentially (approximately)
	// First retry should wait ~1 second (2^0), second retry ~2 seconds (2^1)
	if len(attemptTimes) >= 2 {
		delay1 := attemptTimes[1].Sub(attemptTimes[0])
		expectedDelay1 := 1 * time.Second
		if delay1 < expectedDelay1-time.Millisecond*200 || delay1 > expectedDelay1+time.Millisecond*500 {
			t.Errorf("First retry delay should be approximately 1s, got %v", delay1)
		}
	}
	if len(attemptTimes) >= 3 {
		delay2 := attemptTimes[2].Sub(attemptTimes[1])
		expectedDelay2 := 2 * time.Second
		if delay2 < expectedDelay2-time.Millisecond*200 || delay2 > expectedDelay2+time.Millisecond*500 {
			t.Errorf("Second retry delay should be approximately 2s, got %v", delay2)
		}
	}
}

func TestRetry_ZeroRetries(t *testing.T) {
	attempts := 0
	fn := func() error {
		attempts++
		return errors.New("error")
	}

	err := Retry(fn, 0)
	if err == nil {
		t.Error("Retry() should return error with 0 retries")
	}
	if attempts != 0 {
		t.Errorf("Expected 0 attempts with 0 retries, got %d", attempts)
	}
}

func TestRetry_OneRetry(t *testing.T) {
	attempts := 0
	fn := func() error {
		attempts++
		return errors.New("error")
	}

	err := Retry(fn, 1)
	if err == nil {
		t.Error("Retry() should return error after 1 retry")
	}
	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}
