package main

import (
	"fmt"
	"hrsync-backend/internal/db"
	"hrsync-backend/internal/handler"
	"hrsync-backend/internal/middleware"
	"hrsync-backend/internal/repository"
	"hrsync-backend/internal/router"
	"hrsync-backend/internal/service"
	"hrsync-backend/internal/utils"
	"context"

	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	utils.InitMinio()

	url := os.Getenv("DATABASE_URL")
	client := db.NewClient(db.WithDatasourceURL(url))

	if err := client.Prisma.Connect(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			log.Fatalf("failed to disconnect from database: %v", err)
		}
	}()

	// Wire up dependencies

	// Auth
	authRepo := repository.NewAuthRepository(client)
	authSrv := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authSrv)

	// Employee
	employeeRepo := repository.NewEmployeeRepository(client)
	employeeSrv := service.NewEmployeeService(employeeRepo)
	employeeHandler := handler.NewEmployeeHandler(employeeSrv)

	// Leave
	leaveRepo := repository.NewLeaveRepository(client)
	leaveSrv := service.NewLeaveService(leaveRepo)
	leaveHandler := handler.NewLeaveHandler(leaveSrv)

	leaveTypeRepo := repository.NewLeaveTypeRepository(client)
	leaveTypeSrv := service.NewLeaveTypeService(leaveTypeRepo)
	leaveTypeHandler := handler.NewLeaveTypeHandler(leaveTypeSrv)

	// Overtime
	overtimeRepo := repository.NewOvertimeRepository(client)
	overtimeSrv := service.NewOvertimeService(overtimeRepo)
	overtimeHandler := handler.NewOvertimeHandler(overtimeSrv)

	// Reimburse
	reimburseRepo := repository.NewReimburseRepository(client)
	reimburseSrv := service.NewReimburseService(reimburseRepo)
	reimburseHandler := handler.NewReimburseHandler(reimburseSrv)

	departmentHandler := handler.NewDepartmentHandler()

	// Feedback
	feedbackRepo := repository.NewFeedbackRepository(client)
	feedbackSrv := service.NewFeedbackService(feedbackRepo)
	feedbackHandler := handler.NewFeedbackHandler(feedbackSrv)

	// KPI
	kpiRepo := repository.NewTemplateKPIRepository(client)
	kpiSrv := service.NewTemplateKPIService(kpiRepo)
	kpiHandler := handler.NewTemplateKPIHandler(kpiSrv)

	fileHandler := handler.NewFileHandler()

	// Holiday
	holidayRepo := repository.NewHolidayRepository(client)
	holidaySrv := service.NewHolidayService(holidayRepo)
	holidayHandler := handler.NewHolidayHandler(holidaySrv)

	// Payslip
	payslipRepo := repository.NewPayslipRepository(client)
	payslipSrv := service.NewPayslipService(payslipRepo, employeeRepo)
	payslipHandler := handler.NewPayslipHandler(payslipSrv)

	// Background Sync Holidays for current, previous and next years
	go func() {
		holidaySrv.SyncHolidays(context.Background(), 2024)
		holidaySrv.SyncHolidays(context.Background(), 2025)
		holidaySrv.SyncHolidays(context.Background(), 2026)
	}()

	r := router.NewRouter(authHandler, employeeHandler, departmentHandler, leaveHandler, leaveTypeHandler, overtimeHandler, reimburseHandler, feedbackHandler, kpiHandler, fileHandler, holidayHandler, payslipHandler)

	fmt.Println("Server starting on :8080")
	handler := middleware.CORSMiddleware(r)
	if err := http.ListenAndServe(":8080", handler); err != nil {

		log.Fatal(err)
	}
}
