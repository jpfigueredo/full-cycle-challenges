package internal

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

func (r Report) Print() {
	fmt.Println("===== Stress Test Report =====")
	fmt.Printf("Total Requests: %d\n", r.TotalRequests)
	fmt.Printf("Successful (200): %d\n", r.SuccessCount)
	fmt.Printf("Duration: %s\n", r.Duration)

	fmt.Println("\nLatency Metrics:")
	fmt.Printf("  Min:   %s\n", r.Min)
	fmt.Printf("  Max:   %s\n", r.Max)
	fmt.Printf("  Mean:  %s\n", r.Mean)
	fmt.Printf("  P95:   %s\n", r.P95)
	fmt.Printf("  P99:   %s\n", r.P99)

	fmt.Println("\nStatus Codes Distribution:")
	for code, count := range r.StatusCodes {
		fmt.Printf("  %d : %d\n", code, count)
	}
}

// JSON output
func (r Report) PrintJSON() {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		fmt.Println("Error generating JSON:", err)
		return
	}
	fmt.Println(string(data))
}

// CSV output
func (r Report) PrintCSV() {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	// Cabeçalho
	writer.Write([]string{
		"TotalRequests", "Successful", "Duration",
		"Min", "Max", "Mean", "P95", "P99", "StatusCodes",
	})

	// Converte status codes para string
	status := ""
	for code, count := range r.StatusCodes {
		status += fmt.Sprintf("%d:%d ", code, count)
	}

	// Linha de dados
	writer.Write([]string{
		fmt.Sprintf("%d", r.TotalRequests),
		fmt.Sprintf("%d", r.SuccessCount),
		r.Duration.String(),
		r.Min.String(),
		r.Max.String(),
		r.Mean.String(),
		r.P95.String(),
		r.P99.String(),
		status,
	})
}
