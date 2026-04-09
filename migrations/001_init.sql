-- Migration: Create products table
-- Run automatically when MySQL container initializes

CREATE DATABASE IF NOT EXISTS products_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE products_db;

CREATE TABLE IF NOT EXISTS products (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(255)       NOT NULL,
    description TEXT,
    price       DECIMAL(15, 2)     NOT NULL DEFAULT 0.00,
    stock       INT                NOT NULL DEFAULT 0,
    created_at  DATETIME           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME           NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  DATETIME           NULL,

    INDEX idx_deleted_at (deleted_at),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Seed data for RnD testing
INSERT INTO products (name, description, price, stock) VALUES
    ('Laptop Dell XPS 13',       'High-performance ultrabook 13 inch',         18500000.00, 25),
    ('Mechanical Keyboard',      'TKL layout, Cherry MX Red switches',          850000.00, 50),
    ('Monitor LG 27" 4K',        'IPS panel, 144Hz, HDR400',                   7200000.00, 15),
    ('USB-C Hub 7-in-1',         'HDMI 4K, PD 100W, USB 3.0 x3, SD card',      450000.00, 100),
    ('Webcam Logitech C920',     'Full HD 1080p, autofocus, stereo mic',        1200000.00, 30);
