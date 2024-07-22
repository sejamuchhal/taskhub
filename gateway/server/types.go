package server

type CreateTaskRequest struct {
	Title       string `form:"title" json:"title" binding:"required"`
	Description string `form:"description" json:"description"`
	DueDateTime string `form:"due_date_time" json:"due_date_time"`
}

type CreateTaskResponse struct {
	ID string `json:"id"`
}

type ListTasksRequest struct {
	Limit   int  `form:"limit"`
	Offset  int  `form:"offset"`
	Pending bool `form:"pending"`
}

type TaskDetails struct {
	Title       string
	Description string
	Status      string
	UserId      string
	DueDate     string
	CreatedAt   string
	UpdatedAt   string
}

type GetTaskResponse struct {
	Task TaskDetails
}

type ListTasksResponse struct {
	Count int           `json:"count"`
	Tasks []*TaskDetails `json:"tasks"`
}

type UpdateTaskRequest struct {
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
	DueDateTime string `form:"due_date_time" json:"due_date_time"`
}
