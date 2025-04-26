package grouter

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func getRouteTypeStr(route *Route) (string, error) {
	switch route.rtype {
	case GET:
		return "GET", nil
	case POST:
		return "POST", nil
	case PUT:
		return "PUT", nil
	case PATCH:
		return "PATCH", nil
	case DELETE:
		return "DELETE", nil
	}

	return "", &InvalidRouteTypeError{route: route}
}

func printAndExit(data any, exitCode int) {
	fmt.Println(data)
	os.Exit(exitCode)
}

func generateUniqueNumber() int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int()
}
