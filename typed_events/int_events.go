
/*
 * generated by event_generator
 *
 * DO NOT EDIT
 */

package typed_events

import "github.com/joernweissenborn/eventual2go"



type IntCompleter struct {
	*eventual2go.Completer
}

func (c *IntCompleter) Complete(d int) {
	c.Completer.Complete(d)
}

func (c *IntCompleter) Future() *IntFuture {
	return &IntFuture{c.Completer.Future()}
}

type IntFuture struct {
	*eventual2go.Future
}

type IntCompletionHandler func(int) int

func (ch IntCompletionHandler) toCompletionHandler() eventual2go.CompletionHandler {
	return func(d eventual2go.Data) eventual2go.Data {
		return ch(d.(int))
	}
}

func (f *IntFuture) Then(ch IntCompletionHandler) *IntFuture {
	return &IntFuture{f.Future.Then(ch.toCompletionHandler())}
}

type IntStream struct {
	*eventual2go.Stream
}

type IntSuscriber func(int)

func (l IntSuscriber) toSuscriber() eventual2go.Subscriber {
	return func(d eventual2go.Data) { l(d.(int)) }
}

func (s *IntStream) Listen(ss IntSuscriber) *eventual2go.Subscription{
	return s.Stream.Listen(ss.toSuscriber())
}

type IntFilter func(int) bool

func (f IntFilter) toFilter() eventual2go.Filter {
	return func(d eventual2go.Data) bool { return f(d.(int)) }
}

func (s *IntStream) Where(f IntFilter) {
	s.Stream.Where(f.toFilter())
}

func (s *IntStream) WhereNot(f IntFilter) {
	s.Stream.WhereNot(f.toFilter())
}

func (s *IntStream) First() *IntFuture {
	return &IntFuture{s.Stream.First()}
}

func (s *IntStream) FirstWhere(f IntFilter) *IntFuture {
	return &IntFuture{s.Stream.FirstWhere(f.toFilter())}
}

func (s *IntStream) FirstWhereNot(f IntFilter) *IntFuture {
	return &IntFuture{s.Stream.FirstWhereNot(f.toFilter())}
}

func (s *IntStream) AsChan() (c chan int) {
	c = make(chan int)
	s.Listen(pipeToIntChan(c)).Closed().Then(closeIntChan(c))
	return
}

func pipeToIntChan(c chan int) IntSuscriber {
	return func(d int) {
		c<-d
	}
}

func closeIntChan(c chan int) eventual2go.CompletionHandler {
	return func(d eventual2go.Data) eventual2go.Data {
		close(c)
		return nil
	}
}

