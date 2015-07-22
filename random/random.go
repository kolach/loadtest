package random

import (
	"math/rand"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("random")

func init() {
	rand.Seed(1329351249374509812)	
}
