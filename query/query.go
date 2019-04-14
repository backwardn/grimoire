package query

// Query defines information about query generated by query builder.
type Query struct {
	empty        bool // todo: use bit to mark what is updated and use it when building
	Collection   string
	SelectClause SelectClause
	JoinClause   []JoinClause
	WhereClause  FilterClause
	GroupClause  GroupClause
	SortClause   []SortClause
	OffsetClause Offset
	LimitClause  Limit
}

func (q Query) Build(query *Query) {
	if query.empty {
		*query = q
	} else {
		// manual merge
		if q.Collection != "" {
			query.Collection = q.Collection
		}

		if q.SelectClause.Fields != nil {
			query.SelectClause = q.SelectClause
		}

		query.JoinClause = append(query.JoinClause, q.JoinClause...)

		query.WhereClause = query.WhereClause.And(q.WhereClause)

		if q.GroupClause.Fields != nil {
			query.GroupClause = q.GroupClause
		}

		q.SortClause = append(q.SortClause, query.SortClause...)

		if q.OffsetClause != 0 {
			query.OffsetClause = q.OffsetClause
		}

		if q.LimitClause != 0 {
			query.LimitClause = q.LimitClause
		}
	}
}

// Select filter fields to be selected from database.
func (q Query) Select(fields ...string) Query {
	q.SelectClause = NewSelect(fields...)
	return q
}

func (q Query) From(collection string) Query {
	q.Collection = collection

	if len(q.SelectClause.Fields) == 0 {
		q.SelectClause = NewSelect(collection + ".*")
	}

	// if len(q.JoinClause) > 0 {
	// 	for i := range q.JoinClause {

	// 	}
	// }

	return q
}

func (q Query) Distinct() Query {
	q.SelectClause.OnlyDistinct = true
	return q
}

// Join current collection with other collection.
func (q Query) Join(collection string) Query {
	return q.JoinOn(collection, "", "")
}

// Join current collection with other collection.
func (q Query) JoinOn(collection string, from string, to string) Query {
	return q.JoinWith("JOIN", collection, from, to)
}

// JoinWith current collection with other collection with custom join mode.
func (q Query) JoinWith(mode string, collection string, from string, to string) Query {
	NewJoinWith(mode, collection, from, to).Build(&q) // TODO: ensure this always called last

	return q
}

func (q Query) JoinFragment(expr string, args ...interface{}) Query {
	NewJoinFragment(expr, args...).Build(&q) // TODO: ensure this always called last

	return q
}

func (q Query) Where(filters ...FilterClause) Query {
	q.WhereClause = q.WhereClause.And(filters...)
	return q
}

func (q Query) OrWhere(filters ...FilterClause) Query {
	q.WhereClause = q.WhereClause.Or(FilterAnd(filters...))
	return q
}

func (q Query) Group(fields ...string) Query {
	q.GroupClause.Fields = fields
	return q
}

func (q Query) Having(filters ...FilterClause) Query {
	q.GroupClause.Filter = q.GroupClause.Filter.And(filters...)
	return q
}

func (q Query) OrHaving(filters ...FilterClause) Query {
	q.GroupClause.Filter = q.GroupClause.Filter.Or(FilterAnd(filters...))
	return q
}

func (q Query) Sort(fields ...string) Query {
	return q.SortAsc(fields...)
}

func (q Query) SortAsc(fields ...string) Query {
	sorts := make([]SortClause, len(fields))
	for i := range fields {
		sorts[i] = NewSortAsc(fields[i])
	}

	q.SortClause = append(q.SortClause, sorts...)
	return q
}

func (q Query) SortDesc(fields ...string) Query {
	sorts := make([]SortClause, len(fields))
	for i := range fields {
		sorts[i] = NewSortDesc(fields[i])
	}

	q.SortClause = append(q.SortClause, sorts...)
	return q
}

// Offset the result returned by database.
func (q Query) Offset(offset Offset) Query {
	q.OffsetClause = offset
	return q
}

// Limit result returned by database.
func (q Query) Limit(limit Limit) Query {
	q.LimitClause = limit
	return q
}

// From create query for collection.
func From(collection string) Query {
	return Query{
		Collection:   collection,
		SelectClause: NewSelect(collection + ".*"),
	}
}

// Join current collection with other collection.
func Join(collection string) Query {
	return JoinOn(collection, "", "")
}

// JoinOn current collection with other collection.
func JoinOn(collection string, from string, to string) Query {
	return JoinWith("JOIN", collection, from, to)
}

// JoinWith current collection with other collection with custom join mode.
func JoinWith(mode string, collection string, from string, to string) Query {
	return Query{
		JoinClause: []JoinClause{
			NewJoinWith(mode, collection, from, to),
		},
	}
	// var q Query
	// NewJoinWith(mode, collection, from, to).Build(&q) // TODO: ensure this always called last

	// return q
}

func JoinFragment(expr string, args ...interface{}) Query {
	return Query{
		JoinClause: []JoinClause{
			NewJoinFragment(expr, args...),
		},
	}

	// var q Query
	// NewJoinFragment(expr, args...).Build(&q) // TODO: ensure this always called last

	// return q
}

func Where(filters ...FilterClause) Query {
	return Query{
		WhereClause: FilterAnd(filters...),
	}
}

func Group(fields ...string) Query {
	return Query{
		GroupClause: NewGroup(fields...),
	}
}