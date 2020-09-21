package signal

import(
    "os"
    "os/signal"
)

func WaitSignal(signals ...os.Signal){
    sigChan := make(chan os.Signal,0)
    signal.Notify(sigChan,signals...)
    <-sigChan
}
