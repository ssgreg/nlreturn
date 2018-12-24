package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"path"

	"github.com/kisielk/gotool"
	"github.com/ssgreg/logf"
	"github.com/ssgreg/logftext"
	"golang.org/x/tools/go/loader"
)

var (
	verbose = flag.Bool("v", false, "verbose logs")
)

func main() {
	flag.Parse()

	writer, writerClose := logf.NewChannelWriter(logf.ChannelWriterConfig{
		Appender: logftext.NewAppender(os.Stderr, logftext.EncoderConfig{}),
	})
	defer writerClose()

	var logger *logf.Logger
	if *verbose {
		logger = logf.NewLogger(logf.LevelDebug, writer)
	} else {
		logger = logf.NewDisabledLogger()
	}

	importPaths := gotool.ImportPaths(flag.Args())
	if len(importPaths) == 0 {
		return
	}
	logger.Debug("parsed import paths", logf.Strings("import-paths", importPaths))

	err := handleImportPaths(logger, os.Stdout, importPaths)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func handleImportPaths(logger *logf.Logger, w io.Writer, importPaths []string) error {
	fset := token.NewFileSet()

	cfg := loader.Config{
		Fset: fset,
	}
	for _, importPath := range importPaths {
		cfg.Import(importPath)
	}

	program, err := cfg.Load()
	if err != nil {
		return err
	}
	handleProgram(logger, w, program, fset)

	return nil
}

func handleProgram(logger *logf.Logger, w io.Writer, program *loader.Program, fset *token.FileSet) {
	logger = logger.WithName("package")

	for _, pkg := range program.InitialPackages() {
		for _, file := range pkg.Files {
			packageLogger := logger.With(logf.Stringer("package", file.Name))
			packageLogger.Debug("handle package")

			handleFile(packageLogger, w, file, fset)
		}
	}
}

func handleFile(logger *logf.Logger, w io.Writer, file *ast.File, fset *token.FileSet) {
	var prevNode ast.Node
	logger = logger.WithName("file").With(logf.String("file", path.Base(fset.Position(file.Pos()).Filename)))
	logger.Debug("handle file")

	confirmedReturns := []int{}
	linesWithBlockEnd := []int{}

	ast.Inspect(file, func(node ast.Node) bool {
		if node != nil {
			defer func() {
				prevNode = node
			}()
			nodeLogger := logger.WithName("node").With(
				logf.String("type", fmt.Sprintf("%T", node)),
				logf.Int("pos", fset.Position(node.Pos()).Line),
				logf.Int("pos", fset.Position(node.End()).Line),
			)
			nodeLogger.Debug("got node")

			switch c := node.(type) {
			case *ast.CaseClause:
				if len(c.Body) > 0 {
					switch c.Body[0].(type) {
					case *ast.BranchStmt, *ast.ReturnStmt:
						confirmedReturns = append(confirmedReturns, fset.Position(c.Body[0].Pos()).Line)

						return true
					}
				}

			case *ast.CommClause:
				if len(c.Body) > 0 {
					switch c.Body[0].(type) {
					case *ast.BranchStmt, *ast.ReturnStmt:
						confirmedReturns = append(confirmedReturns, fset.Position(c.Body[0].Pos()).Line)

						return true
					}
				}

			case *ast.BlockStmt:
				linesWithBlockEnd = append(linesWithBlockEnd, fset.Position(node.End()).Line)
				logger.Debug("update line with block ends", logf.Ints("block-ends", linesWithBlockEnd))

				if len(c.List) > 0 {
					switch c.List[0].(type) {
					case *ast.BranchStmt, *ast.ReturnStmt:
						confirmedReturns = append(confirmedReturns, fset.Position(c.List[0].Pos()).Line)

						return true
					}
				}

			case *ast.AssignStmt, *ast.CallExpr, *ast.DeferStmt, *ast.ExprStmt, *ast.GoStmt:
				linesWithBlockEnd = append(linesWithBlockEnd, fset.Position(node.End()).Line)
				logger.Debug("update line with block ends", logf.Ints("block-ends", linesWithBlockEnd))

			case *ast.BranchStmt, *ast.ReturnStmt:
				for _, confirmedReturn := range confirmedReturns {
					if confirmedReturn == fset.Position(node.Pos()).Line {
						return true
					}
				}

				pos := fset.Position(node.Pos()).Line
				prevEnd := fset.Position(prevNode.End()).Line

				if pos-prevEnd > 1 {
					// Check if line before return is occupied by ending curly bracket.
					found := false
					for _, blockEnd := range linesWithBlockEnd {
						if pos-1 == blockEnd {
							found = true
						}
					}
					if !found {
						return true
					}
				}

				printErrorMessage(w, node, fset)
			}
		}

		return true
	})
}

func printErrorMessage(w io.Writer, node ast.Node, fset *token.FileSet) {
	nodeName := "return"

	switch c := node.(type) {
	case *ast.BranchStmt:
		nodeName = c.Tok.String()
	}

	w.Write([]byte(fmt.Sprintf("%s: %s with no blank line before\n", fset.Position(node.Pos()), nodeName)))
}
