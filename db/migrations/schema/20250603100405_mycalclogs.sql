-- +goose Up
CREATE TABLE my_calculation_logs (
    id INT IDENTITY(1,1) PRIMARY KEY,
    request_time   DATETIME2,
    response_time  DATETIME2,
    duration_ms    FLOAT,
    input          NVARCHAR(2000),
    operation      NVARCHAR(20),
    request_data   NVARCHAR(MAX),
    response_data  NVARCHAR(MAX),
    error          NVARCHAR(MAX)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE my_calculation_logs;

-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
