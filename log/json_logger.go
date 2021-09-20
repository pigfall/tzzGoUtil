package log

import(
	"io"
	kit_log "github.com/go-kit/kit/log"
	timeHelper "github.com/pigfall/tzzGoUtil/time"
	"time"
)

func NewJsonLogger(out io.Writer) Logger_Log{
	 return kit_log.With(kit_log.NewJSONLogger(out),"caller",kit_log.Caller(4),"ts",kit_log.TimestampFormat(func ()time.Time{
		 return time.Now()
	 },timeHelper.Format_yyyy_mm_dd_hh_mm_ss))
	  // return kit_log.With(kit_log.NewJSONLogger(out),"caller",kit_log.Caller(4))
}
