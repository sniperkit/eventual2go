package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/codegangsta/cli"
)

//flag vars

var pkgName string
var typeNames []string
var names []string
var export bool
var verbose bool
var print bool
var outputFile string

func main() {

	app := cli.NewApp()

	app.Name = "event_generator"

	app.Usage = "Generate typed events for eventual2go"

	app.Version = "0.1"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "export, e",
			Usage:       "exports the events by rasing the first letter of the typename, use for unexported types or internals, like string",
			Destination: &export,
		},
		cli.BoolFlag{
			Name:        "verbose",
			Usage:       "verbose output",
			Destination: &verbose,
		},
		cli.StringFlag{
			Name:        "package, p",
			Usage:       "name of package",
			Destination: &pkgName,
		},
		cli.BoolFlag{
			Name:        "print, pr",
			Usage:       "outputs to the commandline instead writing to file",
			Destination: &print,
		},
		cli.StringFlag{
			Name:        "output, o",
			Usage:       "file to output, leave empty for setting output to %CURRENT_DIR%/TYPE_NAME_events.go. If multiple types are specified, the first is used for the filename",
			Destination: &outputFile,
		},
		cli.StringSliceFlag{
			Name:  "type, t",
			Usage: "name of the type(s), set multiple for multiple types",
		},
		cli.StringSliceFlag{
			Name:  "name, n",
			Usage: "name of the type(s) to be used by methods and type declartion, set multiple for multiple types. Usful if generate events for slices.",
		},
	}

	app.Action = run

	app.Run(os.Args)

}

func run(c *cli.Context) {

	logger := log.New(os.Stdout, "event_generator: ", 0)

	if pkgName == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("error retrieving pkgName from working dir", err)
		}
		pkgName = filepath.Base(cwd)
	}

	typeNames = c.StringSlice("type")
	names = c.StringSlice("name")

	if len(typeNames) == 0 {
		cli.ShowAppHelp(c)
		logger.Fatal("No types defined")
	}
	if verbose {
		logger.Println("Generating for Package", pkgName)
	}
	out := &bytes.Buffer{}

	generateHeader(out)

	for i, typeName := range typeNames {
		if typeName == "" {
			logger.Fatal("Empty type name")
		}

		name := typeName

		if i < len(names) {
			name = names[i]
		} else if export {
			name = strings.ToUpper(string(name[0])) + name[1:]
		}
		generateType(out, typeName, name)
	}
	if print {

		fmt.Print(out)
	} else {
		if outputFile == "" {
			var name string
			if len(names) != 0 {
				name = names[0]
			} else {
				name = typeNames[0]
			}
			outputFile = fmt.Sprintf("%s_events.go", name)
		}
		outputFile = strings.ToLower(outputFile)
		path, err := filepath.Abs(outputFile)
		if err != nil {
			logger.Fatalf("error creating file", err)
		}

		err = ioutil.WriteFile(path, out.Bytes(), os.ModePerm)
		if err != nil {
			logger.Fatalf("error creating file", err)
		}

	}

}

func generateHeader(out io.Writer) {
	t, _ := template.New("").Parse(tmplHeader)

	header := struct{ PkgName string }{pkgName}

	t.Execute(out, header)
}

func generateType(out io.Writer, typeName, name string) {
	t, _ := template.New("").Parse(tmplType)

	tname := struct {
		TypeName string
		Name     string
	}{typeName, name}

	t.Execute(out, tname)
}

var tmplHeader = `
/*
 * generated by event_generator
 *
 * DO NOT EDIT
 */

package {{.PkgName}}

import "github.com/joernweissenborn/eventual2go"

`

