--
-- Clearing batch statuses
--

DROP TABLE IF EXISTS `[[[.DB_SCHEMA]]]`.`clearing_batch_statuses`;

CREATE TABLE `[[[.DB_SCHEMA]]]`.`clearing_batch_statuses` (
  `id` SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `label` TINYTEXT NOT NULL,
  `description` TINYTEXT DEFAULT NULL,
  `clearing_ledger_status_id` SMALLINT UNSIGNED DEFAULT NULL,
  `created_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(),
  `updated_at`  DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
  `deleted_at` DATETIME(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `clearing_batch_types_deleted_at_idx` (`deleted_at`) USING BTREE,
  CONSTRAINT `clearing_batch_statuses_clearing_ledger_status_id_fk` FOREIGN KEY (`clearing_ledger_status_id`) REFERENCES `clearing_ledger_statuses` (`id`) ON UPDATE CASCADE
) [[[.DB_TABLE_OPTIONS]]];

INSERT INTO `[[[.DB_SCHEMA]]]`.`clearing_batch_statuses`
  (label, description, clearing_ledger_status_id)
VALUES
  ('Unpublished', 'The batch is being prepared and/or awaiting deposit payment.', NULL),
  ('In progress', 'The batch is currently being processed, the deposit has been paid.', NULL),
  ('Suspended', 'Deposit paid, but the processing of the batch is suspended.', NULL),
  ('Canceled', 'The processing of the batch has been canceled.', NULL),
  ('Approved', 'All tasks are customer approved or canceled status.', 1),
  ('Cleared', 'Settled in full, the total cost of the bundle is covered.', 2);

--
-- Clearing batch types
--

DROP TABLE IF EXISTS `[[[.DB_SCHEMA]]]`.`clearing_batch_types`;

CREATE TABLE `[[[.DB_SCHEMA]]]`.`clearing_batch_types` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `label` TINYTEXT NOT NULL,
  `description` TINYTEXT DEFAULT NULL,
  `multiplier` FLOAT NOT NULL DEFAULT 1,
  `created_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(),
  `updated_at`  DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
  `deleted_at` DATETIME(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `clearing_batch_types_deleted_at_idx` (`deleted_at`) USING BTREE
) [[[.DB_TABLE_OPTIONS]]];

INSERT INTO [[[.DB_SCHEMA]]].clearing_batch_types 
  (label,description,multiplier)
VALUES
  ('Exam', 'To be used during the testing or training of a Grinder', 0.0),
  ('Test', 'To be Used When Examining the efficiency and/or clarity of an annotation task', 1.0),
  ('Production', 'Production batch', 1.0),
  ('Training', 'Training batch', 0.0);

--
-- Clearing batches
--

DROP TABLE IF EXISTS `[[[.DB_SCHEMA]]]`.`clearing_batches`;

CREATE TABLE `[[[.DB_SCHEMA]]]`.`clearing_batches` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `project_id` INT(11) NOT NULL,
  `user_id` INT(11) NOT NULL,
  `clearing_batch_type_id` BIGINT UNSIGNED NOT NULL,
  `clearing_batch_status_id` SMALLINT UNSIGNED NOT NULL,
  `label` TINYTEXT NOT NULL,
  `description` TEXT DEFAULT NULL,
  `created_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(),
  `updated_at`  DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
  `deleted_at` DATETIME(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `clearing_batches_deleted_at_idx` (`deleted_at`) USING BTREE,
  CONSTRAINT `clearing_batches_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)  ON UPDATE CASCADE,
  CONSTRAINT `clearing_batches_clearing_batch_types_fk` FOREIGN KEY (`clearing_batch_type_id`) REFERENCES `clearing_batch_types` (`id`) ON UPDATE CASCADE,
  CONSTRAINT `clearing_batches_clearing_batch_statuses_fk` FOREIGN KEY (`clearing_batch_status_id`) REFERENCES `clearing_batch_statuses` (`id`) ON UPDATE CASCADE,
  CONSTRAINT `clearing_batches_projects_fk` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON UPDATE CASCADE
) [[[.DB_TABLE_OPTIONS]]];

