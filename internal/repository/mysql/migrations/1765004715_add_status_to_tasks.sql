-- +migrate Up
ALTER TABLE `tasks`
ADD COLUMN `status` ENUM('pending', 'inprogress', 'done')
DEFAULT 'pending'
AFTER `priority`;

-- +migrate Down
ALTER TABLE `tasks`
DROP COLUMN `status`;