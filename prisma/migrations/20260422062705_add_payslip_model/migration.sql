-- CreateTable
CREATE TABLE "Payslip" (
    "id" TEXT NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "fileUrl" VARCHAR(255) NOT NULL,
    "month" VARCHAR(255) NOT NULL,
    "year" VARCHAR(255) NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "Payslip_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "Payslip" ADD CONSTRAINT "Payslip_email_fkey" FOREIGN KEY ("email") REFERENCES "Employee"("email") ON DELETE RESTRICT ON UPDATE CASCADE;
