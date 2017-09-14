
/*
 * generated by event_generator
 *
 * DO NOT EDIT
 */

package typed_events

import "github.com/joernweissenborn/eventual2go"



type StringSliceCompleter struct {
	*eventual2go.Completer
}

func NewStringSliceCompleter() *StringSliceCompleter {
	return &StringSliceCompleter{eventual2go.NewCompleter()}
}

func (c *StringSliceCompleter) Complete(d []string) {
	c.Completer.Complete(d)
}

func (c *StringSliceCompleter) Future() *StringSliceFuture {
	return &StringSliceFuture{c.Completer.Future()}
}

type StringSliceFuture struct {
	*eventual2go.Future
}

func (f *StringSliceFuture) Result() []string {
	return f.Future.Result().([]string)
}

type StringSliceCompletionHandler func([]string) []string

func (ch StringSliceCompletionHandler) toCompletionHandler() eventual2go.CompletionHandler {
	return func(d eventual2go.Data) eventual2go.Data {
		return ch(d.([]string))
	}
}

func (f *StringSliceFuture) Then(ch StringSliceCompletionHandler) *StringSliceFuture {
	return &StringSliceFuture{f.Future.Then(ch.toCompletionHandler())}
}

func (f *StringSliceFuture) AsChan() chan []string {
	c := make(chan []string, 1)
	cmpl := func(d chan []string) StringSliceCompletionHandler {
		return func(e []string) []string {
			d <- e
			close(d)
			return e
		}
	}
	ecmpl := func(d chan []string) eventual2go.ErrorHandler {
		return func(error) (eventual2go.Data, error) {
			close(d)
			return nil, nil
		}
	}
	f.Then(cmpl(c))
	f.Err(ecmpl(c))
	return c
}

type StringSliceStreamController struct {
	*eventual2go.StreamController
}

func NewStringSliceStreamController() *StringSliceStreamController {
	return &StringSliceStreamController{eventual2go.NewStreamController()}
}

func (sc *StringSliceStreamController) Add(d []string) {
	sc.StreamController.Add(d)
}

func (sc *StringSliceStreamController) Join(s *StringSliceStream) {
	sc.StreamController.Join(s.Stream)
}

func (sc *StringSliceStreamController) JoinFuture(f *StringSliceFuture) {
	sc.StreamController.JoinFuture(f.Future)
}

func (sc *StringSliceStreamController) Stream() *StringSliceStream {
	return &StringSliceStream{sc.StreamController.Stream()}
}

type StringSliceStream struct {
	*eventual2go.Stream
}

type StringSliceSubscriber func([]string)

func (l StringSliceSubscriber) toSubscriber() eventual2go.Subscriber {
	return func(d eventual2go.Data) { l(d.([]string)) }
}

func (s *StringSliceStream) Listen(ss StringSliceSubscriber) *eventual2go.Completer {
	return s.Stream.Listen(ss.toSubscriber())
}

func (s *StringSliceStream) ListenNonBlocking(ss StringSliceSubscriber) *eventual2go.Completer {
	return s.Stream.ListenNonBlocking(ss.toSubscriber())
}

type StringSliceFilter func([]string) bool

func (f StringSliceFilter) toFilter() eventual2go.Filter {
	return func(d eventual2go.Data) bool { return f(d.([]string)) }
}

func toStringSliceFilterArray(f ...StringSliceFilter) (filter []eventual2go.Filter){

	filter = make([]eventual2go.Filter, len(f))
	for i, el := range f {
		filter[i] = el.toFilter()
	}
	return
}

func (s *StringSliceStream) Where(f ...StringSliceFilter) *StringSliceStream {
	return &StringSliceStream{s.Stream.Where(toStringSliceFilterArray(f...)...)}
}

func (s *StringSliceStream) WhereNot(f ...StringSliceFilter) *StringSliceStream {
	return &StringSliceStream{s.Stream.WhereNot(toStringSliceFilterArray(f...)...)}
}

func (s *StringSliceStream) TransformWhere(t eventual2go.Transformer, f ...StringSliceFilter) *eventual2go.Stream {
	return s.Stream.TransformWhere(t, toStringSliceFilterArray(f...)...)
}

func (s *StringSliceStream) Split(f StringSliceFilter) (*StringSliceStream, *StringSliceStream)  {
	return s.Where(f), s.WhereNot(f)
}

func (s *StringSliceStream) First() *StringSliceFuture {
	return &StringSliceFuture{s.Stream.First()}
}

func (s *StringSliceStream) FirstWhere(f... StringSliceFilter) *StringSliceFuture {
	return &StringSliceFuture{s.Stream.FirstWhere(toStringSliceFilterArray(f...)...)}
}

func (s *StringSliceStream) FirstWhereNot(f ...StringSliceFilter) *StringSliceFuture {
	return &StringSliceFuture{s.Stream.FirstWhereNot(toStringSliceFilterArray(f...)...)}
}

func (s *StringSliceStream) AsChan() (c chan []string, stop *eventual2go.Completer) {
	c = make(chan []string)
	stop = s.Listen(pipeToStringSliceChan(c))
	stop.Future().Then(closeStringSliceChan(c))
	return
}

func pipeToStringSliceChan(c chan []string) StringSliceSubscriber {
	return func(d []string) {
		c <- d
	}
}

func closeStringSliceChan(c chan []string) eventual2go.CompletionHandler {
	return func(d eventual2go.Data) eventual2go.Data {
		close(c)
		return nil
	}
}

type StringSliceCollector struct {
	*eventual2go.Collector
}

func NewStringSliceCollector() *StringSliceCollector {
	return &StringSliceCollector{eventual2go.NewCollector()}
}

func (c *StringSliceCollector) Add(d []string) {
	c.Collector.Add(d)
}

func (c *StringSliceCollector) AddFuture(f *StringSliceFuture) {
	c.Collector.Add(f.Future)
}

func (c *StringSliceCollector) AddStream(s *StringSliceStream) {
	c.Collector.AddStream(s.Stream)
}

func (c *StringSliceCollector) Get() []string {
	return c.Collector.Get().([]string)
}

func (c *StringSliceCollector) Preview() []string {
	return c.Collector.Preview().([]string)
}
