package services

import (
	"ozMadeBack/internal/dto"
	"testing"
)

func TestCalculateIntercityFareValidation(t *testing.T) {
	service := &CDEKService{}

	tests := []struct {
		name    string
		req     *dto.IntercityEstimateRequest
		wantErr bool
	}{
		{
			name: "Same city",
			req: &dto.IntercityEstimateRequest{
				FromAddress: dto.AddressDetails{City: "Almaty", FullAddress: "A"},
				ToAddress:   dto.AddressDetails{City: "Almaty", FullAddress: "B"},
			},
			wantErr: true,
		},
		{
			name: "Zero weight",
			req: &dto.IntercityEstimateRequest{
				FromAddress: dto.AddressDetails{City: "Almaty", FullAddress: "A"},
				ToAddress:   dto.AddressDetails{City: "Astana", FullAddress: "B"},
				Package:     dto.PackageDetails{WeightGrams: 0, HeightCm: 10, WidthCm: 10, DepthCm: 10},
			},
			wantErr: true,
		},
		{
			name: "Missing address",
			req: &dto.IntercityEstimateRequest{
				FromAddress: dto.AddressDetails{City: "Almaty", FullAddress: ""},
				ToAddress:   dto.AddressDetails{City: "Astana", FullAddress: "B"},
				Package:     dto.PackageDetails{WeightGrams: 100, HeightCm: 10, WidthCm: 10, DepthCm: 10},
			},
			wantErr: true,
		},
		{
			name: "Invalid dimensions",
			req: &dto.IntercityEstimateRequest{
				FromAddress: dto.AddressDetails{City: "Almaty", FullAddress: "A"},
				ToAddress:   dto.AddressDetails{City: "Astana", FullAddress: "B"},
				Package:     dto.PackageDetails{WeightGrams: 100, HeightCm: -1, WidthCm: 10, DepthCm: 10},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.CalculateIntercityFare(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateIntercityFare() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
