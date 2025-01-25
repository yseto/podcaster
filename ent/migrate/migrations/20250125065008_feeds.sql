-- Create "feeds" table
CREATE TABLE `feeds` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `title` text NOT NULL, `url` text NOT NULL, `users_feeds` integer NULL, CONSTRAINT `feeds_users_feeds` FOREIGN KEY (`users_feeds`) REFERENCES `users` (`id`) ON DELETE SET NULL);
