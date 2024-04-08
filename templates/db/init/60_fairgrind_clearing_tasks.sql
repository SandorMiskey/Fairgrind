--
-- Clearing task statuses
--
DROP TABLE IF EXISTS `[[[.DB_SCHEMA]]]`.`clearing_task_statuses`;

CREATE TABLE `[[[.DB_SCHEMA]]]`.`clearing_task_statuses` (
  `id` SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `label` TINYTEXT NOT NULL,
  `description` TINYTEXT DEFAULT NULL,
  `payable` TINYINT(1) DEFAULT 1,
  `parent_payable` TINYINT(1) DEFAULT 1,
  `created_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(),
  `updated_at`  DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
  `deleted_at` DATETIME(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `clearing_task_statuses_deleted_at_idx` (`deleted_at`) USING BTREE
) [[[.DB_TABLE_OPTIONS]]];

INSERT INTO `[[[.DB_SCHEMA]]]`.`clearing_task_statuses`
  (label, description, payable, parent_payable)
VALUES
  ('Awaiting review', 'Grinder has completed the annotation or uploaded raw data and is awaiting review.', 0, 0),
  ('Reviewer approved', 'The annotation or data upload performed by Grinder has been accepted by reviewer.', 0, 0),
  ('Customer approved', 'The annotation or data upload performed by Grinder has been accepted by the customer.', 1, 1),
  ('Reviewer rejected', 'The annotation or data upload performed by Grinder has been rejected by reviewer.', 0, 0),
  ('Customer rejected', 'The annotation or data upload performed by Grinder has been rejected by the customer.', 0, 0);

--
-- Clearing task types
--

DROP TABLE IF EXISTS `[[[.DB_SCHEMA]]]`.`clearing_task_types`;

CREATE TABLE `[[[.DB_SCHEMA]]]`.`clearing_task_types` (
  `id` SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `label` TINYTEXT NOT NULL,
  `description` TINYTEXT DEFAULT NULL,
  `task_payable` TINYINT(1) DEFAULT 1,
  `subtask_payable` TINYINT(1) DEFAULT 1,
  -- `review_payable` TINYINT(1) DEFAULT 0,
  `created_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(),
  `updated_at`  DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
  `deleted_at` DATETIME(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `clearing_task_types_deleted_at_idx` (`deleted_at`) USING BTREE
) [[[.DB_TABLE_OPTIONS]]];

INSERT INTO `[[[.DB_SCHEMA]]]`.`clearing_task_types`
  (label, description, task_payable, subtask_payable)
VALUES
  ('Annotation', 'Grinder annotated a raw data.', 1, 1),
  ('Data collection', 'Grinder collected and uploaded raw data.', 1, 1),
  ('Review', 'Grinder reviewed an annotation or upload.', 1, 1);

--
-- Clearing task fees
--

DROP TABLE IF EXISTS `[[[.DB_SCHEMA]]]`.`clearing_task_fees`;

CREATE TABLE `[[[.DB_SCHEMA]]]`.`clearing_task_fees` (
  -- `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` INT(11) NOT NULL,
  `project_id` INT(11) NOT NULL,
  `clearing_task_type_id` SMALLINT UNSIGNED NOT NULL,
  `task_fee` FLOAT NOT NULL DEFAULT 0,
  `subtask_fee` FLOAT NOT NULL DEFAULT 0,
  `clearing_token_id` MEDIUMINT UNSIGNED NOT NULL,
  `created_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(),
  `updated_at`  DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
  `deleted_at` DATETIME(3) NULL DEFAULT NULL,
  -- PRIMARY KEY (`id`),
  PRIMARY KEY (`user_id`, `project_id`, `clearing_task_type_id`, `clearing_token_id`),
  KEY `clearing_task_fees_deleted_at_idx` (`deleted_at`) USING BTREE,
  CONSTRAINT `clearing_task_fees_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)  ON UPDATE CASCADE,
  CONSTRAINT `clearing_task_fees_project_id_fk` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON UPDATE CASCADE,
  CONSTRAINT `clearing_task_fees_clearing_task_type_id_fk` FOREIGN KEY (`clearing_task_type_id`) REFERENCES `clearing_task_types` (`id`) ON UPDATE CASCADE,
  CONSTRAINT `clearing_task_fees_clearing_token_id_fk` FOREIGN KEY (`clearing_token_id`) REFERENCES `clearing_tokens` (`id`) ON UPDATE CASCADE
) [[[.DB_TABLE_OPTIONS]]];

--
-- Clearing tasks
--

DROP TABLE IF EXISTS `[[[.DB_SCHEMA]]]`.`clearing_tasks`;

CREATE TABLE `fairgrind`.`clearing_tasks` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `clearing_batch_id` BIGINT UNSIGNED NOT NULL, 
  `clearing_task_id` BIGINT UNSIGNED DEFAULT NULL, 
  `clearing_task_type_id` SMALLINT UNSIGNED NOT NULL,
  `clearing_task_status_id` SMALLINT UNSIGNED NOT NULL,
  -- `input` JSON DEFAULT NULL,
  -- `output` JSON NOT NULL DEFAULT '[]',
  `output` JSON DEFAULT NULL, 
  `task_id` INT(11) NOT NULL,
  `user_id` INT(11) NOT NULL,
  `task_reject_issues_id` INT(11) DEFAULT NULL,
  `reference` MEDIUMTEXT DEFAULT NULL,
  `started_at` DATETIME(3) DEFAULT NULL,
  `finished_at` DATETIME(3) DEFAULT NULL,
	`cleared_at` DATETIME(3) DEFAULT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT current_timestamp(),
  `updated_at` DATETIME(3) DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `clearing_tasks_deleted_at_idx` (`deleted_at`) USING BTREE,
  CONSTRAINT `clearing_tasks_clearing_batch_id_fk` FOREIGN KEY (`clearing_batch_id`) REFERENCES `clearing_batches` (`id`) ON UPDATE CASCADE,
  CONSTRAINT `clearing_tasks_clearing_task_id_fk` FOREIGN KEY (`clearing_task_id`) REFERENCES `clearing_tasks` (`id`) ON UPDATE CASCADE,
  CONSTRAINT `clearing_tasks_clearing_task_type_id_fk` FOREIGN KEY (`clearing_task_type_id`) REFERENCES `clearing_task_types` (`id`) ON UPDATE CASCADE,
  CONSTRAINT `clearing_tasks_clearing_task_status_id_fk` FOREIGN KEY (`clearing_task_status_id`) REFERENCES `clearing_task_statuses` (`id`) ON UPDATE CASCADE,
  CONSTRAINT `clearing_tasks_task_id_fk` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`)  ON UPDATE CASCADE,
  CONSTRAINT `clearing_tasks_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)  ON UPDATE CASCADE,
  CONSTRAINT `task_reject_issues_id_fk` FOREIGN KEY (`task_reject_issues_id`) REFERENCES `task_reject_issues` (`id`)  ON UPDATE CASCADE
) [[[.DB_TABLE_OPTIONS]]];

