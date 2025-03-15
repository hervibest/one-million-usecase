package usecase

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"github.com/hervibest/one-million-usecase/internal/repository"
)

var (
	totalWorker = 90
	dataHeaders = []string{
		"GlobalRank",
		"TldRank",
		"Domain",
		"TLD",
		"RefSubNets",
		"RefIPs",
		"IDN_Domain",
		"IDN_TLD",
		"PrevGlobalRank",
		"PrevTldRank",
		"PrevRefSubNets",
		"PrevRefIPs",
	}
)

type UploadUseCase interface {
	UploadFile(file string) error
}

type uploadUseCase struct {
	domainRepository repository.DomainRepository
}

func NewUploadUseCase(domainRepository repository.DomainRepository) UploadUseCase {
	return &uploadUseCase{domainRepository: domainRepository}
}

func (u *uploadUseCase) UploadFile(filePath string) error {
	start := time.Now()
	csvReader, csvFile, err := openCsvFile(filePath)
	if err != nil {
		return fmt.Errorf("gagal membuka file CSV: %w", err)
	}
	defer csvFile.Close()

	jobs := make(chan []interface{})
	wg := new(sync.WaitGroup)

	go u.dispatchWorker(wg, jobs)
	u.readPerLineThenSendToWorker(csvReader, wg, jobs)

	wg.Wait()

	duration := time.Since(start)
	fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")
	os.Remove(filePath)
	return nil

}

// Fungsi openCsvFile yang diperbaiki
func openCsvFile(filePath string) (*csv.Reader, *os.File, error) {
	log.Println("=> open csv file")

	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("file majestic_million.csv tidak ditemukan. silakan download terlebih dahulu di https://blog.majestic.com/development/majestic-million-csv-daily")
		}

		return nil, nil, err
	}

	reader := csv.NewReader(f)
	return reader, f, nil
}

func (u *uploadUseCase) readPerLineThenSendToWorker(reader *csv.Reader, wg *sync.WaitGroup, jobs chan<- []interface{}) {
	isHeader := true
	for {
		value := make([]interface{}, 0)
		rows, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
		}

		if isHeader {
			isHeader = false
			continue
		}

		for _, row := range rows {
			value = append(value, row)
		}

		wg.Add(1)
		jobs <- value
	}

	close(jobs)
}

func (u *uploadUseCase) dispatchWorker(wg *sync.WaitGroup, jobs <-chan []interface{}) {
	for workerIndex := 0; workerIndex <= totalWorker; workerIndex++ {
		go func(workerIndex int, wg *sync.WaitGroup, jobs <-chan []interface{}) {
			counter := 0

			for job := range jobs {
				u.doTheJob(workerIndex, job, counter)
				wg.Done()
				counter++
			}

		}(workerIndex, wg, jobs)
	}

}

func (u *uploadUseCase) doTheJob(workerIndex int, values []interface{}, counter int) {

	for {
		var outerError error
		func(outerError *error) {
			defer func() {
				if err := recover(); err != nil {
					*outerError = fmt.Errorf("%v", err)
				}
			}()

			ctx := context.Background()
			if err := u.domainRepository.Insert(ctx, dataHeaders, values); err != nil {
				log.Fatal(err.Error())
			}

		}(&outerError)
		if outerError == nil {
			break
		}
	}

	if counter%100 == 0 {
		log.Println("=> worker", workerIndex, "inserted", counter, "data")
	}
}
