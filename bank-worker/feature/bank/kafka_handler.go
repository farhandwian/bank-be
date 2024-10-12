package bank

import (
	"bank-worker/feature/shared"
	"bank-worker/pkg"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type NewTransferEventHandler struct {
}

func (*NewTransferEventHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) {
	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_message_status"

		lvState2       = shared.LogEventStateInsertDB
		lfState2Status = "state_2_insert_db_status"

		lf = []slog.Attr{
			pkg.LogEventName("Transfer-Worker"),
		}
	)

	fmt.Println("testes kafka")
	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))

	var payload TransferEvent
	err := json.Unmarshal(msg.Value, &payload)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(payload),
	)
	/*------------------------------------
	| Step 2 : Insert transfer Transaction
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	user := User{
		PhoneNumber: payload.PhoneNumberOriginUser,
		Balance:     payload.Amount,
	}
	parse, err := uuid.Parse(payload.TargetUser)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}
	layout := "2006-01-02 15:04:05.000000"
	t, err := time.Parse(layout, payload.CreatedAt)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	_, _, _, _, err = transferTX(ctx, user, parse, payload.Remarks, t, payload.Transfer)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState2Status))

	pkg.LogInfoWithContext(ctx, "success insert user", lf)
}
