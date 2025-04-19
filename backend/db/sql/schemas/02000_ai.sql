-- +goose Up
-- Crear tabla ai.conversation
CREATE    TABLE IF NOT EXISTS ai.conversation (
          id SERIAL PRIMARY KEY,
          UUID UUID NOT NULL DEFAULT gen_random_uuid (),
          NAME TEXT NOT NULL,
          user_id INTEGER NOT NULL REFERENCES auth."user" (id) ON DELETE CASCADE,
          created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
          updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
          );

CREATE    INDEX IF NOT EXISTS idx_ai_conversation_user_id ON ai.conversation (user_id);

-- Crear tabla ai.message
CREATE    TABLE IF NOT EXISTS ai.message (
          id SERIAL PRIMARY KEY,
          CONTENT TEXT NOT NULL,
          ROLE TEXT NOT NULL,
          conversation_id INTEGER NOT NULL REFERENCES ai.conversation (id) ON DELETE CASCADE,
          tokens INTEGER,
          finish_reason TEXT,
          created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
          updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
          );

CREATE    INDEX IF NOT EXISTS idx_ai_message_conversation_id ON ai.message (conversation_id);