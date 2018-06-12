package main

import (
	"fmt"
	"github.com/AlpherJang/inkchain/core/chaincode/shim"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
	"strings"
	"strconv"
	"encoding/json"
)

const (
	SCHOOL_TYPE_FULL_TIME = 1
	SCHOOL_TYPE_PART_TIME = 2
	SEX_MALE              = 1
	SEX_FAMEL             = 2
	SCHOOL_PREFIX         = "school"
	STUDENT_PREFIX        = "student"
)

const (
	SPECIALTY_DEGREE = iota
	BACHELOR_DEGREE
	MASTER_DEGREE
	DOCTOR_DEGREE
)

const (
	STUDY_LOG_KEY      = "studyLog"
	GRADUATION_LOG_KEY = "graduationLog"
)

type School struct {
	Number       string `json:"number"`      //school number
	Name         string `json:"name"`        //school name
	Address      string `json:"address"`     //school address
	CreateAt     string `json:"createAt"`    //school create date
	Manager      string `json:"manager"`     //school manager
	SchoolLevel  int    `json:"schoolLevel"` //school level
	RegisterTime string `json:"registerTime"`
}

type Student struct {
	Number               string `json:"number"`               //personal card number
	StudentNumber        string `json:"studentNumber"`        //student number
	Name                 string `json:"name"`                 //student name
	Age                  int    `json:"age"`                  //student age
	Sex                  int    `json:"sex"`                  //student sex
	Credit               int    `json:"credit"`               //academic credit number
	Grade                string `json:"grade"`                //student grade
	Class                string `json:"class"`                //student class
	CurrentSchool        string `json:"currentSchool"`        //student current school
	CurrentSchoolName    string `json:"currentSchoolName"`    //student current school name
	CurrentLevel         int    `json:"currentLevel"`         //student current level
	AdmissionTime        string `json:"admissionTime"`        //student admission time
	GraduationSchool     string `json:"graduationSchool"`     //graduation school number
	GraduationSchoolName string `json:"graduationSchoolName"` //graduation school name
	GraduationTime       string `json:"graduationTime"`       //student graduation time
	GraduationLevel      int    `json:"graduationLevel"`      //student graduation level
	RegisterTime         string `json:"registerTime"`
}

type StudyLog struct {
	SchoolName    string `json:"schoolName"`    //study school name
	SchoolNumber  string `json:"schoolNumber"`  //study school number
	StudentName   string `json:"studentName"`   //student name
	StudentNumber string `json:"studentNumber"` //student number
	Level         int    `json:"level"`         //student study level
	Class         string `json:"class"`         //student study class
	Grade         string `json:"grade"`         //student study grade
	DateTime      string `json:"dateTime"`
}

type GraduationInfo struct {
	SchoolName      string `json:"schoolName"`      //school name
	SchoolNumber    string `json:"schoolNumber"`    //school number
	StudentName     string `json:"studentName"`     //student name
	StudentNumber   string `json:"studentNumber"`   //student number
	GraduationLevel int    `json:"graduationLevel"` //student graduation level
	Description     string `json:"description"`     //description
	Class           string `json:"class"`           //student class
	Grade           string `json:"grade"`           //student grade
	Credit          string `json:"credit"`          //credit number
	DateTime        string `json:"dateTime"`
}

const (
	RegisterSchool       = "registerSchool"
	RegisterStudent      = "registerStudent"
	Enrolment            = "enrolment"
	Graduate             = "graduate"
	QuerySchoolInfo      = "querySchoolInfo"
	QueryStudentInfo     = "queryStudentInfo"
	QueryStudentStudyLog = "queryStudentStudyLog"
	QueryStudentGraduationLog="queryStudentGraduationLog"
)

type StudentChaincode struct {
}

