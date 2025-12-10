-- +migrate Up
CREATE TABLE `tasks` (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `user_id` INT NOT NULL,
    `category_id` INT NOT NULL,
    `title` VARCHAR(191) NOT NULL,
    `description` TEXT,
    `due_date` DATETIME DEFAULT NULL,
    `priority` ENUM('low', 'medium', 'high') DEFAULT 'low',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT `fk_tasks_user` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`),
    CONSTRAINT `fk_tasks_category` FOREIGN KEY (`category_id`) REFERENCES `categories`(`id`)
);

-- +migrate Down
DROP TABLE `tasks`;
