package services

import (
	"encoding/json"
	"testing"
)

func TestBuildSearchBodyIncludesQueryAndFilters(t *testing.T) {
	minCost := 10.0
	maxCost := 50.0
	service := &ProductSearchService{}

	body, err := service.buildSearchBody(ProductSearchParams{
		Query:    "handmade lamp",
		Type:     "home",
		Category: "decor",
		MinCost:  &minCost,
		MaxCost:  &maxCost,
		Limit:    12,
		Offset:   24,
	})
	if err != nil {
		t.Fatalf("buildSearchBody returned error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("failed to decode search body: %v", err)
	}

	if payload["size"].(float64) != 12 {
		t.Fatalf("expected size 12, got %v", payload["size"])
	}
	if payload["from"].(float64) != 24 {
		t.Fatalf("expected offset 24, got %v", payload["from"])
	}

	query := payload["query"].(map[string]any)["bool"].(map[string]any)
	if len(query["should"].([]any)) == 0 {
		t.Fatal("expected full-text query in should clause")
	}
	if query["minimum_should_match"].(float64) != 1 {
		t.Fatalf("expected minimum_should_match 1, got %v", query["minimum_should_match"])
	}
	if len(query["filter"].([]any)) != 3 {
		t.Fatalf("expected 3 filters, got %d", len(query["filter"].([]any)))
	}
}

func TestBuildSearchBodyEdgeCases(t *testing.T) {
	service := &ProductSearchService{}

	t.Run("Negative offset and limit", func(t *testing.T) {
		body, _ := service.buildSearchBody(ProductSearchParams{Limit: -1, Offset: -5})
		var payload map[string]any
		json.Unmarshal(body, &payload)
		if payload["size"].(float64) != 20 {
			t.Errorf("expected default size 20, got %v", payload["size"])
		}
		if payload["from"].(float64) != 0 {
			t.Errorf("expected offset 0, got %v", payload["from"])
		}
	})

	t.Run("Extreme limit", func(t *testing.T) {
		body, _ := service.buildSearchBody(ProductSearchParams{Limit: 1000})
		var payload map[string]any
		json.Unmarshal(body, &payload)
		if payload["size"].(float64) != 100 {
			t.Errorf("expected max size 100, got %v", payload["size"])
		}
	})

	t.Run("Only filters no query", func(t *testing.T) {
		body, _ := service.buildSearchBody(ProductSearchParams{Type: "furniture"})
		var payload map[string]any
		json.Unmarshal(body, &payload)
		query := payload["query"].(map[string]any)["bool"].(map[string]any)
		if _, ok := query["should"]; ok {
			t.Fatal("did not expect should clause when query is empty")
		}
		filters := query["filter"].([]any)
		if len(filters) != 1 {
			t.Errorf("expected 1 filter, got %d", len(filters))
		}
	})
}

func TestParseSearchFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected *float64
		wantErr  bool
	}{
		{"10.5", float64Ptr(10.5), false},
		{" 20 ", float64Ptr(20.0), false},
		{"", nil, false},
		{"  ", nil, false},
		{"invalid", nil, true},
	}

	for _, tt := range tests {
		got, err := ParseSearchFloat(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseSearchFloat(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if tt.expected == nil {
			if got != nil {
				t.Errorf("ParseSearchFloat(%q) = %v, want nil", tt.input, *got)
			}
		} else {
			if got == nil || *got != *tt.expected {
				t.Errorf("ParseSearchFloat(%q) = %v, want %v", tt.input, got, *tt.expected)
			}
		}
	}
}

func float64Ptr(f float64) *float64 {
	return &f
}
