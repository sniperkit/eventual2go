
/*
 * generated by event_generator
 *
 * DO NOT EDIT
 */

package typed_events

import "github.com/joernweissenborn/eventual2go"



type ErrorCompleter struct {
	*eventual2go.Completer
}

func NewErrorCompleter() *ErrorCompleter {
	return &ErrorCompleter{eventual2go.NewCompleter()}
}

func (c *ErrorCompleter) Complete(d error) {
	c.Completer.Complete(d)
}

func (c *ErrorCompleter) Future() *ErrorFuture {
	return &ErrorFuture{c.Completer.Future()}
}

type ErrorFuture struct {
	*eventual2go.Future
}

func (f *ErrorFuture) GetResult() error {
	return f.Future.GetResult().(error)
}

type ErrorCompletionHandler func(error) error

func (ch ErrorCompletionHandler) toCompletionHandler() eventual2go.CompletionHandler {
	return func(d eventual2go.Data) eventual2go.Data {
		return ch(d.(error))
	}
}

func (f *ErrorFuture) Then(ch ErrorCompletionHandler) *ErrorFuture {
	return &ErrorFuture{f.Future.Then(ch.toCompletionHandler())}
}

func (f *ErrorFuture) AsChan() chan error {
	c := make(chan error, 1)
	cmpl := func(d chan error) ErrorCompletionHandler {
		return func(e error) error {
			d <- e
			close(d)
			return e
		}
	}
	ecmpl := func(d chan error) eventual2go.ErrorHandler {
		return func(error) (eventual2go.Data, error) {
			close(d)
			return nil, nil
		}
	}
	f.Then(cmpl(c))
	f.Err(ecmpl(c))
	return c
}

type ErrorStreamController struct {
	*eventual2go.StreamController
}

func NewErrorStreamController() *ErrorStreamController {
	return &ErrorStreamController{eventual2go.NewStreamController()}
}

func (sc *ErrorStreamController) Add(d error) {
	sc.StreamController.Add(d)
}

func (sc *ErrorStreamController) Join(s *ErrorStream) {
	sc.StreamController.Join(s.Stream)
}

func (sc *ErrorStreamController) JoinFuture(f *ErrorFuture) {
	sc.StreamController.JoinFuture(f.Future)
}

func (sc *ErrorStreamController) Stream() *ErrorStream {
	return &ErrorStream{sc.StreamController.Stream()}
}

type ErrorStream struct {
	*eventual2go.Stream
}

type ErrorSubscriber func(error)

func (l ErrorSubscriber) toSubscriber() eventual2go.Subscriber {
	return func(d eventual2go.Data) { l(d.(error)) }
}

func (s *ErrorStream) Listen(ss ErrorSubscriber) *eventual2go.Subscription {
	return s.Stream.Listen(ss.toSubscriber())
}

type ErrorFilter func(error) bool

func (f ErrorFilter) toFilter() eventual2go.Filter {
	return func(d eventual2go.Data) bool { return f(d.(error)) }
}

func toErrorFilterArray(f ...ErrorFilter) (filter []eventual2go.Filter){

	filter = make([]eventual2go.Filter, len(f))
	for i, el := range f {
		filter[i] = el.toFilter()
	}
	return
}

func (s *ErrorStream) Where(f ...ErrorFilter) *ErrorStream {
	return &ErrorStream{s.Stream.Where(toErrorFilterArray(f...)...)}
}

func (s *ErrorStream) WhereNot(f ...ErrorFilter) *ErrorStream {
	return &ErrorStream{s.Stream.WhereNot(toErrorFilterArray(f...)...)}
}

func (s *ErrorStream) Split(f ErrorFilter) (*ErrorStream, *ErrorStream)  {
	return s.Where(f), s.WhereNot(f)
}

func (s *ErrorStream) First() *ErrorFuture {
	return &ErrorFuture{s.Stream.First()}
}

func (s *ErrorStream) FirstWhere(f... ErrorFilter) *ErrorFuture {
	return &ErrorFuture{s.Stream.FirstWhere(toErrorFilterArray(f...)...)}
}

func (s *ErrorStream) FirstWhereNot(f ...ErrorFilter) *ErrorFuture {
	return &ErrorFuture{s.Stream.FirstWhereNot(toErrorFilterArray(f...)...)}
}

func (s *ErrorStream) AsChan() (c chan error) {
	c = make(chan error)
	s.Listen(pipeToErrorChan(c)).Closed().Then(closeErrorChan(c))
	return
}

func pipeToErrorChan(c chan error) ErrorSubscriber {
	return func(d error) {
		c <- d
	}
}

func closeErrorChan(c chan error) eventual2go.CompletionHandler {
	return func(d eventual2go.Data) eventual2go.Data {
		close(c)
		return nil
	}
}

type ErrorCollector struct {
	*eventual2go.Collector
}

func NewErrorCollector() *ErrorCollector {
	return &ErrorCollector{eventual2go.NewCollector()}
}

func (c *ErrorCollector) Add(d error) {
	c.Collector.Add(d)
}

func (c *ErrorCollector) AddFuture(f *ErrorFuture) {
	c.Collector.Add(f.Future)
}

func (c *ErrorCollector) AddStream(s *ErrorStream) {
	c.Collector.AddStream(s.Stream)
}

func (c *ErrorCollector) Get() error {
	return c.Collector.Get().(error)
}

func (c *ErrorCollector) Preview() error {
	return c.Collector.Preview().(error)
}
