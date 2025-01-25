-- Create "entries" table
CREATE TABLE `entries` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `title` text NOT NULL, `description` text NOT NULL, `url` text NOT NULL, `feeds_entries` integer NULL, CONSTRAINT `entries_feeds_entries` FOREIGN KEY (`feeds_entries`) REFERENCES `feeds` (`id`) ON DELETE SET NULL);
-- Create index "entries_url_feeds_entries" to table: "entries"
CREATE UNIQUE INDEX `entries_url_feeds_entries` ON `entries` (`url`, `feeds_entries`);
