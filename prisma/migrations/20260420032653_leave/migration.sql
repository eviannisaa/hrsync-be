-- AddForeignKey
ALTER TABLE "Leave" ADD CONSTRAINT "Leave_email_fkey" FOREIGN KEY ("email") REFERENCES "Employee"("email") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Overtime" ADD CONSTRAINT "Overtime_email_fkey" FOREIGN KEY ("email") REFERENCES "Employee"("email") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Reimburse" ADD CONSTRAINT "Reimburse_email_fkey" FOREIGN KEY ("email") REFERENCES "Employee"("email") ON DELETE RESTRICT ON UPDATE CASCADE;
