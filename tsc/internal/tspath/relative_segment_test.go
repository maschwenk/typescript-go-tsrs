package tspath

import (
	"strings"
	"testing"
)

func referenceRelativeSegment(p string) bool {
	parts := strings.Split(p, "/")
	for i, part := range parts {
		if part == "." || part == ".." || part == "" && i > 0 && i < len(parts)-1 {
			return true
		}
	}
	return false
}

func TestRelativePathSegmentWordScan(t *testing.T) {
	var walk func(string, int)
	walk = func(s string, n int) {
		for _, p := range []string{s, "long/path/prefix/" + s, s + "/trailing/path", "" + s + "é"} {
			if hasRelativePathSegment(p) != referenceRelativeSegment(p) {
				t.Fatalf("different answer for %q", p)
			}
		}
		if n == 0 {
			return
		}
		for _, b := range []byte{'/', '.', 'a'} {
			walk(s+string(b), n-1)
		}
	}
	walk("", 11)
	for _, p := range []string{"/pkg/node_modules/.pnpm/a@1.0/node_modules/a/file.ts", "é/../b", "é/./b", "a/\xff/../b"} {
		if hasRelativePathSegment(p) != referenceRelativeSegment(p) {
			t.Fatalf("different answer for %q", p)
		}
	}
}

func BenchmarkRelativePathSegmentWordScan(b *testing.B) {
	p := "/repo/node_modules/.pnpm/typescript@7.0/node_modules/typescript/lib/lib.esnext.d.ts"
	b.ReportAllocs()
	for b.Loop() {
		hasRelativePathSegment(p)
	}
}
func hasRelativePathSegmentOriginal(p string) bool {
	n := len(p)
	if n == 0 {
		return false
	}

	if p == "." || p == ".." {
		return true
	}

	// Leading "./" OR "../"
	if p[0] == '.' {
		if n >= 2 && p[1] == '/' {
			return true
		}
		// Leading "../"
		if n >= 3 && p[1] == '.' && p[2] == '/' {
			return true
		}
	}
	// Trailing "/." OR "/.."
	if p[n-1] == '.' {
		if n >= 2 && p[n-2] == '/' {
			return true
		}
		if n >= 3 && p[n-2] == '.' && p[n-3] == '/' {
			return true
		}
	}

	// Now look for any `//` or `/./` or `/../`

	prevSlash := false
	segLen := 0   // length of current segment since last slash
	dotCount := 0 // consecutive dots at start of the current segment; -1 => not only dots

	for i := range n {
		c := p[i]
		if c == '/' {
			// "//"
			if prevSlash {
				return true
			}
			// "/./" or "/../"
			if (segLen == 1 && dotCount == 1) || (segLen == 2 && dotCount == 2) {
				return true
			}
			prevSlash = true
			segLen = 0
			dotCount = 0
			continue
		}

		if c == '.' {
			if dotCount >= 0 {
				dotCount++
			}
		} else {
			dotCount = -1
		}
		segLen++
		prevSlash = false
	}

	// Trailing "/." or "/.."
	return (segLen == 1 && dotCount == 1) || (segLen == 2 && dotCount == 2)
}

func BenchmarkRelativePathSegmentOriginal(b *testing.B) {
	p := "/repo/node_modules/.pnpm/typescript@7.0/node_modules/typescript/lib/lib.esnext.d.ts"
	b.ReportAllocs()
	for b.Loop() {
		hasRelativePathSegmentOriginal(p)
	}
}
