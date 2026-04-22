-- CreateTable
CREATE TABLE "EmployeeOrganization" (
    "id" TEXT NOT NULL,
    "organizationImage" VARCHAR(255),
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "EmployeeOrganization_pkey" PRIMARY KEY ("id")
);
