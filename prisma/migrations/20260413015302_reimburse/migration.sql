-- CreateTable
CREATE TABLE "Reimburse" (
    "id" TEXT NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "amount" INTEGER NOT NULL,
    "description" VARCHAR(255) NOT NULL,
    "status" VARCHAR(255) NOT NULL DEFAULT 'SUBMITTED',
    "attachBill" VARCHAR(255) NOT NULL,
    "paymentProof" VARCHAR(255) NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "Reimburse_pkey" PRIMARY KEY ("id")
);