var tmplType = `

type {{.Name}}Completer struct {
	*eventual2go.Completer
}

func New{{.Name}}Completer() *{{.Name}}Completer {
	return &{{.Name}}Completer{eventual2go.NewCompleter()}
}

func (c *{{.Name}}Completer) Complete(d {{.TypeName}}) {
	c.Completer.Complete(d)
}

func (c *{{.Name}}Completer) Future() *{{.Name}}Future {
	return &{{.Name}}Future{c.Completer.Future()}
}

type {{.Name}}Future struct {
	*eventual2go.Future
}

func (f *{{.Name}}Future) Result() {{.TypeName}} {
	return f.Future.Result().({{.TypeName}})
}

type {{.Name}}CompletionHandler func({{.TypeName}}) {{.TypeName}}

func (ch {{.Name}}CompletionHandler) toCompletionHandler() eventual2go.CompletionHandler {
	return func(d eventual2go.Data) eventual2go.Data {
		return ch(d.({{.TypeName}}))
	}
}

func (f *{{.Name}}Future) Then(ch {{.Name}}CompletionHandler) *{{.Name}}Future {
	return &{{.Name}}Future{f.Future.Then(ch.toCompletionHandler())}
}

func (f *{{.Name}}Future) AsChan() chan {{.TypeName}} {
	c := make(chan {{.TypeName}}, 1)
	cmpl := func(d chan {{.TypeName}}) {{.Name}}CompletionHandler {
		return func(e {{.TypeName}}) {{.TypeName}} {
			d <- e
			close(d)
			return e
		}
	}
	ecmpl := func(d chan {{.TypeName}}) eventual2go.ErrorHandler {
		return func(error) (eventual2go.Data, error) {
			close(d)
			return nil, nil
		}
	}
	f.Then(cmpl(c))
	f.Err(ecmpl(c))
	return c
}

type {{.Name}}StreamController struct {
	*eventual2go.StreamController
}

func New{{.Name}}StreamController() *{{.Name}}StreamController {
	return &{{.Name}}StreamController{eventual2go.NewStreamController()}
}

func (sc *{{.Name}}StreamController) Add(d {{.TypeName}}) {
	sc.StreamController.Add(d)
}

func (sc *{{.Name}}StreamController) Join(s *{{.Name}}Stream) {
	sc.StreamController.Join(s.Stream)
}

func (sc *{{.Name}}StreamController) JoinFuture(f *{{.Name}}Future) {
	sc.StreamController.JoinFuture(f.Future)
}

func (sc *{{.Name}}StreamController) Stream() *{{.Name}}Stream {
	return &{{.Name}}Stream{sc.StreamController.Stream()}
}

type {{.Name}}Stream struct {
	*eventual2go.Stream
}

type {{.Name}}Subscriber func({{.TypeName}})

func (l {{.Name}}Subscriber) toSubscriber() eventual2go.Subscriber {
	return func(d eventual2go.Data) { l(d.({{.TypeName}})) }
}

func (s *{{.Name}}Stream) Listen(ss {{.Name}}Subscriber) *eventual2go.Completer {
	return s.Stream.Listen(ss.toSubscriber())
}

func (s *{{.Name}}Stream) ListenNonBlocking(ss {{.Name}}Subscriber) *eventual2go.Completer {
	return s.Stream.ListenNonBlocking(ss.toSubscriber())
}

type {{.Name}}Filter func({{.TypeName}}) bool

func (f {{.Name}}Filter) toFilter() eventual2go.Filter {
	return func(d eventual2go.Data) bool { return f(d.({{.TypeName}})) }
}

func to{{.Name}}FilterArray(f ...{{.Name}}Filter) (filter []eventual2go.Filter){

	filter = make([]eventual2go.Filter, len(f))
	for i, el := range f {
		filter[i] = el.toFilter()
	}
	return
}

func (s *{{.Name}}Stream) Where(f ...{{.Name}}Filter) *{{.Name}}Stream {
	return &{{.Name}}Stream{s.Stream.Where(to{{.Name}}FilterArray(f...)...)}
}

func (s *{{.Name}}Stream) WhereNot(f ...{{.Name}}Filter) *{{.Name}}Stream {
	return &{{.Name}}Stream{s.Stream.WhereNot(to{{.Name}}FilterArray(f...)...)}
}

func (s *{{.Name}}Stream) TransformWhere(t eventual2go.Transformer, f ...{{.Name}}Filter) *eventual2go.Stream {
	return s.Stream.TransformWhere(t, to{{.Name}}FilterArray(f...)...)
}

func (s *{{.Name}}Stream) Split(f {{.Name}}Filter) (*{{.Name}}Stream, *{{.Name}}Stream)  {
	return s.Where(f), s.WhereNot(f)
}

func (s *{{.Name}}Stream) First() *{{.Name}}Future {
	return &{{.Name}}Future{s.Stream.First()}
}

func (s *{{.Name}}Stream) FirstWhere(f... {{.Name}}Filter) *{{.Name}}Future {
	return &{{.Name}}Future{s.Stream.FirstWhere(to{{.Name}}FilterArray(f...)...)}
}

func (s *{{.Name}}Stream) FirstWhereNot(f ...{{.Name}}Filter) *{{.Name}}Future {
	return &{{.Name}}Future{s.Stream.FirstWhereNot(to{{.Name}}FilterArray(f...)...)}
}

func (s *{{.Name}}Stream) AsChan() (c chan {{.TypeName}}, stop *eventual2go.Completer) {
	c = make(chan {{.TypeName}})
	stop = s.Listen(pipeTo{{.Name}}Chan(c))
	stop.Future().Then(close{{.Name}}Chan(c))
	return
}

func pipeTo{{.Name}}Chan(c chan {{.TypeName}}) {{.Name}}Subscriber {
	return func(d {{.TypeName}}) {
		c <- d
	}
}

func close{{.Name}}Chan(c chan {{.TypeName}}) eventual2go.CompletionHandler {
	return func(d eventual2go.Data) eventual2go.Data {
		close(c)
		return nil
	}
}

type {{.Name}}Collector struct {
	*eventual2go.Collector
}

func New{{.Name}}Collector() *{{.Name}}Collector {
	return &{{.Name}}Collector{eventual2go.NewCollector()}
}

func (c *{{.Name}}Collector) Add(d {{.TypeName}}) {
	c.Collector.Add(d)
}

func (c *{{.Name}}Collector) AddFuture(f *{{.Name}}Future) {
	c.Collector.Add(f.Future)
}

func (c *{{.Name}}Collector) AddStream(s *{{.Name}}Stream) {
	c.Collector.AddStream(s.Stream)
}

func (c *{{.Name}}Collector) Get() {{.TypeName}} {
	return c.Collector.Get().({{.TypeName}})
}

func (c *{{.Name}}Collector) Preview() {{.TypeName}} {
	return c.Collector.Preview().({{.TypeName}})
}
`
