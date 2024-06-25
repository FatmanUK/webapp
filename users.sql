PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS `users` (`name` text,`nick` text,`groups` text);
INSERT INTO users VALUES('fatman@dreamtrack.net','Fatman','meatbags;authors;wizards');
INSERT INTO users VALUES('adam.richardson@cgi.com','Adam','wizards;meatbags;authors');
-- DELETE FROM sqlite_sequence;
-- INSERT INTO sqlite_sequence VALUES('users',2);
CREATE INDEX IF NOT EXISTS `idx_users_name` ON `users`(`name`);
COMMIT;
