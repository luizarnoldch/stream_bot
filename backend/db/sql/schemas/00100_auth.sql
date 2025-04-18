-- +goose Up
CREATE    TABLE auth."user" (
          id SERIAL PRIMARY KEY,
          NAME TEXT NOT NULL,
          email TEXT UNIQUE NOT NULL,
          created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
          updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
          );