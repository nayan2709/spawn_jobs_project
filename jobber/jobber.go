package jobber

import (
	"errors"
	"fmt"
	"github.com/dunzoit/projects/nirwana_spwanJob/entities"
	"github.com/dunzoit/projects/nirwana_spwanJob/mockServer"
	"sync"
	"time"
)

type Jobber struct {
	jobId        int
	JobberStatus map[int]chan bool
	mutex        sync.Mutex
}

func NewJobber(jobberStatus map[int]chan bool) *Jobber {
	var mutex sync.Mutex
	return &Jobber{
		jobId:        0,
		JobberStatus: jobberStatus,
		mutex:        mutex,
	}
}

type JobberMethods interface {
	SpawnJob(payload entities.Payload) (jobId int, err error)
	WaitForJobCompletion(jobId int, timeout time.Duration) error
	CancelJob(jobId int)
}

func (j *Jobber) SpawnJob(payload entities.Payload) (jobId int, err error) {
	j.mutex.Lock() // Acquire the lock before accessing the variable
	id := j.jobId
	j.jobId++
	j.JobberStatus[id] = make(chan bool, 1)
	j.mutex.Unlock() // Release the lock after modifying the variable

	// Do api call or any other execution
	err = mockServer.NirwanaMockServer(payload)
	if err != nil {
		fmt.Println("Error while executing Job:", id, " err:", err)
		return id, err
	}

	//fmt.Println("Job completed", id)
	j.JobberStatus[id] <- true
	return id, nil
}

func (j *Jobber) WaitForJobCompletion(jobId int, timeout time.Duration) error {
	if _, ok := j.JobberStatus[jobId]; !ok {
		fmt.Println("Job not found-> either it's executed or timeout", jobId)
		return errors.New("job not found")
	}
	t := time.NewTicker(1 * time.Second)
	timeoutChan := make(chan int, 1)
	var c int
	for {
		select {
		case <-j.JobberStatus[jobId]:
			j.mutex.Lock()
			delete(j.JobberStatus, jobId)
			j.mutex.Unlock()
			return nil
		case <-t.C:
			c++
		case <-timeoutChan:
			fmt.Println("-->> timeout: ", jobId, " c:", c)
			j.mutex.Lock()
			delete(j.JobberStatus, jobId)
			j.mutex.Unlock()
			return errors.New("timeout")
		}
		if c == int(timeout.Seconds()) {
			timeoutChan <- 1
			t.Stop()
		}
	}
}

func (j *Jobber) CancelJob(jobId int) {
	if _, ok := j.JobberStatus[jobId]; !ok {
		fmt.Println("Job not found-> either it's executed or timeout")
		return
	}
	j.mutex.Lock()
	delete(j.JobberStatus, jobId)
	j.mutex.Unlock()
}
