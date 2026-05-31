package database

import "slices"

type Zset struct {
	Keys map[string]*Zvalue
	Set []Zvalue
}

type Zvalue struct {
	Rank float64
	Value string
}

func NewZset() *Zset {
	return &Zset{
		Keys: make(map[string]*Zvalue),
	}
}

func (z *Zset) Add(rank float64, value string) int {
	var v *Zvalue
	var isNewValue int = 0
	v, ok := z.Keys[value]
	if !ok {
		v = &Zvalue{Rank: rank, Value: value}
		isNewValue = 1
	}
	
	z.Keys[value] = v
	z.Set = append(z.Set, *v)
	slices.SortFunc(z.Set, func(a, b Zvalue) int {
		if a.Rank < b.Rank {
			return -1
		} else if a.Rank == b.Rank {
			if a.Value < b.Value {
				return -1
			} else {
				return 1
			}
		} else {
			return 1
		}
	})

	return isNewValue
}

func (z *Zset) GetIndex(value string) (int, bool) {
	v, ok := z.Keys[value]
	if !ok {
		return 0, false
	}

	return binarySearch(z.Set, *v), true
}

func binarySearch(set []Zvalue, searchValue Zvalue) int {
	middle := len(set) / 2
	middleValue := set[middle]
	if middleValue == searchValue && middleValue.Value == searchValue.Value {
		return middle
	}
	var index int
	if middleValue.Rank < searchValue.Rank {
		index = binarySearch(set[middle:], searchValue) + middle
	} else if (middleValue.Rank > searchValue.Rank) {
		index = binarySearch(set[:middle], searchValue) 
	} else if (middleValue.Value < searchValue.Value) {
		index = binarySearch(set[middle:], searchValue) + middle
	} else if (middleValue.Value > searchValue.Value) {
		index = binarySearch(set[:middle], searchValue) 
	}

	return index
}