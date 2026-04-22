package main
import (
	"context"
	"fmt"
	"hrsync-backend/internal/db"
	"os"
)
func main() {
	client := db.NewClient(db.WithDatasourceURL(os.Getenv("DATABASE_URL")))
	if err := client.Connect(); err != nil { panic(err) }
	defer client.Disconnect()
	employees, _ := client.Employee.FindMany().Exec(context.Background())
	fmt.Println("--- Employees ---")
	for _, e := range employees {
		fmt.Printf("Name: [%s], Email: [%s]\n", e.Name, e.Email)
	}
	payslips, _ := client.Payslip.FindMany().Exec(context.Background())
	fmt.Println("--- Payslips ---")
	for _, p := range payslips {
		fmt.Printf("Email: [%s], File: [%s]\n", p.Email, p.FileURL)
	}
}
