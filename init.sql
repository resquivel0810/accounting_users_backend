-- Script SQL inicial para la base de datos local
-- Base de datos: preview
-- Usuario: preview_usr
-- Contraseña: localdev

-- Crear base de datos si no existe
CREATE DATABASE IF NOT EXISTS preview CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE preview;

-- Tabla de usuarios
CREATE TABLE IF NOT EXISTS users (
    id CHAR(36) PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    pwd VARCHAR(255) NOT NULL,
    profile_picture_url VARCHAR(500) DEFAULT '',
    status INT DEFAULT 1 COMMENT '0=inactivo, 1=activo',
    token VARCHAR(500) DEFAULT '',
    role INT DEFAULT 0 COMMENT 'Rol del usuario',
    account INT DEFAULT 1 COMMENT 'Tipo de cuenta: 1=free, 2=premium',
    email_conf INT DEFAULT 0 COMMENT '0=no confirmado, 1=confirmado',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    active INT DEFAULT 0 COMMENT '0=inactivo, 1=activo (sesión)',
    INDEX idx_email (email),
    INDEX idx_username (username),
    INDEX idx_status (status),
    INDEX idx_token (token(255))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla de códigos de usuario (confirmación email, reset password, etc.)
CREATE TABLE IF NOT EXISTS user_code (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_user CHAR(36) NOT NULL,
    code VARCHAR(255) NOT NULL,
    status INT DEFAULT 1 COMMENT '0=usado/inactivo, 1=activo',
    expiration TIMESTAMP NULL DEFAULT NULL COMMENT 'Fecha de expiración del código',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_user) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_code (code),
    INDEX idx_id_user (id_user),
    INDEX idx_status_expiration (status, expiration)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla de códigos de libros
CREATE TABLE IF NOT EXISTS bookcode (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(255) NOT NULL COMMENT 'MD5 del código',
    status INT DEFAULT 1 COMMENT '0=usado, 1=disponible',
    id_user CHAR(36) NULL COMMENT 'Usuario que usó el código',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_code (code),
    INDEX idx_status (status),
    INDEX idx_id_user (id_user)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla de feedback/comentarios
CREATE TABLE IF NOT EXISTS feedback (
    id INT AUTO_INCREMENT PRIMARY KEY,
    iduser CHAR(36) NOT NULL,
    rate INT NOT NULL COMMENT 'Calificación 1-5',
    comment TEXT NOT NULL,
    consent INT DEFAULT 0 COMMENT '0=sin consentimiento, 1=con consentimiento',
    display INT DEFAULT 0 COMMENT '0=no mostrar, 1=mostrar públicamente',
    watched INT DEFAULT 0 COMMENT '0=no visto, 1=visto por admin',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (iduser) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_iduser (iduser),
    INDEX idx_display (display),
    INDEX idx_consent (consent),
    INDEX idx_watched (watched)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla de razones de eliminación de cuenta
CREATE TABLE IF NOT EXISTS DelAccount (
    id INT AUTO_INCREMENT PRIMARY KEY,
    userid CHAR(36) NOT NULL,
    reason TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_userid (userid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla de métricas de términos buscados
CREATE TABLE IF NOT EXISTS metricsterm (
    id INT AUTO_INCREMENT PRIMARY KEY,
    idterm INT NOT NULL COMMENT 'ID del término buscado',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_idterm (idterm)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla de términos sugeridos/perdidos
CREATE TABLE IF NOT EXISTS sugestedterms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    term VARCHAR(500) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_term (term(255))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla de tiempo usado (métricas)
CREATE TABLE IF NOT EXISTS time_used (
    id INT AUTO_INCREMENT PRIMARY KEY,
    time DECIMAL(10,2) NOT NULL COMMENT 'Tiempo en segundos',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_time (time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Usuario de ejemplo para pruebas (opcional)
-- Contraseña: "password123" (MD5: 482c811da5d5b4bc6d497ffa98491e38)
-- INSERT INTO users (username, name, last_name, email, pwd, status, email_conf, role, account) 
-- VALUES ('testuser', 'Test', 'User', 'test@example.com', MD5('password123'), 1, 1, 0, 1);
