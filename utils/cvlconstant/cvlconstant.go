package cvlconstant

// special chars
const (
	FORWARD_SLASH = "/"
	UNDERSCORE    = "_"
)

// HTTP Methods
const (
	POST_METHOD = "POST"
	GET_METHOD  = "GET"
)

// module and app const
const (
	MODULE_MASTER = "master"
	APP_NAME      = "cvl-kra"
)

// logger priority key
const (
	LOGGER_PRIORITY_KEY = "logger_priority"
)

// middleware consts
const (
	IP = "remoteIp"
)

type ApiName string

func (api ApiName) String() string {
	return string(api)
}

// api label const
const (
	MyHello         ApiName = "myhello"
	MyCalc          ApiName = "mycalc"
	MyBulkCalc      ApiName = "mybulkcalc"
	MyBulkCalcBatch ApiName = "mybulkcalcbatch"
	AddEmp          ApiName = "addemp"
	UpdateEmp       ApiName = "updateemp"
	ListEmp         ApiName = "listemp"
	AddProof        ApiName = "addproof"
	GetPassword     ApiName = "/getpassword"
	PanValidation   ApiName = "/panvalidation"
	PanDownload     ApiName = "/pandownload"
	PanDataUpload   ApiName = "/pandataupload"
)

// logger module
const (
	MODULE_MYHELLO          = "MYHELLO"
	MODULE_MYCALC           = "MYCALC"
	MODULE_MYCALC_V1_2      = "MYCALC V1.2"
	MODULE_MYBULKCALC       = "MYBULKCALC"
	MODULE_MYBULKCALC_BATCH = "MYBULKCALC_BATCH"
	DOBATCHJOB              = "DO BATCH JOB"
	ADD_EMPLOYEE            = "ADD_EMPLOYEE"
	MODULE_ADD_EMPLOYEE     = "ADD_EMPLOYEE"
	MODULE_ADD_PROOF        = "ADD_PROOF"
	MODULE_LIST_EMPLOYEES   = "MODULE_LIST_EMPLOYEES"
	MODULE_UPDATE_EMPLOYEE  = "MODULE_UPDATE_EMPLOYEE"
	MODULE_GET_PASSWORD     = "MODULE_GET_PASSWORD"
	MODULE_PAN_VALIDATION   = "MODULE_PAN_VALIDATION"
	MODULE_PAN_DOWNLOAD     = "MODULE_PAN_DOWNLOAD"
	MODULE_PAN_DATA_UPLOAD  = "MODULE_PAN_DATA_UPLOAD"
)

// common constant
const (
	INVALID_JSON_FORMAT                 = "invalid JSON"
	CALCULATION_ERROR                   = "calculation error"
	CALCULATION_SUCCESS                 = "calculation successful"
	EMPLOYEE_SUCCESS                    = "employee added successfuly"
	DB_STORAGE_FAILED                   = "DB storage failed"
	DB_CHECK_FAILED                     = "DB check failed"
	DB_STORAGE_SUCCESS                  = "DB storage successful"
	UNSUPPORTED_OPERATION               = "unsupported operation"
	EMPTY_NUMBERS_LIST                  = "empty numbers list"
	ERROR_INSERTING_LOG_ENTRY           = "error inserting log entry into database"
	ERROR_IN_MARSHALING                 = "error marshaling response"
	INVALID_BODY                        = "invalid Body"
	INVALID_CSV                         = "invalid base64 CSV"
	INPUTFILE                           = "inputfile"
	INVALID_INPUT_DOBATCH               = "invalid input in DoBatchJob"
	FAILED_DECODE_CSV                   = "Failed to decode base64 input:"
	BULK_OPERATION                      = "bulk_operation"
	EMP_REGEX_EXP                       = `^[A-Za-z]-\d{1,5}$`
	NAME_REGEX_EXP                      = `^[A-Za-z]+(?:'[A-Za-z]+)?(?: [A-Za-z]+(?:'[A-Za-z]+)?)*$`
	PAN_REGEX_EXP                       = `^[A-Z]{5}[0-9]{4}[A-Z]{1}$`
	PAN_TYPE_INDIVIDUAL                 = "P" // Individual
	PAN_TYPE_COMPANY                    = "C" // Company
	PAN_TYPE_HUF                        = "H" // Hindu Undivided Family
	PAN_TYPE_AOP                        = "A" // Association of Persons
	PAN_TYPE_BOI                        = "B" // Body of Individuals
	PAN_TYPE_GOVT                       = "G" // Government Agency
	PAN_TYPE_JURIDICAL_PERSON           = "J" // Artificial Juridical Person
	PAN_TYPE_LOCAL_AUTHORITY            = "L" // Local Authority
	PAN_TYPE_FIRM                       = "F" // Firm
	PAN_TYPE_TRUST                      = "T" // Trust
	PAN_FOURTH_PLACE                    = "PHABGJLFTC"
	PAN_FOURTH_PLACE_MUST               = "4th place must be"
	DATE_FORMAT                         = "2006-01-02"
	DATE_FORMAT_EMP                     = "e.g., 2006-01-02"
	INVALID_DATE_FORMAT                 = "invalid_format"
	INVALID_AGE_RANGE                   = "invalid_range"
	SYSTEM                              = "system"
	INVALID_REPORTSTO                   = "invalid_reportsto"
	DUPLICATE_EMPCODE                   = "duplicate"
	SELF_REPORT_NOT_ALLOWED             = "self_report_not_allowed"
	INVALID_AGE_AT_JOINING              = "Invalid age at joining"
	INVALID_GENDER                      = "invalid gender"
	INVALID_DESIGNATION                 = "invalid designation"
	INVALID_EMP_CODE                    = "invalid_emp_code"
	VALIDATION_ERRORS_FOUND             = "validation_errors_found"
	ERRORCODE_FETCHING_DB_RECORD_FAILED = "emp_list fetch failed"
	STRUCT_VALIDATION_FAILED            = "Struct validation failed"
	VALIDATION_FAILED                   = "validation failed"
	INVALID_XML_FORMAT                  = "invalid XML"
	VALIDATION_PASS                     = "validation pass"
	REQUEST_DATA_IS_VALID               = "Request data is valid"
	RETURNING_VALIDATION_ERRORS         = "Returning validation errors"
	MARSHALING_ERROR                    = "Marshaling error"
	MOCKOON_PANVALID_ERROR              = "Mockoon pan validation error"
	MOCKOON_GETPASSWORD_ERROR           = "Mockoon getpassword error"
	ERROR_IN_SENDING_XML_REQUEST        = "Error in sending XML request"
	READING_RESPONSE_BODY_FAILED        = "Reading response failed"
	EXTERNAL_SERVICE_ERROR              = "Internal Service Error"
)