func (s *StudentChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (s *StudentChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	switch function {
	case RegisterSchool:
		if len(args) < 4 {
			return shim.Error("RegisterSchool, Incorrect number of arguments. Expecting 4")
		}
		return s.registerSchool(stub, args)
	case RegisterStudent:
		if len(args) < 3 {
			return shim.Error("RegisterStudent, Incorrect number of arguments. Expecting 3")
		}
		return s.registerStudent(stub, args)
	case Enrolment:
		if len(args) < 6 {
			return shim.Error("Enrolment, Incorrect number of arguments. Expecting 6")
		}
		return s.enrolment(stub, args)
	case Graduate:
		if len(args) < 5 {
			return shim.Error("Graduate, Incorrect number of arguments. Expecting 5")
		}
		return s.graduate(stub, args)
	case QuerySchoolInfo:
		if len(args) < 1 {
			return shim.Error("QuerySchoolInfo, Incorrect number of arguments. Expecting 1")
		}
		return s.querySchoolInfo(stub, args)
	case QueryStudentInfo:
		if len(args) < 1 {
			return shim.Error("QueryStudentInfo, Incorrect number of arguments. Expecting 1")
		}
		return s.queryStudentInfo(stub, args)
	case QueryStudentStudyLog:
		if len(args) < 1 {
			return shim.Error("QueryStudentStudyLog, Incorrect number of arguments. At least 1")
		}
		return s.queryStudentStudyLog(stub, args)
	case QueryStudentGraduationLog:
		if len(args)<1{
			return shim.Error("QueryStudentGraduationLog, Incorrect number of arguments. At least 1")
		}
		return s.queryStudentGraduationLog(stub,args)
	}
	return shim.Error("Function not found")
}

//school register
func (s *StudentChaincode) registerSchool(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//get transaction time
	timeStamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Failed to get transaction timestamp : " + err.Error())
	}
	registerTime := fmt.Sprintf("%d", timeStamp.GetSeconds())
	//get sender
	sender, err := stub.GetSender()
	if err != nil {
		return shim.Error("Failed to get sender info : " + err.Error())
	}
	//validate args
	number := strings.TrimSpace(strings.ToLower(args[0]))
	if len(number) <= 0 {
		return shim.Error("1st arg must be non-empty string")
	}
	name := strings.TrimSpace(strings.ToLower(args[1]))
	if len(name) <= 0 {
		return shim.Error("2st arg must be non-empty string")
	}
	address := strings.TrimSpace(strings.ToLower(args[2]))
	if len(address) <= 0 {
		return shim.Error("3st arg must be non-empty string")
	}
	level, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("3st arg failed to parse int : " + err.Error())
	}
	if level != SCHOOL_TYPE_FULL_TIME && level != SCHOOL_TYPE_PART_TIME {
		return shim.Error("3st arg must be 1 or 2")
	}
	createTime := strings.TrimSpace(strings.ToLower(args[4]))
	if len(createTime) <= 0 {
		return shim.Error("4st arg must be non-empty string")
	}
	//build school record
	school := &School{
		Number:       number,
		Name:         name,
		Address:      address,
		SchoolLevel:  level,
		Manager:      sender,
		CreateAt:     createTime,
		RegisterTime: registerTime,
	}
	schoolKey := SCHOOL_PREFIX + number
	//validate school exists
	old, err := stub.GetState(schoolKey)
	if err != nil {
		return shim.Error("Failed to check school info : " + err.Error())
	} else if old != nil {
		return shim.Error("School exists")
	}
	//marshal to json
	schoolJson, err := json.Marshal(school)
	if err != nil {
		return shim.Error("Failed to marshal school : " + err.Error())
	}
	err = stub.PutState(schoolKey, schoolJson)
	if err != nil {
		return shim.Error("Failed to save school data : " + err.Error())
	}
	return shim.Success(nil)
}

