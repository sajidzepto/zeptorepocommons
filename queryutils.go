package zeptobaserepo

type Query struct {
}

type QueryConfig struct {
	fields string
}

type AssoicationConfigs struct {
}

type SearcCondition struct {
	operatorClause   string // = , < , >
	inClause         string // in
	betweenClause    string // between
	orderBy          OrderBy
	paginationClause string
}

type OrderBy struct {

	//  order by
}

type PaginatorQueryResult struct {
	values     interface{}
	nextOffset int
}
