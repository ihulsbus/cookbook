package models

import (
	"testing"

	"github.com/shopspring/decimal"
)

func d(s string) decimal.Decimal {
	return decimal.RequireFromString(s)
}

func dp(s string) *decimal.Decimal {
	v := d(s)
	return &v
}

var (
	gram      = Unit{ShortName: "g", Dimension: DimensionMass, ToBase: decimal.NewFromInt(1)}
	kilogram  = Unit{ShortName: "kg", Dimension: DimensionMass, ToBase: decimal.NewFromInt(1000)}
	milligram = Unit{ShortName: "mg", Dimension: DimensionMass, ToBase: d("0.001")}

	milliliter = Unit{ShortName: "ml", Dimension: DimensionVolume, ToBase: decimal.NewFromInt(1)}
	liter      = Unit{ShortName: "l", Dimension: DimensionVolume, ToBase: decimal.NewFromInt(1000)}

	piece = Unit{ShortName: "pcs", Dimension: DimensionCount, ToBase: decimal.NewFromInt(1)}
	dozen = Unit{ShortName: "doz", Dimension: DimensionCount, ToBase: decimal.NewFromInt(12)}

	water     = Ingredient{Name: "water", DensityGPerMl: dp("1")}
	oliveOil  = Ingredient{Name: "olive oil", DensityGPerMl: dp("0.92")}
	egg       = Ingredient{Name: "egg", GramsPerPiece: dp("50")}
	eggAndOil = Ingredient{Name: "egg-like oil", DensityGPerMl: dp("1"), GramsPerPiece: dp("50")}
	unknown   = Ingredient{Name: "mystery powder"}
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name    string
		qty     decimal.Decimal
		from    Unit
		to      Unit
		ing     Ingredient
		want    decimal.Decimal
		wantErr bool
	}{
		{
			name: "mass to mass, same dimension",
			qty:  decimal.NewFromInt(2000),
			from: gram, to: kilogram, ing: unknown,
			want: decimal.NewFromInt(2),
		},
		{
			name: "mass to mass, fractional factor",
			qty:  decimal.NewFromInt(5),
			from: gram, to: milligram, ing: unknown,
			want: decimal.NewFromInt(5000),
		},
		{
			name: "volume to volume, same dimension",
			qty:  d("1.5"),
			from: liter, to: milliliter, ing: unknown,
			want: decimal.NewFromInt(1500),
		},
		{
			name: "count to count, same dimension",
			qty:  decimal.NewFromInt(24),
			from: piece, to: dozen, ing: unknown,
			want: decimal.NewFromInt(2),
		},
		{
			name: "volume to mass via density 1 (water)",
			qty:  decimal.NewFromInt(500),
			from: milliliter, to: gram, ing: water,
			want: decimal.NewFromInt(500),
		},
		{
			name: "mass to volume via density 1 (water)",
			qty:  decimal.NewFromInt(500),
			from: gram, to: milliliter, ing: water,
			want: decimal.NewFromInt(500),
		},
		{
			name: "volume to mass via density < 1 (olive oil)",
			qty:  decimal.NewFromInt(100),
			from: milliliter, to: gram, ing: oliveOil,
			want: d("92"),
		},
		{
			name: "mass to volume via density < 1 (olive oil)",
			qty:  d("92"),
			from: gram, to: milliliter, ing: oliveOil,
			want: decimal.NewFromInt(100),
		},
		{
			name: "volume to mass, missing density",
			qty:  decimal.NewFromInt(100),
			from: milliliter, to: gram, ing: unknown,
			wantErr: true,
		},
		{
			name: "count to mass via grams-per-piece (egg)",
			qty:  decimal.NewFromInt(3),
			from: piece, to: gram, ing: egg,
			want: decimal.NewFromInt(150),
		},
		{
			name: "mass to count via grams-per-piece (egg)",
			qty:  decimal.NewFromInt(300),
			from: gram, to: piece, ing: egg,
			want: decimal.NewFromInt(6),
		},
		{
			name: "mass to count, missing grams-per-piece",
			qty:  decimal.NewFromInt(300),
			from: gram, to: piece, ing: unknown,
			wantErr: true,
		},
		{
			name: "count to volume needs both density and grams-per-piece",
			qty:  decimal.NewFromInt(2),
			from: piece, to: milliliter, ing: eggAndOil,
			want: decimal.NewFromInt(100),
		},
		{
			name: "volume to count needs both density and grams-per-piece",
			qty:  decimal.NewFromInt(100),
			from: milliliter, to: piece, ing: eggAndOil,
			want: decimal.NewFromInt(2),
		},
		{
			name: "count to volume, missing density only",
			qty:  decimal.NewFromInt(2),
			from: piece, to: milliliter, ing: egg,
			wantErr: true,
		},
		{
			name: "count to volume, missing grams-per-piece only",
			qty:  decimal.NewFromInt(2),
			from: piece, to: milliliter, ing: oliveOil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert(tt.qty, tt.from, tt.to, tt.ing)

			if tt.wantErr {
				if err != ErrNoConversion {
					t.Fatalf("expected ErrNoConversion, got %v (result %s)", err, got)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Equal(tt.want) {
				t.Fatalf("got %s, want %s", got, tt.want)
			}
		})
	}
}

// TestConvert_DefaultUnits sanity-checks the real seed data (KIT-37) -
// not exact real-world unit equivalences (those are empirically rounded
// constants and won't divide evenly), just that converting there and
// back lands within a cent's worth of rounding.
func TestConvert_DefaultUnits(t *testing.T) {
	byShortName := func(short string) Unit {
		for _, u := range DefaultUnits {
			if u.ShortName == short {
				return u
			}
		}
		t.Fatalf("no default unit with short name %q", short)
		return Unit{}
	}

	tbsp := byShortName("tbsp")
	ml := byShortName("ml")

	got, err := Convert(decimal.NewFromInt(2), tbsp, ml, unknown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := d("29.5736")
	diff := got.Sub(want).Abs()
	if diff.GreaterThan(d("0.01")) {
		t.Fatalf("2 tbsp in ml: got %s, want approximately %s", got, want)
	}
}
