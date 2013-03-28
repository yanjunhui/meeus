// Copyright 2013 Sonia Keys
// License MIT: http://www.opensource.org/licenses/MIT

// Precession: Chapter 21, Precession.
//
// Functions in this package take Julian epoch argurments rather than Julian
// days.  Use base.JDEToJulianYear() to convert.
//
// Also in package base are some definitions related to the Besselian and
// Julian Year.
//
// Proper motion units
//
// Meeus gives some example annual proper motions in units of seconds of
// right ascension and seconds of declination.  To make units clear,
// functions in this package take proper motions with argument types of
// base.HourAngle and base.Angle respectively.  Error-prone conversion
// factors conversion factors can be avoided by using the constructors
// for these base types.
//
// For example, rather than
//
//	mra := -0.03847 / 13751 // as Meeus suggests
//
// or
//
//	mra := -0.03847 * (15/3600) * (pi/180) // less magic
//
// use
//
//	mra := base.NewHourAngle(false, 0, 0, -0.03847)
//
// In all cases, proper motions are expected to be annual proper motions
// so the unit denominator is years.  (The code, following Meeus's example,
// technically treats it as Julian years.)
package precess

import (
	"math"

	"github.com/soniakeys/meeus/base"
	"github.com/soniakeys/meeus/coord"
	"github.com/soniakeys/meeus/nutation"
)

// ApproxAnnualPrecession returns approximate annual precision in right
// ascension and declination.
//
// The two epochs should be within a few hundred years.
// The declinations should not be too close to the poles.
func ApproxAnnualPrecession(eq *coord.Equatorial, epochFrom, epochTo float64) (Δα base.HourAngle, Δδ base.Angle) {
	m, na, nd := mn(epochFrom, epochTo)
	sa, ca := math.Sincos(eq.RA)
	Δαs := m + na*sa*math.Tan(eq.Dec) // seconds of RA
	Δδs := nd * ca                    // seconds of Dec
	return base.NewHourAngle(false, 0, 0, Δαs), base.NewAngle(false, 0, 0, Δδs)
}

// mn as separate function for testing purposes
func mn(epochFrom, epochTo float64) (m, na, nd float64) {
	T := (epochTo - epochFrom) * .01
	m = 3.07496 + 0.00186*T
	na = 1.33621 - 0.00057*T
	nd = 20.0431 - 0.0085*T
	return
}

// ApproxPosition uses ApproxAnnualPrecession to compute a simple and quick
// precession while still considering proper motion.
//
// Both eqFrom and eqTo must be non-nil, although they may point to the same
// struct.  EqTo is returned for convenience.
func ApproxPosition(eqFrom, eqTo *coord.Equatorial, epochFrom, epochTo float64, mα base.HourAngle, mδ base.Angle) *coord.Equatorial {
	Δα, Δδ := ApproxAnnualPrecession(eqFrom, epochFrom, epochTo)
	dy := epochTo - epochFrom
	eqTo.RA = eqFrom.RA + (Δα+mα).Rad()*dy
	eqTo.Dec = eqFrom.Dec + (Δδ+mδ).Rad()*dy
	return eqTo
}

// Precessor represents precession from one epoch to another.
//
// Construct with NewPrecessor, then call method Precess.
// After construction, Precess may be called multiple times to precess
// different coordinates with the same initial and final epochs.
type Precessor struct {
	ζ, z, sθ, cθ float64
}

// Package variables allow these slices to be reused.  (As composite
// literals inside of NewPrecessor they would be reallocated on every
// function call.)
var (
	ζt = []float64{2306.2181, 0.30188, 0.017998}
	zt = []float64{2306.2181, 1.09468, 0.018203}
	θt = []float64{2004.3109, -0.42665, -0.041833}

	ζT = []float64{2306.2181, 1.39656, -0.000139}
	zT = []float64{2306.2181, 1.39656, -0.000139}
	θT = []float64{2004.3109, -0.8533, -0.000217}
)

