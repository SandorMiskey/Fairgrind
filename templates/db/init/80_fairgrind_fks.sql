--
-- leger_clearing_task_id -> tasks_id
--

-- ALTER TABLE `{{.DB_SCHEMA}}`.`clearing_ledger`
-- ADD CONSTRAINT `clearing_ledger_clearing_task_id_fk`
-- FOREIGN KEY (`clearing_task_id`)
-- REFERENCES `clearing_tasks` (`id`)
-- ON UPDATE CASCADE;
