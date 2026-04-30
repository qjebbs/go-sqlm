package util

// List is a helper type for options that can be specified multiple times, such as SelectTags and SelectNullZeroTables.
type List[T comparable] []T

// NewList creates a new List from the given items.
func NewList[T comparable](items []T) List[T] {
	return List[T](items)
}

// Add adds items to the List.
func (l *List[T]) Add(items ...T) {
	*l = append(*l, items...)
}

// Contains checks if the List contains the specified item.
func (l List[T]) Contains(item T) bool {
	for _, v := range l {
		if v == item {
			return true
		}
	}
	return false
}

// ContainsAny checks if the List contains any of the specified items.
func (l List[T]) ContainsAny(items []T) bool {
	for _, item := range items {
		if l.Contains(item) {
			return true
		}
	}
	return false
}

// Index returns the index of the specified item in the List, or -1 if not found.
func (l List[T]) Index(v T) int {
	for i := range l {
		if v == l[i] {
			return i
		}
	}
	return -1
}

// Unique returns a new List containing only unique items from the original List.
func (l List[T]) Unique() List[T] {
	seen := make(map[T]struct{})
	var unique List[T]
	for _, v := range l {
		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			unique = append(unique, v)
		}
	}
	return unique
}

// Merge merges another List into the current List, ensuring uniqueness.
func (l *List[T]) Merge(other List[T]) {
	unique := l.Unique()
	for _, item := range other {
		if !unique.Contains(item) {
			unique = append(unique, item)
		}
	}
	*l = unique
}

// ToSlice converts the List to a slice.
func (l List[T]) ToSlice() []T {
	return []T(l)
}
