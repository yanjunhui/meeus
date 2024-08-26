// Copyright 2013 Sonia Keys
// License: MIT

// Parallax: Chapter 40, Correction for Parallax.
package parallax

import (
	"math"

	"github.com/yanjunhui/meeus/globe"
	"github.com/yanjunhui/meeus/sidereal"
	"github.com/yanjunhui/unit"
)

// constant for Horizontal.  p. 279.
var hp = unit.AngleFromSec(8.794)

// Horizontal returns equatorial horizontal parallax of a body.
//
// Argument Δ is distance in AU.
//
// Meeus mentions use of this function for the Moon, Sun, planet, or comet.
// That is, for relatively distant objects.  For parallax of the Moon (or
// other relatively close object) see moonposition.Parallax.
func Horizontal(Δ float64) (π unit.Angle) {
	return hp.Div(Δ) // (40.1) p. 279
}

// Topocentric returns topocentric positions including parallax.
//
// Arguments α, δ are geocentric right ascension and declination in radians.
// Δ is distance to the observed object in AU.  ρsφʹ, ρcφʹ are parallax
// constants (see package globe.) L is geographic longitude of the observer,
// jde is time of observation.
//
// Results are observed topocentric ra and dec in radians.
func Topocentric(α unit.RA, δ unit.Angle, Δ, ρsφʹ, ρcφʹ float64, L unit.Angle, jde float64) (αʹ unit.RA, δʹ unit.Angle) {
	π := Horizontal(Δ)
	θ0 := sidereal.Apparent(jde)
	H := (θ0.Angle() - L - unit.Angle(α)).Mod1()
	sπ := π.Sin()
	sH, cH := H.Sincos()
	sδ, cδ := δ.Sincos()
	// (40.2) p. 279
	Δα := unit.HourAngle(math.Atan2(-ρcφʹ*sπ*sH, cδ-ρcφʹ*sπ*cH))
	αʹ = α.Add(Δα)
	// (40.3) p. 279
	δʹ = unit.Angle(math.Atan2((sδ-ρsφʹ*sπ)*Δα.Cos(), cδ-ρcφʹ*sπ*cH))
	return
}

// Topocentric2 returns topocentric corrections including parallax.
//
// This function implements the "non-rigorous" method descripted in the text.
//
// Note that results are corrections, not corrected coordinates.
func Topocentric2(α unit.RA, δ unit.Angle, Δ, ρsφʹ, ρcφʹ float64, L unit.Angle, jde float64) (Δα unit.HourAngle, Δδ unit.Angle) {
	π := Horizontal(Δ)
	θ0 := sidereal.Apparent(jde)
	H := (θ0.Angle() - L - unit.Angle(α)).Mod1()
	sH, cH := H.Sincos()
	sδ, cδ := δ.Sincos()
	Δα = unit.HourAngle(-π.Mul(ρcφʹ * sH / cδ)) // (40.4) p. 280
	Δδ = -π.Mul(ρsφʹ*cδ - ρcφʹ*cH*sδ)           // (40.5) p. 280
	return
}

// Topocentric3 returns topocentric hour angle and declination including parallax.
//
// This function implements the "alternative" method described in the text.
// The method should be similarly rigorous to that of Topocentric() and results
// should be virtually consistent.
func Topocentric3(α unit.RA, δ unit.Angle, Δ, ρsφʹ, ρcφʹ float64, L unit.Angle, jde float64) (Hʹ unit.HourAngle, δʹ unit.Angle) {
	π := Horizontal(Δ)
	θ0 := sidereal.Apparent(jde)
	H := (θ0.Angle() - L - unit.Angle(α)).Mod1()
	sπ := π.Sin()
	sH, cH := H.Sincos()
	sδ, cδ := δ.Sincos()
	A := cδ * sH
	B := cδ*cH - ρcφʹ*sπ
	C := sδ - ρsφʹ*sπ
	q := math.Sqrt(A*A + B*B + C*C)
	Hʹ = unit.HourAngle(math.Atan2(A, B))
	δʹ = unit.Angle(math.Asin(C / q))
	return
}

// TopocentricEcliptical returns topocentric ecliptical coordinates including parallax.
//
// Arguments λ, β are geocentric ecliptical longitude and latitude of a body,
// s is its geocentric semidiameter. φ, h are the observer's latitude and
// and height above the ellipsoid in meters.  ε is the obliquity of the
// ecliptic, θ is local sidereal time, π is equatorial horizontal parallax
// of the body (see Horizonal()).
//
// Results are observed topocentric coordinates and semidiameter.
func TopocentricEcliptical(λ, β, s, φ unit.Angle, h float64, ε unit.Angle, θ unit.Time, π unit.Angle) (λʹ, βʹ, sʹ unit.Angle) {
	S, C := globe.Earth76.ParallaxConstants(φ, h)
	sλ, cλ := λ.Sincos()
	sβ, cβ := β.Sincos()
	sε, cε := ε.Sincos()
	sθ, cθ := θ.Angle().Sincos()
	sπ := π.Sin()
	N := cλ*cβ - C*sπ*cθ
	λʹ = unit.Angle(math.Atan2(sλ*cβ-sπ*(S*sε+C*cε*sθ), N))
	if λʹ < 0 {
		λʹ += 2 * math.Pi
	}
	cλʹ := λʹ.Cos()
	βʹ = unit.Angle(math.Atan(cλʹ * (sβ - sπ*(S*cε-C*sε*sθ)) / N))
	sʹ = unit.Angle(math.Asin(cλʹ * βʹ.Cos() * s.Sin() / N))
	return
}
