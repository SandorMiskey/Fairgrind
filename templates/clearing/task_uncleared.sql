-- Project and users --------------------------------------------------------- {{{
-- Clearing template user #1: [[[.DB_TEMPLATE_USER1]]]
-- Clearing template user #2: [[[.DB_TEMPLATE_USER2]]]

-- Clearing template project: [[[.DB_TEMPLATE_PROJECT]]]

-- --------------------------------------------------------------------------- }}}
-- clearing_tasks update ----------------------------------------------------- {{{

-- INSERT INTO [[[.DB_SCHEMA]]].clearing_tasks
--   (project_id, user_id, clearing_task_type_id, clearing_task_status_id, label, description)
-- VALUES
--   ([[[.DB_TEMPLATE_PROJECT]]], [[[.DB_TEMPLATE_USER1]]], 1, 1, 'Clearing test', 'multipier=1, withdrawable=NULL');
-- SET @last_id = LAST_INSERT_ID();

-- --------------------------------------------------------------------------- }}}
-- clearing_batch_type_id update --------------------------------------------- {{{

-- batch type {{{

INSERT INTO [[[.DB_SCHEMA]]].clearing_batches
  (project_id, user_id, clearing_batch_type_id, clearing_batch_status_id, label, description)
VALUES
  ([[[.DB_TEMPLATE_PROJECT]]], [[[.DB_TEMPLATE_USER1]]], 3, 2, 'Clearing test', 'multipier=1, withdrawable=NULL');
SET @last_id = LAST_INSERT_ID();

-- expected: MSG_SUBROUTING_MULTIPLIERS_MISMATCH
UPDATE [[[.DB_SCHEMA]]].clearing_batches
SET clearing_batch_type_id = 2  -- multiplier=1
WHERE id = @last_id;

-- expected: MSG_SUBROUTING_MULTIPLIERS_MISMATCH
UPDATE [[[.DB_SCHEMA]]].clearing_batches
SET clearing_batch_type_id = 4  -- multiplier=0
WHERE id = @last_id;

-- }}}
-- batch status {{{

-- INSERT INTO [[[.DB_SCHEMA]]].clearing_batches
--   (project_id, user_id, clearing_batch_type_id, clearing_batch_status_id, label, description)
-- VALUES
--   ([[[.DB_TEMPLATE_PROJECT]]], [[[.DB_TEMPLATE_USER1]]], 3, 6, 'Clearing test', 'multipier=1, withdrawable=1');
-- SET @last_id = LAST_INSERT_ID();

-- }}}

-- --------------------------------------------------------------------------- }}}

