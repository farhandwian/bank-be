package module

// import (
// 	"bank-backend/module/handler/bank"
// 	"bank-backend/module/handler/user"

// 	"github.com/IBM/sarama"
// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// type InitialInjection struct {
// 	User user.Handler
// 	Bank bank.Handler
// }

// func NewInitialInjection(PGx *pgxpool.Pool, Producer *sarama.SyncProducer) InitialInjection {
// 	return InitialInjection{
// 		User: user.NewHandler(PGx),
// 		Bank: bank.NewHandler(PGx, Producer),
// 	}
// }
