package param


type GetTaskResponse struct {
	Task TaskInfo `json:"task"`
}

type GetTasksResponse struct {
	Tasks []TaskInfo `json:"tasks"`
}