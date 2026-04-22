-- CreateTable
CREATE TABLE "Feedback" (
    "id" TEXT NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "employeeName" VARCHAR(255) NOT NULL,
    "employeeEmail" VARCHAR(255) NOT NULL,
    "employeeDepartment" VARCHAR(255) NOT NULL,
    "month" VARCHAR(255) NOT NULL,
    "isAnonymouse" BOOLEAN NOT NULL DEFAULT true,
    "positiveExperience" VARCHAR(255) NOT NULL,
    "suggestion" VARCHAR(255) NOT NULL DEFAULT 'SUBMITTED',
    "workEnvironment" INTEGER NOT NULL,
    "workQualityReliability" INTEGER NOT NULL,
    "collaborationCommunication" INTEGER NOT NULL,
    "workLifeBalance" INTEGER NOT NULL,
    "criticalThinking" INTEGER NOT NULL,
    "overallSatisfaction" INTEGER NOT NULL,
    "score" INTEGER NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "Feedback_pkey" PRIMARY KEY ("id")
);
