package entity

type CategoryFilter struct {
	SortBy  string
	Order  string
	Name   string
	Limit  int
	Page   int
}

type CategoryFilterCount struct {
	Name string
}

func (f *CategoryFilter) GetLimit() int {
	if f.Limit < 1 {
		return 10
	}
	return f.Limit
}

func (f *CategoryFilter) GetOffset() int {
	if f.Page < 1 {
		return 0
	}
	return (f.Page - 1) * f.GetLimit()
}
