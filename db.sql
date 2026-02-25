-- =============================================
-- STORE JV - Schéma MySQL
-- =============================================

CREATE DATABASE IF NOT EXISTS partiel CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE partiel;

-- ---------------------------------------------
-- Table : games
-- Contient le catalogue de jeux
-- ---------------------------------------------
CREATE TABLE games (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    title       VARCHAR(255)   NOT NULL,
    description TEXT,
    price       DECIMAL(10, 2) NOT NULL,
    image_url   VARCHAR(500),
    stock       INT            NOT NULL DEFAULT 0,
    created_at  TIMESTAMP      DEFAULT CURRENT_TIMESTAMP
);

-- ---------------------------------------------
-- Table : cart
-- Un panier est lié à un session_id (UUID généré côté front)
-- Pas d'utilisateur = pas de FK vers une table users
-- ---------------------------------------------
CREATE TABLE cart (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    session_id  VARCHAR(36)  NOT NULL UNIQUE,   -- UUID format : xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
    created_at  TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- ---------------------------------------------
-- Table : cart_items
-- Ligne de panier : un jeu + une quantité pour un cart donné
-- ON DELETE CASCADE : si le cart est supprimé, ses items le sont aussi
-- UNIQUE (cart_id, game_id) : évite les doublons, on incrémente quantity plutôt
-- ---------------------------------------------
CREATE TABLE cart_items (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    cart_id     INT NOT NULL,
    game_id     INT NOT NULL,
    quantity    INT NOT NULL DEFAULT 1,
    FOREIGN KEY (cart_id) REFERENCES cart(id) ON DELETE CASCADE,
    FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
    UNIQUE KEY uq_cart_game (cart_id, game_id)
);

-- ---------------------------------------------
-- Données de test
-- ---------------------------------------------
INSERT INTO games (title, description, price, image_url, stock) VALUES
('The Witcher 3', 'RPG open world épique dans un monde dark fantasy.', 19.99, 'https://placehold.co/300x400?text=Witcher+3', 50),
('Hollow Knight', 'Metroidvania atmosphérique et difficile.', 14.99, 'https://placehold.co/300x400?text=Hollow+Knight', 30),
('Celeste', 'Platformer narratif sur le dépassement de soi.', 19.99, 'https://placehold.co/300x400?text=Celeste', 25),
('Disco Elysium', 'RPG narratif policier unique en son genre.', 29.99, 'https://placehold.co/300x400?text=Disco+Elysium', 15),
('Hades', 'Roguelite action avec une narration excellente.', 24.99, 'https://placehold.co/300x400?text=Hades', 40),
('Stardew Valley', 'Simulation de ferme relaxante et addictive.', 13.99, 'https://placehold.co/300x400?text=Stardew', 60);