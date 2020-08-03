
package student	
import (
	"errors"
)

type ErrStudent uint32

const (
	Err_Name_Cannot_Print ErrStudent = 0
	Err_Transgender ErrStudent = 1
	Err_Grade_Zero ErrStudent = 2

)


func (e ErrStudent) Error() error {
	switch e {
	case Err_Name_Cannot_Print:
		return errors.New("0: Err_Name_Cannot_Print")
	case Err_Transgender:
		return errors.New("1: Err_Transgender")
	case Err_Grade_Zero:
		return errors.New("2: Err_Grade_Zero")

	}
	return nil
}
