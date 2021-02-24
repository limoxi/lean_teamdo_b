package task

import (
	"github.com/kfchen81/beego/vanilla"
	b_task "teamdo/business/task"
)

type TaskExecutor struct {
	vanilla.RestResource
}

func (this *TaskExecutor) Resource() string {
	return "project.task_executor"
}

func (this *TaskExecutor) GetParameters() map[string][]string {
	return map[string][]string{
		"POST": []string{
			"task_id:int",
			"user_id:int",
		},
	}
}

func (this *TaskExecutor) Post()  {
	bCtx := this.GetBusinessContext()

	id, _ := this.GetInt("task_id")
	userId, _ := this.GetInt("user_id")
	task := b_task.NewTaskRepository(bCtx).GetTaskById(id)
	task.UpdateExecutor(userId)

	response := vanilla.MakeResponse(vanilla.Map{})
	this.ReturnJSON(response)
}