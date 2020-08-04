package errgen
import (
	"flag"
	"github.com/tallstoat/pbparser"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	errPkgName = flag.String("pkg_name", "", "Specify a new package name for generated error file, proto package for default")
	//errFilename = flag.String("err_file", "", "Specify err filename, error enum element name for default")
	errFilePath = flag.String("err_path", "./", "Specify a the path for output err file, current directory for default")
)

const ERR_SIGN string = "Err"

func GenerateErrFileFromProto(inputFilename string){
	pf, err := pbparser.ParseFile(inputFilename)
	if err != nil {
		fmt.Println(err)
		return
	}
	GenerateErrFile(pf)
}

func GenerateErrFile(pf pbparser.ProtoFile) {
	var code string

	errEnums := getErrEnum(pf)
	if len(errEnums) == 0 {
		fmt.Println("文件中没有以Err为前缀或后缀的enum类型")
		return
	}

	// 设置 package 和 import 部分
	var pkgName string
	if *errPkgName == "" {  // 如果命令行参数中没指定生成文件的路径和名称
		pkgName = pf.PackageName // 默认为 .proto 文件中的package名
	} else {
		pkgName = *errPkgName
	}
	header := fmt.Sprintf(headerTemplate, pkgName)
	code += header

	// error 类型定义 和 枚举定义
	for _, enum := range errEnums {
		errDefCode := getErrDefCode(enum)
		errFuncCode := getErrFuncCode(enum)
		code += errDefCode + "\n" + errFuncCode
	}

	// 输出文件的名称
	var filename string
	//if *errFilename == "" {
	//	filename = strings.ToLower(errEnums[0].Name)  // 使用其中第一个enum类型的名称
	//} else {
	//	filename = *errFilename
	//}

	filename = strings.ToLower(errEnums[0].Name)
	pathLen := len(*errFilePath)
	if '/' != (*errFilePath)[pathLen-1]{   // 如果命令行参数中输入的路径后面没带 / ，就在后面给补充一个
		*errFilePath += "/"
	}
	filename = *errFilePath + filename + ".go"
	file, _  := os.Create(filename)
	file.WriteString(code)

}

// GetErrEnum 从ProtoFile中读取enum类型，返回前缀或后缀为 Err 的enum列表
func getErrEnum(pf pbparser.ProtoFile) (errEnums []pbparser.EnumElement){
	for _, enum := range pf.Enums{ // []EnumElements
		if strings.HasPrefix(enum.Name, ERR_SIGN){
			errEnums = append(errEnums, enum)
		}
	}
	return errEnums
}

// getErrDefCode 生成error类型定义和相应错误常量
// type ERR_TYPE uint32
// const ( ERR_CONST_1 ERR_TYPE = 1, ...)
func getErrDefCode(errEnum pbparser.EnumElement) string{
	constCodeBlock := getErrConstCodeBlock(errEnum)
	errDefCode := fmt.Sprintf(errDefTemplate, errEnum.Name, constCodeBlock)
	return errDefCode
}

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


// getErrFuncCode 生成一个 根据自身数值返回相应error 的函数代码
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
	// 如果在 .proto 文件中使用 注释 描述了错误，就直接使用该注释
	if errConst.Documentation != "" {
		return errConst.Documentation
	}
	// 如果没有使用注释，就把tag和错误名组合成错误信息
	index := strings.LastIndex(errConst.Name, ERR_SIGN)
	var errName string
	if index != -1 {
		errName = errConst.Name[index:]
	} else {
		errName = errConst.Name
	}
	return strconv.Itoa(errConst.Tag) + ": " + errName
}