package test

type ReportEventType uint8

const (
	REPORT_START_SUITE ReportEventType = iota
	REPORT_FINISH_SUITE
	REPORT_START_CASE
	REPORT_FINISH_CASE
)

type ReportEvent struct {
	SuiteReport *SuiteReport
	CaseReport  *CaseReport
	Type        ReportEventType
}

func NewSuiteReportEvent(suiteReport *SuiteReport, typ ReportEventType) *ReportEvent {
	return &ReportEvent{
		SuiteReport: suiteReport,
		Type:        typ,
	}
}

func NewCaseReportEvent(caseReport *CaseReport, typ ReportEventType) *ReportEvent {
	return &ReportEvent{
		CaseReport: caseReport,
		Type:       typ,
	}
}
