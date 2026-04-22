/*
  Warnings:

  - You are about to drop the column `time` on the `Overtime` table. All the data in the column will be lost.
  - Added the required column `endDate` to the `Overtime` table without a default value. This is not possible if the table is not empty.
  - Added the required column `endTime` to the `Overtime` table without a default value. This is not possible if the table is not empty.
  - Added the required column `startDate` to the `Overtime` table without a default value. This is not possible if the table is not empty.
  - Added the required column `startTime` to the `Overtime` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "Overtime" DROP COLUMN "time",
ADD COLUMN     "duration" INTEGER NOT NULL DEFAULT 0,
ADD COLUMN     "endDate" TIMESTAMP(3) NOT NULL,
ADD COLUMN     "endTime" VARCHAR(255) NOT NULL,
ADD COLUMN     "startDate" TIMESTAMP(3) NOT NULL,
ADD COLUMN     "startTime" VARCHAR(255) NOT NULL;

-- AlterTable
ALTER TABLE "Reimburse" ADD COLUMN     "creatorRole" VARCHAR(255) DEFAULT 'EMPLOYEE',
ADD COLUMN     "updatedByRole" VARCHAR(255);
