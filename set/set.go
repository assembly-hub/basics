// Package set
package set

// Set string set
type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(keyList ...T) Set[T] {
	for _, k := range keyList {
		s[k] = struct{}{}
	}
	return s
}

func (s Set[K]) Has(key K) bool {
	_, has := s[key]

	return has
}

func (s Set[T]) Del(keyList ...T) Set[T] {
	for _, k := range keyList {
		if s.Has(k) {
			delete(s, k)
		}
	}
	return s
}

func (s Set[T]) Size() int {
	return len(s)
}

func (s Set[T]) Empty() bool {
	return s.Size() == 0
}

func (s Set[T]) ToList() []T {
	if len(s) <= 0 {
		return nil
	}

	list := make([]T, 0, len(s))
	for k := range s {
		list = append(list, k)
	}
	return list
}

// Union 并集
func (s Set[T]) Union(other Set[T]) Set[T] {
	newSet := Set[T]{}
	for k := range s {
		newSet.Add(k)
	}
	for k := range other {
		newSet.Add(k)
	}
	return newSet
}

// Intersection 交集
func (s Set[T]) Intersection(other Set[T]) Set[T] {
	s1, s2 := s.Size(), other.Size()
	newSet := Set[T]{}
	if s1 < s2 {
		for k := range s {
			if other.Has(k) {
				newSet.Add(k)
			}
		}
	} else {
		for k := range other {
			if s.Has(k) {
				newSet.Add(k)
			}
		}
	}

	return newSet
}

// Difference 差集
func (s Set[T]) Difference(other Set[T]) Set[T] {
	newSet := Set[T]{}
	for k := range s {
		if !other.Has(k) {
			newSet.Add(k)
		}
	}
	return newSet
}
