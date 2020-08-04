
package myerr2	
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
		return errors.New("注释1")
	case Err_Transgender:
		return errors.New("注释2")
	case Err_Grade_Zero:
		return errors.New("注释3")

	}
	return nil
}
