PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS `users` (`name` text,`nick` text,`groups` text,`created` datetime,`last_login` datetime,`last_request` datetime);
INSERT INTO users VALUES('fatman@dreamtrack.net','Fatman','meatbags;authors;wizards;stewards',NULL,NULL,NULL);
INSERT INTO users VALUES('adam.richardson@cgi.com','Adam','wizards;meatbags;authors',NULL,NULL,NULL);
CREATE INDEX IF NOT EXISTS `idx_users_name` ON `users`(`name`);
COMMIT;
