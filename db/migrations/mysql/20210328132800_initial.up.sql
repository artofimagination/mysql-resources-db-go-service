-- +migrate Up
CREATE TABLE IF NOT EXISTS categories(
   id smallint PRIMARY KEY AUTO_INCREMENT,
   name VARCHAR (50),
   description VARCHAR (300),
   created_at DATETIME NOT NULL DEFAULT NOW(),
   updated_at DATETIME NOT NULL DEFAULT NOW()
);

-- +migrate Up
CREATE TABLE IF NOT EXISTS resources(
   id binary(16) PRIMARY KEY,
   category smallint,
   FOREIGN KEY (category) REFERENCES categories(id),
   content json,
   created_at DATETIME NOT NULL DEFAULT NOW(),
   updated_at DATETIME NOT NULL DEFAULT NOW()
);

-- +migrate Up
INSERT INTO categories (name, description) VALUES ('News feed', 'Resource marked as news feed item');
INSERT INTO categories (name, description) VALUES ('Content', 'All resource that has been uploaded as an attachement in another resource. For example, news feed image for news feed resource item');
