// Code generated by "stringer -type=AlignHorizontal"; DO NOT EDIT.

package gi

import (
	"fmt"
	"strconv"
)

const _AlignHorizontal_name = "AlignLeftAlignHCenterAlignRightAlignJustify"

var _AlignHorizontal_index = [...]uint8{0, 9, 21, 31, 43}

func (i AlignHorizontal) String() string {
	if i < 0 || i >= AlignHorizontal(len(_AlignHorizontal_index)-1) {
		return "AlignHorizontal(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _AlignHorizontal_name[_AlignHorizontal_index[i]:_AlignHorizontal_index[i+1]]
}

func StringToAlignHorizontal(s string) (AlignHorizontal, error) {
	for i := 0; i < len(_AlignHorizontal_index)-1; i++ {
		if s == _AlignHorizontal_name[_AlignHorizontal_index[i]:_AlignHorizontal_index[i+1]] {
			return AlignHorizontal(i), nil
		}
	}
	return 0, fmt.Errorf("String %v is not a valid option for type AlignHorizontal", s)
}