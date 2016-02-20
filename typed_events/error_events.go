
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

func (c *ErrorCompleter) Complete(d error) {
	c.Completer.Complete(d)
}

func (c *ErrorCompleter) Future() *ErrorFuture {
	return &ErrorFuture{c.Completer.Future()}
}

type ErrorFuture struct {
	*eventual2go.Future
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

type ErrorStream struct {
	*eventual2go.Stream
}

type ErrorSuscriber func(error)

func (l ErrorSuscriber) toSuscriber() eventual2go.Subscriber {
	return func(d eventual2go.Data) { l(d.(error)) }
}

func (s *ErrorStream) Listen(ss ErrorSuscriber) *eventual2go.Subscription{
	return s.Stream.Listen(ss.toSuscriber())
}

type ErrorFilter func(error) bool

func (f ErrorFilter) toFilter() eventual2go.Filter {
	return func(d eventual2go.Data) bool { return f(d.(error)) }
}

func (s *ErrorStream) Where(f ErrorFilter) {
	s.Stream.Where(f.toFilter())
}

func (s *ErrorStream) WhereNot(f ErrorFilter) {
	s.Stream.WhereNot(f.toFilter())
}

func (s *ErrorStream) First() *ErrorFuture {
	return &ErrorFuture{s.Stream.First()}
}

func (s *ErrorStream) FirstWhere(f ErrorFilter) *ErrorFuture {
	return &ErrorFuture{s.Stream.FirstWhere(f.toFilter())}
}

func (s *ErrorStream) FirstWhereNot(f ErrorFilter) *ErrorFuture {
	return &ErrorFuture{s.Stream.FirstWhereNot(f.toFilter())}
}

func (s *ErrorStream) AsChan() (c chan error) {
	c = make(chan error)
	s.Listen(pipeToErrorChan(c)).Closed().Then(closeErrorChan(c))
	return
}

func pipeToErrorChan(c chan error) ErrorSuscriber {
	return func(d error) {
		c<-d
	}
}

func closeErrorChan(c chan error) eventual2go.CompletionHandler {
	return func(d eventual2go.Data) eventual2go.Data {
		close(c)
		return nil
	}
}

