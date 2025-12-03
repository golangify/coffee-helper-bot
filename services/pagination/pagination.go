package pagination

type Pagination struct {
	CurrentPage int  `json:"current_page"`
	Total       int  `json:"total"`
	TotalPages  int  `json:"total_pages"`
	HasPrev     bool `json:"has_prev"`
	HasNext     bool `json:"has_next"`
	PrevPage    int  `json:"prev_page"`
	NextPage    int  `json:"next_page"`
}

func Paginate(currentPage, limit, total int) (*Pagination, int, int) {
	if currentPage < 1 {
		currentPage = 1
	}

	if limit <= 0 {
		limit = 1 // предотвращаем деление на ноль
	}

	pagination := &Pagination{
		CurrentPage: currentPage,
		Total:       total,
		PrevPage:    0, // 0 означает отсутствие страницы
		NextPage:    0, // 0 означает отсутствие страницы
	}

	// Вычисляем общее количество страниц
	if total == 0 {
		pagination.TotalPages = 0
	} else {
		pagination.TotalPages = (total + limit - 1) / limit
	}

	// Корректируем текущую страницу, если она превышает общее количество
	if pagination.TotalPages > 0 && pagination.CurrentPage > pagination.TotalPages {
		pagination.CurrentPage = pagination.TotalPages
	}

	// Проверяем наличие предыдущей страницы
	if pagination.CurrentPage > 1 {
		pagination.HasPrev = true
		pagination.PrevPage = pagination.CurrentPage - 1
	}

	// Проверяем наличие следующей страницы
	if pagination.CurrentPage < pagination.TotalPages {
		pagination.HasNext = true
		pagination.NextPage = pagination.CurrentPage + 1
	}

	// Вычисляем offset для SQL запросов
	offset := (pagination.CurrentPage - 1) * limit

	return pagination, offset, limit
}
