package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/ashish-dhakane/curlsc/internal/models"
)

type Renderer interface {
	Render([]models.CheckResult, models.Summary) error
}

type TableRenderer struct {
	Out     io.Writer
	NoColor bool
}

func (r *TableRenderer) Render(results []models.CheckResult, summary models.Summary) error {
	w := tabwriter.NewWriter(r.Out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "STATUS\tCODE\tTIME(ms)\tURL")

	for _, res := range results {
		statusStr := ""
		colorCode := ""
		resetCode := ""

		if !r.NoColor {
			resetCode = "\033[0m"
			switch {
			case res.Status == models.StatusUp:
				colorCode = "\033[32m" // Green
				statusStr = "✓ UP"
			case res.Status == models.StatusDown:
				colorCode = "\033[33m" // Yellow
				statusStr = "✗ DOWN"
			case res.Status == models.StatusError:
				colorCode = "\033[31m" // Red
				statusStr = "✗ ERR"
			}
		} else {
			switch res.Status {
			case models.StatusUp:
				statusStr = "✓ UP"
			case models.StatusDown:
				statusStr = "✗ DOWN"
			case models.StatusError:
				statusStr = "✗ ERR"
			}
		}

		codeStr := "---"
		if res.StatusCode > 0 {
			codeStr = strconv.Itoa(res.StatusCode)
		}

		errSuffix := ""
		if res.ErrorType != "" {
			errSuffix = fmt.Sprintf("  [%s]", res.ErrorType)
		}

		fmt.Fprintf(w, "%s%s%s\t%s\t%d\t%s%s\n", colorCode, statusStr, resetCode, codeStr, res.ResponseTimeMs, res.URL, errSuffix)
	}
	w.Flush()

	fmt.Fprintf(r.Out, "\nSummary: %d checked | %d up | %d failed | %s elapsed\n", 
		summary.Total, summary.Success, summary.Failure, summary.ElapsedTime.Round(time.Millisecond))
	
	return nil
}

type JSONRenderer struct {
	Out io.Writer
}

func (r *JSONRenderer) Render(results []models.CheckResult, summary models.Summary) error {
	encoder := json.NewEncoder(r.Out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

type CSVRenderer struct {
	Out io.Writer
}

func (r *CSVRenderer) Render(results []models.CheckResult, summary models.Summary) error {
	w := csv.NewWriter(r.Out)
	defer w.Flush()

	w.Write([]string{"url", "status", "statusCode", "responseTimeMs", "finalUrl", "errorType"})
	for _, res := range results {
		w.Write([]string{
			res.URL,
			string(res.Status),
			strconv.Itoa(res.StatusCode),
			strconv.FormatInt(res.ResponseTimeMs, 10),
			res.FinalURL,
			res.ErrorType,
		})
	}
	return nil
}

func GetRenderer(format string, out io.Writer, noColor bool) Renderer {
	switch format {
	case "json":
		return &JSONRenderer{Out: out}
	case "csv":
		return &CSVRenderer{Out: out}
	default:
		return &TableRenderer{Out: out, NoColor: noColor}
	}
}
