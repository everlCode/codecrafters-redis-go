package database

import "slices"

type Zset struct {
	Keys map[string]*Zvalue
	Set  []Zvalue
}

type Zvalue struct {
	Score float64
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
		v = &Zvalue{Score: rank, Value: value}
		z.Keys[value] = v
		z.Set = append(z.Set, *v)
		isNewValue = 1
	}
	v.Score = rank

	slices.SortFunc(z.Set, func(a, b Zvalue) int {
		if a.Score < b.Score {
			return -1
		} else if a.Score == b.Score {
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

func (z *Zset) Remove(key string) int  {
	_, ok := z.Keys[key]
	if !ok {
		return 0
	}
	delete(z.Keys, key)	
	var set []Zvalue
	for _, v := range z.Set {
		if v.Value != key {
			set = append(set, v)
		}
	}
	z.Set = set

	return 1
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
	if middleValue.Score < searchValue.Score {
		index = binarySearch(set[middle:], searchValue) + middle
	} else if middleValue.Score > searchValue.Score {
		index = binarySearch(set[:middle], searchValue)
	} else if middleValue.Value < searchValue.Value {
		index = binarySearch(set[middle:], searchValue) + middle
	} else if middleValue.Value > searchValue.Value {
		index = binarySearch(set[:middle], searchValue)
	}

	return index
}
