--
-- Clearing token types
--

DROP TABLE IF EXISTS `[[[.DB_SCHEMA]]]`.`clearing_token_types`;

CREATE TABLE `[[[.DB_SCHEMA]]]`.`clearing_token_types` (
	`id` MEDIUMINT UNSIGNED NOT NULL AUTO_INCREMENT,
	`label` TINYTEXT NOT NULL,
	`description` TINYTEXT DEFAULT NULL,
	`created_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(),
	`updated_at`  DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
	`deleted_at` DATETIME(3) NULL DEFAULT NULL,
	PRIMARY KEY (`id`),
	KEY `clearing_token_types_deleted_at_idx` (`deleted_at`) USING BTREE
) [[[.DB_TABLE_OPTIONS]]];

INSERT INTO `[[[.DB_SCHEMA]]]`.`clearing_token_types`
	(label, description)
VALUES
	('Fiat', 'Government-issued, legal tender with no intrinsic value; value relies on public trust.'),
	('Crypto', 'Decentralized digital assets; blockchain-based peer-to-peer currencies.'),
	('TrustChain', 'Fungible token on the TrustChain network.'),
	('ERC-20', 'Ethereum fungible token'),
	('SPL', 'Fungible tokens on the Solana network.'),
	('Virtual', 'Project token with no actual implementation yet.');

--
-- Clearing tokens
--

DROP TABLE IF EXISTS `[[[.DB_SCHEMA]]]`.`clearing_tokens`;

CREATE TABLE `[[[.DB_SCHEMA]]]`.`clearing_tokens` (
	`id` MEDIUMINT UNSIGNED NOT NULL AUTO_INCREMENT,
	`label` TINYTEXT NOT NULL,
	`symbol` VARCHAR(16) NOT NULL,
	`description` TINYTEXT DEFAULT NULL,
	`clearing_token_type_id` MEDIUMINT UNSIGNED NOT NULL,
	`created_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(),
	`updated_at`  DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
	`deleted_at` DATETIME(3) NULL DEFAULT NULL,
	PRIMARY KEY (`id`),
	KEY `clearing_tokens_deleted_at_idx` (`deleted_at`) USING BTREE,
	UNIQUE KEY `clearing_tokens_symbol_clearing_token_type_id_idx` (symbol, clearing_token_type_id),
	CONSTRAINT `clearing_tokens_clearing_token_type_id_fk` FOREIGN KEY (`clearing_token_type_id`) REFERENCES `clearing_token_types` (`id`) ON UPDATE CASCADE
) [[[.DB_TABLE_OPTIONS]]];

INSERT INTO `[[[.DB_SCHEMA]]]`.`clearing_tokens`
	(label, symbol, description, clearing_token_type_id)
VALUES
	('US dollar', 'USD', 'United States Dollar: Official U.S. fiat currency', 1),
	('Euro', 'EUR', 'Official fiat currency in the Eurozone', 1),
	('Ether', 'ETH', 'Native cryptocurrency of the Ethereum network.', 2),
	('Bitcoin', 'BTC', 'The digital gold.', 2),
	('FairGrind', 'FGDT', 'FairGrind token to be minted on the TrustChain mainnet.', 3),
	('TE-FOOD/TONE', 'TONE', 'ERC-20 utility token of TE-FOOD.', 4),
	('Template', 'FOOBAR', 'Virtual project token with no blockchain implementation yet.', 6);


