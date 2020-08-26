package errgen


var (

	headerTemplate string = `
package %s	
`
	// PackageName, EnumElement.Name/ErrTypeName,  const block(enum中定义的常量)
	errDefTemplate string = `
type %s uint32

const (
%s
)
`
	// EnumElement.Name/ErrTypeName， case block
	errFuncTemplate string = `
func (e %s) Error() string {
	switch e {
%s
	}
	return ""
}
`
	// ErrConst, ErrConst
	caseTemplate string = `	case %ERR_CONST%:
		return "%ERR_MESSAGE%"
`
	// ErrConst,  ErrType, ErrTag
	constTemplate string = `	%ERR_CONST% %ERR_TYPE% = %ERR_TAG%
`

)