-- AlterTable
ALTER TABLE "Feedback" ADD COLUMN     "createdBy" VARCHAR(255);

-- AlterTable
ALTER TABLE "Overtime" ADD COLUMN     "isAdditionalLeave" BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN     "isDeployment" BOOLEAN NOT NULL DEFAULT false;

-- AlterTable
ALTER TABLE "TemplateKPI" ALTER COLUMN "attachment" DROP NOT NULL;

-- CreateTable
CREATE TABLE "KPIItem" (
    "id" TEXT NOT NULL,
    "templateId" TEXT NOT NULL,
    "nameResult" VARCHAR(255) NOT NULL,
    "kpiResult" TEXT NOT NULL,
    "weight" DOUBLE PRECISION NOT NULL,
    "target" DOUBLE PRECISION NOT NULL,
    "actual" DOUBLE PRECISION NOT NULL,
    "score" DOUBLE PRECISION NOT NULL,
    "finalScore" DOUBLE PRECISION NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "KPIItem_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "KPIItem" ADD CONSTRAINT "KPIItem_templateId_fkey" FOREIGN KEY ("templateId") REFERENCES "TemplateKPI"("id") ON DELETE CASCADE ON UPDATE CASCADE;
