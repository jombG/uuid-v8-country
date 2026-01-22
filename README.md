# UUID v8 with Country Codes

[![Go Reference](https://pkg.go.dev/badge/github.com/jombG/uuid-v8-country.svg)](https://pkg.go.dev/github.com/jombG/uuid-v8-country)
[![Go Report Card](https://goreportcard.com/badge/github.com/jombG/uuid-v8-country)](https://goreportcard.com/report/github.com/jombG/uuid-v8-country)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go implementation of UUID version 8 with embedded country codes. This library allows you to generate UUIDs that contain geographical information, making it useful for distributed systems that need to track the origin of entities.

## Features

- **UUID v8 Compliant**: Fully compliant with RFC 4122 UUID version 8 specification
- **Country Code Embedding**: Embeds country codes from [biter777/countries](https://github.com/biter777/countries)
- **Timestamp Support**: Includes nanosecond-precision timestamps for temporal ordering
- **Cryptographically Secure**: Uses `crypto/rand` for random number generation
- **Zero External Dependencies**: Minimal dependencies, only standard library plus UUID and countries packages
- **High Performance**: Optimized for speed with minimal allocations
- **Thread-Safe**: Safe for concurrent use across multiple goroutines

## Motivation

In distributed systems, it's often useful to know where an entity originated. Traditional UUIDs don't carry this information. By embedding country codes directly into UUIDs, we can:

- Track geographical distribution of data
- Implement geo-aware sharding strategies
- Comply with data residency requirements
- Debug issues related to specific regions
- Analyze usage patterns by country

## Installation

```bash
go get github.com/jombG/uuid-v8-country
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    uuidcountry "github.com/jombG/uuid-v8-country"
    "github.com/biter777/countries"
)

func main() {
    // Generate a UUID v8 with a country code
    u, err := uuidcountry.CountryUUIDv8(countries.USA)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Generated UUID:", u)
    // Output: xxxxxxxx-xxxx-8xxx-xxxx-xxxxxxxxxxxx
}
```

## Usage Examples

### Basic Usage

```go
import (
    uuidcountry "github.com/jombG/uuid-v8-country"
    "github.com/biter777/countries"
)

// Generate UUID for a specific country
u, err := uuidcountry.CountryUUIDv8(countries.Japan)
if err != nil {
    log.Fatal(err)
}
fmt.Println(u) // Example: 018d1234-5678-8abc-bdef-0123456789ab
```

### Extracting Country Information

```go
// Create UUID
u, _ := uuidcountry.CountryUUIDv8(countries.Germany)

// Extract country code back
country, err := uuidcountry.ExtractCountry(u)
if err != nil {
    log.Fatal(err)
}

fmt.Println(country) // Output: Germany
fmt.Println(country.Alpha2()) // Output: DE
fmt.Println(country.Alpha3()) // Output: DEU
```

### Working with Timestamps

```go
// Generate UUID
u, _ := uuidcountry.CountryUUIDv8(countries.France)

// Extract timestamp
timestamp := uuidcountry.GetTimestamp(u)
fmt.Println(timestamp.Format(time.RFC3339))
// Output: 2026-01-22T10:30:45.123456789Z
```

### Complete Example

See [examples/main.go](examples/main.go) for a complete working example.

```go
package main

import (
    "fmt"
    "time"

    "github.com/biter777/countries"
    uuidcountry "github.com/jombG/uuid-v8-country"
)

func main() {
    // Generate UUIDs for different countries
    countriesTest := []countries.CountryCode{
        countries.USA,
        countries.Russia,
        countries.China,
        countries.India,
    }

    for _, country := range countriesTest {
        // Create UUID
        u, err := uuidcountry.CountryUUIDv8(country)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }

        // Extract information
        extractedCountry, _ := uuidcountry.ExtractCountry(u)
        timestamp := uuidcountry.GetTimestamp(u)

        // Display results
        fmt.Printf("UUID: %s\n", u)
        fmt.Printf("  Country: %s (%s)\n", extractedCountry, extractedCountry.Alpha2())
        fmt.Printf("  Time: %s\n", timestamp.Format(time.RFC3339))
        fmt.Printf("  Version: %d\n\n", (u[6]&0xf0)>>4)
    }
}
```

To run the example:

```bash
cd examples
go run main.go
```

## UUID Structure

The UUID v8 format used by this library has the following structure:

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                    unix_ts_ns (bytes 0-3)                     |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                    unix_ts_ns (bytes 4-7)                     |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|  ver  |       rand (byte 6)     | var | country_code (byte 8)|
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|              country_code (bytes 9-10)        | rand (byte 11)|
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                       rand (bytes 12-15)                      |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

### Field Descriptions

- **unix_ts_ns** (bytes 0-7): Unix timestamp in nanoseconds (64-bit, big-endian)
- **ver** (4 bits): UUID version, always `8`
- **var** (2 bits): UUID variant, always `10` (RFC 4122)
- **country_code** (22 bits): Country code from biter777/countries package
- **rand**: Cryptographically secure random data

## API Reference

### CountryUUIDv8

```go
func CountryUUIDv8(country countries.CountryCode) (uuid.UUID, error)
```

Generates a new UUID v8 with the specified country code embedded.

**Parameters:**
- `country`: A country code from the `github.com/biter777/countries` package

**Returns:**
- `uuid.UUID`: The generated UUID
- `error`: Error if random number generation fails

### ExtractCountry

```go
func ExtractCountry(u uuid.UUID) (countries.CountryCode, error)
```

Extracts the country code from a UUID v8 generated by `CountryUUIDv8`.

**Parameters:**
- `u`: The UUID to extract from

**Returns:**
- `countries.CountryCode`: The extracted country code
- `error`: Error if the UUID is not version 8

### GetTimestamp

```go
func GetTimestamp(u uuid.UUID) time.Time
```

Extracts the timestamp from a UUID generated by `CountryUUIDv8`.

**Parameters:**
- `u`: The UUID to extract from

**Returns:**
- `time.Time`: The timestamp embedded in the UUID

## Performance

Benchmarks run on Apple M1:

```
BenchmarkCountryUUIDv8-8      2841036      422.3 ns/op      16 B/op      1 allocs/op
BenchmarkExtractCountry-8    185304114     6.458 ns/op       0 B/op      0 allocs/op
BenchmarkGetTimestamp-8      194682717     6.153 ns/op       0 B/op      0 allocs/op
```

## Testing

Run tests:

```bash
go test -v
```

Run tests with coverage:

```bash
go test -v -cover
```

Run benchmarks:

```bash
go test -bench=. -benchmem
```

## Dependencies

- [github.com/google/uuid](https://github.com/google/uuid) - UUID generation and parsing
- [github.com/biter777/countries](https://github.com/biter777/countries) - Country codes and information

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [ash3in/uuidv8](https://github.com/ash3in/uuidv8)
- Uses country codes from [biter777/countries](https://github.com/biter777/countries)
- Built on top of [google/uuid](https://github.com/google/uuid)

## Related Projects

- [UUID v7 RFC](https://datatracker.ietf.org/doc/html/rfc9562) - Official UUID specification
- [UUIDv8 Draft](https://datatracker.ietf.org/doc/html/draft-ietf-uuidrev-rfc4122bis) - UUID version 8 specification

## Support

If you have any questions or issues, please open an issue on GitHub.
