package queue

import (
	"bank-backend/module/bank/entity"
	"bank-backend/pkg"
	"bank-backend/utils"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type ProcessTransferQueue struct {
	Producer sarama.SyncProducer
	Topic    string
}

func NewProcessTransferQueue(producer sarama.SyncProducer, topic string) *ProcessTransferQueue {
	return &ProcessTransferQueue{Producer: producer, Topic: topic}
}

func (q *ProcessTransferQueue) PublishProcessTransferJob(ctx context.Context, request entity.TransferRequest, userPhoneNumber string) (uuid.UUID, string, error) {
	fmt.Println("q.Topic nih:")
	fmt.Println(q.Topic)
	var (
		lvState3       = utils.LogEventStateKafkaPublish
		lfState3Status = "state_2_kafka_publish_status"

		lf = []slog.Attr{
			pkg.LogEventName("bank-service"),
		}
	)
	/*------------------------------------
	| Step 3 : Publish TransferEvent
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState3))

	id, err := pkg.GenerateId()
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogWarnWithContext(ctx, "generate uuid error", err, lf)
		return uuid.UUID{}, "", err
	}
	now := time.Now()
	// Format the time
	formatted := now.Format("2006-01-02 15:04:05.000000")

	event := entity.TransferEvent{
		Transfer:              id.String(),
		Amount:                request.Amount,
		PhoneNumberOriginUser: userPhoneNumber,
		TargetUser:            request.TargetUser,
		Remarks:               request.Remarks,
		CreatedAt:             formatted,
	}
	messageByte, err := json.Marshal(event)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogWarnWithContext(ctx, "kafka publish error", err, lf)
		return uuid.UUID{}, "", err
	}
	err = pkg.PublishMessage(q.Producer, q.Topic, string(messageByte))
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogWarnWithContext(ctx, "kafka publish error", err, lf)
		return uuid.UUID{}, "", err
	}
	return id, event.CreatedAt, nil
}