// NewPrecessor constructs a Precessor object and initializes it to precess
// coordinates from epochFrom to epochTo.
func NewPrecessor(epochFrom, epochTo float64) *Precessor {
	ζCoeff := ζt
	zCoeff := zt
	θCoeff := θt
	if epochFrom != 2000 {
		T := (epochFrom - 2000) * .01
		ζCoeff = []float64{
			base.Horner(T, ζT...),
			0.30188 - 0.000344*T,
			0.017998}
		zCoeff = []float64{
			base.Horner(T, zT...),
			1.09468 - 0.000066*T,
			0.018203}
		θCoeff = []float64{
			base.Horner(T, θT...),
			-0.42665 - 0.000217*T,
			-0.041833}
	}
	t := (epochTo - epochFrom) * .01
	p := &Precessor{
		ζ: base.Horner(t, ζCoeff...) * t * math.Pi / 180 / 3600,
		z: base.Horner(t, zCoeff...) * t * math.Pi / 180 / 3600,
	}
	θ := base.Horner(t, θCoeff...) * t * math.Pi / 180 / 3600
	p.sθ, p.cθ = math.Sincos(θ)
	return p
}

// Precess precesses coordinates eqFrom, leaving result in eqTo.
//
// The same struct may be used for eqFrom and eqTo.
// EqTo is returned for convenience.
func (p *Precessor) Precess(eqFrom, eqTo *coord.Equatorial) *coord.Equatorial {
	sδ, cδ := math.Sincos(eqFrom.Dec)
	sαζ, cαζ := math.Sincos(eqFrom.RA + p.ζ)
	A := cδ * sαζ
	B := p.cθ*cδ*cαζ - p.sθ*sδ
	C := p.sθ*cδ*cαζ + p.cθ*sδ
	eqTo.RA = math.Atan2(A, B) + p.z
	if C < base.CosSmallAngle {
		eqTo.Dec = math.Asin(C)
	} else {
		eqTo.Dec = math.Acos(math.Sqrt(A*A + B*B)) // near pole
	}
	return eqTo
}

// Position precesses equatorial coordinates from one epoch to another,
// including proper motions.
//
// If proper motions are not to be considered or are not applicable, pass 0, 0
// for mα, mδ
//
// Both eqFrom and eqTo must be non-nil, although they may point to the same
// struct.  EqTo is returned for convenience.
func Position(eqFrom, eqTo *coord.Equatorial, epochFrom, epochTo float64, mα base.HourAngle, mδ base.Angle) *coord.Equatorial {
	p := NewPrecessor(epochFrom, epochTo)
	t := epochTo - epochFrom
	eqTo.RA = eqFrom.RA + mα.Rad()*t
	eqTo.Dec = eqFrom.Dec + mδ.Rad()*t
	return p.Precess(eqTo, eqTo)
}

// EclipticPrecessor represents precession from one epoch to another.
//
// Construct with NewEclipticPrecessor, then call method Precess.
// After construction, Precess may be called multiple times to precess
// different coordinates with the same initial and final epochs.
type EclipticPrecessor struct {
	sη, cη, π, p float64
}

var (
	ηt = []float64{0, 47.0029, -0.03302, 0.000060}
	πt = []float64{3600 * 174.876384, -869.8089, 0.03536}
	pt = []float64{0, 5029.0966, 1.11113, -0.000006}

	ηT = []float64{74.0029, -0.06603, 0.000598}
	πT = []float64{3600 * 174.876384, 3289.4789, 0.60622}
	pT = []float64{5029.0966, 2.22226, -0.000042}
)

// NewEclipticPrecessor constructs an EclipticPrecessor object and initializes
// it to precess coordinates from epochFrom to epochTo.
func NewEclipticPrecessor(epochFrom, epochTo float64) *EclipticPrecessor {
	ηCoeff := ηt
	πCoeff := πt
	pCoeff := pt
	if epochFrom != 2000 {
		T := (epochFrom - 2000) * .01
		ηCoeff = []float64{0,
			base.Horner(T, ηT...),
			-0.03302 + 0.000598*T,
			0.000060}
		πCoeff = []float64{
			base.Horner(T, πT...),
			-869.8089 - 0.50491*T,
			0.03536}
		pCoeff = []float64{0,
			base.Horner(T, pT...),
			1.11113 - 0.000042*T,
			-0.000006}
	}
	t := (epochTo - epochFrom) * .01
	p := &EclipticPrecessor{
		π: base.NewAngle(false, 0, 0, base.Horner(t, πCoeff...)).Rad(),
		p: base.NewAngle(false, 0, 0, base.Horner(t, pCoeff...)).Rad(),
	}
	η := base.NewAngle(false, 0, 0, base.Horner(t, ηCoeff...)).Rad()
	p.sη, p.cη = math.Sincos(η)
	return p
}

