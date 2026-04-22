-- CreateTable
CREATE TABLE "Leave" (
    "id" TEXT NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "reason" VARCHAR(255) NOT NULL,
    "leaveType" VARCHAR(255) NOT NULL,
    "requestPeriod" TIMESTAMP(3) NOT NULL,
    "status" VARCHAR(255) NOT NULL DEFAULT 'SUBMITTED',

    CONSTRAINT "Leave_pkey" PRIMARY KEY ("id")
);
