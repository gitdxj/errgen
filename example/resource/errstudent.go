
package student	

type ErrStudent uint32

const (
	Err_Name_Cannot_Print ErrStudent = 0
	Err_Transgender ErrStudent = 1
	Err_Grade_Zero ErrStudent = 2

)


func (e ErrStudent) Error() string {
	switch e {
	case Err_Name_Cannot_Print:
		return "0: Err_Name_Cannot_Print"
	case Err_Transgender:
		return "注释2"
	case Err_Grade_Zero:
		return "注释3"

	}
	return ""
}
