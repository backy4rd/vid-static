package util

type TaskQueue struct {
    Queue chan func()
}

func NewTaskQueue() TaskQueue {
    return TaskQueue{
        make(chan func()),
    }
}

func (tq *TaskQueue) Push(task func()) {
    tq.Queue <- task
}

func (tq *TaskQueue) Run() {
    for i := 0; i < 3; i++ {
        go func () {
            for {
                select {
                case task := <-tq.Queue:
                    task();
                }
            }
        }()
    }
}
