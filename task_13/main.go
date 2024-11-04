package main

import (
	"encoding/json"

	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

type ReadWorker struct {
	file_name     string
	output_chanel chan []float64
}

type FilterWorker struct {
	window_size   int
	input_chanel  chan []float64
	output_chanel chan []float64
}

type WriteWorker struct {
	file_name    string
	input_chanel chan []float64
}

func padArray(array []float64, window_size int) []float64 {
	size := len(array)
	if window_size > size {
		log.Fatal("Window size is shorter than array size")
	}

	half_window := int(window_size / 2)
	data := make([]float64, size+2*half_window)
	copy(data[half_window:(half_window+size)], array[:])
	copy(data[:half_window], array[:half_window])
	copy(data[(size+half_window):], array[(size-half_window):])
	log.Debug("Pad array ready!")
	return data
}

func getMeadValue(data []float64) float64 {
	sum := 0.0
	for _, d := range data {
		sum += d
	}
	return sum / float64(len(data))
}

func workReading(worker ReadWorker, wg *sync.WaitGroup) {
	defer wg.Done()
	var array []float64

	content, err := os.ReadFile(worker.file_name)
	if err != nil {
		log.Info(err)
		return
	}

	err = json.Unmarshal(content, &array)
	if err != nil {
		log.Info(err)
		return
	}

	worker.output_chanel <- array
	close(worker.output_chanel)
}

func workFiltering(worker FilterWorker, wg *sync.WaitGroup) {
	defer wg.Done()

	window_size := worker.window_size

	for array := range worker.input_chanel {
		size := len(array)

		if window_size > size {
			continue
		}
		padded_data := padArray(array, window_size)

		filtered_data := make([]float64, size)
		log.Debug("Padded array size:", len(padded_data))
		for i := 0; i < size; i++ {
			filtered_data[i] = getMeadValue(padded_data[i:(i + window_size)])
		}
		worker.output_chanel <- filtered_data
	}
	log.Debug("Filtering ready!")
	close(worker.output_chanel)
}

func workWriting(worker WriteWorker, wg *sync.WaitGroup) {
	defer wg.Done()

	for array := range worker.input_chanel {
		data, err := json.Marshal(array)
		if err != nil {
			continue
		}
		err = os.WriteFile(worker.file_name, data, 0644)
		if err != nil {
			log.Info(err)
			continue
		}
		log.Debug("Writing ready!")
	}
}

func main() {
	log.SetLevel(log.DebugLevel)

	var wg_main sync.WaitGroup
	entries, err := os.ReadDir("./inputs")
	if err != nil {
		log.Info(err)
		return
	}
	for _, e := range entries {
		log.Debug("Entry:", e.Name())
		read_worker := ReadWorker{"./inputs/" + e.Name(), make(chan []float64)}
		filter_worker := FilterWorker{5, read_worker.output_chanel, make(chan []float64)}
		write_worker := WriteWorker{"./outputs/" + e.Name(), filter_worker.output_chanel}
		wg_main.Add(3)
		go workReading(read_worker, &wg_main)
		go workFiltering(filter_worker, &wg_main)
		go workWriting(write_worker, &wg_main)
	}

	wg_main.Wait()
}
