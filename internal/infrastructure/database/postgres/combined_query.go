package postgres

import (
	repository "doan/internal/repositories"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

func BuildCommonCombinedQuery(db *gorm.DB, condition *repository.CommonCombinedCondition) (*gorm.DB, error) {
	db, err := BuildCombinedCondition(db, condition.Conditions)
	if err != nil {
		return db, err
	}
	db = BuildSorting(db, condition.Sorting)
	if condition.Paging != nil && condition.Paging.Limit > 0 {
		db = BuildPaging(db, condition.Paging)
	}
	return db, nil
}

func BuildCombinedCondition(db *gorm.DB, conditions repository.CombinedCondition) (*gorm.DB, error) {
	if conditions.IsEmpty() {
		return db, nil
	}
	whereArgs := []interface{}{}
	whereQuery, err := buildCombinedWhereCondition(conditions, &whereArgs)
	if err != nil {
		return db, err
	}
	if whereQuery != nil {
		db = db.Where(*whereQuery, whereArgs...)
	}
	return db, nil
}

func buildCombinedWhereCondition(condition repository.CombinedCondition, whereArgs *[]interface{}) (*string, error) {
	if condition.IsEmpty() {
		return nil, nil
	}
	if condition.IsSingleCondition() {
		return buildSingleWhereCondition(*condition.SingleCondition, whereArgs)
	}
	result := ""
	for index, combinedCondition := range condition.CombinedConditions {
		whereQuery, err := buildCombinedWhereCondition(combinedCondition, whereArgs)
		if err != nil {
			return nil, err
		}
		if whereQuery != nil {
			result = fmt.Sprintf("%s %s", result, *whereQuery)
		}
		if index != len(condition.CombinedConditions)-1 {
			if combinedCondition.RelationKey == nil {
				return nil, fmt.Errorf("relation key is nil")
			}
			result = fmt.Sprintf("%s %s", result, strings.ToUpper(*combinedCondition.RelationKey))
		}
	}
	result = fmt.Sprintf("(%s)", result)
	return &result, nil
}

func buildSingleWhereCondition(cond repository.SingleCondition, whereArgs *[]interface{}) (*string, error) {
	result := ""
	if whereArgs == nil {
		whereArgs = &[]interface{}{}
	}
	switch strings.ToLower(cond.Operator) {
	case repository.Equal:
		result = fmt.Sprintf(" %s = ? ", cond.FieldName)
		*whereArgs = append(*whereArgs, cond.Value)
	case repository.EqualIgnore:
		result = fmt.Sprintf("LOWER(%s) = LOWER(?)", cond.FieldName)
		*whereArgs = append(*whereArgs, cond.Value)
	case repository.NotEqual:
		result = fmt.Sprintf("%s != ?", cond.FieldName)
		*whereArgs = append(*whereArgs, cond.Value)
	case repository.LessThan:
		result = fmt.Sprintf("%s < ?", cond.FieldName)
		*whereArgs = append(*whereArgs, cond.Value)
	case repository.GreaterThan:
		result = fmt.Sprintf("%s > ?", cond.FieldName)
		*whereArgs = append(*whereArgs, cond.Value)
	case repository.LessThanOrEqual:
		result = fmt.Sprintf("%s <= ?", cond.FieldName)
		*whereArgs = append(*whereArgs, cond.Value)
	case repository.GreaterThanOrEqual:
		result = fmt.Sprintf("%s >= ?", cond.FieldName)
		*whereArgs = append(*whereArgs, cond.Value)
	case repository.In:
		result = fmt.Sprintf("%s IN (?)", cond.FieldName)
		*whereArgs = append(*whereArgs, cond.Value)
	case repository.LikeContains:
		if isCastField(cond.FieldName) {
			result = fmt.Sprintf("%s::text LIKE ?", cond.FieldName)
			*whereArgs = append(*whereArgs, fmt.Sprintf("%%%s%%", fmt.Sprint(cond.Value)))
		} else {
			result = fmt.Sprintf("%s LIKE ? ESCAPE '%s'", cond.FieldName, repository.EscapeLike)
			*whereArgs = append(*whereArgs, fmt.Sprintf("%%%s%%", EscapeLike(fmt.Sprintf("%s", cond.Value))))
		}
	case repository.NotLikeContains:
		result = fmt.Sprintf("%s NOT LIKE ? ESCAPE '%s'", cond.FieldName, repository.EscapeLike)
		*whereArgs = append(*whereArgs, fmt.Sprintf("%%%s%%", EscapeLike(fmt.Sprintf("%s", cond.Value))))
	case repository.ILikeContains:
		result = fmt.Sprintf("%s ILIKE ? ESCAPE '%s'", cond.FieldName, repository.EscapeLike)
		*whereArgs = append(*whereArgs, fmt.Sprintf("%%%s%%", EscapeLike(fmt.Sprintf("%s", cond.Value))))
	case repository.NotILikeContains:
		result = fmt.Sprintf("%s NOT ILIKE ? ESCAPE '%s'", cond.FieldName, repository.EscapeLike)
		*whereArgs = append(*whereArgs, fmt.Sprintf("%%%s%%", EscapeLike(fmt.Sprintf("%s", cond.Value))))
	case repository.Like:
		result = fmt.Sprintf("%s LIKE ? ESCAPE '%s'", cond.FieldName, repository.EscapeLike)
		*whereArgs = append(*whereArgs, EscapeLike(fmt.Sprintf("%s", cond.Value)))
	case repository.NotLike:
		result = fmt.Sprintf("%s NOT LIKE ? ESCAPE '%s'", cond.FieldName, repository.EscapeLike)
		*whereArgs = append(*whereArgs, EscapeLike(fmt.Sprintf("%s", cond.Value)))
	case repository.ILike:
		result = fmt.Sprintf("%s ILIKE ? ESCAPE '%s'", cond.FieldName, repository.EscapeLike)
		*whereArgs = append(*whereArgs, EscapeLike(fmt.Sprintf("%s", cond.Value)))
	case repository.NotILike:
		result = fmt.Sprintf("%s NOT ILIKE ? ESCAPE '%s'", cond.FieldName, repository.EscapeLike)
		*whereArgs = append(*whereArgs, EscapeLike(fmt.Sprintf("%s", cond.Value)))
	case repository.NotIn:
		result = fmt.Sprintf("%s not IN (?)", cond.FieldName)
		*whereArgs = append(*whereArgs, cond.Value)
	default:
		return nil, fmt.Errorf("unsupported operator: %s", cond.Operator)
	}
	return &result, nil
}
