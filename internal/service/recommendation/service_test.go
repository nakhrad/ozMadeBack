package recommendation

import (
	"testing"
	"time"

	"ozMadeBack/internal/models"
)

func TestGlobalScorePrefersRecentPopularProducts(t *testing.T) {
	now := time.Date(2026, time.April, 8, 12, 0, 0, 0, time.UTC)

	recentPopular := globalScore(models.Product{
		ViewCount:     200,
		AverageRating: 4.5,
		CreatedAt:     now.Add(-6 * time.Hour),
	}, 10, 3, now)

	oldUnpopular := globalScore(models.Product{
		ViewCount:     10,
		AverageRating: 1.0,
		CreatedAt:     now.Add(-20 * 24 * time.Hour),
	}, 1, 0, now)

	if recentPopular <= oldUnpopular {
		t.Fatalf("expected recent popular product to score higher, got recent=%f old=%f", recentPopular, oldUnpopular)
	}
}

func TestGlobalScoreEdgeCases(t *testing.T) {
	now := time.Date(2026, time.April, 8, 12, 0, 0, 0, time.UTC)

	t.Run("Zero values", func(t *testing.T) {
		score := globalScore(models.Product{CreatedAt: now}, 0, 0, now)
		if score <= 0 {
			t.Errorf("expected positive score for new product even with 0 views, got %f", score)
		}
	})

	t.Run("Future creation date", func(t *testing.T) {
		score := globalScore(models.Product{CreatedAt: now.Add(1 * time.Hour)}, 0, 0, now)
		expected := globalScore(models.Product{CreatedAt: now}, 0, 0, now)
		if score != expected {
			t.Errorf("expected future creation date to be treated as now, got %f expected %f", score, expected)
		}
	})

	t.Run("High views low rating vs low views high rating", func(t *testing.T) {
		highViews := globalScore(models.Product{ViewCount: 1000, AverageRating: 1.0, CreatedAt: now}, 0, 0, now)
		highRating := globalScore(models.Product{ViewCount: 10, AverageRating: 5.0, CreatedAt: now}, 0, 0, now)
		// Verification of actual balance depends on weights, but let's ensure they are distinguishable
		if highViews == highRating {
			t.Errorf("scores should likely be different for high views vs high rating")
		}
	})
}

func TestUniqueUint(t *testing.T) {
	tests := []struct {
		name     string
		input    []uint
		expected []uint
	}{
		{"Empty", []uint{}, []uint{}},
		{"No duplicates", []uint{1, 2, 3}, []uint{1, 2, 3}},
		{"Duplicates", []uint{1, 2, 2, 3, 1}, []uint{1, 2, 3}},
		{"Single element", []uint{5}, []uint{5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uniqueUint(tt.input)
			if len(got) != len(tt.expected) {
				t.Fatalf("expected length %d, got %d", len(tt.expected), len(got))
			}
			for i, v := range got {
				if v != tt.expected[i] {
					t.Errorf("at index %d: expected %d, got %d", i, tt.expected[i], v)
				}
			}
		})
	}
}

func TestNormalizePreferenceKeyExtra(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  MixedCase  ", "mixedcase"},
		{"\tTabsAndNewlines\n", "tabsandnewlines"},
		{"", ""},
		{"AlreadyLower", "alreadylower"},
	}

	for _, tt := range tests {
		if got := normalizePreferenceKey(tt.input); got != tt.expected {
			t.Errorf("normalizePreferenceKey(%q) = %q; expected %q", tt.input, got, tt.expected)
		}
	}
}
