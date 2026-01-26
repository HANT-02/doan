package repositories

type Condition struct {
	Field string
	Value interface{}
	Op    string // "eq", "ne", "lt", "gt", "lte", "gte", "in", "like", ...
}

type CommonCondition struct {
	Conditions   []Condition
	OrConditions []Condition
	Columns      []string // Specific columns to select
	Preloads     []string // Fields to preload (e.g., "User", "Profile")
	Sorting      []Sorting
	Paging       *Paging
}

func NewCommonCondition() *CommonCondition {
	return &CommonCondition{
		Conditions: []Condition{},
		Sorting:    []Sorting{},
		Paging:     &Paging{},
	}
}
func (cc *CommonCondition) AddColumns(columns []string) {
	cc.Columns = append(cc.Columns, columns...)
}
func (cc *CommonCondition) AddCondition(field string, value interface{}, op string) {
	cc.Conditions = append(cc.Conditions, Condition{
		Field: field,
		Value: value,
		Op:    op,
	})
}
func (cc *CommonCondition) AddOrCondition(input []Condition) {
	for _, cond := range input {
		cc.OrConditions = append(cc.OrConditions, cond)
	}
}
func (cc *CommonCondition) SetPaging(limit, page uint64) {
	cc.Paging.Limit = limit
	cc.Paging.Page = page
}
func (cc *CommonCondition) SetPreload(preloads []string) {
	cc.Preloads = preloads
}
func (cc *CommonCondition) AddSorting(field, order string) {
	cc.Sorting = append(cc.Sorting, Sorting{
		Field: field,
		Order: order,
	})
}

func (cc *CommonCondition) WithPaging(limit, page uint64) *CommonCondition {
	cc.Paging = &Paging{
		Limit: limit,
		Page:  page,
	}
	return cc
}

func (cc *CommonCondition) WithCondition(field string, value interface{}, op string) *CommonCondition {
	condition := Condition{
		Field: field,
		Value: value,
		Op:    op,
	}
	cc.Conditions = append(cc.Conditions, condition)
	return cc
}

func (cc *CommonCondition) WithSorting(field string, order string) *CommonCondition {
	sorting := Sorting{
		Field: field,
		Order: order,
	}
	cc.Sorting = append(cc.Sorting, sorting)
	return cc
}