// EclipticPrecess precesses coordinates eclFrom, leaving result in eclTo.
//
// The same struct may be used for eclFrom and eclTo.
// EclTo is returned for convenience.
func (p *EclipticPrecessor) Precess(eclFrom, eclTo *coord.Ecliptic) *coord.Ecliptic {
	sβ, cβ := math.Sincos(eclFrom.Lat)
	sd, cd := math.Sincos(p.π - eclFrom.Lon)
	A := p.cη*cβ*sd - p.sη*sβ
	B := cβ * cd
	C := p.cη*sβ + p.sη*cβ*sd
	eclTo.Lon = p.p + p.π - math.Atan2(A, B)
	if C < base.CosSmallAngle {
		eclTo.Lat = math.Asin(C)
	} else {
		eclTo.Lat = math.Acos(math.Sqrt(A*A + B*B)) // near pole
	}
	return eclTo
}

// EclipticPosition precesses ecliptic coordinates from one epoch to another,
// including proper motions.
//
// While eclFrom is given as ecliptic coordinates, proper motions mα, mδ are
// still expected to be equatorial.  If proper motions are not to be considered
// or are not applicable, pass 0, 0.
//
// Both eclFrom and eclTo must be non-nil, although they may point to the same
// struct.  EclTo is returned for convenience.
func EclipticPosition(eclFrom, eclTo *coord.Ecliptic, epochFrom, epochTo float64, mα base.HourAngle, mδ base.Angle) *coord.Ecliptic {
	p := NewEclipticPrecessor(epochFrom, epochTo)
	*eclTo = *eclFrom
	if mα != 0 || mδ != 0 {
		mλ, mβ := eqProperMotionToEcl(mα.Rad(), mδ.Rad(), epochFrom, eclFrom)
		t := epochTo - epochFrom
		eclTo.Lon += mλ * t
		eclTo.Lat += mβ * t
	}
	return p.Precess(eclTo, eclTo)
}

func eqProperMotionToEcl(mα, mδ, epoch float64, pos *coord.Ecliptic) (mλ, mβ float64) {
	ε := nutation.MeanObliquity(base.JulianYearToJDE(epoch))
	sε, cε := math.Sincos(ε)
	eqPos := new(coord.Equatorial).EclToEq(pos, sε, cε)
	sα, cα := math.Sincos(eqPos.RA)
	sδ, cδ := math.Sincos(eqPos.Dec)
	cβ := math.Cos(pos.Lat)
	mλ = (mδ*sε*cα + mα*cδ*(cε*cδ+sε*sδ*sα)) / (cβ * cβ)
	mβ = (mδ*(cε*cδ+sε*sδ*sα) - mα*sε*cα*cδ) / cβ
	return
}

// ProperMotion3D takes the 3D equatorial coordinates of an object
// at one epoch and computes its coordinates at a new epoch, considering
// proper motion and radial velocity.
//
// Radial distance (r) must be in parsecs, radial velocitiy (mr) in
// parsecs per year.
//
// Both eqFrom and eqTo must be non-nil, although they may point to the same
// struct.  EqTo is returned for convenience.
func ProperMotion3D(eqFrom, eqTo *coord.Equatorial, epochFrom, epochTo, r, mr float64, mα base.HourAngle, mδ base.Angle) *coord.Equatorial {
	sα, cα := math.Sincos(mα.Rad())
	sδ, cδ := math.Sincos(mδ.Rad())
	x := r * cδ * cα
	y := r * cδ * sα
	z := r * sδ
	mrr := mr / r
	zmδ := z * mδ.Rad()
	mx := x*mrr - zmδ*cα - y*mα.Rad()
	my := y*mrr - zmδ*sα + x*mα.Rad()
	mz := z*mrr + r*mδ.Rad()*cδ
	t := epochTo - epochFrom
	xp := x + t*mx
	yp := y + t*my
	zp := z + t*mz
	eqTo.RA = eqFrom.RA + math.Atan2(yp, xp)
	eqTo.Dec = eqFrom.Dec + math.Atan2(zp, math.Hypot(xp, yp))
	return eqTo
}
