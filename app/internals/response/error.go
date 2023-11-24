package response

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/muasx88/stealth_commerce/app/domain"
)

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckInt(value *int, name string) error {
	if value == nil {
		return nil
	}

	fmt.Println(reflect.ValueOf(*value).Kind())
	if reflect.ValueOf(*value).Kind() != reflect.Int {
		return fmt.Errorf("'%s' value is not of type integer or number", name)
	}

	return nil
}

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if code, ok := domain.ErrorMap[err]; ok {
		return code
	} else {
		return http.StatusInternalServerError
	}
}
