-- Crear una nueva operación de pago (sender_name opcional)
-- name: CreatePaymentOperation :one
INSERT    INTO payment.operations (
          destination,
          operation_number,
          operation_date,
          sender_name,
          amount_sent,
          currency,
          additional_notes
          )
VALUES    ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- Obtener operación por ID
-- name: GetPaymentOperationByID :one
SELECT    *
FROM      payment.operations
WHERE     id = $1;

-- Obtener operación por número de operación
-- name: GetPaymentOperationByNumber :one
SELECT    *
FROM      payment.operations
WHERE     operation_number = $1;

-- Listar todas las operaciones ordenadas por fecha de creación
-- name: ListPaymentOperations :many
SELECT    *
FROM      payment.operations
ORDER BY  created_at;

-- Actualizar operación (campos opcionales, sender_name puede mantenerse o establecerse a NULL)
-- name: UpdatePaymentOperation :one
UPDATE    payment.operations
SET       destination = COALESCE($2, destination),
          operation_number = COALESCE($3, operation_number),
          operation_date = COALESCE($4, operation_date),
          sender_name = COALESCE($5, sender_name),
          amount_sent = COALESCE($6, amount_sent),
          currency = COALESCE($7, currency),
          additional_notes = COALESCE($8, additional_notes)
WHERE     id = $1
RETURNING *;

-- Eliminar operación
-- name: DeletePaymentOperation :one
DELETE    FROM payment.operations
WHERE     id = $1
RETURNING *;