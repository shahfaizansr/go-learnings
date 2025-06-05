package hrmsmodels

import (
	"github.com/remiges-tech/alya/wscutils"
)

type EmployeeRequestModel struct {
	Name          string `json:"name" validate:"required,min=5,max=100"`
	EmpCode       string `json:"empcode" validate:"required,max=7"`
	Gender        string `json:"gender" validate:"required"`
	DOB           string `json:"dob" validate:"required"`
	DOJ           string `json:"doj" validate:"required"`
	Salary        string `json:"salary" validate:"required"`
	Reportsto     string `json:"reportsto"`
	Designation   string `json:"designation" validate:"required,oneof=Admin Manager TeamLead Developer Tester HR"`
	BriefBio      string `json:"briefbio" validate:"required"`
	Interests     string `json:"interests" validate:"omitempty"`
	AddedFromIP   string `json:"added_from_ip"`
	UpdatedFromIP string `json:"updated_from_ip"`
}

type EmployeeResponseModel struct {
	Name          string `json:"name"`
	EmpCode       string `json:"empcode"`
	Gender        string `json:"gender"`
	DOB           string `json:"dob"`
	DOJ           string `json:"doj"`
	Salary        string `json:"salary"`
	Reportsto     string `json:"reportsto"`
	Designation   string `json:"designation"`
	BriefBio      string `json:"briefbio"`
	Interests     string `json:"interests"`
	AddedAt       string `json:"added_at"`
	AddedBy       string `json:"added_by"`
	UpdatedAt     string `json:"updated_at"`
	UpdatedBy     string `json:"updated_by"`
	AddedFromIP   string `json:"added_from_ip"`
	UpdatedFromIP string `json:"updated_from_ip"`
	IsActive      bool   `json:"is_active"`
}

type UpdateEmployeeRequestModel struct {
	Name          wscutils.Optional[string] `json:"name"`
	EmpCode       string                    `json:"empcode" validate:"required,max=7"`
	Gender        wscutils.Optional[string] `json:"gender"`
	DOB           wscutils.Optional[string] `json:"dob"`
	DOJ           wscutils.Optional[string] `json:"doj"`
	Salary        wscutils.Optional[string] `json:"salary"`
	Reportsto     wscutils.Optional[string] `json:"reportsto"`
	Designation   wscutils.Optional[string] `json:"designation"`
	BriefBio      wscutils.Optional[string] `json:"briefbio"`
	Interests     wscutils.Optional[string] `json:"interests"`
	UpdatedFromIP wscutils.Optional[string] `json:"updated_from_ip"`
}

type EmployeeDynamic map[string]interface{}
type ListEmployeesRequestModel struct {
	Gender  wscutils.Optional[string] `json:"gender,omitempty"`
	Salary  wscutils.Optional[string] `json:"salary,omitempty"`
	DOB     wscutils.Optional[string] `json:"dob,omitempty"`
	DOJ     wscutils.Optional[string] `json:"doj,omitempty"`
	Name    wscutils.Optional[string] `json:"name,omitempty"`
	EmpCode wscutils.Optional[string] `json:"empcode,omitempty"`
}

type ListEmployeesModel struct {
	AddedAt      wscutils.Optional[string] `json:"added_at"`
	UpdatedAt    wscutils.Optional[string] `json:"updated_at"`
	Fields       []string                  `json:"fields" validate:"required"`
	FilterParams ListEmployeesRequestModel `json:"filter_param" validate:"required"`
	Start        wscutils.Optional[int]    `json:"start"`
	Length       wscutils.Optional[int]    `json:"length"`
	IsActive     wscutils.Optional[bool]   `json:"is_active"`
	SortBy       string                    `json:"sortBy"`  // e.g. "name", "dob", "salary"
	SortDir      string                    `json:"sortDir"` // "A" or "D"
}

type ProofRequestModel struct {
	EmpCode   string            `json:"empcode" validate:"required,max=10"`
	Type      string            `json:"type" validate:"required,oneof=poi poc por"`
	Size      int64             `json:"size,omitempty"`
	MimeType  string            `json:"mimetype" validate:"required,oneof=image/jpeg application/pdf"`
	IssuedAt  string            `json:"issuedat" validate:"required"`
	ValidFrom string            `json:"validfrom" validate:"required"`
	ValidTo   string            `json:"validto" validate:"required"`
	IsActive  bool              `json:"isactive"`
	Version   int               `json:"version"`
	Details   ProofDetailsModel `json:"details" validate:"required"`
}

type ProofDetailsModel struct {
	Identifier string `json:"identifier" validate:"required"`
	Address    string `json:"address,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	Country    string `json:"country,omitempty"`
	Pincode    string `json:"pincode,omitempty" validate:"len=6,numeric"`
	Name       string `json:"name"`
	DOB        string `json:"dob"`
}

type ProofResponseModel struct {
	EmpCode   string            `json:"empcode"`
	Type      string            `json:"type"`
	Size      int64             `json:"size"`
	MimeType  string            `json:"mimetype"`
	IssuedAt  string            `json:"issuedat"`
	ValidFrom string            `json:"validfrom"`
	ValidTo   string            `json:"validto"`
	IsActive  bool              `json:"isactive"`
	Version   int               `json:"version"`
	Added     string            `json:"added,omitempty"`
	Updated   string            `json:"updated,omitempty"`
	Details   ProofDetailsModel `json:"details"`
}

type Item struct {
	Message string `json:"message"`
	EmpCode string `json:"empcode"`
}
