--
-- Clearing wallet view - sum
--

DROP VIEW IF EXISTS `{{.DB_SCHEMA}}`.`clearing_wallets_summed_view`;

CREATE VIEW `{{.DB_SCHEMA}}`.`clearing_wallets_summed_view` AS
SELECT
	SUM(clearing_ledger.amount) AS clearing_ledger_amount_sum,
	MAX(clearing_ledger.updated_at) AS clearing_ledger_updated_at_max,
	clearing_ledger.user_id AS clearing_ledger_user_id,
	clearing_ledger_statuses.withdrawable AS clearing_ledger_status_withdrawable,
	clearing_tokens.id AS clearing_token_id,
	clearing_tokens.symbol AS clearing_token_symbol
FROM
    `{{.DB_SCHEMA}}`.`clearing_ledger`
JOIN
    `{{.DB_SCHEMA}}`.`clearing_ledger_statuses` ON clearing_ledger.clearing_ledger_status_id = clearing_ledger_statuses.id
JOIN
	`{{.DB_SCHEMA}}`.`clearing_tokens` ON clearing_ledger.clearing_token_id = clearing_tokens.id
GROUP BY
	clearing_ledger.user_id, clearing_tokens.symbol, clearing_ledger_statuses.withdrawable;

--
-- Clearing wallet view - by project
--

DROP VIEW IF EXISTS `{{.DB_SCHEMA}}`.`clearing_wallets_detailed_view`;


CREATE VIEW `{{.DB_SCHEMA}}`.`clearing_wallets_detailed_view` AS
SELECT
	clearing_ledger.user_id AS clearing_ledger_user_id,
	SUM(clearing_ledger.amount) AS clearing_ledger_amount_sum,
	MAX(clearing_ledger.updated_at) AS clearing_ledger_updated_at_max,
	clearing_ledger_labels.id AS clearing_ledger_label_id,
	clearing_ledger_labels.label clearing_ledger_label_label,
	clearing_ledger_statuses.withdrawable AS clearing_ledger_status_withdrawable,
	COUNT(clearing_tasks.id) AS clearing_tasks_id_count,
	clearing_tokens.symbol AS clearing_token_symbol,
	projects.id AS project_id,
	projects.name AS project_name
FROM
    `{{.DB_SCHEMA}}`.`clearing_ledger` 
JOIN
	`{{.DB_SCHEMA}}`.`clearing_ledger_labels` ON clearing_ledger.clearing_ledger_label_id = clearing_ledger_labels.id
JOIN
	`{{.DB_SCHEMA}}`.`clearing_ledger_statuses` ON clearing_ledger.clearing_ledger_status_id = clearing_ledger_statuses.id
JOIN
	`{{.DB_SCHEMA}}`.`clearing_tokens` ON clearing_ledger.clearing_token_id = clearing_tokens.id
LEFT JOIN
	`{{.DB_SCHEMA}}`.`clearing_tasks` ON clearing_ledger.clearing_task_id = clearing_tasks.id
LEFT JOIN
	`{{.DB_SCHEMA}}`.`clearing_batches` ON clearing_tasks.clearing_batch_id = clearing_batches.id
-- LEFT JOIN
-- 	`{{.DB_SCHEMA}}`.`clearing_batch_types` ON clearing_batches.clearing_batch_type_id = clearing_batch_types.id
LEFT JOIN
	`{{.DB_SCHEMA}}`.`projects` ON clearing_batches.project_id = projects.id
GROUP BY
	clearing_ledger.user_id, clearing_ledger_labels.id, projects.id, clearing_tokens.symbol, clearing_ledger_statuses.withdrawable;
