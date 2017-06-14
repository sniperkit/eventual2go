
/*
 * generated by event_generator
 *
 * DO NOT EDIT
 */

package typed_events

import "github.com/joernweissenborn/eventual2go"



type StringCompleter struct {
	*eventual2go.Completer
}

func NewStringCompleter() *StringCompleter {
	return &StringCompleter{eventual2go.NewCompleter()}
}

func (c *StringCompleter) Complete(d string) {
	c.Completer.Complete(d)
}

func (c *StringCompleter) Future() *StringFuture {
	return &StringFuture{c.Completer.Future()}
}

type StringFuture struct {
	*eventual2go.Future
}

func (f *StringFuture) Result() string {
	return f.Future.Result().(string)
}

type StringCompletionHandler func(string) string

func (ch StringCompletionHandler) toCompletionHandler() eventual2go.CompletionHandler {
	return func(d eventual2go.Data) eventual2go.Data {
		return ch(d.(string))
	}
}

func (f *StringFuture) Then(ch StringCompletionHandler) *StringFuture {
	return &StringFuture{f.Future.Then(ch.toCompletionHandler())}
}

func (f *StringFuture) AsChan() chan string {
	c := make(chan string, 1)
	cmpl := func(d chan string) StringCompletionHandler {
		return func(e string) string {
			d <- e
			close(d)
			return e
		}
	}
	ecmpl := func(d chan string) eventual2go.ErrorHandler {
		return func(error) (eventual2go.Data, error) {
			close(d)
			return nil, nil
		}
	}
	f.Then(cmpl(c))
	f.Err(ecmpl(c))
	return c
}

type StringStreamController struct {
	*eventual2go.StreamController
}

func NewStringStreamController() *StringStreamController {
	return &StringStreamController{eventual2go.NewStreamController()}
}

func (sc *StringStreamController) Add(d string) {
	sc.StreamController.Add(d)
}

func (sc *StringStreamController) Join(s *StringStream) {
	sc.StreamController.Join(s.Stream)
}

func (sc *StringStreamController) JoinFuture(f *StringFuture) {
	sc.StreamController.JoinFuture(f.Future)
}

func (sc *StringStreamController) Stream() *StringStream {
	return &StringStream{sc.StreamController.Stream()}
}

type StringStream struct {
	*eventual2go.Stream
}

type StringSubscriber func(string)

func (l StringSubscriber) toSubscriber() eventual2go.Subscriber {
	return func(d eventual2go.Data) { l(d.(string)) }
}

func (s *StringStream) Listen(ss StringSubscriber) *eventual2go.Completer {
	return s.Stream.Listen(ss.toSubscriber())
}

type StringFilter func(string) bool

func (f StringFilter) toFilter() eventual2go.Filter {
	return func(d eventual2go.Data) bool { return f(d.(string)) }
}

func toStringFilterArray(f ...StringFilter) (filter []eventual2go.Filter){

	filter = make([]eventual2go.Filter, len(f))
	for i, el := range f {
		filter[i] = el.toFilter()
	}
	return
}

func (s *StringStream) Where(f ...StringFilter) *StringStream {
	return &StringStream{s.Stream.Where(toStringFilterArray(f...)...)}
}

func (s *StringStream) WhereNot(f ...StringFilter) *StringStream {
	return &StringStream{s.Stream.WhereNot(toStringFilterArray(f...)...)}
}

func (s *StringStream) TransformWhere(t eventual2go.Transformer, f ...StringFilter) *eventual2go.Stream {
	return s.Stream.TransformWhere(t, toStringFilterArray(f...)...)
}

func (s *StringStream) Split(f StringFilter) (*StringStream, *StringStream)  {
	return s.Where(f), s.WhereNot(f)
}

func (s *StringStream) First() *StringFuture {
	return &StringFuture{s.Stream.First()}
}

func (s *StringStream) FirstWhere(f... StringFilter) *StringFuture {
	return &StringFuture{s.Stream.FirstWhere(toStringFilterArray(f...)...)}
}

func (s *StringStream) FirstWhereNot(f ...StringFilter) *StringFuture {
	return &StringFuture{s.Stream.FirstWhereNot(toStringFilterArray(f...)...)}
}

func (s *StringStream) AsChan() (c chan string, stop *eventual2go.Completer) {
	c = make(chan string)
	stop = s.Listen(pipeToStringChan(c))
	stop.Future().Then(closeStringChan(c))
	return
}

func pipeToStringChan(c chan string) StringSubscriber {
	return func(d string) {
		c <- d
	}
}

func closeStringChan(c chan string) eventual2go.CompletionHandler {
	return func(d eventual2go.Data) eventual2go.Data {
		close(c)
		return nil
	}
}

type StringCollector struct {
	*eventual2go.Collector
}

func NewStringCollector() *StringCollector {
	return &StringCollector{eventual2go.NewCollector()}
}

func (c *StringCollector) Add(d string) {
	c.Collector.Add(d)
}

func (c *StringCollector) AddFuture(f *StringFuture) {
	c.Collector.Add(f.Future)
}

func (c *StringCollector) AddStream(s *StringStream) {
	c.Collector.AddStream(s.Stream)
}

func (c *StringCollector) Get() string {
	return c.Collector.Get().(string)
}

func (c *StringCollector) Preview() string {
	return c.Collector.Preview().(string)
}
