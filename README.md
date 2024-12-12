# GTaskMg

A Goroutine and Channel Management System.

```go
package main

import (
    "flag"
    "fmt"

    gtask "github.com/lion187chen/gTaskMg"
)

type MGCore struct {
    *gtask.GManager
    *gtask.GTask
}

var Core MGCore

func main() {
    Core.GManager = new(gtask.GManager).Init()
    Core.GTask = Core.CreateTask(Core.demoTask, "Demo.Task", 16)

    Core.GManager.EnQueueSync("Demo.Task", "Message")                   // Send message to named queue
    Core.GTask.GQueue.EnQueue(gtask.GMSG_EXIT, 100*time.Millisecond)    // OR use my.GQueue
    go Core.Run().(func(string))("Hello")

    Core.GManager.Join()
    fmt.Println("All task done!")
}

func (my *MGCore) demoTask(name string) {
    fmt.Println(name)
    for {
        msg := <-my.GQueue          // OR use my.DeQueue()/my.DeQueueSync()
        switch tmsg := msg.(type) {
        case string:
            switch tmsg {
            case gtask.GMSG_EXIT:
                Core.GTask.Exit()
                return
            default:
                fmt.Println(tmsg)
            }
        }
    }
}
```
