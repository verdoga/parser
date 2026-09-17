package app

import (
	"testing"
	"time"
)

// TestNewProcessingAttempt проверяет обязательные поля и нормализацию времени.
func TestNewProcessingAttempt(t *testing.T) {
	started := time.Date(2025, 1, 2, 3, 4, 5, 0, time.FixedZone("test", 2*60*60))
	attempt, err := newProcessingAttempt("id-1", "tool-1", started)
	if err != nil {
		t.Fatalf("newProcessingAttempt() error: %v", err)
	}
	if attempt.id != "id-1" || attempt.version != "tool-1" || !attempt.startedAt.Equal(started.UTC()) || attempt.startedAt.Location() != time.UTC {
		t.Fatalf("attempt = %#v", attempt)
	}
	for _, test := range []struct {
		name, id, version string
	}{
		{name: "empty id", version: "v"},
		{name: "empty version", id: "id"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, err := newProcessingAttempt(test.id, test.version, started); err == nil {
				t.Fatal("newProcessingAttempt() error = nil")
			}
		})
	}
}

// TestProcessingAttemptFinish проверяет UTC-время и неотрицательную длительность.
func TestProcessingAttemptFinish(t *testing.T) {
	started := time.Date(2025, 1, 2, 3, 0, 0, 0, time.UTC)
	attempt := processingAttempt{id: "id", version: "v", startedAt: started}
	finished := started.Add(1500 * time.Millisecond).In(time.FixedZone("later", -4*60*60))
	processing := attempt.finish(finished)
	if processing.ID != "id" || processing.Tool != "dslparser" || processing.Version == nil || *processing.Version != "v" {
		t.Fatalf("processing identity = %#v", processing)
	}
	if processing.StartedAt == nil || *processing.StartedAt != started.Format(time.RFC3339Nano) || processing.DurationMS == nil || *processing.DurationMS != 1500 {
		t.Fatalf("processing timing = %#v", processing)
	}
	if got := attempt.finish(started.Add(-time.Second)); got.DurationMS == nil || *got.DurationMS != 0 {
		t.Fatalf("negative duration = %v, want pointer to 0", got.DurationMS)
	}
}
