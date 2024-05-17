## Find min, mean & max temperature for weather stations

### Generate the data
  - Run the below command to generate the csv file. File will be created in ./data folder
  - python3 create_measurements.py 10000000

### Generate the output
  - Run the below command. output file will be created in ./data
  - go run main.go 

### Approach
1. Create a map with key as station name and value as channel
2. Iterate through the input file and push the temperature values into a particaluar station channel(**Fan Out**)
3. A goroutine will be created for each channel and it consumes the temperature values and find the min, mean and max temperature and publish this into aggregate channel(**Fan In**)
4. Single goroutine will consume from aggregate channel and create a slice, sort it and write to output file

```
go run main.go -input ./data/measurements_100M.txt -output ./data/sorted_records_100M.txt

2024-05-17 09:03:06.142743122 +0200 CEST m=+0.000046738
records: 413
Sorted records written to sorted_records.txt
Execution time: 1m37.671871738s
```