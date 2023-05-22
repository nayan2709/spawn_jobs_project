package spawn_jobs

import (
	"errors"
	"fmt"
	"github.com/dunzoit/projects/nirwana_spwanJob/entities"
	"github.com/dunzoit/projects/nirwana_spwanJob/jobber"
	"math"
	"time"
)

type SpawnJobs struct {
	Jobber     *jobber.Jobber
	spawnChan  chan int
	errChan    chan entities.ErrorChanStruct
	CurrSpawns map[int]int
}

type SpawnJobsClient interface {
	SpawnJobs(payloads []entities.Payload, fanout int) error
}

func NewSpawnJobsClient() SpawnJobs {
	jobberStatus := make(map[int]chan bool)
	spawChan := make(chan int, 1)
	errChan := make(chan entities.ErrorChanStruct, 1)
	return SpawnJobs{
		Jobber:     jobber.NewJobber(jobberStatus),
		spawnChan:  spawChan,
		errChan:    errChan,
		CurrSpawns: make(map[int]int),
	}
}

func (j *SpawnJobs) SpawnJobs(payload []entities.Payload, fanout int) error {
	for i := 0; i < len(payload); i++ {
		// 0 means not spawned, 1 means spawned, 2 means completed, 3 means cancel due to error
		j.CurrSpawns[i] = 0
	}

	fanout = int(math.Min(float64(fanout), float64(len(payload))))
	for i := 0; i < fanout; i++ {
		fmt.Println("Spawning: ", i, " payload: ", payload[i])
		j.CurrSpawns[i] = 1
		go j.spawnAndWaitForJobCompletion(payload[i])
	}
	if fanout != len(payload) {
		fmt.Println("fanout: ", fanout)
		i := fanout
		for {
			select {
			case err := <-j.errChan:
				fmt.Println(fmt.Sprintf("Error while executing Job:%v, err:%v", err.JobId, err.Error))
				return errors.New("error scheduling job")
			case jobId := <-j.spawnChan:
				fmt.Println(fmt.Sprintf("Job:%v completed successfully", jobId))
				fmt.Println("Spawning: ", i, " payload: ", payload[i])
				j.CurrSpawns[jobId] = 2
				j.CurrSpawns[i] = 1
				go j.spawnAndWaitForJobCompletion(payload[i])
				i++
			}
			if i == len(payload) {
				break
			}
		}
	}
	return nil
}

func (j *SpawnJobs) spawnAndWaitForJobCompletion(payload entities.Payload) {
	fmt.Println("start spawnAndWaitForJobCompletion for payload", payload)
	jobID, err := j.Jobber.SpawnJob(payload)
	if err != nil {
		j.errChan <- entities.ErrorChanStruct{Error: err, JobId: jobID}
		j.CancelJobs()
		return
	}
	err = j.Jobber.WaitForJobCompletion(jobID, 2*time.Second)
	if err != nil {
		j.errChan <- entities.ErrorChanStruct{Error: err, JobId: jobID}
		j.CancelJobs()
		return
	}
	j.spawnChan <- jobID
	fmt.Println("end spawnAndWaitForJobCompletion for payload", payload)
}

func (j *SpawnJobs) CancelJobs() {
	for i := 0; i < len(j.CurrSpawns); i++ {
		if j.CurrSpawns[i] == 1 {
			j.Jobber.CancelJob(i)
			j.CurrSpawns[i] = 3
		}
	}
}
