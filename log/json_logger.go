package log

import(
	"io"
	kit_log "github.com/go-kit/kit/log"
)

func NewJsonLogger(out io.Writer, ) Logger_Log{
	return kit_log.With(kit_log.NewJSONLogger(out),"caller",kit_log.Caller(4))
}