//student register
func (s *StudentChaincode) registerStudent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//get transaction time
	timeStamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Failed to get transaction timestamp : " + err.Error())
	}
	registerTime := fmt.Sprintf("%d", timeStamp.GetSeconds())
	//get sender
	sender, err := stub.GetSender()
	if err != nil {
		return shim.Error("Failed to get sender info : " + err.Error())
	}
	//validate args
	number := strings.TrimSpace(strings.ToLower(args[0]))
	if len(number) <= 0 {
		return shim.Error("1st arg must be non-empty")
	}
	name := strings.TrimSpace(strings.ToLower(args[1]))
	if len(name) <= 0 {
		return shim.Error("2st arg must be non-empty")
	}
	age, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("2st arg parse to int failed : " + err.Error())
	}
	if age <= 0 {
		return shim.Error("2st arg must be more than 0")
	}
	sex, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("3st arg parse to int failed : " + err.Error())
	}
	if sex != SEX_FAMEL && sex != SEX_MALE {
		return shim.Error("3st arg must be 0 or 1")
	}
	//build student struct
	student := &Student{
		Number:       number,
		Name:         name,
		Age:          age,
		Sex:          sex,
		RegisterTime: registerTime,
	}
	//check sender exists
	studentKey := STUDENT_PREFIX + sender
	old, err := stub.GetState(studentKey)
	if err != nil {
		return shim.Error("Failed to validate sender register status : " + err.Error())
	} else if old != nil {
		return shim.Error("Sender has registered")
	}
	//marshal json
	studentJson, err := json.Marshal(student)
	if err != nil {
		return shim.Error("Marshal student info error : " + err.Error())
	}
	err = stub.PutState(studentKey, studentJson)
	if err != nil {
		return shim.Error("Failed to save student info : " + err.Error())
	}
	return shim.Success(nil)
}

//student enrolment
func (s *StudentChaincode) enrolment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//get transaction time
	timeStamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Failed to get transaction timestamp : " + err.Error())
	}
	admissionTime := fmt.Sprintf("%d", timeStamp.GetSeconds())
	//get sender
	sender, err := stub.GetSender()
	if err != nil {
		return shim.Error("Failed to get sender info : " + err.Error())
	}
	//validate args
	schoolNumber := strings.TrimSpace(strings.ToLower(args[0]))
	if len(schoolNumber) <= 0 {
		return shim.Error("1st arg must be non-empty string")
	}
	studentAddr := strings.TrimSpace(strings.ToLower(args[1]))
	if len(studentAddr) <= 0 {
		return shim.Error("2st arg must be non-empty string")
	}
	grade := strings.TrimSpace(strings.ToLower(args[2]))
	if len(grade) <= 0 {
		return shim.Error("3st arg must be non-empty string")
	}
	class := strings.TrimSpace(strings.ToLower(args[3]))
	if len(class) <= 0 {
		return shim.Error("4st arg must be non-empty string")
	}
	studentNumber := strings.TrimSpace(strings.ToLower(args[4]))
	if len(studentNumber) <= 0 {
		return shim.Error("5st arg must be non-empty string")
	}
	level, err := strconv.Atoi(args[5])
	if err != nil {
		return shim.Error("6st parse int failed : " + err.Error())
	}
	if level < SPECIALTY_DEGREE && level > DOCTOR_DEGREE {
		return shim.Error("6st must be 0 , 1 , 2 , 3")
	}
	//check school info
	schoolJson, err := stub.GetState(SCHOOL_PREFIX + schoolNumber)
	if err != nil {
		return shim.Error("Failed to get school info : " + err.Error())
	} else if schoolJson == nil {
		return shim.Error("School not exists")
	}
	//unmarshal school info
	school := &School{}
	err = json.Unmarshal(schoolJson, school)
	if err != nil {
		return shim.Error("Failed unmarshal school info : " + err.Error())
	}
	//check sender
	if school.Manager != sender {
		return shim.Error("Authority failed")
	}
	//check student info
	studentKey := STUDENT_PREFIX + studentAddr
	studentJson, err := stub.GetState(studentKey)
	if err != nil {
		return shim.Error("Failed to get student info : " + err.Error())
	} else if studentJson == nil {
		return shim.Error("Student not exists")
	}
	//unmarshal student info
	student := &Student{}
	err = json.Unmarshal(studentJson, student)
	if err != nil {
		return shim.Error("Failed unmarshal student info : " + err.Error())
	}
	if student.CurrentSchool != "" {
		return shim.Error("Student has study at another school")
	}
	//update school info
	student.StudentNumber = studentNumber
	student.CurrentSchool = schoolNumber
	student.CurrentSchoolName = school.Name
	student.AdmissionTime = admissionTime
	student.Class = class
	student.Grade = grade
	student.CurrentLevel = level
	studentJson, err = json.Marshal(student)
	if err != nil {
		return shim.Error("Marshal student info to json failed : " + err.Error())
	}
	//update student info
	err = stub.PutState(studentKey, studentJson)
	if err != nil {
		return shim.Error("Failed to save student info : " + err.Error())
	}
	//save study log
	err = s.saveStudyLog(stub, student, school, admissionTime)
	if err != nil {
		return shim.Error("Failed to save study log : " + err.Error())
	}
	return shim.Success(nil)
}

