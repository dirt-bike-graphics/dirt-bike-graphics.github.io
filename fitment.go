// Package fitment provides dirt bike plastic and graphics kit fitment lookup helpers
// across 175 model years, 76 models and 11 brands (motocross, enduro and electric).
//
// Maintained from the open fitment catalog powering [Custom Dirt Bike Graphics] and its
// interactive [dirt bike graphics designer].
//
// [Custom Dirt Bike Graphics]: https://customdirtbikegraphics.com
// [dirt bike graphics designer]: https://customdirtbikegraphics.com/design?utm_source=pkg_go_dev&utm_medium=referral&utm_campaign=go
package fitment

import (
	_ "embed"
	"encoding/json"
	"strings"
	"sync"
)

// RawJSON holds the embedded fitment.json catalog bytes.
//
//go:embed fitment.json
var RawJSON []byte

// BikeFitment describes one OEM plastic generation shared across sister models and years.
type BikeFitment struct {
	ID            string   `json:"id"`
	Brand         string   `json:"brand"`
	Model         string   `json:"model"`
	ReferenceYear *string  `json:"reference_year"`
	Category      string   `json:"category"`
	Fits          []string `json:"fits"`
}

type catalogFile struct {
	Updated string        `json:"updated"`
	Bikes   []BikeFitment `json:"bikes"`
}

var (
	once    sync.Once
	catalog catalogFile
)

func load() {
	once.Do(func() {
		_ = json.Unmarshal(RawJSON, &catalog)
	})
}

func norm(s string) string {
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return strings.ToUpper(b.String())
}

// Bikes returns all 17 reference dirt bike fitment groups.
func Bikes() []BikeFitment {
	load()
	out := make([]BikeFitment, len(catalog.Bikes))
	copy(out, catalog.Bikes)
	return out
}

// FindKit returns the reference bike whose graphics kit fits model and optional year (pass "" for any year).
func FindKit(model, year string) *BikeFitment {
	load()
	m := norm(model)
	if m == "" {
		return nil
	}
	y := strings.TrimSpace(year)
	for i := range catalog.Bikes {
		bike := &catalog.Bikes[i]
		for _, fit := range bike.Fits {
			idx := strings.IndexByte(fit, ' ')
			if idx <= 0 {
				continue
			}
			fYear, fModel := fit[:idx], fit[idx+1:]
			if norm(fModel) == m && (y == "" || fYear == y) {
				copyBike := *bike
				return &copyBike
			}
		}
	}
	return nil
}

// SharedFits returns every model year sharing the same plastics with model and optional year.
func SharedFits(model, year string) []string {
	b := FindKit(model, year)
	if b == nil {
		return nil
	}
	out := make([]string, len(b.Fits))
	copy(out, b.Fits)
	return out
}

// SameKit reports whether two model/year pairs share the exact same graphics kit layout.
func SameKit(modelA, yearA, modelB, yearB string) bool {
	a := FindKit(modelA, yearA)
	b := FindKit(modelB, yearB)
	return a != nil && b != nil && a.ID == b.ID
}
