package time

import(
		stdtime "time"
)


func InTheRange(input, from,to stdtime.Time)(bool){
	return (input.After(from) && input.Before(to))
}
