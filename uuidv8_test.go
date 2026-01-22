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
				t.Fatalf("CountryUUIDv8() ошибка = %v", err)
			}

			// Проверяем, что UUID не пустой
			if u == uuid.Nil {
				t.Error("UUID не должен быть пустым")
			}

			// Проверяем версию UUID (должна быть 8)
			version := (u[6] & 0xf0) >> 4
			if version != 8 {
				t.Errorf("версия UUID = %d, ожидалось 8", version)
			}

			// Проверяем вариант RFC 4122 (должен быть 10x в двоичном)
			variant := (u[8] & 0xc0) >> 6
			if variant != 2 { // 10 в двоичном = 2 в десятичном
				t.Errorf("вариант UUID = %02b, ожидалось 10", variant)
			}
		})
	}
}

func TestCountryUUIDv8_Uniqueness(t *testing.T) {
	// Генерируем множество UUID и проверяем, что все уникальные
	country := countries.Russia
	uuidMap := make(map[uuid.UUID]bool)
	count := 1000

	for i := 0; i < count; i++ {
		u, err := CountryUUIDv8(country)
		if err != nil {
			t.Fatalf("CountryUUIDv8() ошибка = %v", err)
		}

		if uuidMap[u] {
			t.Fatalf("Найден дубликат UUID: %s", u)
		}
		uuidMap[u] = true
	}

	if len(uuidMap) != count {
		t.Errorf("Сгенерировано уникальных UUID = %d, ожидалось %d", len(uuidMap), count)
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
			// Генерируем UUID
			u, err := CountryUUIDv8(tt.country)
			if err != nil {
				t.Fatalf("CountryUUIDv8() ошибка = %v", err)
			}

			// Извлекаем страну обратно
			extractedCountry, err := ExtractCountry(u)
			if err != nil {
				t.Fatalf("ExtractCountry() ошибка = %v", err)
			}

			// Проверяем совпадение
			if extractedCountry != tt.country {
				t.Errorf("ExtractCountry() = %v (код: %d), ожидалось %v (код: %d)",
					extractedCountry, extractedCountry, tt.country, tt.country)
			}
		})
	}
}

func TestExtractCountry_WrongVersion(t *testing.T) {
	// Создаем UUID v4 (не v8)
	u := uuid.New()

	_, err := ExtractCountry(u)
	if err == nil {
		t.Error("ExtractCountry() должен вернуть ошибку для не-v8 UUID")
	}
}

func TestGetTimestamp(t *testing.T) {
	beforeTime := time.Now()
	time.Sleep(1 * time.Millisecond)

	u, err := CountryUUIDv8(countries.Russia)
	if err != nil {
		t.Fatalf("CountryUUIDv8() ошибка = %v", err)
	}

	time.Sleep(1 * time.Millisecond)
	afterTime := time.Now()

	extractedTime := GetTimestamp(u)

	// Проверяем, что время находится в ожидаемом диапазоне
	if extractedTime.Before(beforeTime) {
		t.Errorf("Извлеченное время %v раньше чем время до создания %v", extractedTime, beforeTime)
	}

	if extractedTime.After(afterTime) {
		t.Errorf("Извлеченное время %v позже чем время после создания %v", extractedTime, afterTime)
	}
}

func TestCountryUUIDv8_RoundTrip(t *testing.T) {
	// Тест полного цикла: создание -> извлечение -> проверка
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
			t.Fatalf("CountryUUIDv8(%v) ошибка = %v", originalCountry, err)
		}

		extractedCountry, err := ExtractCountry(u)
		if err != nil {
			t.Fatalf("ExtractCountry() ошибка = %v", err)
		}

		if extractedCountry != originalCountry {
			t.Errorf("Страна не совпадает: оригинал %v (код: %d), извлечено %v (код: %d)",
				originalCountry, originalCountry, extractedCountry, extractedCountry)
		}
	}
}

func TestCountryUUIDv8_TimestampProgression(t *testing.T) {
	// Проверяем, что timestamp увеличивается со временем
	u1, err := CountryUUIDv8(countries.Russia)
	if err != nil {
		t.Fatalf("CountryUUIDv8() ошибка = %v", err)
	}
	time1 := GetTimestamp(u1)

	time.Sleep(10 * time.Millisecond)

	u2, err := CountryUUIDv8(countries.Russia)
	if err != nil {
		t.Fatalf("CountryUUIDv8() ошибка = %v", err)
	}
	time2 := GetTimestamp(u2)

	if !time2.After(time1) {
		t.Errorf("Второй timestamp (%v) должен быть позже первого (%v)", time2, time1)
	}
}

// Бенчмарки для оценки производительности
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

// Тест параллельной генерации
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
		}(countries.CountryCode(i % 250)) // Перебираем разные коды стран
	}

	// Собираем результаты
	uuidMap := make(map[uuid.UUID]bool)
	for i := 0; i < goroutines*uuidsPerGoroutine; i++ {
		select {
		case u := <-results:
			if uuidMap[u] {
				t.Fatalf("Найден дубликат UUID при параллельной генерации: %s", u)
			}
			uuidMap[u] = true
		case err := <-errors:
			t.Fatalf("Ошибка при параллельной генерации: %v", err)
		}
	}

	if len(uuidMap) != goroutines*uuidsPerGoroutine {
		t.Errorf("Сгенерировано уникальных UUID = %d, ожидалось %d",
			len(uuidMap), goroutines*uuidsPerGoroutine)
	}
}
