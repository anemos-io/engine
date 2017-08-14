package anemos

import "fmt"

func (uri Uri) String() string {

	return fmt.Sprintf("%s:%s:%s:%s:%s:%s",
		uri.Kind,
		uri.Provider,
		uri.Operation,
		uri.Name,
		uri.Id,
		uri.Status)
}
