/*
  Warnings:

  - You are about to drop the column `memberIds` on the `Group` table. All the data in the column will be lost.
  - You are about to drop the column `waitingIds` on the `Group` table. All the data in the column will be lost.
  - You are about to drop the `_group_member` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `_waiting_members` table. If the table is not empty, all the data it contains will be lost.
  - A unique constraint covering the columns `[name]` on the table `Group` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[loginCode]` on the table `Group` will be added. If there are existing duplicate values, this will fail.

*/
-- DropForeignKey
ALTER TABLE "Group" DROP CONSTRAINT "Group_memberIds_fkey";

-- DropForeignKey
ALTER TABLE "Group" DROP CONSTRAINT "Group_waitingIds_fkey";

-- DropForeignKey
ALTER TABLE "_group_member" DROP CONSTRAINT "_group_member_A_fkey";

-- DropForeignKey
ALTER TABLE "_group_member" DROP CONSTRAINT "_group_member_B_fkey";

-- DropForeignKey
ALTER TABLE "_waiting_members" DROP CONSTRAINT "_waiting_members_A_fkey";

-- DropForeignKey
ALTER TABLE "_waiting_members" DROP CONSTRAINT "_waiting_members_B_fkey";

-- AlterTable
ALTER TABLE "Group" DROP COLUMN "memberIds",
DROP COLUMN "waitingIds";

-- DropTable
DROP TABLE "_group_member";

-- DropTable
DROP TABLE "_waiting_members";

-- CreateTable
CREATE TABLE "GroupMember" (
    "groupId" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "waiting" BOOLEAN NOT NULL,

    CONSTRAINT "GroupMember_pkey" PRIMARY KEY ("groupId","userId")
);

-- CreateIndex
CREATE UNIQUE INDEX "Group_name_key" ON "Group"("name");

-- CreateIndex
CREATE UNIQUE INDEX "Group_loginCode_key" ON "Group"("loginCode");

-- AddForeignKey
ALTER TABLE "GroupMember" ADD CONSTRAINT "GroupMember_groupId_fkey" FOREIGN KEY ("groupId") REFERENCES "Group"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "GroupMember" ADD CONSTRAINT "GroupMember_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
