package reports

import (
	"encoding/json"
	"net/http"
)

type DateRange struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type ReportHandler struct {
	ReportUsecase ReportUsecase
}

func NewReportHandler(reportUsecase ReportUsecase) ReportHandler {
	return ReportHandler{
		ReportUsecase: reportUsecase,
	}
}

func (h ReportHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	// Generate Excel report
	switch r.Method {
	case http.MethodPost:
		{
			// Get start date and end date
			var dateRange DateRange
			var path string
			err := json.NewDecoder(r.Body).Decode(&dateRange)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Generate Excel report
			path, err = h.ReportUsecase.GenerateReport(dateRange.StartDate, dateRange.EndDate)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Store path in response
			err = json.NewEncoder(w).Encode(map[string]string{"path": path})
			if err != nil {
				return
			}

			// Send path to client
			w.Header().Set("Content-Type", "application/json")

		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}
