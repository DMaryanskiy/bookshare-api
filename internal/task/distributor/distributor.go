package distributor

import (
	"context"
	"encoding/json"
	"log"

	"github.com/DMaryanskiy/bookshare-api/internal/task"
	"github.com/hibiken/asynq"
)

type TaskDistributor struct {
	Client *asynq.Client
}

func NewTaskDistributor(redisAddr string) *TaskDistributor {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	return &TaskDistributor{Client: client}
}

func (d *TaskDistributor) DistributeVerificationEmail(ctx context.Context, payload task.PayloadSendVerificationEmail) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(task.TaskSendVerificationEmail, data)
	info, err := d.Client.EnqueueContext(ctx, task)
	if err != nil {
		return err
	}

	log.Println("Task was enqueued:")
	log.Printf("UUID: %+v\nType: %+v\nPayload: %+v\n", info.ID, info.Type, string(info.Payload))
	return nil
}