//student graduate
func (s *StudentChaincode) graduate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//get transaction time
	timeStamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Failed to get transaction timestamp : " + err.Error())
	}
	graduateTime := fmt.Sprintf("%d", timeStamp.GetSeconds())
	//get sender
	sender, err := stub.GetSender()
	if err != nil {
		return shim.Error("Failed to get sender info : " + err.Error())
	}
	//check args
	schoolNumber := strings.TrimSpace(strings.ToLower(args[0]))
	if len(schoolNumber) <= 0 {
		return shim.Error("1st arg must be non-empty string")
	}
	studentAddr := strings.TrimSpace(strings.ToLower(args[1]))
	if len(studentAddr) <= 0 {
		return shim.Error("2st arg must be non-empty string")
	}
	description := strings.TrimSpace(strings.ToLower(args[2]))
	if len(description) <= 0 {
		return shim.Error("3st arg must be non-empty string")
	}
	creditNum := strings.TrimSpace(strings.ToLower(args[3]))
	if len(creditNum) <= 0 {
		return shim.Error("4st arg must be non-empty string")
	}
	credit, err := strconv.Atoi(args[4])
	if err != nil {
		return shim.Error("5st arg parse int failed : " + err.Error())
	}
	if credit <= 0 {
		return shim.Error("5st arg must be more than 0")
	}
	//check school
	schoolJson, err := stub.GetState(SCHOOL_PREFIX + schoolNumber)
	if err != nil {
		return shim.Error("Failed to get school info : " + err.Error())
	}
	school := &School{}
	err = json.Unmarshal(schoolJson, school)
	if err != nil {
		return shim.Error("Unmarshal school info failed : " + err.Error())
	}
	//check sender
	if school.Manager != sender {
		return shim.Error("Authority failed")
	}
	//check student info
	studentKey := STUDENT_PREFIX + studentAddr
	studentJson, err := stub.GetState(studentKey)
	if err != nil {
		return shim.Error("Failed to get student info : " + err.Error())
	}
	student := &Student{}
	err = json.Unmarshal(studentJson, student)
	if err != nil {
		return shim.Error("Unmarshal student info error : " + err.Error())
	}
	if student.CurrentSchool != school.Number {
		return shim.Error("Student not study at your school")
	}
	//update student info
	student.GraduationLevel = student.CurrentLevel
	student.GraduationSchool = school.Number
	student.GraduationSchoolName = school.Name
	student.GraduationTime = graduateTime
	student.Credit = credit
	student.CurrentSchool = ""
	studentJson, err = json.Marshal(student)
	if err != nil {
		return shim.Error("Student marshal to json failed : " + err.Error())
	}
	//update student state
	err = stub.PutState(studentKey, studentJson)
	if err != nil {
		return shim.Error("Failed to update student info : " + err.Error())
	}
	//save graduate log
	err = s.saveGraduateLog(stub, school, student, description, creditNum)
	if err != nil {
		return shim.Error("Failed to save graduate log : " + err.Error())
	}
	return shim.Success(nil)
}

