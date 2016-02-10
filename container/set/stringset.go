package set

// An unordered collection of unique elements which supports lookups, insertions, deletions,
// iteration, and common binary set operations.  It is not guaranteed to be thread-safe.
type Set interface {
	// Returns a new Set that contains exactly the same elements as this set.
	Copy() Set

	// Returns the cardinality of this set.
	Len() int

	// Returns true if and only if this set contains v (according to Go equality rules).
	Contains(v string) bool
	// Inserts v into this set.
	Add(v string)
	// Removes v from this set, if it is present.  Returns true if and only if v was present.
	Remove(v string) bool
	// Return All Keys
	ToArray() []string

	// Executes f(v) for every element v in this set.  If f mutates this set, behavior is undefined.
	Do(f func(string))
	// Executes f(v) once for every element v in the set, aborting if f ever returns false. If f
	// mutates this set, behavior is undefined.
	DoWhile(f func(string) bool)
	// Returns a channel from which each element in the set can be read exactly once.  If this set
	// is mutated before the channel is emptied, the exact data read from the channel is undefined.
	Iter() <-chan string

	// Adds every element in s into this set.
	Union(s Set)
	// Removes every element not in s from this set.
	Intersect(s Set)
	// Removes every element in s from this set.
	Subtract(s Set)
	// Removes all elements from the set.
	Init()
	// Returns true if and only if all elements in this set are elements in s.
	IsSubset(s Set) bool
	// Returns true if and only if all elements in s are elements in this set.
	IsSuperset(s Set) bool
	// Returns true if and only if this set and s contain exactly the same elements.
	IsEqual(s Set) bool
	// Removes all elements v from this set that satisfy f(v) == true.
	RemoveIf(f func(string) bool)
}

// Returns a new set which is the union of s1 and s2.  s1 and s2 are unmodified.
func Union(s1 Set, s2 Set) Set {
	s3 := s1.Copy()
	s3.Union(s2)
	return s3
}

// Returns a new set which is the intersect of s1 and s2.  s1 and s2 are
// unmodified.
func Intersect(s1 Set, s2 Set) Set {
	s3 := s1.Copy()
	s3.Intersect(s2)
	return s3
}

// Returns a new set which is the difference between s1 and s2.  s1 and s2 are
// unmodified.
func Subtract(s1 Set, s2 Set) Set {
	s3 := s1.Copy()
	s3.Subtract(s2)
	return s3
}

// Returns a new Set pre-populated with the given items
func NewStringSet(items ...string) Set {
	res := setImpl{
		data: make(map[string]struct{}),
	}
	for _, item := range items {
		res.Add(item)
	}
	return res
}

type setImpl struct {
	data map[string]struct{}
}

func (s setImpl) Len() int {
	return len(s.data)
}

func (s setImpl) Copy() Set {
	res := NewStringSet()
	res.Union(s)
	return res
}

func (s setImpl) Init() {
	s.data = make(map[string]struct{})
}

func (s setImpl) Contains(v string) bool {
	_, ok := s.data[v]
	return ok
}

func (s setImpl) Add(v string) {
	s.data[v] = struct{}{}
}

func (s setImpl) Remove(v string) bool {
	_, ok := s.data[v]
	if ok {
		delete(s.data, v)
	}
	return ok
}

func (s setImpl) Do(f func(string)) {
	for key := range s.data {
		f(key)
	}
}

func (s setImpl) DoWhile(f func(string) bool) {
	for key := range s.data {
		if !f(key) {
			break
		}
	}
}

func (s setImpl) Iter() <-chan string {
	iter := make(chan string)
	go func() {
		for key := range s.data {
			iter <- key
		}
		close(iter)
	}()
	return iter
}

func (s setImpl) Union(s2 Set) {
	s2.Do(func(item string) { s.Add(item) })
}

func (s setImpl) Intersect(s2 Set) {
	var toRemove []string = nil
	for key := range s.data {
		if !s2.Contains(key) {
			toRemove = append(toRemove, key)
		}
	}

	for _, key := range toRemove {
		s.Remove(key)
	}
}

func (s setImpl) Subtract(s2 Set) {
	s2.Do(func(item string) { s.Remove(item) })
}

func (s setImpl) IsSubset(s2 Set) (isSubset bool) {
	isSubset = true
	s.DoWhile(func(item string) bool {
		if !s2.Contains(item) {
			isSubset = false
		}
		return isSubset
	})
	return
}

func (s setImpl) IsSuperset(s2 Set) bool {
	return s2.IsSubset(s)
}

func (s setImpl) IsEqual(s2 Set) bool {
	if s.Len() != s2.Len() {
		return false
	}

	return s.IsSubset(s2)
}

func (s setImpl) RemoveIf(f func(string) bool) {
	var toRemove []string
	for item := range s.data {
		if f(item) {
			toRemove = append(toRemove, item)
		}
	}

	for _, item := range toRemove {
		s.Remove(item)
	}
}

func (s setImpl) ToArray() []string {
	result := make([]string, len(s.data))
	pos := 0
	for key, _ := range s.data {
		result[pos] = key
		pos++
	}

	return result
}
