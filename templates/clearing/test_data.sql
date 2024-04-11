-- Project and users --------------------------------------------------------- {{{

-- Clearing template user #1: [[[.DB_TEMPLATE_USER1]]]
-- Clearing template user #2: [[[.DB_TEMPLATE_USER2]]]

-- Clearing template project: [[[.DB_TEMPLATE_PROJECT]]]

-- --------------------------------------------------------------------------- }}}
-- clearing_batch_type_id update (clearing/clearing.go:240) ------------------ {{{

-- batch type {{{

INSERT INTO [[[.DB_SCHEMA]]].clearing_batches
  (project_id, user_id, clearing_batch_type_id, clearing_batch_status_id, label, description)
VALUES
  ([[[.DB_TEMPLATE_PROJECT]]], [[[.DB_TEMPLATE_USER1]]], 3, 2, 'CLR #1', 'multipier=1, ledger_status_id=NULL');
SET @last_id = LAST_INSERT_ID();

-- test #1, expected: MSG_SUBROUTING_MULTIPLIERS_MISMATCH (clearing/clearing.go:284)
UPDATE [[[.DB_SCHEMA]]].clearing_batches
SET clearing_batch_type_id = 2  -- multiplier=1
WHERE id = @last_id;

-- test #2, expected: MSG_SUBROUTING_MULTIPLIERS_MISMATCH (clearing/clearing.go:284)
UPDATE [[[.DB_SCHEMA]]].clearing_batches
SET clearing_batch_type_id = 4  -- multiplier=0
WHERE id = @last_id;

-- }}}
-- batch status {{{

INSERT INTO [[[.DB_SCHEMA]]].clearing_batches
  (project_id, user_id, clearing_batch_type_id, clearing_batch_status_id, label, description)
VALUES
  ([[[.DB_TEMPLATE_PROJECT]]], [[[.DB_TEMPLATE_USER1]]], 4, 2, 'CLR #2', 'multipier=0, ledger_status_id=NULL');
SET @last_id = LAST_INSERT_ID();

-- test #3, expected: MSG_SUBROU (clearing/clearing.go:284)

UPDATE [[[.DB_SCHEMA]]].clearing_batches
SET clearing_batch_type_id = 2  -- multiplier=1
WHERE id = @last_id;


-- }}}

-- --------------------------------------------------------------------------- }}}

