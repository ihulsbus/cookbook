package models

import (
	"errors"

	"github.com/shopspring/decimal"
)

// ErrNoConversion is returned when crossing dimensions needs an
// Ingredient field (DensityGPerMl or GramsPerPiece) that is nil. The
// caller's fallback is to keep showing the original unit.
var ErrNoConversion = errors.New("no conversion available")

// Convert converts qty of from into the equivalent quantity of to.
// Within one Dimension, each Unit's ToBase factor (its conversion to
// that Dimension's base unit - g, ml or piece) is enough. Crossing
// dimensions (mass<->volume<->count) additionally routes through ing's
// DensityGPerMl and/or GramsPerPiece; if the one needed is nil, this
// returns ErrNoConversion.
func Convert(qty decimal.Decimal, from, to Unit, ing Ingredient) (decimal.Decimal, error) {
	// qty in from's Dimension's base unit - the common starting point for
	// both paths below.
	baseQty := qty.Mul(from.ToBase)

	if from.Dimension == to.Dimension {
		return baseQty.Div(to.ToBase), nil
	}

	var baseToQty decimal.Decimal
	switch {
	case from.Dimension == DimensionMass && to.Dimension == DimensionVolume:
		if ing.DensityGPerMl == nil {
			return decimal.Zero, ErrNoConversion
		}
		baseToQty = baseQty.Div(*ing.DensityGPerMl)
	case from.Dimension == DimensionVolume && to.Dimension == DimensionMass:
		if ing.DensityGPerMl == nil {
			return decimal.Zero, ErrNoConversion
		}
		baseToQty = baseQty.Mul(*ing.DensityGPerMl)
	case from.Dimension == DimensionCount && to.Dimension == DimensionMass:
		if ing.GramsPerPiece == nil {
			return decimal.Zero, ErrNoConversion
		}
		baseToQty = baseQty.Mul(*ing.GramsPerPiece)
	case from.Dimension == DimensionMass && to.Dimension == DimensionCount:
		if ing.GramsPerPiece == nil {
			return decimal.Zero, ErrNoConversion
		}
		baseToQty = baseQty.Div(*ing.GramsPerPiece)
	case from.Dimension == DimensionCount && to.Dimension == DimensionVolume:
		if ing.GramsPerPiece == nil || ing.DensityGPerMl == nil {
			return decimal.Zero, ErrNoConversion
		}
		baseToQty = baseQty.Mul(*ing.GramsPerPiece).Div(*ing.DensityGPerMl)
	case from.Dimension == DimensionVolume && to.Dimension == DimensionCount:
		if ing.GramsPerPiece == nil || ing.DensityGPerMl == nil {
			return decimal.Zero, ErrNoConversion
		}
		baseToQty = baseQty.Mul(*ing.DensityGPerMl).Div(*ing.GramsPerPiece)
	default:
		return decimal.Zero, ErrNoConversion
	}

	return baseToQty.Div(to.ToBase), nil
}
