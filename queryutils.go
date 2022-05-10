package zeptorepocommons

import "strings"

type Query struct {
	queryCondition *QueryCondition
	pageConfig     *PageConfig
}

type QueryCondition struct {
	orConditions []OrConditions
}

type PageConfig struct {
	orderBy map[string]string
	offset  int
	limit   int
}

type PaginatorQueryResult struct {
	values     interface{}
	nextOffset int
}

type Condition interface {
	getPreparedStatement() (string, interface{})
}

type AndConditions struct {
	conditions []Condition
}

type OrConditions struct {
	andConditions []AndConditions
}

type SearchOperatorCondition struct {
	field    string
	operator string
	value    interface{}
}

type SearchInCondition struct {
	field  string
	values []interface{}
}

type SearchLikeCondition struct {
	field string
	regex string
}

type SearchBetweenCondition struct {
	field      string
	lowerValue interface{}
	upperValue interface{}
}

func (s SearchOperatorCondition) getPreparedStatement() (string, interface{}) {
	sb := strings.Builder{}
	sb.WriteString(s.field + " ")
	sb.WriteString(s.operator + " ? ")
	return sb.String(), s.value
}

func (s SearchInCondition) getPreparedStatement() (string, interface{}) {
	sb := strings.Builder{}
	sb.WriteString(s.field + " ")
	sb.WriteString("IN" + " ? ")
	return sb.String(), s.values
}

func (s SearchLikeCondition) getPreparedStatement() (string, interface{}) {

	sb := strings.Builder{}
	sb.WriteString(s.field + " ")
	sb.WriteString("LIKE" + " ? ")
	return sb.String(), s.regex
}

func (s SearchBetweenCondition) getPreparedStatement() (string, interface{}) {
	sb := strings.Builder{}
	sb.WriteString(s.field + " ")
	sb.WriteString("BETWEEN" + " ? AND ? ")
	return sb.String(), []interface{}{s.lowerValue, s.upperValue}
}

func (s *QueryCondition) getPreparedStatement() (string, interface{}) {
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
