package handlers

type ListQuery struct {
	CurrentPage int    `form:"currentPage" binding:"required"`
	PerPage     int    `form:"perPage" binding:"required"`
	Name        string `form:"name"`
}

