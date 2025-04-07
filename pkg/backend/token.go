package backend

import(
	"github.com/matoous/go-nanoid/v2"
)

func CreateToken() (string, error) {
	return gonanoid.New()
}