// message ids
const (
	INVALID_MESSAGE_ID           = 201
	EMPTY_MESSAGE_ID             = 202
	INVALID_BASE64               = 203
	DUPLICATE_MESSAGE_ID         = 204
	DB_STORAGE_FAILED_MESSAGE_ID = 205
	SELF_REPORT_NOT_ALLOWED_ID   = 206
	FETCHING_DB_RECORD_FAILED    = 207
	INVALID_JSON_MESSAGE_ID      = 208
)

// errcodes
const (
	ERRORCODE_EMPTY              = "empty"
	ERRORCODE_REQUIRED           = "required"
	ERRORCODE_INVALID            = "invalid"
	GENDER_REQUIRED              = "gender"
	ERRORCODE_NOT_ADULT          = "not adult"
	ERRORCODE_INVALID_DEPENDANCY = "invalid_dependency"
	ERRORCODE_BASE64             = "invalid base64 file"
	BATCH_SUBMISSION_FAILED      = "batch failed"
	CSV_READ_ERROR               = "csv read error"
	CSV_PARSE_ERROR              = "csv parse error"
	CSV_EMPTY_RECORD             = "csv empty record"
	JSON_MARSHAL_ERROR           = "JSON marshal error"
	GENERATED_JSON               = "generated JSON"
	INVALID_JSON_CONVEN          = "invalid JSON convension"
	INVALID_SALARY_RANGE         = "invalid_range"
	INVALID_SALARY_FORMAT        = "invalid_format"
	INVALID_NAME_FORMAT          = "invalid_name_format"
	DOB_VALIDATION               = "DOB validation"
	DOJ_VALIDATION               = "DOJ validation"
	SALARY_VALIDATION            = "salary validation"
	REPORSTO_VALIDATION          = "reporting validation"
	FILTER_VALIDATION_ERROR      = "FILTER_VALIDATION_ERROR"
	INVALID_JSON_MESSAGE_ERROR   = "invalid_json"
	REQUEST_BODY                 = "request_body"
)

const (
	INPUT           = "input"
	OPERATION       = "op"
	RESULT          = "result"
	INPUT_OPERATION = "input & operation"
	SUM             = "sum"
	MEAN            = "mean"
	SORT            = "sort"
	AVERAGE         = "average"
	VALID_EMPCODE   = "e.g., A-1, A-12"
	VALID_NAME      = "e.g., Allan D'costa"
	ISSUE_DATE      = "issuedat"
	VALID_FROM      = "validfrom"
	VALIDTO         = "validto"
	ASC             = "asc"
	DESC            = "desc"
	VALID_SORTBY    = "asc or desc"
	ADMIN           = "Admin"
	MANAGER         = "Manager"
	TEAMLEAD        = "TeamLead"
	DEVELOPER       = "Developer"
	TESTER          = "Tester"
	HR              = "HR"
)

const (
	UNDER_AGE   = "under_age"
	PAST_DATE   = "past_date"
	DOB         = "dob"
	EMPCODE     = "empcode"
	NAME        = "name"
	DOJ         = "doj"
	SALARY      = "salary"
	REPORTSTO   = "reportsto"
	DESIGNATION = "designation"
	GENDER      = "gender"
	MALE        = "M"
	FEMALE      = "F"
	OTHER       = "O"
	BRIEF_BIO   = "briefbio"
	IDENTIFIER  = "identifier"
	SORTDIR     = "sortDir"
	SORTBY      = "sortBy"
)

const (
	MIN_SALARY = 0
	MAX_SALARY = 1000000000
)

// env constant
const (
	CALC_LOG_INSERT_QUERY = "CALC_LOG_INSERT_QUERY"
)
