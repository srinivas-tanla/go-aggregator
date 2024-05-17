package main

// measurements.txt 

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"flag"
)


type Record struct {
	Station     string
	AverageTemp float64
	MinValue    float64
	MaxValue    float64
}

var stations = make(map[string]chan float64)
var aggregates = make(chan Record)
var wg sync.WaitGroup
var cg sync.WaitGroup



func main() {
	start := time.Now()
	fmt.Println(start)
	// read the file
	input_file := flag.String("input", "./data/measurements.txt", "input file to process" )
	output_file := flag.String("output", "./data/sorted_records.txt", "output file to process" )
	// Parse the flags
	flag.Parse()

	file, err := os.Open(*input_file)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	go func() {
		cg.Add(1)
		defer cg.Done()

		// fanout 
		var records []Record
		for output := range aggregates {
			records = append(records, output)
		}

		// Sort the records slice based on AverageTemp
		sort.Slice(records, func(i, j int) bool {
			return records[i].AverageTemp < records[j].AverageTemp
		})

		fmt.Println("records:", len(records))

		// Write the sorted records to a file
		file, err := os.Create(*output_file)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		for _, record := range records {
			fmt.Fprintf(file, "%s;%.2f;%.2f;%.2f\n", record.Station, record.MinValue, record.AverageTemp, record.MaxValue)
		}

		fmt.Println("Sorted records written to sorted_records.txt")
	}()
  
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				// close all the channels
				go func() {
					for _, channel := range stations {
						close(channel)
					}
				}()
				break
			}
			fmt.Println("Error:", err)
			return
		}
		parts := strings.Split(line, ";")
		parts[1] = strings.Replace(parts[1], "\n", "", -1)
		floatValue, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			fmt.Println(err)
		}

		if _, ok := stations[parts[0]]; !ok {
			stations[parts[0]] = make(chan float64)
			go func() {
				defer wg.Done()
				wg.Add(1)
				var sumValue float64
				var counter int
				var minValue float64 = 1000
				var maxValue float64 = -1000
				for value := range stations[parts[0]] {
					sumValue += value
					counter++
					minValue = math.Min(minValue, value)
					maxValue = math.Max(maxValue, value)
				}
				aggregates <- Record{
					Station:     parts[0],
					AverageTemp: sumValue / float64(counter),
					MinValue:    minValue,
					MaxValue:    maxValue,
				}
			}()
		}
		// fan out
		stations[parts[0]] <- floatValue
	}

	wg.Wait()
	close(aggregates)
	cg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)
}
