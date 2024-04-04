package utility

import (
	"fmt"

	"github.com/google/uuid"
)

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func GenerateUUID() string {
	u, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Errorf("error generating uuid: %v", err))
	}
	return u.String()
}
