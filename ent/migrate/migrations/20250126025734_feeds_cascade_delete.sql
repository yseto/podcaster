-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "new_entries" table
CREATE TABLE `new_entries` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `title` text NOT NULL, `description` text NOT NULL, `url` text NOT NULL, `published_at` datetime NULL, `new` bool NOT NULL DEFAULT (true), `feeds_entries` integer NULL, CONSTRAINT `entries_feeds_entries` FOREIGN KEY (`feeds_entries`) REFERENCES `feeds` (`id`) ON DELETE CASCADE);
-- Copy rows from old table "entries" to new temporary table "new_entries"
INSERT INTO `new_entries` (`id`, `title`, `description`, `url`, `published_at`, `new`, `feeds_entries`) SELECT `id`, `title`, `description`, `url`, `published_at`, `new`, `feeds_entries` FROM `entries`;
-- Drop "entries" table after copying rows
DROP TABLE `entries`;
-- Rename temporary table "new_entries" to "entries"
ALTER TABLE `new_entries` RENAME TO `entries`;
-- Create index "entries_url_feeds_entries" to table: "entries"
CREATE UNIQUE INDEX `entries_url_feeds_entries` ON `entries` (`url`, `feeds_entries`);
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
