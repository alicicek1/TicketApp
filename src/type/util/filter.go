package util

type Filter struct {
	Filters          map[string]interface{}
	Page             int64
	PageSize         int64
	SortingField     string
	SortingDirection int
}
