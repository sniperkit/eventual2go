
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

type StringSliceSuscriber func([]string)

func (l StringSliceSuscriber) toSuscriber() eventual2go.Subscriber {
	return func(d eventual2go.Data) { l(d.([]string)) }
}

func (s *StringSliceStream) Listen(ss StringSliceSuscriber) *eventual2go.Subscription {
	return s.Stream.Listen(ss.toSuscriber())
}

type StringSliceFilter func([]string) bool

func (f StringSliceFilter) toFilter() eventual2go.Filter {
	return func(d eventual2go.Data) bool { return f(d.([]string)) }
}

func (s *StringSliceStream) Where(f StringSliceFilter) *StringSliceStream {
	return &StringSliceStream{s.Stream.Where(f.toFilter())}
}

func (s *StringSliceStream) WhereNot(f StringSliceFilter) *StringSliceStream {
	return &StringSliceStream{s.Stream.WhereNot(f.toFilter())}
}

func (s *StringSliceStream) First() *StringSliceFuture {
	return &StringSliceFuture{s.Stream.First()}
}

func (s *StringSliceStream) FirstWhere(f StringSliceFilter) *StringSliceFuture {
	return &StringSliceFuture{s.Stream.FirstWhere(f.toFilter())}
}

func (s *StringSliceStream) FirstWhereNot(f StringSliceFilter) *StringSliceFuture {
	return &StringSliceFuture{s.Stream.FirstWhereNot(f.toFilter())}
}

func (s *StringSliceStream) AsChan() (c chan []string) {
	c = make(chan []string)
	s.Listen(pipeToStringSliceChan(c)).Closed().Then(closeStringSliceChan(c))
	return
}

func pipeToStringSliceChan(c chan []string) StringSliceSuscriber {
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
