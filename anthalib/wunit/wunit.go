// wunit/wunit.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 1 Royal College St, London NW1 0NH UK

package wunit

import (
	"fmt"
)

// structure defining a base unit
type BaseUnit interface {
	// unit name
	Name() string
	// unit symbol
	Symbol() string
	// multiply by this to get SI value
	// nb this should be a function since we actually need
	// an affine transformation
	BaseSIConversionFactor() float64 // this can be calculated in many cases
	// if we convert to the SI units what is the appropriate unit symbol
	BaseSIUnit() string // if we use the above, what unit do we get?
	// print this
	ToString() string
}

// a unit with an SI prefix
type PrefixedUnit interface {
	BaseUnit
	// the prefix of the unit
	Prefix() SIPrefix
	// the symbol including prefix
	PrefixedSymbol() string
	// the symbol excluding prefix
	RawSymbol() string
	// appropriate unit if we ask for SI values
	BaseSISymbol() string
	// returns conversion factor from *this* unit to the other
	ConvertTo(pu PrefixedUnit) float64
}

// fundamental representation of a value in the system
type Measurement interface {
	// the value in base SI units
	SIValue() float64
	// the value in the current units
	RawValue() float64
	// unit plus prefix
	Unit() PrefixedUnit
	// set the value, this must be thread-safe
	// returns old value
	SetValue(v float64) float64
	// convert units
	ConvertTo(p PrefixedUnit) float64
	// add to this measurement
	Add(m Measurement)
	// subtract from this measurement
	Subtract(m Measurement)
	// A nice string representation
	ToString() string
}

// structure implementing the Measurement interface
type ConcreteMeasurement struct {
	// the raw value
	Mvalue float64
	// the relevant units
	Munit *GenericPrefixedUnit
}

// value when converted to SI units
func (cm *ConcreteMeasurement) SIValue() float64 {
	return cm.Mvalue * cm.Munit.BaseSIConversionFactor()
}

// value without conversion
func (cm *ConcreteMeasurement) RawValue() float64 {
	return cm.Mvalue
}

// get unit with prefix
func (cm *ConcreteMeasurement) Unit() PrefixedUnit {
	return cm.Munit
}

// set the value of this measurement
func (cm *ConcreteMeasurement) SetValue(v float64) float64 {
	cm.Mvalue = v
	return v
}

// convert to a different unit
// nb this is NOT destructive
func (cm *ConcreteMeasurement) ConvertTo(p PrefixedUnit) float64 {
	return cm.Unit().ConvertTo(p) * cm.RawValue()
}

// add to this

func (cm *ConcreteMeasurement) Add(m Measurement) {
	// ideally should check these have the same Dimension
	// need to improve this

	cm.SetValue(m.ConvertTo(cm.Unit()) + cm.RawValue())
}

// subtract

func (cm *ConcreteMeasurement) Subtract(m Measurement) {
	// ideally should check these have the same Dimension
	// need to improve this

	cm.SetValue(cm.RawValue() - m.ConvertTo(cm.Unit()))
}

func (cm *ConcreteMeasurement) ToString() string {
	return fmt.Sprintf("%-6.3f%s", cm.RawValue(), cm.Unit().PrefixedSymbol())
}

/**********/

func NewPMeasurement(v float64, pu string) ConcreteMeasurement {
	return ConcreteMeasurement{v, ParsePrefixedUnit(pu)}
}

// helper function for creating a new measurement
func NewMeasurement(v float64, prefix string, unit string) ConcreteMeasurement {
	gpu := NewPrefixedUnit(prefix, unit)
	return ConcreteMeasurement{v, gpu}
}
