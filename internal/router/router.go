package router

import (
	"hrsync-backend/internal/handler"
	"net/http"
)

func NewRouter(
	authHandler *handler.AuthHandler,
	employeeHandler *handler.EmployeeHandler,
	departmentHandler *handler.DepartmentHandler,
	leaveHandler *handler.LeaveHandler,
	leaveTypeHandler *handler.LeaveTypeHandler,
	overtimeHandler *handler.OvertimeHandler,
	reimburseHandler *handler.ReimburseHandler,
	feedbackHandler *handler.FeedbackHandler,
	kpiHandler *handler.TemplateKPIHandler,
	fileHandler *handler.FileHandler,
	holidayHandler *handler.HolidayHandler,
	payslipHandler *handler.PayslipHandler,
) http.Handler {
	mux := http.NewServeMux()

	// Auth routes (public)
	RegisterAuthRoutes(mux, authHandler, employeeHandler)

	// Resource routes
	RegisterTeamRoutes(mux, employeeHandler)
	RegisterDepartmentRoutes(mux, departmentHandler)
	RegisterLeaveRoutes(mux, leaveHandler, leaveTypeHandler)
	RegisterOvertimeRoutes(mux, overtimeHandler)
	RegisterReimburseRoutes(mux, reimburseHandler)
	RegisterFeedbackRoutes(mux, feedbackHandler)
	RegisterKPIRoutes(mux, kpiHandler)
	RegisterFileRoutes(mux, fileHandler)
	RegisterHolidayRoutes(mux, holidayHandler)
	RegisterPayslipRoutes(mux, payslipHandler)

	return mux
}
