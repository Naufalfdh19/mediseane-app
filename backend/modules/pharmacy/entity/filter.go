package entity

type PharmacyFilter struct {
	Field string
	Order string
	Name  string
	City  string
	Limit int
	Page  int
}

type PharmacyFilterCount struct {
	Name string
	City string
}

func (f *PharmacyFilter) GetLimit() int {
	if f.Limit < 1 {
		return 10
	}
	return f.Limit
}

func (f *PharmacyFilter) GetOffset() int {
	if f.Page < 1 {
		return 0
	}
	return (f.Page - 1) * f.GetLimit()
}
