package main

import (
	"github.com/ssgreg/nlreturn/pkg/nlreturn"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(nlreturn.NewAnalyzer())
}
