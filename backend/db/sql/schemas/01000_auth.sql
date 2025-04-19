-- +goose Up
-- Crear tabla auth.user
CREATE    TABLE IF NOT EXISTS auth."user" (
          id SERIAL PRIMARY KEY,
          NAME TEXT NOT NULL,
          phone TEXT NOT NULL UNIQUE,
          created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
          updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
          );

CREATE    INDEX IF NOT EXISTS idx_user_phone ON auth.user (phone);