//query school info
func (s *StudentChaincode) querySchoolInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	number := strings.TrimSpace(strings.ToLower(args[0]))
	if len(number) <= 0 {
		return shim.Error("1st arg must be non-empty string")
	}
	schoolKey := SCHOOL_PREFIX + number
	schoolJson, err := stub.GetState(schoolKey)
	if err != nil {
		return shim.Error("Failed to get school info : " + err.Error())
	} else if schoolJson == nil {
		return shim.Error("School not exists")
	}
	return shim.Success(schoolJson)
}

//query student info
func (s *StudentChaincode) queryStudentInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//check args
	studentAddr := strings.TrimSpace(strings.ToLower(args[0]))
	if len(studentAddr) <= 0 {
		return shim.Error("1st arg must be non-empty string")
	}
	//get student info
	studentKey := STUDENT_PREFIX + studentAddr
	studentJson, err := stub.GetState(studentKey)
	if err != nil {
		return shim.Error("Failed to get student info : " + err.Error())
	} else if studentJson != nil {
		return shim.Error("Student not exists")
	}
	//unmarshal student info
	student := &Student{}
	err = json.Unmarshal(studentJson, student)
	if err != nil {
		return shim.Error("Failed to unmarshal student info : " + err.Error())
	}
	result := map[string]interface{}{
		"studentInfo": student,
	}
	//get study log
	studentStudyLog, err := s.getStudentStudyLog(stub, []string{student.StudentNumber})
	if err == nil {
		result["studyLog"] = studentStudyLog
	}
	//get graduation info log
	graduationInfoLog, err := s.getStudentGraduationLog(stub, []string{student.StudentNumber})
	if err == nil {
		result["graduationLog"] = graduationInfoLog
	}
	//marshal result
	resultJson, err := json.Marshal(graduationInfoLog)
	if err != nil {
		return shim.Error("Failed to marshal result : " + err.Error())
	}
	return shim.Success(resultJson)
}

//query student study log
func (s *StudentChaincode) queryStudentStudyLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	studentAddr := strings.TrimSpace(strings.ToLower(args[0]))
	if len(studentAddr) <= 0 {
		return shim.Error("1st arg must be non-empty string")
	}
	studentJson, err := stub.GetState(STUDENT_PREFIX + studentAddr)
	if err != nil {
		return shim.Error("Failed to get student info : " + err.Error())
	} else if studentJson != nil {
		return shim.Error("Student not exists")
	}
	student := &Student{}
	err = json.Unmarshal(studentJson, student)
	if err != nil {
		return shim.Error("Unmarshal student info failed : " + err.Error())
	}
	keyArg := []string{student.StudentNumber}
	schoolNumber := strings.TrimSpace(strings.ToLower(args[1]))
	if len(schoolNumber) > 0 {
		keyArg = append(keyArg, schoolNumber)
	}
	result, err := s.getStudentStudyLog(stub, keyArg)
	if err != nil {
		return shim.Error("Failed to get student study log : " + err.Error())
	}
	resultJson, err := json.Marshal(result)
	if err != nil {
		return shim.Error("Failed to marshal result : " + err.Error())
	}
	return shim.Success(resultJson)
}

