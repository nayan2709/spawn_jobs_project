package main

import (
	"fmt"
	"github.com/dunzoit/projects/nirwana_spwanJob/entities"
	"github.com/dunzoit/projects/nirwana_spwanJob/spawn_jobs"
)

func main() {
	// // Redirect Output
	//outputFile, err := os.OpenFile("output.txt", 'a', 0666)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer outputFile.Close()
	//os.Stdout = outputFile

	fmt.Println("😎.Starting spawning Jobs.😎")
	fanout := 5
	payloads := []entities.Payload{"m1", "m2", "m3", "m4", "m5", "m6", "m7", "m8", "m9", "m10"}
	spawnJobs := spawn_jobs.NewSpawnJobsClient()
	err := spawnJobs.SpawnJobs(payloads, fanout)
	if err != nil {
		fmt.Println("Error while spawning Jobs 🥲: ", err)
		return
	}
	fmt.Println("😎.Successfully Executed spawning.😎")
}
