package game

import "nechego/item"

// Fund represents a collectable Item received from some Source.
type Fund struct {
	Source string
	Item   *item.Item
}

// Funds is a set of items indexed by their source. Acts like a buffer
// between a user's inventory and the source of a collectable item.
type Funds map[string][]*item.Item

// Add adds the given item to the specified source.
func (f Funds) Add(source string, i *item.Item) {
	f[source] = append(f[source], i)
}

// Collect returns all funds from the specified sources and deletes
// them from the Funds. If sources are not specified, returns and
// deletes all funds. If there is no items at the source or the given
// source is not found, returns empty slice.
func (f Funds) Collect(sources ...string) (retrieved []*Fund) {
	if len(sources) == 0 {
		for src := range f {
			sources = append(sources, src)
		}
	}
	for _, src := range sources {
		defer delete(f, src)
		items := f[src]
		for _, x := range items {
			retrieved = append(retrieved, &Fund{src, x})
		}
	}
	return
}

// Sources returns all fund sources.
func (f Funds) Sources() []string {
	s := []string{}
	for source := range f {
		s = append(s, source)
	}
	return s
}
