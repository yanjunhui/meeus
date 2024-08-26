// Copyright 2013 Sonia Keys
// License: MIT

// Saturnrings: Chapter 45, The Ring of Saturn
package saturnring

import (
	"math"

	"github.com/yanjunhui/meeus/base"
	"github.com/yanjunhui/meeus/coord"
	"github.com/yanjunhui/meeus/nutation"
	pp "github.com/yanjunhui/meeus/planetposition"
	"github.com/yanjunhui/meeus/unit"
)

// Constants for scaling aEdge and bEdge.
const (
	InnerEdgeOfOuter = .8801
	OuterEdgeOfInner = .8599
	InnerEdgeOfInner = .6650
	InnerEdgeOfDusky = .5486
)

// Ring computes quantities of the ring of Saturn.
//
//	B  Saturnicentric latitude of the Earth referred to the plane of the ring.
//	Bʹ  Saturnicentric latitude of the Sun referred to the plane of the ring.
//	ΔU  Difference between Saturnicentric longitudes of the Sun and the Earth.
//	P  Geometric position angle of the northern semiminor axis of the ring.
//	aEdge  Major axis of the out edge of the outer ring.
//	bEdge  Minor axis of the out edge of the outer ring.
func Ring(jde float64, earth, saturn *pp.V87Planet) (B, Bʹ, ΔU, P, aEdge, bEdge unit.Angle) {
	f1, f2 := cl(jde, earth, saturn)
	ΔU, B = f1()
	Bʹ, P, aEdge, bEdge = f2()
	return
}

// UB computes quantities required by illum.Saturn().
//
// Same as ΔU and B returned by Ring().  Results in radians.
func UB(jde float64, earth, saturn *pp.V87Planet) (ΔU, B unit.Angle) {
	f1, _ := cl(jde, earth, saturn)
	return f1()
}

// cl splits the work into two closures.
func cl(jde float64, earth, saturn *pp.V87Planet) (f1 func() (ΔU, B unit.Angle),
	f2 func() (Bʹ, P, aEdge, bEdge unit.Angle)) {
	const p = math.Pi / 180
	var i, Ω unit.Angle
	var l0, b0 unit.Angle
	var R float64
	Δ := 9.
	var λ, β unit.Angle
	var si, ci, sβ, cβ, sB float64
	var sbʹ, cbʹ, slʹΩ, clʹΩ float64
	f1 = func() (ΔU, B unit.Angle) {
		// (45.1), p. 318
		T := base.J2000Century(jde)
		i = unit.AngleFromDeg(base.Horner(T, 28.075216, -.012998, .000004))
		Ω = unit.AngleFromDeg(base.Horner(T, 169.50847, 1.394681, .000412))
		// Step 2.
		l0, b0, R = earth.Position(jde)
		l0, b0 = pp.ToFK5(l0, b0, jde)
		sl0, cl0 := l0.Sincos()
		sb0 := b0.Sin()
		// Steps 3, 4.
		var l, b unit.Angle
		var r, x, y, z float64
		f := func() {
			τ := base.LightTime(Δ)
			l, b, r = saturn.Position(jde - τ)
			l, b = pp.ToFK5(l, b, jde)
			sl, cl := l.Sincos()
			sb, cb := b.Sincos()
			x = r*cb*cl - R*cl0
			y = r*cb*sl - R*sl0
			z = r*sb - R*sb0
			Δ = math.Sqrt(x*x + y*y + z*z)
		}
		f()
		f()
		// Step 5.
		λ = unit.Angle(math.Atan2(y, x))
		β = unit.Angle(math.Atan(z / math.Hypot(x, y)))
		// First part of step 6.
		si, ci = i.Sincos()
		sβ, cβ = β.Sincos()
		sB = si*cβ*(λ-Ω).Sin() - ci*sβ
		B = unit.Angle(math.Asin(sB)) // return value
		// Step 7.
		N := unit.AngleFromDeg(113.6655 + .8771*T)
		lʹ := l - unit.AngleFromDeg(.01759).Div(r)
		bʹ := b - unit.AngleFromDeg(.000764).Mul((l-N).Cos()/r)
		// Setup for steps 8, 9.
		sbʹ, cbʹ = bʹ.Sincos()
		slʹΩ, clʹΩ = (lʹ - Ω).Sincos()
		// Step 9.
		sλΩ, cλΩ := (λ - Ω).Sincos()
		U1 := unit.Angle(math.Atan2(si*sbʹ+ci*cbʹ*slʹΩ, cbʹ*clʹΩ))
		U2 := unit.Angle(math.Atan2(si*sβ+ci*cβ*sλΩ, cβ*cλΩ))
		ΔU = unit.Angle(math.Abs((U1 - U2).Rad())) // return value
		return
	}
	f2 = func() (Bʹ, P, aEdge, bEdge unit.Angle) {
		// Remainder of step 6.
		aEdge = unit.AngleFromSec(375.35).Div(Δ) // return value
		bEdge = aEdge.Mul(math.Abs(sB))          // return value
		// Step 8.
		sBʹ := si*cbʹ*slʹΩ - ci*sbʹ
		Bʹ = unit.Angle(math.Asin(sBʹ)) // return value
		// Step 10.
		Δψ, Δε := nutation.Nutation(jde)
		ε := nutation.MeanObliquity(jde) + Δε
		// Step 11.
		λ0 := Ω - math.Pi/2
		β0 := math.Pi/2 - i
		// Step 12.
		sl0λ, cl0λ := (l0 - λ).Sincos()
		λ += unit.AngleFromDeg(.005693).Mul(cl0λ / cβ)
		β += unit.AngleFromDeg(.005693).Mul(sl0λ * sβ)
		// Step 13.
		λ0 += Δψ
		λ += Δψ
		// Step 14.
		sε, cε := ε.Sincos()
		α0, δ0 := coord.EclToEq(λ0, β0, sε, cε)
		α, δ := coord.EclToEq(λ, β, sε, cε)
		// Step 15.
		sδ0, cδ0 := δ0.Sincos()
		sδ, cδ := δ.Sincos()
		sα0α, cα0α := (α0 - α).Sincos()
		P = unit.Angle(math.Atan2(cδ0*sα0α, sδ0*cδ-cδ0*sδ*cα0α)) // return value
		return
	}
	return
}
