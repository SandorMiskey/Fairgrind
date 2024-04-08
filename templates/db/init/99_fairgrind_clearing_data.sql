--
-- Project and users
--

INSERT INTO fairgrind.users
  (id, first_name, last_name, email, password, birth_year, country_id)
VALUES
  ({{.DB_TEMPLATE_USER1}}, '', 'Clearing template user #1', '', '', 0, 0),
  ({{.DB_TEMPLATE_USER2}}, '', 'Clearing template user #2', '', '', 0, 0);

INSERT INTO fairgrind.projects
  (id, uid, name, description)
VALUES
  ({{.DB_TEMPLATE_PROJECT}}, {{.DB_TEMPLATE_USER1}}, 'Clearing', 'Clearing template project');

--
-- Batches
--

INSERT INTO {{.DB_SCHEMA}}.clearing_batches
  (project_id, user_id, clearing_batch_type_id, clearing_batch_status_id, label, description)
VALUES
  ({{.DB_TEMPLATE_PROJECT}}, {{.DB_TEMPLATE_USER1}}, 1, 2, 'Exam template #1', 'Clearing template batch: type=Exam, status=In progress'),
  ({{.DB_TEMPLATE_PROJECT}}, {{.DB_TEMPLATE_USER1}}, 1, 3, 'Exam template #2', 'Clearing template batch: type=Exam, status=Suspended'),
  ({{.DB_TEMPLATE_PROJECT}}, {{.DB_TEMPLATE_USER1}}, 1, 5, 'Exam template #3', 'Clearing template batch: type=Exam, status=Cleared'),
  ({{.DB_TEMPLATE_PROJECT}}, {{.DB_TEMPLATE_USER1}}, 2, 2, 'Test template #1', 'Clearing template batch: type=Test, status=In progress'),
  ({{.DB_TEMPLATE_PROJECT}}, {{.DB_TEMPLATE_USER1}}, 2, 3, 'Test template #2', 'Clearing template batch: type=Test, status=Suspended'),
  ({{.DB_TEMPLATE_PROJECT}}, {{.DB_TEMPLATE_USER1}}, 2, 5, 'Test template #3', 'Clearing template batch: type=Test, status=Cleared'),
  ({{.DB_TEMPLATE_PROJECT}}, {{.DB_TEMPLATE_USER1}}, 3, 2, 'Prod template #1', 'Clearing template batch: type=Production, status=In progress'),
  ({{.DB_TEMPLATE_PROJECT}}, {{.DB_TEMPLATE_USER1}}, 3, 3, 'Prod template #2', 'Clearing template batch: type=Production, status=Suspended'),
  ({{.DB_TEMPLATE_PROJECT}}, {{.DB_TEMPLATE_USER1}}, 3, 5, 'Prod template #3', 'Clearing template batch: type=Production, status=Cleared');

--
-- Fees
--

-- TODO: Add fees

--
-- Tasks
--

INSERT INTO `fairgrind`.`tasks`
  (batch_id, input_json, status_id, grinder_uid)
VALUES
  (7, '[]', 3, {{.DB_TEMPLATE_USER1}});

INSERT INTO `fairgrind`.`clearing_tasks`
  (clearing_batch_id, clearing_task_type_id, clearing_task_status_id, task_id, user_id)
VALUES
  (7, 1, 2, 511, {{.DB_TEMPLATE_USER1}}),
  (9, 1, 2, 511, {{.DB_TEMPLATE_USER1}});

--
-- Ledger
--

INSERT INTO `{{.DB_SCHEMA}}`.`clearing_ledger`
  (user_id, clearing_task_id, reference, clearing_ledger_status_id, clearing_ledger_label_id, amount, clearing_token_id)
VALUES
  ({{.DB_TEMPLATE_USER1}}, NULL, 'Withdrawable FGDT #1', 2, 3, 500.0, 5),
  ({{.DB_TEMPLATE_USER1}}, NULL, 'Withdrawable FGDT #2', 2, 3, 500.0, 5),
  ({{.DB_TEMPLATE_USER1}}, NULL, 'Withdrew FGDT', 2, 7, -100.0, 5),
  ({{.DB_TEMPLATE_USER1}}, 1, 'Pending FGDT #1', 1, 2, 100.0, 5),
  ({{.DB_TEMPLATE_USER1}}, 2, 'Pending FGDT #2', 1, 2, 50.0, 5),
  ({{.DB_TEMPLATE_USER1}}, 1, 'Pending USD #1', 1, 2, 10.0, 1),
  ({{.DB_TEMPLATE_USER1}}, 1, 'Pending USD #2', 1, 2, 20.0, 1),
  ({{.DB_TEMPLATE_USER2}}, NULL, 'Pending USD #1', 1, 2, 30.0, 1);

