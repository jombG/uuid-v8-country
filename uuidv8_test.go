package uuidv8country

import (
	"testing"
	"time"

	"github.com/biter777/countries"
	"github.com/google/uuid"
)

func TestCountryUUIDv8_Generation(t *testing.T) {
	tests := []struct {
		name    string
		country countries.CountryCode
	}{
		{"Russia", countries.Russia},
		{"USA", countries.USA},
		{"Germany", countries.Germany},
		{"Ukraine", countries.Ukraine},
		{"Japan", countries.Japan},
		{"Brazil", countries.Brazil},
		{"China", countries.China},
		{"India", countries.India},
		{"France", countries.France},
		{"Canada", countries.Canada},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := CountryUUIDv8(tt.country)
			if err != nil {
				t.Fatalf("CountryUUIDv8() error = %v", err)
			}

			// Check that UUID is not nil
			if u == uuid.Nil {
				t.Error("UUID should not be nil")
			}

			// Check UUID version (should be 8)
			version := (u[6] & 0xf0) >> 4
			if version != 8 {
				t.Errorf("UUID version = %d, expected 8", version)
			}

			// Check RFC 4122 variant (should be 10x in binary)
			variant := (u[8] & 0xc0) >> 6
			if variant != 2 { // 10 in binary = 2 in decimal
				t.Errorf("UUID variant = %02b, expected 10", variant)
			}
		})
	}
}

func TestCountryUUIDv8_Uniqueness(t *testing.T) {
	// Generate multiple UUIDs and check that all are unique
	country := countries.Russia
	uuidMap := make(map[uuid.UUID]bool)
	count := 1000

	for i := 0; i < count; i++ {
		u, err := CountryUUIDv8(country)
		if err != nil {
			t.Fatalf("CountryUUIDv8() error = %v", err)
		}

		if uuidMap[u] {
			t.Fatalf("Found duplicate UUID: %s", u)
		}
		uuidMap[u] = true
	}

	if len(uuidMap) != count {
		t.Errorf("Generated unique UUIDs = %d, expected %d", len(uuidMap), count)
	}
}

func TestExtractCountry_Success(t *testing.T) {
	tests := []struct {
		name    string
		country countries.CountryCode
	}{
		{"Russia", countries.Russia},
		{"USA", countries.USA},
		{"Germany", countries.Germany},
		{"Ukraine", countries.Ukraine},
		{"Japan", countries.Japan},
		{"Brazil", countries.Brazil},
		{"China", countries.China},
		{"Unknown", countries.Unknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate UUID
			u, err := CountryUUIDv8(tt.country)
			if err != nil {
				t.Fatalf("CountryUUIDv8() error = %v", err)
			}

			// Extract country back
			extractedCountry, err := ExtractCountry(u)
			if err != nil {
				t.Fatalf("ExtractCountry() error = %v", err)
			}

			// Check match
			if extractedCountry != tt.country {
				t.Errorf("ExtractCountry() = %v (code: %d), expected %v (code: %d)",
					extractedCountry, extractedCountry, tt.country, tt.country)
			}
		})
	}
}

func TestExtractCountry_WrongVersion(t *testing.T) {
	// Create UUID v4 (not v8)
	u := uuid.New()

	_, err := ExtractCountry(u)
	if err == nil {
		t.Error("ExtractCountry() should return error for non-v8 UUID")
	}
}

func TestGetTimestamp(t *testing.T) {
	beforeTime := time.Now()
	time.Sleep(1 * time.Millisecond)

	u, err := CountryUUIDv8(countries.Russia)
	if err != nil {
		t.Fatalf("CountryUUIDv8() error = %v", err)
	}

	time.Sleep(1 * time.Millisecond)
	afterTime := time.Now()

	extractedTime := GetTimestamp(u)

	// Check that time is in expected range
	if extractedTime.Before(beforeTime) {
		t.Errorf("Extracted time %v is before creation time %v", extractedTime, beforeTime)
	}

	if extractedTime.After(afterTime) {
		t.Errorf("Extracted time %v is after creation time %v", extractedTime, afterTime)
	}
}

func TestCountryUUIDv8_RoundTrip(t *testing.T) {
	// Test full cycle: creation -> extraction -> validation
	allCountries := []countries.CountryCode{
		countries.Russia,
		countries.USA,
		countries.Germany,
		countries.Ukraine,
		countries.Japan,
		countries.Brazil,
		countries.China,
		countries.India,
		countries.France,
		countries.Canada,
		countries.Australia,
		countries.Mexico,
		countries.Spain,
		countries.Italy,
		countries.UnitedKingdom,
	}

	for _, originalCountry := range allCountries {
		u, err := CountryUUIDv8(originalCountry)
		if err != nil {
			t.Fatalf("CountryUUIDv8(%v) error = %v", originalCountry, err)
		}

		extractedCountry, err := ExtractCountry(u)
		if err != nil {
			t.Fatalf("ExtractCountry() error = %v", err)
		}

		if extractedCountry != originalCountry {
			t.Errorf("Country mismatch: original %v (code: %d), extracted %v (code: %d)",
				originalCountry, originalCountry, extractedCountry, extractedCountry)
		}
	}
}

func TestCountryUUIDv8_TimestampProgression(t *testing.T) {
	// Check that timestamp increases over time
	u1, err := CountryUUIDv8(countries.Russia)
	if err != nil {
		t.Fatalf("CountryUUIDv8() error = %v", err)
	}
	time1 := GetTimestamp(u1)

	time.Sleep(10 * time.Millisecond)

	u2, err := CountryUUIDv8(countries.Russia)
	if err != nil {
		t.Fatalf("CountryUUIDv8() error = %v", err)
	}
	time2 := GetTimestamp(u2)

	if !time2.After(time1) {
		t.Errorf("Second timestamp (%v) should be after first (%v)", time2, time1)
	}
}

// Benchmarks for performance evaluation
func BenchmarkCountryUUIDv8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = CountryUUIDv8(countries.Russia)
	}
}

func BenchmarkExtractCountry(b *testing.B) {
	u, _ := CountryUUIDv8(countries.Russia)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ExtractCountry(u)
	}
}

func BenchmarkGetTimestamp(b *testing.B) {
	u, _ := CountryUUIDv8(countries.Russia)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = GetTimestamp(u)
	}
}

// Test concurrent generation
func TestCountryUUIDv8_Concurrent(t *testing.T) {
	const goroutines = 100
	const uuidsPerGoroutine = 100

	results := make(chan uuid.UUID, goroutines*uuidsPerGoroutine)
	errors := make(chan error, goroutines*uuidsPerGoroutine)

	for i := 0; i < goroutines; i++ {
		go func(country countries.CountryCode) {
			for j := 0; j < uuidsPerGoroutine; j++ {
				u, err := CountryUUIDv8(country)
				if err != nil {
					errors <- err
					return
				}
				results <- u
			}
		}(countries.CountryCode(i % 250)) // Iterate through different country codes
	}

	// Collect results
	uuidMap := make(map[uuid.UUID]bool)
	for i := 0; i < goroutines*uuidsPerGoroutine; i++ {
		select {
		case u := <-results:
			if uuidMap[u] {
				t.Fatalf("Found duplicate UUID during concurrent generation: %s", u)
			}
			uuidMap[u] = true
		case err := <-errors:
			t.Fatalf("Error during concurrent generation: %v", err)
		}
	}

	if len(uuidMap) != goroutines*uuidsPerGoroutine {
		t.Errorf("Generated unique UUIDs = %d, expected %d",
			len(uuidMap), goroutines*uuidsPerGoroutine)
	}
}
