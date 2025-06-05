-- +goose Up
CREATE TABLE employee (
    id INT IDENTITY(1,1) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    empcode VARCHAR(10) UNIQUE NOT NULL,
    -- gender_enum emulated via CHECK constraint
    gender CHAR(1) NOT NULL CHECK (gender IN ('M', 'F', 'O')),
    dob DATE,
    doj DATE,
    salary BIGINT,
    reportsto VARCHAR(10), -- FK below

    -- designation_enum emulated via CHECK constraint
    designation VARCHAR(20) NOT NULL CHECK (designation IN ('Admin', 'Manager', 'TeamLead', 'Developer', 'Tester', 'HR')),
    briefbio TEXT,
    interests VARCHAR(255),
    addedat DATETIMEOFFSET,
    addedby VARCHAR(50),
    updatedat DATETIMEOFFSET,
    updatedby VARCHAR(50),
    isactive BIT NOT NULL DEFAULT 1,
    added_fromip VARCHAR(45),
    updated_fromip VARCHAR(45),
    CONSTRAINT fk_reportsto FOREIGN KEY (reportsto) REFERENCES employee(empcode)
);

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
