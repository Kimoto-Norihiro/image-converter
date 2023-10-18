CREATE TABLE images (
	`id` int NOT NULL AUTO_INCREMENT,
	`object_name` varchar(255) NOT NULL UNIQUE,
	`resize_width_percent` int NOT NULL,
	`resize_height_percent` int NOT NULL,
	`encode_format_id` int NOT NULL,
	`status_id` int NOT NULL,
	`converted_image_url` varchar(255) NOT NULL,
	PRIMARY KEY (`id`)
);