//query student graduation log
func (s *StudentChaincode) queryStudentGraduationLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	studentAddr := strings.TrimSpace(strings.ToLower(args[0]))
	if len(studentAddr) <= 0 {
		return shim.Error("1st arg must be non-empty string")
	}
	studentJson, err := stub.GetState(STUDENT_PREFIX + studentAddr)
	if err != nil {
		return shim.Error("Failed to get student info : " + err.Error())
	} else if studentJson != nil {
		return shim.Error("Student not exists")
	}
	student := &Student{}
	err = json.Unmarshal(studentJson, student)
	if err != nil {
		return shim.Error("Unmarshal student info failed : " + err.Error())
	}
	keyArg := []string{student.StudentNumber}
	schoolNumber := strings.TrimSpace(strings.ToLower(args[1]))
	if len(schoolNumber) > 0 {
		keyArg = append(keyArg, schoolNumber)
	}
	result, err := s.getStudentGraduationLog(stub, keyArg)
	if err != nil {
		return shim.Error("Failed to get student study log : " + err.Error())
	}
	resultJson, err := json.Marshal(result)
	if err != nil {
		return shim.Error("Failed to marshal result : " + err.Error())
	}
	return shim.Success(resultJson)
}

//get student study log
func (s *StudentChaincode) getStudentStudyLog(stub shim.ChaincodeStubInterface, keyArg []string) ([]StudyLog, error) {
	//get study log key
	resultsIterator, err := stub.GetStateByPartialCompositeKey(STUDY_LOG_KEY, keyArg)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	resultList := make([]StudyLog, 0)
	for i := 0; resultsIterator.HasNext(); i++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		value := responseRange.Value
		tmp := StudyLog{}
		err = json.Unmarshal(value, &tmp)
		if err == nil {
			resultList = append(resultList, tmp)
		}
	}
	return resultList, nil
}

//get student graduation log
func (s *StudentChaincode) getStudentGraduationLog(stub shim.ChaincodeStubInterface, keyArg []string) ([]GraduationInfo, error) {
	//get graduation key
	resultsIterator, err := stub.GetStateByPartialCompositeKey(GRADUATION_LOG_KEY, keyArg)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	resultList := make([]GraduationInfo, 0)
	for i := 0; resultsIterator.HasNext(); i++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		value := responseRange.Value
		tmp := GraduationInfo{}
		err = json.Unmarshal(value, &tmp)
		if err == nil {
			resultList = append(resultList, tmp)
		}
	}
	return resultList, nil
}

//save study log
func (s *StudentChaincode) saveStudyLog(stub shim.ChaincodeStubInterface, student *Student, school *School, enrolTime string) error {
	log := &StudyLog{
		SchoolName:    school.Name,
		SchoolNumber:  school.Number,
		StudentName:   student.Name,
		StudentNumber: student.StudentNumber,
		Level:         student.CurrentLevel,
		Class:         student.Class,
		Grade:         student.Grade,
		DateTime:      enrolTime,
	}
	//marshal study log data
	logJson, err := json.Marshal(log)
	if err != nil {
		return err
	}
	//create composite key
	compositeKey, err := stub.CreateCompositeKey(STUDY_LOG_KEY, []string{student.Number, school.Number, enrolTime})
	if err != nil {
		return err
	}
	err = stub.PutState(compositeKey, logJson)
	if err != nil {
		return err
	}
	return nil
}

//save graduate log
func (s *StudentChaincode) saveGraduateLog(stub shim.ChaincodeStubInterface, school *School, student *Student, description string, creditNumber string) error {
	graduateLog := &GraduationInfo{
		SchoolNumber:    school.Number,
		SchoolName:      school.Name,
		StudentNumber:   student.StudentNumber,
		StudentName:     student.Name,
		GraduationLevel: student.CurrentLevel,
		Description:     description,
		Class:           student.Class,
		Grade:           student.Grade,
		Credit:          creditNumber,
		DateTime:        student.GraduationTime,
	}
	//marshal graduate log data
	graduationJson, err := json.Marshal(graduateLog)
	if err != nil {
		return err
	}
	//create composite key
	compositeKey, err := stub.CreateCompositeKey(GRADUATION_LOG_KEY, []string{student.Number, school.Number, student.GraduationTime})
	if err != nil {
		return err
	}
	err = stub.PutState(compositeKey, graduationJson)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := shim.Start(new(StudentChaincode))
	if err == nil {
		fmt.Printf("Error starting WorkChaincode: %s", err)
	}
}