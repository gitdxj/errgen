package main 
import (
	"github.com/tallstoat/pbparser"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main(){
	filename := "./demo/student.proto"
	pf, err := pbparser.ParseFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Println("package name is", pf.PackageName)
	//
	//for _, enum := range pf.Enums {
	//	fmt.Println("Enum type is", enum.Name)
	//	for _, constant := range enum.EnumConstants {
	//		fmt.Println(constant.Name, constant.Tag)
	//	}
	//}

	GenerateErrFile(pf, "./demo/generated_err.go")
}

const ERR_PREFIX string = "Err"


// GetErrEnum 从ProtoFile中读取enum类型，返回带有Err前缀的enum列表
func GetErrEnum(pf pbparser.ProtoFile) (errEnums []pbparser.EnumElement){
	for _, enum := range pf.Enums{ // []EnumElements
		if strings.HasPrefix(enum.Name, ERR_PREFIX){
			errEnums = append(errEnums, enum)
		}
	}
	return errEnums
}

func GenerateErrFile(pf pbparser.ProtoFile, filename string) {
	var code string

	errEnums := GetErrEnum(pf)

	packageName := pf.PackageName
	header := fmt.Sprintf(headerTemplate, packageName)
	code += header

	for _, enum := range errEnums {
		errDefCode := getErrDefCode(enum)
		errFuncCode := getErrFuncCode(enum)
		code += errDefCode + "\n" + errFuncCode
	}

	file, _  := os.Create(filename)
	file.WriteString(code)

}


func getErrDefCode(errEnum pbparser.EnumElement) string{
	constCodeBlock := getErrConstCodeBlock(errEnum)
	errDefCode := fmt.Sprintf(errDefTemplate, errEnum.Name, constCodeBlock)
	return errDefCode
}

// const block
func getErrConstCodeBlock(errEnum pbparser.EnumElement) string {
	var block string
	errType := errEnum.Name
	for _, errConst := range errEnum.EnumConstants {
		newConst := getErrConstCode(errConst.Name, errType, strconv.Itoa(errConst.Tag))
		block += newConst
	}
	return block
}

func getErrConstCode(errConst, errType, errTag string) string{
	replacer := strings.NewReplacer("%ERR_CONST%", errConst,
				"%ERR_TYPE%", errType,
				"%ERR_TAG%", errTag)
	return replacer.Replace(constTemplate)
}



func getErrFuncCode(errEnum pbparser.EnumElement) string{
	caseCodeBlock := getErrCaseCodeBlock(errEnum)
	errFuncCode := fmt.Sprintf(errFuncTemplate, errEnum.Name, caseCodeBlock)
	return errFuncCode
}

// case block
func getErrCaseCodeBlock(errEnum pbparser.EnumElement) string{
	var block string
	for _, errConst := range errEnum.EnumConstants {
		errMessage := getErrMessage(errConst)
		newConst := getErrCaseCode(errConst.Name, errMessage)
		block += newConst
	}
	return block
}

func getErrCaseCode(errConst, errMessage string) string{
	replacer := strings.NewReplacer("%ERR_CONST%", errConst,
		"%ERR_MESSAGE%", errMessage)
	return replacer.Replace(caseTemplate)
}

// getErrMessage 得到 通过errors.New(string)创建新error时的 错误信息
func getErrMessage(errConst pbparser.EnumConstantElement) string{
	index := strings.LastIndex(errConst.Name, ERR_PREFIX)
	var errName string
	if index != -1 {
		errName = errConst.Name[index:]
	} else {
		errName = errConst.Name
	}
	return strconv.Itoa(errConst.Tag) + ": " + errName
}