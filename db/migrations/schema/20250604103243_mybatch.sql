-- +goose Up
CREATE TABLE dbo.batches (
    id UNIQUEIDENTIFIER NOT NULL DEFAULT NEWID(),
    app VARCHAR(255) NOT NULL,
    op VARCHAR(255) NOT NULL CHECK (op = LOWER(op)),
    context NVARCHAR(MAX) NOT NULL, -- JSONB equivalent
    inputfile VARCHAR(255),
    status VARCHAR(10) NOT NULL CHECK (status IN ('queued', 'inprog', 'success', 'failed', 'aborted', 'wait')),
    reqat DATETIME2 NOT NULL,
    doneat DATETIME2,
    created_at DATETIME2 DEFAULT GETDATE(),
    outputfiles NVARCHAR(MAX), -- JSONB equivalent
    nsuccess INT,
    nfailed INT,
    naborted INT
);


CREATE TABLE dbo.batchrows (
    rowid BIGINT IDENTITY(1,1) PRIMARY KEY,
    batch UNIQUEIDENTIFIER NOT NULL,
    line INT NOT NULL,
    input NVARCHAR(MAX) NOT NULL,     -- JSONB equivalent
    status VARCHAR(10) NOT NULL CHECK (status IN ('queued', 'inprog', 'success', 'failed', 'aborted', 'wait')),
    reqat DATETIME2 NOT NULL,
    doneat DATETIME2,
    res NVARCHAR(MAX),                -- JSONB equivalent
    blobrows NVARCHAR(MAX),           -- JSONB equivalent
    messages NVARCHAR(MAX),           -- JSONB equivalent
    doneby VARCHAR(255)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS dbo.batchrows;
DROP TABLE IF EXISTS dbo.batches;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
