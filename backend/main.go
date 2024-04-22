package main

import (
	"fmt"
	"os"

	"github.com/svenliebig/env"
)

func main() {
	env.Load()
	fmt.Println(os.Getenv("THE_MOVIE_DB_API_KEY"))
}
