package filter

// Page holds the result of applying pagination to a slice of DiffResults.
type Page[T any] struct {
	Items      []T
	Total      int
	Offset     int
	Limit      int
	HasMore    bool
}

// Paginate returns a Page from items using the given offset and limit.
// A limit of 0 means no pagination (all items are returned).
func Paginate[T any](items []T, offset, limit int) Page[T] {
	total := len(items)

	if limit <= 0 {
		return Page[T]{
			Items:   items,
			Total:   total,
			Offset:  0,
			Limit:   0,
			HasMore: false,
		}
	}

	if offset < 0 {
		offset = 0
	}

	if offset >= total {
		return Page[T]{
			Items:   []T{},
			Total:   total,
			Offset:  offset,
			Limit:   limit,
			HasMore: false,
		}
	}

	end := offset + limit
	hasMore := end < total
	if end > total {
		end = total
	}

	return Page[T]{
		Items:   items[offset:end],
		Total:   total,
		Offset:  offset,
		Limit:   limit,
		HasMore: hasMore,
	}
}
