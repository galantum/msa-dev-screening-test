package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

// Dataset represents the structure of the API response
// Updated to match the provided API response format
type DataSet struct {
	Year                   string `json:"year"`
	University             string `json:"university"`
	School                 string `json:"school"`
	Degree                 string `json:"degree"`
	EmploymentRateOverall  string `json:"employment_rate_overall"`
	EmploymentRateFtPerm   string `json:"employment_rate_ft_perm"`
	BasicMonthlyMean       string `json:"basic_monthly_mean"`
	BasicMonthlyMedian     string `json:"basic_monthly_median"`
	GrossMonthlyMean       string `json:"gross_monthly_mean"`
	GrossMonthlyMedian     string `json:"gross_monthly_median"`
	GrossMthly25Percentile string `json:"gross_mthly_25_percentile"`
	GrossMthly75Percentile string `json:"gross_mthly_75_percentile"`
}

// Change limit to 2000 because total data from response api is 1262
const apiURL = "https://api-production.data.gov.sg/v2/internal/api/datasets/d_3c55210de27fcccda2ed0c63fdd2b352/rows?limit=2000"

func fetchData() ([]DataSet, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed fetch data: %s", resp.Status)
	}

	var result struct {
		Code int `json:"code"`
		Data struct {
			Rows []DataSet `json:"rows"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data.Rows, nil
}

func writeCSV(outputDir string, year string, data []DataSet) error {
	filePath := fmt.Sprintf("%s/%s.csv", outputDir, year)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{
		"Year", "University", "School", "Degree", "Employment Rate Overall", "Employment Rate FT Perm", "Basic Monthly Mean", "Basic Monthly Median", "Gross Monthly Mean", "Gross Monthly Median", "Gross Mthly 25 Percentile", "Gross Mthly 75 Percentile",
	})

	// Write data
	for _, record := range data {
		if record.Year == year {
			row := []string{
				record.Year,
				record.University,
				record.School,
				record.Degree,
				record.EmploymentRateOverall,
				record.EmploymentRateFtPerm,
				record.BasicMonthlyMean,
				record.BasicMonthlyMedian,
				record.GrossMonthlyMean,
				record.GrossMonthlyMedian,
				record.GrossMthly25Percentile,
				record.GrossMthly75Percentile,
			}
			writer.Write(row)
		}
	}

	return writer.Error()
}

func worker(years <-chan string, data []DataSet, outputDir string, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()

	for year := range years {
		mu.Lock()
		if err := writeCSV(outputDir, year, data); err != nil {
			fmt.Printf("Error writing CSV for year %s: %v\n", year, err)
		}
		mu.Unlock()
	}
}

func main() {
	// define value concurent limit & output directory
	concurrentLimit := 2
	outputDir := "/home/yourname/csv"

	data, err := fetchData()
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return
	}

	// debug
	//fmt.Printf("%+v\n", data)

	// Define years which in dataset for filtering data
	yearSet := make(map[string]struct{})
	for _, record := range data {
		yearSet[record.Year] = struct{}{}
	}

	// debug
	//fmt.Printf("%+v\n", yearSet)

	// Define years which in dataset for knowing total year
	years := make([]string, 0, len(yearSet))
	for year := range yearSet {
		years = append(years, year)
	}

	// debug
	//fmt.Printf("%+v\n", years)

	yearChan := make(chan string, len(years))
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Start workers
	for i := 0; i < concurrentLimit; i++ {
		wg.Add(1)
		go worker(yearChan, data, outputDir, &wg, &mu)
	}

	// Send years to workers
	for _, year := range years {
		yearChan <- year
	}
	close(yearChan)

	// Wait for all workers to finish
	wg.Wait()

	fmt.Println("Data processing completed.")
}
