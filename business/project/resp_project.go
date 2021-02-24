package project

type RProject struct {
	Id 		int	`json:"id"`
	Name 	string `json:"name"`
	ManagerId	int	`json:"manager_id"`
}

type RTunnel struct {
	Id 	int		`json:"id"`
	Title	string	`json:"title"`
	DisplayIndex	int	`json:"display_index"`
	Tasks 	[]*RTask 	`json:"tasks"`
}

type RTask struct {
	Id int 		`json:"id"`
	ProjectName string	`json:"project_name"`
	Title string	`json:"title"`
	Status bool 	`json:"status"`
	Remark string 	`json:"remark"`
	Priority string `json:"priority"`
	StartDate string	`json:"start_date"`
	EndDate	string	`json:"end_date"`
	ExecutorId int 	`json:"executor_id"`
}