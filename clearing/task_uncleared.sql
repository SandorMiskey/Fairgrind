-- User -------------------------------------------------------------------- {{{

INSERT INTO fairgrind.users
  (first_name, last_name, email, password, birth_year, country_id)
VALUES
  ('Task uncleared #1', 'Clearing test', '', '', 0, 0);
SET @user_id = LAST_INSERT_ID();

-- ------------------------------------------------------------------------- }}}
-- Project ----------------------------------------------------------------- {{{

INSERT INTO fairgrind.projects
  (uid, name, description)
VALUES
  (@user, 'Clearing test', 'Task uncleared #1');
SET @project_id = LAST_INSERT_ID();

-- ------------------------------------------------------------------------- }}}
-- Batches ----------------------------------------------------------------- {{{

-- multiplier=0, withdrawable=NULL
INSERT INTO fairgrind.clearing_batches
  (project_id, user_id, clearing_batch_type_id, clearing_batch_status_id, label, description)
VALUES
  (@project_id, @user_id, 1, 2, 'Clearing test', 'Task uncleard #1, multipier=0, withdrawable=NULL');
SET @batch_id_0n = LAST_INSERT_ID();

-- multiplier=0, withdrawable=false
INSERT INTO fairgrind.clearing_batches
  (project_id, user_id, clearing_batch_type_id, clearing_batch_status_id, label, description)
VALUES
  (@project_id, @user_id, 1, 5, 'Clearing test', 'Task uncleard #2, multipier=0, withdrawable=false');
SET @batch_id_0f = LAST_INSERT_ID();

-- multiplier=0, withdrawable=true
INSERT INTO fairgrind.clearing_batches
  (project_id, user_id, clearing_batch_type_id, clearing_batch_status_id, label, description)
VALUES
  (@project_id, @user_id, 1, 6, 'Clearing test', 'Task uncleard #3, multipier=0, withdrawable=true');
SET @batch_id_0t = LAST_INSERT_ID();

-- multiplier=1, withdrawable=NULL
INSERT INTO fairgrind.clearing_batches
  (project_id, user_id, clearing_batch_type_id, clearing_batch_status_id, label, description)
VALUES
  (@project_id, @user_id, 3, 2, 'Clearing test', 'Task uncleard #4, multipier=1, withdrawable=NULL');
SET @batch_id_1n = LAST_INSERT_ID();

-- multiplier=1, withdrawable=false
INSERT INTO fairgrind.clearing_batches
  (project_id, user_id, clearing_batch_type_id, clearing_batch_status_id, label, description)
VALUES
  (@project_id, @user_id, 3, 5, 'Clearing test', 'Task uncleard #5, multipier=1, withdrawable=false');
SET @batch_id_1f = LAST_INSERT_ID();

-- multiplier=1, withdrawable=true
INSERT INTO fairgrind.clearing_batches
  (project_id, user_id, clearing_batch_type_id, clearing_batch_status_id, label, description)
VALUES
  (@project_id, @user_id, 3, 6, 'Clearing test', 'Task uncleard #6, multipier=1, withdrawable=true');
SET @batch_id_1t = LAST_INSERT_ID();

-- ------------------------------------------------------------------------- }}}
-- Fees -------------------------------------------------------------------- {{{

INSERT INTO fairgrind.clearing_task_fees
  (user_id, project_id, clearing_task_type_id, task_fee, subtask_fee, clearing_token_id)
VALUES
  (@user_id, @project_id, 1, 1.0, 2.0, 2),  -- annotation, eur
  (@user_id, @project_id, 1, 0.1, 0.2, 3),  -- annotation, eth
  (@user_id, @project_id, 2, 3.0, 4.0, 2),  -- data collection, eur
  (@user_id, @project_id, 2, 0.3, 0.4, 3),  -- data collection, eth
  (@user_id, @project_id, 3, 5.0, 6.0, 2),  -- review, eur
  (@user_id, @project_id, 3, 0.5, 0.0, 3);  -- review, eth

-- ------------------------------------------------------------------------- }}}
-- Tasks ------------------------------------------------------------------- {{{

INSERT INTO fairgrind.tasks
  (grinder_uid)
VALUES
  (@user_id);
SET @task_id = LAST_INSERT_ID();

-- batch_id_0n, multiplier=0, withdrawable=NULL {{{

-- payable=true, parent_payable=true
INSERT INTO `fairgrind`.`clearing_tasks`
  (clearing_batch_id, clearing_task_type_id, clearing_task_status_id, clearing_task_id, task_id, user_id)
VALUES
  (@batch_id_0n, 1, 3, NULL, @task_id, @user_id),   -- annotation
  (@batch_id_0n, 2, 3, NULL, @task_id, @user_id),   -- data collection
  (@batch_id_0n, 3, 3, NULL, @task_id, @user_id);   -- review

-- payable=false, parent_payable=false
INSERT INTO `fairgrind`.`clearing_tasks`
  (clearing_batch_id, clearing_task_type_id, clearing_task_status_id, clearing_task_id, task_id, user_id)
VALUES
  (@batch_id_0n, 1, 5, NULL, @task_id, @user_id),   -- annotation
  (@batch_id_0n, 2, 5, NULL, @task_id, @user_id),   -- data collection
  (@batch_id_0n, 3, 5, NULL, @task_id, @user_id);   -- review

-- }}}
-- batch_id_0f, multiplier=0, withdrawable=false {{{

-- payable=true, parent_payable=true
INSERT INTO `fairgrind`.`clearing_tasks`
  (clearing_batch_id, clearing_task_type_id, clearing_task_status_id, clearing_task_id, task_id, user_id)
VALUES
  (@batch_id_0f, 1, 3, NULL, @task_id, @user_id),   -- annotation
  (@batch_id_0f, 2, 3, NULL, @task_id, @user_id),   -- data collection
  (@batch_id_0f, 3, 3, NULL, @task_id, @user_id);   -- review

-- payable=false, parent_payable=false
INSERT INTO `fairgrind`.`clearing_tasks`
  (clearing_batch_id, clearing_task_type_id, clearing_task_status_id, clearing_task_id, task_id, user_id)
VALUES
  (@batch_id_0f, 1, 5, NULL, @task_id, @user_id),   -- annotation
  (@batch_id_0f, 2, 5, NULL, @task_id, @user_id),   -- data collection
  (@batch_id_0f, 3, 5, NULL, @task_id, @user_id);   -- review

-- }}}
-- batch_id_0t, multiplier=0, withdrawable=true {{{

-- payable=true, parent_payable=true
INSERT INTO `fairgrind`.`clearing_tasks`
  (clearing_batch_id, clearing_task_type_id, clearing_task_status_id, clearing_task_id, task_id, user_id)
VALUES
  (@batch_id_0t, 1, 3, NULL, @task_id, @user_id),   -- annotation
  (@batch_id_0t, 2, 3, NULL, @task_id, @user_id),   -- data collection
  (@batch_id_0t, 3, 3, NULL, @task_id, @user_id);   -- review

-- payable=false, parent_payable=false
INSERT INTO `fairgrind`.`clearing_tasks`
  (clearing_batch_id, clearing_task_type_id, clearing_task_status_id, clearing_task_id, task_id, user_id)
VALUES
  (@batch_id_0t, 1, 5, NULL, @task_id, @user_id),   -- annotation
  (@batch_id_0t, 2, 5, NULL, @task_id, @user_id),   -- data collection
  (@batch_id_0t, 3, 5, NULL, @task_id, @user_id);   -- review

-- }}}
-- batch_id_1n, multiplier=1, withdrawable=NULL {{{

-- payable=true, parent_payable=true
INSERT INTO `fairgrind`.`clearing_tasks`
  (clearing_batch_id, clearing_task_type_id, clearing_task_status_id, clearing_task_id, task_id, user_id)
VALUES
  (@batch_id_1n, 1, 3, NULL, @task_id, @user_id),   -- annotation
  (@batch_id_1n, 2, 3, NULL, @task_id, @user_id),   -- data collection
  (@batch_id_1n, 3, 3, NULL, @task_id, @user_id);   -- review

-- payable=false, parent_payable=false
INSERT INTO `fairgrind`.`clearing_tasks`
  (clearing_batch_id, clearing_task_type_id, clearing_task_status_id, clearing_task_id, task_id, user_id)
VALUES
  (@batch_id_1n, 1, 5, NULL, @task_id, @user_id),   -- annotation
  (@batch_id_1n, 2, 5, NULL, @task_id, @user_id),   -- data collection
  (@batch_id_1n, 3, 5, NULL, @task_id, @user_id);   -- review

-- }}}
-- batch_id_1f, multiplier=1, withdrawable=false {{{

-- payable=true, parent_payable=true
INSERT INTO `fairgrind`.`clearing_tasks`
  (clearing_batch_id, clearing_task_type_id, clearing_task_status_id, clearing_task_id, task_id, user_id)
VALUES
  (@batch_id_1f, 1, 3, NULL, @task_id, @user_id),   -- annotation
  (@batch_id_1f, 2, 3, NULL, @task_id, @user_id),   -- data collection
  (@batch_id_1f, 3, 3, NULL, @task_id, @user_id);   -- review

-- payable=false, parent_payable=false
INSERT INTO `fairgrind`.`clearing_tasks`
  (clearing_batch_id, clearing_task_type_id, clearing_task_status_id, clearing_task_id, task_id, user_id)
VALUES
  (@batch_id_1f, 1, 5, NULL, @task_id, @user_id),   -- annotation
  (@batch_id_1f, 2, 5, NULL, @task_id, @user_id),   -- data collection
  (@batch_id_1f, 3, 5, NULL, @task_id, @user_id);   -- review

-- }}}
-- batch_id_1t, multiplier=1, withdrawable=true {{{

-- payable=true, parent_payable=true
INSERT INTO `fairgrind`.`clearing_tasks`
  (clearing_batch_id, clearing_task_type_id, clearing_task_status_id, clearing_task_id, task_id, user_id)
VALUES
  (@batch_id_1t, 1, 3, NULL, @task_id, @user_id),   -- annotation
  (@batch_id_1t, 2, 3, NULL, @task_id, @user_id),   -- data collection
  (@batch_id_1t, 3, 3, NULL, @task_id, @user_id);   -- review

-- payable=false, parent_payable=false
INSERT INTO `fairgrind`.`clearing_tasks`
  (clearing_batch_id, clearing_task_type_id, clearing_task_status_id, clearing_task_id, task_id, user_id)
VALUES
  (@batch_id_1t, 1, 5, NULL, @task_id, @user_id),   -- annotation
  (@batch_id_1t, 2, 5, NULL, @task_id, @user_id),   -- data collection
  (@batch_id_1t, 3, 5, NULL, @task_id, @user_id);   -- review

-- }}}

-- ------------------------------------------------------------------------- }}}

