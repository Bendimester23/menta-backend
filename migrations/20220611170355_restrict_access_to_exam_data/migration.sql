-- AlterTable
ALTER TABLE "Exam" ADD COLUMN     "show" BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN     "showResults" BOOLEAN NOT NULL DEFAULT false;
