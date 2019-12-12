package consume

import (
	jsoniter "github.com/json-iterator/go"

	"go.uber.org/zap"
)

type (
	//Job 工作項目
	Job interface {
		Preprocess() error
		Do() error
	}

	//JobChan 工作通道
	JobChan chan Job

	//WorkPool 工作池
	WorkPool struct {
		Pool chan Job
		Size int
	}

	//WorkerChan 工作者通道
	WorkerChan chan chan Job
	//Worker 工作者
	Worker struct {
		//Task <-chan chan Job
		*WorkPool
		quit chan bool
	}

	//Dispatcher 配適器 ,設定工作者數量
	Dispatcher struct {
		*WorkPool
		MaxWorker int
		quit      chan bool
	}
)

//NewWorkPool work pool instance
//field worker chan , Pool size
func NewWorkPool(size int) *WorkPool {
	return &WorkPool{
		Pool: make(JobChan, size),
		Size: size,
	}
}

//Preprocess preprocess queue message
func (wp *WorkPool) Preprocess(data []byte, job Job) error {

	if err := jsoniter.Unmarshal(data, job); err != nil {
		return err
	}
	if err := job.Preprocess(); err != nil {
		return err
	}

	return nil
}

//Put put int work pool
func (wp *WorkPool) Put(job Job) {
	wp.Pool <- job
}

//NewWorker worker instance
func NewWorker(pool *WorkPool) *Worker {
	return &Worker{
		WorkPool: pool,
		quit:     make(chan bool),
	}
}

//NewDispatcher 建立配適器 , 設定工作者數量
//@numberOfWorkers 工作者數量
func NewDispatcher(numberOfWorkers int) *Dispatcher {

	d := &Dispatcher{
		MaxWorker: numberOfWorkers,
		quit:      make(chan bool),
	}

	return d
}

//Run 運行配適器 , 建立工作池 ,建立worker
//worker 開始監聽workpool
//@poolSize 建立pool size
func (d *Dispatcher) Run(poolSize int) {

	d.WorkPool = NewWorkPool(poolSize)

	for i := 0; i < d.MaxWorker; i++ {
		worker := NewWorker(d.WorkPool)
		worker.Start()
	}
}

//Start 開始工作
func (w *Worker) Start() {
	go func() {
		for job := range w.Pool {
			if err := job.Do(); err != nil {
				zap.L().Error(err.Error() + " of Work Do")
			}
		}
	}()
}
