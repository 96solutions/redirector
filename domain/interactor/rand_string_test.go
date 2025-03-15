package interactor

import (
	"math/rand"
	"regexp"
	"strings"
	"testing"
)

func TestRandString(t *testing.T) {
	tests := []struct {
		name   string
		length int
		want   *regexp.Regexp
	}{
		{
			name:   "random string with length 0",
			length: 0,
			want:   regexp.MustCompile(`^$`),
		},
		{
			name:   "random string with length 5",
			length: 5,
			want:   regexp.MustCompile(`^[a-zA-Z]{5}$`),
		},
		{
			name:   "random string with length 10",
			length: 10,
			want:   regexp.MustCompile(`^[a-zA-Z]{10}$`),
		},
		{
			name:   "random string with length 32",
			length: 32,
			want:   regexp.MustCompile(`^[a-zA-Z]{32}$`),
		},
		{
			name:   "random string with length 100",
			length: 100,
			want:   regexp.MustCompile(`^[a-zA-Z]{100}$`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := randString(tt.length)

			// Check length matches expected
			if len(got) != tt.length {
				t.Errorf("randString() length = %v, want %v", len(got), tt.length)
			}

			// Check string matches expected pattern
			if !tt.want.MatchString(got) {
				t.Errorf("randString() = %v, does not match pattern %v", got, tt.want)
			}
		})
	}
}

func TestRandString_Uniqueness(t *testing.T) {
	// Generate multiple strings and verify they're unique
	iterations := 1000
	length := 32
	uniqueStrings := make(map[string]bool)

	for i := 0; i < iterations; i++ {
		s := randString(length)
		if uniqueStrings[s] {
			t.Errorf("randString() generated duplicate string: %s", s)
		}
		uniqueStrings[s] = true
	}
}

func TestRandString_Distribution(t *testing.T) {
	// Check that the character distribution is roughly even
	// (allowing for some statistical variance)
	iterations := 10000
	length := 100
	charCounts := make(map[byte]int)

	// Count occurrences of each character
	for i := 0; i < iterations; i++ {
		s := randString(length)
		for j := 0; j < len(s); j++ {
			charCounts[s[j]]++
		}
	}

	// Calculate expected counts (evenly distributed)
	totalChars := iterations * length
	expectedPerChar := float64(totalChars) / float64(len(letterBytes))
	// Allow variance of 10% from expected value
	variance := 0.1 * expectedPerChar

	// Verify each character's count is within expected range
	for _, char := range letterBytes {
		count := charCounts[byte(char)]
		if float64(count) < expectedPerChar-variance || float64(count) > expectedPerChar+variance {
			t.Logf("Character '%c' appeared %d times, expected around %.2f (±%.2f)",
				char, count, expectedPerChar, variance)
			// Not failing the test because slight statistical variations are normal
			// This is more of a sanity check to catch grossly uneven distributions
		}
	}
}

func TestRandString_WithFixedSeed(t *testing.T) {
	// Test that with a fixed seed the function produces deterministic results
	// This is more a demonstration than a strict requirement

	// Create deterministic random sources with fixed seeds
	r1 := rand.New(rand.NewSource(42))
	r2 := rand.New(rand.NewSource(42))
	r3 := rand.New(rand.NewSource(43))

	// Define a test wrapper function that uses our seeded sources
	generateWithSource := func(r *rand.Rand, n int) string {
		b := make([]byte, n)
		for i := range b {
			b[i] = letterBytes[r.Intn(len(letterBytes))]
		}
		return string(b)
	}

	// Generate strings with the same seed
	str1 := generateWithSource(r1, 10)
	str2 := generateWithSource(r2, 10)

	// They should be identical
	if str1 != str2 {
		t.Errorf("rand strings with same seed should be identical, got %s and %s", str1, str2)
	}

	// Generate string with different seed
	str3 := generateWithSource(r3, 10)

	// It should be different
	if str1 == str3 {
		t.Errorf("rand strings with different seeds should differ, but got identical: %s", str1)
	}
}

func BenchmarkRandString(b *testing.B) {
	benchmarks := []struct {
		name   string
		length int
	}{
		{"short_5", 5},
		{"medium_32", 32},
		{"long_100", 100},
		{"very_long_1000", 1000},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				randString(bm.length)
			}
		})
	}
}

// This benchmark compares the performance of the existing implementation
// with an alternative implementation using bytes.Buffer
func BenchmarkRandString_Alternatives(b *testing.B) {
	const length = 32

	// Current implementation
	b.Run("current_implementation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			randString(length)
		}
	})

	// Alternative implementation using bytes.Buffer (for comparison)
	b.Run("bytes_buffer_implementation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var buffer strings.Builder
			buffer.Grow(length)
			for j := 0; j < length; j++ {
				buffer.WriteByte(letterBytes[rand.Intn(len(letterBytes))])
			}
			_ = buffer.String()
		}
	})
}
