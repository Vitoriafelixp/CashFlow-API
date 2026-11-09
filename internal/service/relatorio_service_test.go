package service

import (
	"testing"
	"time"
)

func TestInicioDoDiaPreservaDataLocal(t *testing.T) {
	local := time.FixedZone("UTC-03", -3*60*60)
	data := time.Date(2026, time.October, 1, 21, 15, 0, 0, local)

	inicio := inicioDoDia(data)
	if inicio.Year() != 2026 || inicio.Month() != time.October || inicio.Day() != 1 {
		t.Fatalf("inicioDoDia() = %s; want local date 2026-10-01", inicio.Format(time.RFC3339))
	}
	if inicio.Hour() != 0 || inicio.Minute() != 0 || inicio.Second() != 0 {
		t.Fatalf("inicioDoDia() = %s; want local midnight", inicio.Format(time.RFC3339))
	}
}
