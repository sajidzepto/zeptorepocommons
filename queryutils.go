package zeptobaserepo

import "strings"

type QueryConfig struct {
	fields string
}

type AssociationConfigs struct {
}

type Condition interface {
	getPreparedStatement() (string, interface{})
}

type SearchCondition struct {
	orConditions []OrConditions
	orderBy      *OrderBy
	offset       int
}

func (s SearchCondition) getPreparedStatement() (string, interface{}) {
	sb := strings.Builder{}
	var arguments []interface{}
	for i, orC := range s.orConditions {
		for j, andC := range orC.andConditions {
			for _, condition := range andC.conditions {
				stm, args := condition.getPreparedStatement()
				sb.WriteString(stm)
				arguments = append(arguments, args)
			}
			if j < len(orC.andConditions)-1 {
				sb.WriteString(" AND ")
			}

		}
		if i < len(s.orConditions)-1 {
			sb.WriteString(" OR ")
		}
	}
	return sb.String(), arguments
}

type AndConditions struct {
	conditions []Condition
}

type OrConditions struct {
	andConditions []AndConditions
}

// db.Where("name = ?", "jinzhu").First(&user)
type SearchOperatorCondition struct {
	field    string
	operator string
	value    interface{}
}

func (s SearchOperatorCondition) getPreparedStatement() (string, interface{}) {
	sb := strings.Builder{}
	sb.WriteString(s.field + " ")
	sb.WriteString(s.operator + " ? ")
	return sb.String(), s.value
}

// db.Where("name IN ?", []string{"jinzhu", "jinzhu 2"}).Find(&users)
type SearchInCondition struct {
	field  string
	values []interface{}
}

func (s SearchInCondition) getPreparedStatement() (string, interface{}) {
	sb := strings.Builder{}
	sb.WriteString(s.field + " ")
	sb.WriteString("IN" + " ? ")
	return sb.String(), s.values
}

// db.Where("name LIKE ?", "%jin%").Find(&users)
type SearchLikeCondition struct {
	field string
	regex string
}

func (s SearchLikeCondition) getPreparedStatement() (string, interface{}) {

	sb := strings.Builder{}
	sb.WriteString(s.field + " ")
	sb.WriteString("LIKE" + " ? ")
	return sb.String(), s.regex
}

//db.Where("created_at BETWEEN ? AND ?", lastWeek, today).Find(&users)
type SearchBetweenCondition struct {
	field      string
	lowerValue interface{}
	upperValue interface{}
}

func (s SearchBetweenCondition) getPreparedStatement() (string, interface{}) {
	sb := strings.Builder{}
	sb.WriteString(s.field + " ")
	sb.WriteString("BETWEEN" + " ? AND ? ")
	return sb.String(), []interface{}{s.lowerValue, s.upperValue}
}

type OrderBy struct {
	order map[string]string
}

type PaginatorQueryResult struct {
	values     interface{}
	nextOffset int
}
