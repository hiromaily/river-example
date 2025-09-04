package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/riverqueue/river"
)

// EmailArgs はメール送信ジョブの引数です。
// River では JobArgs は Kind() string を実装します。
type EmailArgs struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	// 例として送信リトライのためのメタ情報なども持てます
	RequestedAt time.Time `json:"requested_at"`
}

func (EmailArgs) Kind() string { return "email.send" }

// EmailWorker は EmailArgs を処理します。
type EmailWorker struct {
	river.WorkerDefaults[EmailArgs]
}

func (w *EmailWorker) Work(ctx context.Context, job *river.Job[EmailArgs]) error {
	// 本来はここで実メール送信（SMTP / API）を行う
	fmt.Printf("[EmailWorker] to=%s subject=%q body=%q\n", job.Args.To, job.Args.Subject, job.Args.Body)
	return nil
}
