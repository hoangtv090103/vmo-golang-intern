package reports

type ReportUsecase struct {
	ReportRepo ReportRepository
}

func NewReportUsecase(reportRepo ReportRepository) *ReportUsecase {
	return &ReportUsecase{
		ReportRepo: reportRepo,
	}
}

func (r *ReportUsecase) GenerateReport(startDate, endDate string) (string, error) {
	return r.ReportRepo.GenerateReport(startDate, endDate)
}
