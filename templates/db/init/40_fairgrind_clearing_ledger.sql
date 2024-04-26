--
-- Clearing ledger statuses
--

DROP TABLE IF EXISTS `[[[.DB_SCHEMA]]]`.`clearing_ledger_statuses`;

CREATE TABLE `[[[.DB_SCHEMA]]]`.`clearing_ledger_statuses` (
	`id` SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
	`label` TINYTEXT NOT NULL,
	`description` TINYTEXT DEFAULT NULL,
	`withdrawable` TINYINT(1) NOT NULL DEFAULT 0,
	`created_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(),
	`updated_at`  DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
	`deleted_at` DATETIME(3) NULL DEFAULT NULL,
	PRIMARY KEY (`id`),
	KEY `clearing_ledger_statuses_deleted_at_idx` (`deleted_at`) USING BTREE
) [[[.DB_TABLE_OPTIONS]]];

INSERT INTO `[[[.DB_SCHEMA]]]`.`clearing_ledger_statuses`
	(label, description, withdrawable)
VALUES
	('Pending', 'Transaction credited but not yet withdrawable.', 0),
	('Withdrawable', 'Credit is withdrawable', 1);

--
-- Clearing ledger types
--

DROP TABLE IF EXISTS `[[[.DB_SCHEMA]]]`.`clearing_ledger_labels`;

CREATE TABLE `[[[.DB_SCHEMA]]]`.`clearing_ledger_labels` (
	`id` SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
	`label` TINYTEXT NOT NULL,
	`description` TINYTEXT DEFAULT NULL,
	`created_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(),
	`updated_at`  DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
	`deleted_at` DATETIME(3) NULL DEFAULT NULL,
	PRIMARY KEY (`id`),
	KEY `clearing_ledger_labels_deleted_at_idx` (`deleted_at`) USING BTREE
) [[[.DB_TABLE_OPTIONS]]];

INSERT INTO `[[[.DB_SCHEMA]]]`.`clearing_ledger_labels`
	(label, description)
VALUES
	('Undefined', 'No label given to the transaction.'),
	('Task', 'Credit processed due to annotation or review.'),
	('Airdrop', 'Registration reward'),
	('Referral', 'Registration referral reward'),
	('Swap', 'Credit has been swapped'),
	('Bridge', 'Credit has been bridged to external wallet.'),
	('Withdraw', 'Credit has been withdrew.');

--
-- Clearing ledger
--

DROP TABLE IF EXISTS `[[[.DB_SCHEMA]]]`.`clearing_ledger`;

CREATE TABLE `[[[.DB_SCHEMA]]]`.`clearing_ledger` (
	`id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
	`user_id` INT(11) NOT NULL,
	-- Handling of references is TBD
	-- `referred_user_id` INT(11) DEFAULT NULL,
	`clearing_task_id` BIGINT UNSIGNED DEFAULT NULL,
	`reference` TEXT DEFAULT NULL,
	`clearing_ledger_status_id` SMALLINT UNSIGNED NOT NULL,
	`clearing_ledger_label_id` SMALLINT UNSIGNED NOT NULL DEFAULT 1,
	`amount` DECIMAL(36,18) NOT NULL DEFAULT 0,
	`clearing_token_id` MEDIUMINT UNSIGNED NOT NULL,
	`created_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(),
	`updated_at`  DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
	`deleted_at` DATETIME(3) NULL DEFAULT NULL,
	PRIMARY KEY (`id`),
	KEY `clearing_ledger_deleted_at_idx` (`deleted_at`) USING BTREE,
	CONSTRAINT `clearing_ledger_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE,
	-- CONSTRAINT `clearing_ledger_referred_users_id_fk` FOREIGN KEY (`referred_user_id`) REFERENCES `users` (`id`)  ON UPDATE CASCADE,
	-- CONSTRAINT `clearing_ledger_clearing_task_id_fk` FOREIGN KEY (`clearing_task_id`) REFERENCES `clearing_tasks` (`id`) ON UPDATE CASCADE,
	CONSTRAINT `clearing_ledger_clearing_token_id_fk` FOREIGN KEY (`clearing_token_id`) REFERENCES `clearing_tokens` (`id`) ON UPDATE CASCADE,
	CONSTRAINT `clearing_ledger_clearing_ledger_status_id_fk` FOREIGN KEY (`clearing_ledger_status_id`) REFERENCES `clearing_ledger_statuses` (`id`) ON UPDATE CASCADE,
	CONSTRAINT `clearing_ledger_clearing_ledger_label_id_fk` FOREIGN KEY (`clearing_ledger_label_id`) REFERENCES `clearing_ledger_labels` (`id`) ON UPDATE CASCADE
) [[[.DB_TABLE_OPTIONS]]];

