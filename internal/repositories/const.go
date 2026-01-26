package repositories

const (
	// Equal operator
	Equal = "eq"
	// EqualIgnore operator
	EqualIgnore = "eq_ignore"
	// NotEqual operator
	NotEqual = "ne"
	// LessThan operator
	LessThan = "lt"
	// GreaterThan operator
	GreaterThan = "gt"
	// LessThanOrEqual operator
	LessThanOrEqual = "lte"
	// GreaterThanOrEqual operator
	GreaterThanOrEqual = "gte"
	// In operator
	In = "in"
	// NotIn operator
	NotIn = "not_in"
	// Like operator
	Like      = "like"
	StartWith = "start_with"
	// NotLike operator
	NotLike = "not_like"
	// ILike operator
	ILike = "ilike"
	// NotILike operator
	NotILike = "not_ilike"
	// LikeContains operator
	LikeContains = "like_contains"
	// NotLikeContains operator
	NotLikeContains = "not_like_contains"
	// ILikeContains operator
	ILikeContains = "ilike_contains"
	// NotILikeContains operator
	NotILikeContains = "not_ilike_contains"
	// IsNotNullContains operator
	IsNotNull = "is_not_null"
	// JsonContains operator
	JSONContains = "json_contains"
)

const (
	EscapeLike = `\`
)
const (
	// Asc order
	Asc = "asc"
	// Desc order
	Desc = "desc"
)
const (
	RelationKeyAnd = "AND"
	RelationKeyOr  = "OR"
)
