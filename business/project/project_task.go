package project

import (
	"context"
	"fmt"
	"github.com/kfchen81/beego"
	"github.com/kfchen81/beego/vanilla"
	m_project "teamdo/models/project"
)

type Task struct {
	vanilla.EntityBase
	Id 			int
	Title       string
	TunnelId	int
	ExecutorId  int
	projectName string

	Status 		bool
	IsDelete	bool
	Remark		string
	Priority	string
}

func NewTask(
	ctx context.Context,
	title string,
	tunnelId int,
	executorId int,
	status 	bool,
	remark	string,
	priority string,
	projectName string,
	) *Task {

	o := vanilla.GetOrmFromContext(ctx)

	model := m_project.Task{}
	model.Title = title
	model.ExecutorId = executorId
	model.TunnelId = tunnelId
	model.Status = status
	model.Remark = remark
	model.ProjectName = projectName
	model.Priority = priority

	id, err := o.Insert(&model)
	if err != nil {
		beego.Error(err)
		panic(vanilla.NewBusinessError("create_task_fail", fmt.Sprintf("创建任务失败")))
	}
	model.Id = int(id)

	return NewTaskForModel(ctx, &model)

}

func NewTaskForModel(ctx context.Context, dbModel *m_project.Task) *Task {
	instance := new(Task)
	instance.Ctx = ctx
	instance.Id = dbModel.Id
	instance.Title = dbModel.Title
	instance.Status = dbModel.Status
	instance.Priority = dbModel.Priority
	instance.Remark = dbModel.Remark
	instance.projectName = dbModel.ProjectName
	return instance
}
