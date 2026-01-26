package postgres

import (
	repository "doan/internal/repositories"
	"doan/pkg/utils"
	"fmt"
	"gorm.io/gorm/clause"
	"strings"

	"gorm.io/gorm"
)

var castFields = map[string]struct{}{
	"year":              {},
	"end_year":          {},
	"duration_of_study": {},
}

func isCastField(field string) bool {
	_, ok := castFields[field]
	return ok
}

func BuildQuery(db *gorm.DB, condition *repository.CommonCondition) (*gorm.DB, error) {
	db, err := BuildConditions(db, condition.Conditions)
	if err != nil {
		return db, err
	}
	db, err = BuildOrConditions(db, condition.OrConditions)
	if err != nil {
		return db, err
	}
	if len(condition.Preloads) > 0 {
		db, err = BuildPreloads(db, condition.Preloads)
		if err != nil {
			return db, err
		}
	}
	db = BuildSorting(db, condition.Sorting)
	if condition.Paging != nil && condition.Paging.Limit > 0 {
		db = BuildPaging(db, condition.Paging)

	}
	return db, nil

}

func BuildConditions(db *gorm.DB, conditions []repository.Condition) (*gorm.DB, error) {
	for _, cond := range conditions {
		switch strings.ToLower(cond.Op) {
		case repository.Equal:
			if utils.IsValueNil(cond.Value) {
				db = db.Where(fmt.Sprintf("%s IS NULL", cond.Field))
			} else {
				db = db.Where(fmt.Sprintf("%s = ?", cond.Field), cond.Value)
			}
		case repository.EqualIgnore:
			db = db.Where(fmt.Sprintf("LOWER(%s) = LOWER(?)", cond.Field), cond.Value)
		case repository.NotEqual:
			if utils.IsValueNil(cond.Value) {
				db = db.Where(fmt.Sprintf("%s IS NOT NULL", cond.Field))
			} else {
				db = db.Where(fmt.Sprintf("%s != ?", cond.Field), cond.Value)
			}
		case repository.LessThan:
			db = db.Where(fmt.Sprintf("%s < ?", cond.Field), cond.Value)
		case repository.GreaterThan:
			db = db.Where(fmt.Sprintf("%s > ?", cond.Field), cond.Value)
		case repository.LessThanOrEqual:
			db = db.Where(fmt.Sprintf("%s <= ?", cond.Field), cond.Value)
		case repository.GreaterThanOrEqual:
			db = db.Where(fmt.Sprintf("%s >= ?", cond.Field), cond.Value)
		case repository.In:
			db = db.Where(fmt.Sprintf("%s IN (?)", cond.Field), cond.Value)
		case repository.NotIn:
			db = db.Where(fmt.Sprintf("%s NOT IN (?)", cond.Field), cond.Value)
		case repository.LikeContains:
			if isCastField(cond.Field) {
				db = db.Where(fmt.Sprintf("%s::text LIKE ?", cond.Field), fmt.Sprintf("%%%s%%", fmt.Sprint(cond.Value)))
				continue
			} else {
				db = db.Where(fmt.Sprintf("%s LIKE ? ESCAPE '%s'", cond.Field, repository.EscapeLike), fmt.Sprintf("%%%s%%", EscapeLike(fmt.Sprintf("%s", cond.Value))))
			}
		case repository.NotLikeContains:
			db = db.Where(fmt.Sprintf("%s NOT LIKE ? ESCAPE '%s'", cond.Field, repository.EscapeLike), fmt.Sprintf("%%%s%%", EscapeLike(fmt.Sprintf("%s", cond.Value))))
		case repository.ILikeContains:
			db = db.Where(fmt.Sprintf("%s ILIKE ? ESCAPE '%s'", cond.Field, repository.EscapeLike), fmt.Sprintf("%%%s%%", EscapeLike(fmt.Sprintf("%s", cond.Value))))
		case repository.StartWith:
			db = db.Where(fmt.Sprintf("%s LIKE ? ESCAPE '%s'", cond.Field, repository.EscapeLike), fmt.Sprintf("%s%%", EscapeLike(fmt.Sprintf("%s", cond.Value))))
		case repository.NotILikeContains:
			db = db.Where(fmt.Sprintf("%s NOT ILIKE ? ESCAPE '%s'", cond.Field, repository.EscapeLike), fmt.Sprintf("%%%s%%", EscapeLike(fmt.Sprintf("%s", cond.Value))))
		case repository.Like:
			db = db.Where(fmt.Sprintf("%s LIKE ? ESCAPE '%s'", cond.Field, repository.EscapeLike), EscapeLike(fmt.Sprintf("%s", cond.Value)))
		case repository.NotLike:
			db = db.Where(fmt.Sprintf("%s NOT LIKE ? ESCAPE '%s'", cond.Field, repository.EscapeLike), EscapeLike(fmt.Sprintf("%s", cond.Value)))
		case repository.ILike:
			db = db.Where(fmt.Sprintf("%s ILIKE ? ESCAPE '%s'", cond.Field, repository.EscapeLike), EscapeLike(fmt.Sprintf("%s", cond.Value)))
		case repository.NotILike:
			db = db.Where(fmt.Sprintf("%s NOT ILIKE ? ESCAPE '%s'", cond.Field, repository.EscapeLike), EscapeLike(fmt.Sprintf("%s", cond.Value)))
		case repository.IsNotNull:
			db = db.Where(fmt.Sprintf("%s IS NOT NULL", cond.Field))
		case repository.JSONContains:
			db = db.Where(fmt.Sprintf("%s @> ?", cond.Field), cond.Value)
		default:
			return db, fmt.Errorf("unsupported operator: %s", cond.Op)
		}
	}
	return db, nil
}
func BuildOrConditions(db *gorm.DB, orConditions []repository.Condition) (*gorm.DB, error) {
	if len(orConditions) == 0 {
		return db, nil
	}

	var ors []clause.Expression

	for _, cond := range orConditions {
		switch strings.ToLower(cond.Op) {
		case repository.Equal:
			ors = append(ors,
				clause.Expr{SQL: fmt.Sprintf("%s = ?", cond.Field), Vars: []interface{}{cond.Value}},
			)

		case repository.ILike:
			ors = append(ors,
				clause.Expr{SQL: fmt.Sprintf("%s ILIKE ?", cond.Field), Vars: []interface{}{cond.Value}},
			)

		case repository.LikeContains:
			ors = append(ors,
				clause.Expr{
					SQL:  fmt.Sprintf("%s LIKE ? ESCAPE '%s'", cond.Field, repository.EscapeLike),
					Vars: []interface{}{fmt.Sprintf("%%%s%%", EscapeLike(fmt.Sprint(cond.Value)))},
				},
			)

		case repository.ILikeContains:
			ors = append(ors,
				clause.Expr{
					SQL:  fmt.Sprintf("%s ILIKE ? ESCAPE '%s'", cond.Field, repository.EscapeLike),
					Vars: []interface{}{fmt.Sprintf("%%%s%%", EscapeLike(fmt.Sprint(cond.Value)))},
				},
			)
		}
	}

	return db.Clauses(clause.Or(ors...)), nil
}

func BuildPreloads(db *gorm.DB, preloads []string) (*gorm.DB, error) {
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	return db, nil
}
func BuildSorting(db *gorm.DB, sorting []repository.Sorting) *gorm.DB {
	for _, sort := range sorting {
		if sort.Order == repository.Asc {
			db = db.Order(sort.Field)
		} else {
			db = db.Order(fmt.Sprintf("%s DESC", sort.Field))
		}
	}
	return db
}

func BuildPaging(db *gorm.DB, paging *repository.Paging) *gorm.DB {
	if paging.Limit == 0 {
		paging.Limit = 10 // Default limit
	}
	if paging.Page == 0 {
		paging.Page = 1 // Default page
	}
	limit := paging.Limit
	// Validate page
	page := paging.Page
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	return db.Limit(int(limit)).Offset(int(offset))
}

func EscapeLike(s string) string {
	s = strings.ReplaceAll(s, repository.EscapeLike, fmt.Sprintf(`%s%s`, repository.EscapeLike, repository.EscapeLike))
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}
