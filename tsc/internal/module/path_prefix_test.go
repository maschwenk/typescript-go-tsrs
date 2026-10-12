package module

import (
	"fmt"
	"testing"

	"github.com/microsoft/TypeScript/tsc/internal/collections"
	"github.com/microsoft/TypeScript/tsc/internal/core"
)

func TestIndexedPathPrefixesPreserveLongestPrefixAndFirstSuffixMatch(t *testing.T) {
	paths := collections.NewOrderedMapWithSizeHint[string, []string](128)
	for _, pattern := range []string{"*", "a*one", "a*two", "a/b*", "é/*", "a/b/exact", "invalid**", "a/b/*end"} {
		paths.Set(pattern, nil)
	}
	for i := 0; i < 96; i++ {
		paths.Set(fmt.Sprintf("prefix%d/*", i), nil)
	}
	parsed := TryParsePatterns(paths)
	if parsed.byPrefix == nil {
		t.Fatal("large mapping did not build its prefix index")
	}
	candidates := []string{"", "a", "axone", "axtwo", "a/b/exact", "a/b/valueend", "a/b/value", "é/value", "é", "unmatched", "invalidxxx"}
	for i := 0; i < 96; i++ {
		candidates = append(candidates, fmt.Sprintf("prefix%d/file.ts", i), fmt.Sprintf("prefix%d", i))
	}
	for _, candidate := range candidates {
		expected := core.FindBestPatternMatch(parsed.patterns, core.Identity, candidate)
		if parsed.matchableStringSet.Has(candidate) {
			expected = core.Pattern{Text: candidate, StarIndex: -1}
		}
		if actual := MatchPatternOrExact(parsed, candidate); actual != expected {
			t.Fatalf("%q: indexed %#v, original %#v", candidate, actual, expected)
		}
	}
}
