-- +goose Up
-- CREATE    TYPE validation_status AS ENUM('pending', 'approved', 'rejected');
-- Payment Operations Table (Outgoing Payments)
CREATE    TABLE payment.operations (
          id SERIAL PRIMARY KEY,
          destination TEXT NOT NULL,
          operation_number TEXT NOT NULL UNIQUE,
          operation_date TIMESTAMP NOT NULL,
          sender_name TEXT NULL,
          amount_sent DECIMAL(15, 2) NOT NULL,
          currency CHAR(3) NOT NULL,
          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
          additional_notes TEXT
          );

-- Payment Receipts Table (Incoming Payments)
-- CREATE TABLE payment_receipts (
--     receipt_id SERIAL PRIMARY KEY,
--     payment_operation_id INT REFERENCES payment_operations(operation_id),
--     receiver_name VARCHAR(255) NOT NULL,
--     receipt_number VARCHAR(50) NOT NULL UNIQUE,
--     receipt_date TIMESTAMP NOT NULL,
--     amount_received DECIMAL(15,2) NOT NULL,
--     received_currency CHAR(3) NOT NULL,
--     receipt_status VARCHAR(20) DEFAULT 'pending',
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     receipt_proof_url TEXT,
--     validation_notes TEXT
-- );
-- -- Payment Validations Table
-- CREATE TABLE payment_validations (
--     validation_id SERIAL PRIMARY KEY,
--     operation_id INT NOT NULL REFERENCES payment_operations(operation_id),
--     receipt_id INT NOT NULL REFERENCES payment_receipts(receipt_id),
--     validator_id INT, -- Could reference users table if exists
--     validation_date TIMESTAMP NOT NULL,
--     validation_status validation_status NOT NULL DEFAULT 'pending',
--     validation_reason TEXT,
--     validated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     UNIQUE(operation_id, receipt_id) -- Ensure one validation per pair
-- );