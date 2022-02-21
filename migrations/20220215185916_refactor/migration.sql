/*
  Warnings:

  - You are about to drop the column `leaderId` on the `Group` table. All the data in the column will be lost.
  - Added the required column `leader` to the `GroupMember` table without a default value. This is not possible if the table is not empty.

*/
-- DropForeignKey
ALTER TABLE "Group" DROP CONSTRAINT "Group_leaderId_fkey";

-- AlterTable
ALTER TABLE "Group" DROP COLUMN "leaderId";

-- AlterTable
ALTER TABLE "GroupMember" ADD COLUMN     "leader" BOOLEAN NOT NULL;
