--
-- Template project and user
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
