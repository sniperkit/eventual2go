package eventual2go

import (
	"sync"
	"fmt"
)

// Creates a new future.
func NewFuture() (F *Future) {
	F = new(Future)
	F.m = new(sync.Mutex)
	F.fcs = []futurecompleter{}
	F.fces = []futurecompletererror{}
	return
}

// Future is thread-safe struct that can be completed with arbitrary data or failed with an error. Handler functions can
// be registered for both events and get invoked after completion..
type Future struct {
	m *sync.Mutex
	fcs []futurecompleter
	fces []futurecompletererror
	c bool
	r interface {}
	e error
}

// Completes the future with the given data and triggers al registered completion handlers. Panics if the future is already
// complete.
func (f *Future) Complete(d Data){
	f.m.Lock()
	defer f.m.Unlock()
	if f.c {
		panic(fmt.Sprint("Completed complete future with",d))
	}
	f.r = d
	for _, fc := range f.fcs {
		go deliverData(fc,d)
	}
	f.c = true
}

// Returns the completion state.
func (f *Future) IsComplete() bool {
	return f.c
}
// Completes the future with the given error and triggers al registered error handlers. Panics if the future is already
// complete.
func (f *Future) CompleteError(err error){
	f.m.Lock()
	defer f.m.Unlock()
	if f.c {
		panic(fmt.Sprint("Errorcompleted complete future with",err))
	}
	f.e = err
	for _, fce := range f.fces {
		deliverErr(fce,f.e)
	}
	f.c = true
}

func deliverData(fc futurecompleter, d interface {}){
		fc.f.Complete(fc.cf(d))
}

func deliverErr(fce futurecompletererror, e error){
	go func(){
		d, err := fce.ef(e)
		if err == nil {
			fce.f.Complete(d)
		} else {
			fce.f.CompleteError(err)
		}
	}()
}

// Then registers a completion handler. If the future is already complete, the handler gets executed immediately.
// Returns a future that gets completed with result of the handler.
func (f *Future) Then(ch CompletionHandler) (nf *Future) {
	f.m.Lock()
	defer f.m.Unlock()

	nf = NewFuture()
	fc := futurecompleter{ch,nf}
	if f.c && f.e == nil {
		deliverData(fc,f.r)
	} else if !f.c  {
		f.fcs = append(f.fcs,fc)
	}
	return
}

// Blocks until the future is complete.
func (f *Future) WaitUntilComplete() {
	<-f.AsChan()
}

// Returns the result of the future, nil called before completion or after error completion.
func (f *Future) GetResult() Data{
	return f.r
}

// Then registers an error handler. If the future is already completed with an error, the handler gets executed
// immediately.
// Returns a future that either gets completed with result of the handler or error completed with the error from handler,
// if not nil.
func (f *Future) Err(eh ErrorHandler) (nf *Future) {
	f.m.Lock()
	defer f.m.Unlock()

	nf = NewFuture()
	fce := futurecompletererror{eh, nf}
	if f.e != nil {
		deliverErr(fce, f.e)
	} else if !f.c{
		f.fces = append(f.fces, fce)
	}
	return
}

// AsChan returns a channel which will receive either the result or the error after completion of the future.
func (f *Future) AsChan() chan Data{
	c := make(chan Data,1)
	cmpl := func(d chan Data) CompletionHandler{
		return func(e Data)Data{
			d<-e
			close(d)
			return nil
		}
	}
	ecmpl :=func(d chan Data) ErrorHandler{
		return func(e error)(Data, error){
			d<-e
			close(d)
			return nil, nil
		}
	}
	f.Then(cmpl(c))
	f.Err(ecmpl(c))
	return c
}


type futurecompleter struct {
	cf CompletionHandler
	f *Future
}

type futurecompletererror struct {
	ef ErrorHandler
	f *Future
}
