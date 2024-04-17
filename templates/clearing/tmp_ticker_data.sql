-- Project and users --------------------------------------------------------- {{{
-- Clearing template user #1: [[[.DB_TEMPLATE_USER1]]]
-- Clearing template user #2: [[[.DB_TEMPLATE_USER2]]]

-- Clearing template project: [[[.DB_TEMPLATE_PROJECT]]]

-- --------------------------------------------------------------------------- }}}
-- Indexes ------------------------------------------------------------------- {{{

ALTER TABLE [[[.DB_SCHEMA]]].clearing_tasks ADD INDEX `clearing_tasks_cleared_at_idx` (`cleared_at`);

-- --------------------------------------------------------------------------- }}}

