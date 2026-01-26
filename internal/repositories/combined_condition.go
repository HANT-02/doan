package repositories

type CommonCombinedCondition struct {
	Conditions CombinedCondition
	Sorting    []Sorting
	Paging     *Paging
}

type SingleCondition struct {
	FieldName string
	Value     interface{}
	Operator  string // "eq", "ne", "lt", "gt", "lte", "gte", "in", "like", ...
}

type CombinedCondition struct {
	SingleCondition    *SingleCondition
	RelationKey        *string // "AND" or "OR"
	CombinedConditions []CombinedCondition
}

func NewCombinedCondition() *CombinedCondition {
	return &CombinedCondition{
		CombinedConditions: []CombinedCondition{},
	}
}

func (c *CombinedCondition) IsEmpty() bool {
	return c.SingleCondition == nil && len(c.CombinedConditions) == 0
}

func (c *CombinedCondition) IsSingleCondition() bool {
	return c.SingleCondition != nil && len(c.CombinedConditions) == 0
}

func (c *CombinedCondition) IsValid() bool {
	return !(c.SingleCondition != nil && len(c.CombinedConditions) > 0)
}

func NewCommonCombinedCondition() *CommonCombinedCondition {
	return &CommonCombinedCondition{
		Conditions: *NewCombinedCondition(),
		Sorting:    []Sorting{},
		Paging:     &Paging{},
	}
}
func (a *CommonCombinedCondition) SetPaging(limit, page uint64) {
	a.Paging.Limit = limit
	a.Paging.Page = page
}

func (a *CommonCombinedCondition) AddSorting(field, order string) {
	a.Sorting = append(a.Sorting, Sorting{
		Field: field,
		Order: order,
	})
}

// AddExternalCondition Bọc toàn bộ điều kiện cũ vào dấu ngoặc và thêm điều kiện mới vào bên ngoài dấu ngoặc
// OldCondition := A x B x ... x C
// NewSingleCondition := D
// RelationKey := "AND" or "OR"
// Result: (A x B x ... x C) RelationKey D
func (a *CommonCombinedCondition) AddExternalCondition(relationKey *string, singleCondition SingleCondition) {
	if a.Conditions.IsEmpty() {
		a.Conditions.SingleCondition = &singleCondition
		return
	}
	// Tạo nhóm điều kiện mới
	newCombinedCondition := NewCombinedCondition()
	newCombinedCondition.SingleCondition = &singleCondition

	// Thêm relation key vào điều kiện cũ
	a.Conditions.RelationKey = relationKey
	wrapCombinedCondition := NewCombinedCondition()
	wrapCombinedCondition.CombinedConditions = []CombinedCondition{a.Conditions, *newCombinedCondition}
	// Cập nhật lại điều kiện
	a.Conditions = *wrapCombinedCondition
}

// JoinInternalCondition Nối điều kiện mới vào sau điều kiện cũ
// OldCondition := A x B x ... x C
// NewSingleCondition := D
// RelationKey := "AND" or "OR"
// Result: A x B x ... x C RelationKey D
func (a *CommonCombinedCondition) JoinInternalCondition(relationKey *string, singleCondition SingleCondition) {
	if a.Conditions.IsEmpty() {
		a.Conditions.SingleCondition = &singleCondition
		return
	}
	if a.Conditions.SingleCondition != nil && len(a.Conditions.CombinedConditions) == 0 {
		a.AddExternalCondition(relationKey, singleCondition)
		return
	}

	if len(a.Conditions.CombinedConditions) > 0 {
		a.Conditions.CombinedConditions[len(a.Conditions.CombinedConditions)-1].RelationKey = relationKey
	}
	a.Conditions.CombinedConditions = append(a.Conditions.CombinedConditions, CombinedCondition{
		SingleCondition: &singleCondition,
	})
}
