package vfsmatch

func (p *globPattern) mayMatchBelow(prefix string) bool {
	offset, component := 0, 0
	for {
		part, next, ok := nextPathPartSingle(prefix, offset)
		if !ok || component == len(p.components) {
			return true
		}
		c := p.components[component]
		if c.kind != kindLiteral {
			return true
		}
		if !p.stringsEqual(c.literal, part) {
			return false
		}
		offset, component = next, component+1
	}
}

// A nil selection means all excludes; an empty non-nil selection means none.
func (m *globMatcher) liveExcludes(parent []int, prefix string) []int {
	var result []int
	changed := false
	count := len(m.excludes)
	if parent != nil {
		count = len(parent)
	}
	for position := 0; position < count; position++ {
		i := position
		if parent != nil {
			i = parent[position]
		}
		if !m.excludes[i].mayMatchBelow(prefix) {
			if !changed {
				changed = true
				result = make([]int, 0, count-1)
				for previous := 0; previous < position; previous++ {
					if parent != nil {
						result = append(result, parent[previous])
					} else {
						result = append(result, previous)
					}
				}
			}
		} else if changed {
			result = append(result, i)
		}
	}
	if !changed {
		return parent
	}
	return result
}
