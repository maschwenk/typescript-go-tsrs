package vfsmatch

import "testing"

func TestLiveExcludesKeepOriginalFileAndDirectoryAnswers(t *testing.T) {
	specs := []string{"unrelated/**/*.ts", "src/private", "src/**/*.generated.ts", "**/node_modules", "Équipe/cache", "other/deep/file.ts"}
	for _, sensitive := range []bool{true, false} {
		files := newGlobMatcher(nil, specs, "/work", sensitive, UsageFiles)
		directories := newGlobMatcher(nil, specs, "/work", sensitive, UsageDirectories)
		for _, prefix := range []string{"/work/", "/work/src/", "/work/src/public/", "/work/équipe/", "/work/unrelated/", "/work/other/deep/"} {
			liveFiles := files.liveExcludes(nil, prefix)
			liveDirs := directories.liveExcludes(nil, prefix)
			for _, entry := range []string{"a.ts", "a.generated.ts", "file.ts", "node_modules", "cache", "private", "child"} {
				i1, b1 := files.matchesFileParts(prefix, entry)
				i2, b2 := files.matchesFilePartsIn(liveFiles, prefix, entry)
				if i1 != i2 || b1 != b2 || directories.matchesDirectoryParts(prefix, entry) != directories.matchesDirectoryPartsIn(liveDirs, prefix, entry) {
					t.Fatalf("case=%v %q+%q changed exclude matching", sensitive, prefix, entry)
				}
			}
			child := prefix + "child/"
			carried := files.liveExcludes(liveFiles, child)
			fresh := files.liveExcludes(nil, child)
			for _, entry := range []string{"a.ts", "a.generated.ts", "node_modules"} {
				a, b := files.matchesFilePartsIn(carried, child, entry)
				c, d := files.matchesFilePartsIn(fresh, child, entry)
				if a != c || b != d {
					t.Fatal("parent exclusion pruning changed a descendant answer")
				}
			}
		}
	}
}
