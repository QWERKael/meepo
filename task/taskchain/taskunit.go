package taskchain

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"meepo/abandon/tasks"
)

type TaskUnit struct {
	TaskId      string            // 任务Id
	TaskName    string            // 用户定义的任务名
	TaskInfo    string            // 任务备注
	TaskType    string            // 任务类型
	TaskFromPID *actor.PID        // 任务来源的PID
	Task        *tasks.PluginTask // 任务
}

func NewTaskUnit(TaskId string, TaskName string, TaskInfo string, TaskType string) *TaskUnit {
	return &TaskUnit{TaskId,
		TaskName,
		TaskInfo,
		TaskType,
		nil,
		&tasks.PluginTask{},
	}
}

func (tu *TaskUnit) TaskUnitDisplay() {
	fmt.Printf("TaskId: %s\nTaskName: %s\nTaskInfo: %s\nTaskType: %s\n", tu.TaskId, tu.TaskName, tu.TaskInfo, tu.TaskType)
}

func (tu *TaskUnit) Exec() error {
	err := tu.Task.Exec()
	if err != nil {
		return err
	}
	return nil
}
