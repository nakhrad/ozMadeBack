package product

import (
	"testing"
	"time"
)

func TestTrendingScorePrefersRecentProducts(t *testing.T) {
	now := time.Date(2026, time.April, 5, 12, 0, 0, 0, time.UTC)
	recent := TrendingScore(100, now.Add(-2*time.Hour), now)
	old := TrendingScore(100, now.Add(-72*time.Hour), now)

	if recent <= old {
		t.Fatalf("expected recent product to have higher score, got recent=%f old=%f", recent, old)
	}
}

func TestTrendingScoreClampsFutureCreationTime(t *testing.T) {
	now := time.Date(2026, time.April, 5, 12, 0, 0, 0, time.UTC)
	score := TrendingScore(50, now.Add(3*time.Hour), now)
	expected := TrendingScore(50, now, now)

	if score != expected {
		t.Fatalf("expected future timestamps to be clamped, got score=%f expected=%f", score, expected)
	}
}

func TestTrendingScoreEdgeCases(t *testing.T) {
	now := time.Date(2026, time.April, 5, 12, 0, 0, 0, time.UTC)

	t.Run("Zero views", func(t *testing.T) {
		score := TrendingScore(0, now.Add(-1*time.Hour), now)
		if score != 0 {
			t.Errorf("expected 0 score for 0 views, got %f", score)
		}
	})

	t.Run("Extreme age", func(t *testing.T) {
		old := TrendingScore(1000, now.Add(-1000*time.Hour), now)
		recent := TrendingScore(1, now.Add(-1*time.Hour), now)
		if recent <= old {
			t.Errorf("expected 1 view on recent product to beat 1000 views on very old product, got recent=%f old=%f", recent, old)
		}
	})

	t.Run("Very new product", func(t *testing.T) {
		score := TrendingScore(100, now, now)
		if score <= 0 {
			t.Errorf("expected positive score for new product with views, got %f", score)
		}
	})
}
