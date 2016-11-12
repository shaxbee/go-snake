package snake

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLineContainsPoint(t *testing.T) {
	contained := []struct {
		l Line
		p Vec
	}{
		{Line{Vec{0, 0}, Vec{1, 1}}, Vec{0.0, 0.0}},
		{Line{Vec{0, 0}, Vec{1, 1}}, Vec{1.0, 1.0}},
		{Line{Vec{0, 0}, Vec{1, 1}}, Vec{0.5, 0.5}},
	}
	for i, v := range contained {
		assert.True(t, v.l.ContainsPoint(v.p), "Expected line %d to contain point", i)
	}

	excluded := []struct {
		l Line
		p Vec
	}{
		{Line{Vec{0, 0}, Vec{1, 1}}, Vec{1.1, 1.1}},
		{Line{Vec{0, 0}, Vec{1, 1}}, Vec{-1.0, -1.0}},
	}
	for i, v := range excluded {
		assert.False(t, v.l.ContainsPoint(v.p), "Expected line %d to exclude point", i)
	}

}

func TestLineIntersectLine(t *testing.T) {
	intersecting := []struct {
		a, b Line
		p    Vec
	}{
		{Line{Vec{0.0, 0.0}, Vec{1.0, 1.0}}, Line{Vec{0.5, 0.0}, Vec{0.5, 1.0}}, Vec{0.5, 0.5}},
		{Line{Vec{0.0, 0.0}, Vec{1.0, 1.0}}, Line{Vec{1.0, 0.0}, Vec{1.0, 1.0}}, Vec{1.0, 1.0}},
		{Line{Vec{0.0, 0.0}, Vec{1.0, 1.0}}, Line{Vec{1.0, 0.0}, Vec{0.0, 1.0}}, Vec{0.5, 0.5}},
	}
	for i, v := range intersecting {
		p, ok := v.a.IntersectLine(v.b)
		if assert.True(t, ok, "Expected lines %d to intersect", i) {
			assert.Equal(t, v.p, p, "Expected lines %d to intersect at point %v", i, v.p)
		}
	}

	disjoined := []struct {
		a, b Line
	}{
		{Line{Vec{0.0, 1.0}, Vec{1.0, 1.0}}, Line{Vec{0.0, 0.0}, Vec{1.0, 0.0}}},
	}
	for i, v := range disjoined {
		_, ok := v.a.IntersectLine(v.b)
		assert.False(t, ok, "Expected lines %i to be disjoined", i)
	}
}

func TestArcIntersectLine(t *testing.T) {

}

func TestArcIntersectArc(t *testing.T) {

}
