CREATE TABLE images (
  "id" int NOT NULL AUTO_INCREMENT PRIMARY KEY,
  "object_name" varchar NOT NULL UNIQUE,
	"resize_width_percent" int NOT NULL,
  "resize_height_percent" int NOT NULL,
	"encode_format_id" int NOT NULL,
	"status_id" int NOT NULL,
	"converted_image_url" varchar NOT NULL,
)