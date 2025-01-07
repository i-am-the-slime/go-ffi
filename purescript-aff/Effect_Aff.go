package purescript_aff

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	. "github.com/purescript-native/go-runtime"
)

var (
	Info = Teal
	Warn = Yellow
	Fata = Red
)

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

// Pure a
type Pure struct {
	Pure  bool
	value Any
}

func (p Pure) String() string {
	return fmt.Sprintf("Pure %s", p.value)
}

// Throw Error
type Throw struct {
	Throw bool
	err   error
}

func (t Throw) String() string {
	return fmt.Sprintf("Throw %s", t.err)
}

// Catch (Aff a) (Error -> Aff a)
type Catch struct {
	Catch      bool
	aff        Any
	errorToAff func(Any) Any
}

func (c Catch) String() string {
	return fmt.Sprintf("Catch %s", c.aff)
}

// Sync (Effect a)
type Sync struct {
	Sync bool
	eff  func() Any
}

func (s Sync) String() string {
	return fmt.Sprintf("Sync")
}

type AsyncCallback = func(Any) func()
type Canceler = Any
type AsyncFn = func(AsyncCallback) func() Canceler

// Async ((Either Error a -> Effect Unit) -> Effect (Canceler))
type Async struct {
	Async   bool
	asyncFn AsyncFn
}

func (a Async) String() string {
	return fmt.Sprintf("Async")
}

// forall b. Bind (Aff b) (b -> Aff a)
type Bind struct {
	Bind   bool
	affOfB Any
	bToAff func(Any) Any
}

func (b Bind) String() string {
	return fmt.Sprintf("Bind %d", b.affOfB)
}

// forall b. Bracket (Aff b) (BracketConditions b) (b -> Aff a)
type Bracket struct {
	Bracket           bool
	acquire           Any
	bracketConditions Dict
	withResource      func(Any) Any
}

// forall b. Fork Boolean (Aff b) ?(Fiber b -> a)
type Fork struct {
	Fork             bool
	questionableBool bool //Uhh
	affOfB           Any
	fiberBToA        func(Any) Any
}

// Sequential (ParAff a)
type Sequential struct {
	Sequential bool
	parAff     Any
}

type Return struct {
	Return bool
}

type Resume struct {
	Resume bool
	b      Cons
}

type Release struct {
	Release           bool
	bracketConditions Dict
	result            Any
}

type Finalizer struct {
	Finalizer bool
	finalizer Any
}

type Finalized struct {
	Finalizer bool
	step      Any
	fail      Any
}

// forall b. Map (b -> a) (ParAff b)
type ParMap struct {
	ParMap    bool
	bToA      func(Any) Any
	parAffOfB Any
}

// forall b. Apply (ParAff (b -> a)) (ParAff b)
type ParApply struct {
	ParApply     bool
	parAffOfBToA Any
	parAffOfB    Any
}

// Alt (ParAff a) (ParAff a)
type ParAlt struct {
	ParAlt  bool
	option1 Any
	option2 Any
}

type OnComplete struct {
	rethrow bool
	handler func(Any /* Either Error a -> Effect Unit */) func() Any
}

type Cons struct {
	head Any
	tail *Cons
}

type InterruptCons struct {
	head      Any
	interrupt bool
	tail      *InterruptCons
}

const SUSPENDED = 0   // Suspended, pending a join.
const CONTINUE = 1    // Interpret the next instruction.
const STEP_BIND = 2   // Apply the next bind.
const STEP_RESULT = 3 // Handle potential failure from a result.
const PENDING = 4     // An async effect is running.
const RETURN = 5      // The current stack has returned.
const COMPLETED = 6   // The entire fiber has completed.

func nonCanceler(error Any) Any {
	return Pure{Pure: true, value: nil}
}

func runSync(left func(Any) Any, right func(Any) Any, eff func() Any) Any {
	fmt.Println(Teal("func: runSync"))
	var res Any
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("\tAn error, error: ", err)
			res = left(err)
		}
	}()
	res = right(eff())
	return res
}

func runAsync(left func(Any) Any, eff AsyncFn, k AsyncCallback) Any {
	fmt.Println(Teal("func: runAsync", k))
	// catch...
	var canceler Canceler = nonCanceler
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("\tAn error, error: ", err)
			k(left(err))()
		}
	}()
	canceler = eff(k)()
	return canceler
}

func sequential(util Any, supervisor Any, par Any) Async {
	var asyncFn AsyncFn = func(cb AsyncCallback) func() Any {
		return func() Any {
			return runPar(util, supervisor, par, cb)
		}
	}
	return Async{Async: true, asyncFn: asyncFn}
}

func runPar(util Any, supervisor Any, par Any, cb AsyncCallback) Any {
	return errors.New("Implement parallel shit")
}

type Scheduler struct {
	isDraining func() bool
	enqueue    func(cb func())
}

var scheduler Scheduler = (func() Scheduler {
	const limit = 1024
	size := 0
	ix := 0
	var queue [limit](func())
	draining := false

	drain := func() {
		fmt.Println("Draining...", size)
		draining = true
		for size != 0 {
			size--
			thunk := queue[ix]
			ix = (ix + 1) % limit
			thunk()
		}
		draining = false
	}

	isDraining := func() bool { return draining }
	enqueue := func(cb func()) {
		fmt.Println("Enqueuing...")
		if size == limit {
			tmp := draining
			drain()
			draining = tmp
		}
		queue[(ix+size)%limit] = cb
		size++

		if !draining {
			drain()
		}
	}

	return Scheduler{isDraining: isDraining, enqueue: enqueue}
})()

func Fiber(util_ Any, supervisor Any, aff_ Any) Dict {
	var util map[string]Any = util_.(map[string]Any)
	var isLeft func(Any) Any = util["isLeft"].(func(Any) Any)
	var left func(Any) Any = util["left"].(func(Any) Any)
	var right func(Any) Any = util["right"].(func(Any) Any)
	var fromRight func(Any) Any = util["fromRight"].(func(Any) Any)
	var fromLeft func(Any) Any = util["fromLeft"].(func(Any) Any)
	runTick := 0

	step := aff_       // Successful step
	var fail Any       // Failure step
	interrupt := false // Asynchronouse interrupt

	// Stack of continuations for the current fiber.
	var b Cons = Cons{}

	// Stack of attempts and finalizers for error recovery. Every `Cons` is also
	// tagged with current `interrupt` state. We use this to track which items
	// should be ignored or evaluated as a result of a kill.
	var attempts *InterruptCons = &InterruptCons{}

	// A special state is needed for Bracket, because it cannot be killed. When
	// we enter a bracket acquisition or finalizer, we increment the counter,
	// and then decrement once complete.
	bracketCount := 0

	// Each join gets a new id so they can be revoked.
	// joinId := 0
	joins := map[int]OnComplete{}
	var rethrow = true

	status := SUSPENDED

	var run func(int) Any
	run = func(localRunTick int) Any {
		fmt.Println("New round", "localRunTick", localRunTick, "bracketCount", bracketCount)
		fmt.Println("Step", step)
		var tmp Any
		var result Any
		var attempt Any
		for {
			tmp = nil
			result = nil
			attempt = nil
			switch status {
			case STEP_BIND:
				fmt.Println("STEP_BIND", b.head)
				status = CONTINUE // next step
				headFn := b.head.(func(Any) Any)
				// Makeshift try catch block
				func() {
					defer func() {
						if err := recover(); err != nil {
							fmt.Println(Red("\tError oh shit ", err))
							// early return on error
							status = RETURN
							fail = left(err)
							step = nil
						}
					}()
					step = headFn(step)
					if b.tail == nil {
						b.head = nil
					} else {
						b.head = b.tail.head
						b.tail = b.tail.tail
					}
				}()
			case STEP_RESULT:
				fmt.Println("STEP_RESULT")
				if isLeft(step).(bool) {
					// early return on error
					status = RETURN
					fail = step
					step = nil
				} else if b.head == nil {
					// happy case done
					status = RETURN
					fmt.Println("My work here is done")
				} else {
					// happy case work left
					status = STEP_BIND
					step = fromRight(step)
					fmt.Println("Next step after happy result", step)
				}

			case CONTINUE:
				fmt.Println("CONTINUE")
				switch currentStep := step.(type) {

				case Bind:
					fmt.Println("\tBind")
					if b.head != nil {
						b.tail = &Cons{head: b.head, tail: b.tail}
					}
					b.head = currentStep.bToAff
					status = CONTINUE
					step = currentStep.affOfB

				case Pure:
					fmt.Println("\tPure")
					if b.head == nil {

						fmt.Println("\t> Head nil")
						// we're done
						status = RETURN
						step = right(currentStep.value)
					} else {
						fmt.Println("\t> Head exists", currentStep.value)
						// this happens after a bind
						status = STEP_BIND
						step = currentStep.value
					}

				case Sync:
					fmt.Println("\tSync")
					status = STEP_RESULT
					step = runSync(left, right, currentStep.eff)

				case Async:
					fmt.Println("\tAsync")
					status = PENDING
					step = runAsync(left, currentStep.asyncFn, func(theResult Any) func() {
						return func() {
							if runTick != localRunTick {
								fmt.Println("\t\tRun tick != localRunTick")
								return
							}
							runTick++
							scheduler.enqueue(func() {
								if runTick != localRunTick+1 {
									fmt.Println("\t\tRun tick != localRunTick + 1")
									return
								}
								status = STEP_RESULT
								step = theResult
								run(runTick)
							})
						}
					})
					fmt.Println("\t\tNext step after Async", step)
					return nil

				case Throw:
					fmt.Println("\tThrow")
					status = RETURN
					fail = left(currentStep.err)
					step = nil

				case Catch:
					fmt.Println("\tCatch", currentStep.aff)
					if b.head == nil {
						fmt.Println("\t\tHead is nil")
						attempts = &InterruptCons{interrupt: interrupt, head: step, tail: attempts}
					} else {
						fmt.Println("\t\tHead is not nil", b.head)
						attempts = &InterruptCons{
							interrupt: interrupt,
							head:      step,
							tail: &InterruptCons{
								interrupt: interrupt,
								head:      Resume{b: b},
								tail:      attempts}}
					}
					b.head = nil
					b.tail = nil
					status = CONTINUE
					step = currentStep.aff
				case Bracket:
					fmt.Println("\tBracket")
					bracketCount++
					if b.head == nil {
						attempts = &InterruptCons{head: step, tail: attempts, interrupt: interrupt}
					} else {
						attempts = &InterruptCons{head: step, tail: &InterruptCons{head: Resume{b: b}, tail: attempts, interrupt: interrupt}, interrupt: interrupt}
					}
					b.head = nil
					b.tail = nil
					status = CONTINUE
					step = currentStep.acquire

				case Fork:
					fmt.Println("\tFork")
					status = STEP_RESULT
					tmp = Fiber(util, supervisor, currentStep.affOfB)
					if supervisor != nil {
						supervisor.(Dict)["register"].(func(Any))(tmp)
					}
					if currentStep.questionableBool {
						tmp.(Dict)["run"].(func())()
					}
					step = right(tmp)
				case Sequential:
					fmt.Println("\tSequential")
					status = CONTINUE
					step = sequential(util, supervisor, currentStep.parAff)

				}

			case RETURN:
				fmt.Println("RETURN")
				b.head = nil
				b.tail = nil
				if attempts == nil || attempts.head == nil {
					fmt.Println("\tAttempts are nil", interrupt, fail)
					status = COMPLETED
					if interrupt {
						step = interrupt
					} else if fail != nil {
						step = fail
					}
				} else {
					fmt.Println("\tAttempts aren't nil", attempts.head, attempts.tail)
					tmp = attempts.interrupt
					attempt = attempts.head
					attempts = attempts.tail

					switch currentAttempt := attempt.(type) {
					case Catch:
						fmt.Println("\tReturn Catch")
						if interrupt && interrupt != tmp && bracketCount == 0 {
							fmt.Println("\t\tGonna RETURN")
							status = RETURN
						} else if fail != nil {
							fmt.Println("\t\tGonna CONTINUE", fromLeft(fail), currentAttempt.errorToAff(fromLeft(fail)))
							status = CONTINUE
							step = currentAttempt.errorToAff(fromLeft(fail))
							fail = nil
						}
						// We cannot resume from an unmasked interrupt or exception.
					case Resume:
						fmt.Println("\tResume")
						if interrupt && interrupt != tmp && bracketCount == 0 || fail != nil {
							status = RETURN
						} else {
							b.head = currentAttempt.b.head
							b.tail = currentAttempt.b.tail
							status = STEP_BIND
							step = fromRight(step)
							fmt.Println("\t\tNext step", step)
						}

					case Bracket:
						fmt.Println("\tBracket")
						bracketCount--
						if fail == nil {
							result = fromRight(step)
							attempts = &InterruptCons{
								head: Release{
									bracketConditions: currentAttempt.bracketConditions,
									result:            result,
								},
								tail:      attempts,
								interrupt: tmp.(bool),
							}

							if interrupt == tmp || bracketCount > 0 {
								status = CONTINUE
								step = currentAttempt.withResource(result)
							}
						}
					case Release:
						fmt.Println("\tRelease")
						attempts = &InterruptCons{
							head:      Finalized{step: step, fail: fail},
							tail:      attempts,
							interrupt: interrupt,
						}
						status = CONTINUE
						// It has only been killed if the interrupt status has changed
						// since we enqueued the item, and the bracket count is 0. If the
						// bracket count is non-zero then we are in a masked state so it's
						// impossible to be killed.
						if interrupt && interrupt != tmp && bracketCount == 0 {
							step = currentAttempt.bracketConditions["killed"].(func(Any) Any)(fromLeft(interrupt)).(func(Any) Any)(currentAttempt.result)
						} else if fail != nil {
							step = currentAttempt.bracketConditions["failed"].(func(Any) Any)(fromLeft(fail)).(func(Any) Any)(currentAttempt.result)
						} else {
							step = currentAttempt.bracketConditions["completed"].(func(Any) Any)(fromRight(step)).(func(Any) Any)(currentAttempt.result)
						}
						fail = nil
						bracketCount++
					case Finalizer:
						fmt.Println("\tFinalizer")
						bracketCount++
						attempts.tail = attempts
						attempts.head = Finalized{step: step, fail: fail}
						attempts.interrupt = interrupt
						status = CONTINUE
						step = currentAttempt.finalizer
					case Finalized:
						fmt.Println("\tFinalized")
						bracketCount--
						status = RETURN
						step = currentAttempt.step
						fail = currentAttempt.fail
					default:
						fmt.Println(Fata("\tUnknown Attempt type hit", reflect.TypeOf(currentAttempt)))
					}
				}
			case COMPLETED:
				fmt.Println("COMPLETED", joins)
				for _, join := range joins {
					if rethrow {
						rethrow = join.rethrow
						// Run eff
						join.handler(step)()
					}
				}
				joins = nil
				if interrupt && fail != nil {
					go panic(fromLeft(fail))
				} else if isLeft(step).(bool) && rethrow {
					if rethrow {
						go panic(fromLeft(fail))
					}
				}
				return nil

			case SUSPENDED:
				fmt.Println("SUSPENDED")
				status = CONTINUE
			case PENDING:
				fmt.Println("PENDING")
				return nil

			default:
				fmt.Println(Fata("Unknown branch hit", step))
			}
		}

	}
	var realRun func() Any
	realRun = func() Any {
		if status == SUSPENDED {
			if !scheduler.isDraining() {
				scheduler.enqueue(func() {
					run(runTick)
				})
			} else {
				run(runTick)
			}
		}
		return nil // Important to avoid panic
	}
	return Dict{
		"run":        realRun,
		"kill":       func() { fmt.Println("[TODO]: kill") },
		"join":       func() { fmt.Println("[TODO]: join") },
		"onComplete": func() { fmt.Println("[TODO]: onComplete") },
		"onSuspend":  func() { fmt.Println("[TODO]: onSuspend") },
	}
}

func init() {
	exports := Foreign("Effect.Aff")

	// ∷ ∀	 a b. Aff a → (a → Aff b) → Aff b
	exports["_pure"] = func(value Any) Any {
		return Pure{Pure: true, value: value}
	}
	// ∷ ∀ a. Error → Aff a
	exports["_throwError"] = func(e_ Any) Any {
		e := e_.(error)
		return Throw{Throw: true, err: e}
	}
	// ∷ ∀ a. Aff a → (Error → Aff a) → Aff a
	exports["_catchError"] = func(aff Any) Any {
		return func(k Any) Any {
			return Catch{aff: aff, errorToAff: k.(func(Any) Any)}
		}
	}
	// ∷ ∀ a b. (a → b) → Aff a → Aff b
	exports["_map"] = func(f_ Any) Any {
		fmt.Println(Yellow("map"))
		f := f_.(func(Any) Any)
		return func(aff Any) Any {
			switch a := aff.(type) {
			case Pure:
				return Pure{Pure: true, value: f(a.value)}
			default:
				return Bind{
					Bind:   true,
					affOfB: aff,
					bToAff: func(b Any) Any { return Pure{Pure: true, value: f(b)} }}
			}
		}
	}

	// ∷ ∀ a b. Aff a → (a → Aff b) → Aff b
	exports["_bind"] = func(aff Any) Any {
		fmt.Println(Yellow("bind"))
		return func(f_ Any) Any {
			f := f_.(func(Any) Any)
			bToAff := func(b Any) Any { return f(b) }
			return Bind{Bind: true, affOfB: aff, bToAff: bToAff}
		}
	}

	// type Util = interface {
	// 	isLeft(Any) bool
	// 	fromLeft(Any) Any
	// 	fromRight(Any) Any
	// 	left(Any) Any
	// 	right(Any) Any
	// }

	// ∷ ∀ a. Fn.Fn2 FFIUtil (Aff a) (Effect (Fiber a))
	exports["_makeFiber"] = func(util Any, aff Any) Any {
		return func() Any {
			return Fiber(util, nil, aff)
		}
	}

	// ∀ a. Effect a → Aff a
	exports["_liftEffect"] = func(effect_ Any) Any {
		effect := func() Any { return Run(effect_) }
		return Sync{Sync: true, eff: effect}
	}

	// ∀ a. Fn.Fn2 (Unit → Either a Unit) Number (Aff Unit)
	exports["_delay"] = func(right_ Any, millis_ Any) Any {
		var asyncFn AsyncFn = func(asyncCb AsyncCallback) func() Canceler {
			return func() Canceler {
				// [TODO] Cancel like this: https://stackoverflow.com/questions/55135239/how-can-i-sleep-with-responsive-context-cancelation
				right := right_.(func(Any) Any)
				var millis int = int(millis_.(float64))
				duration := time.Duration(millis) * time.Millisecond
				sleepContext := func(ctx context.Context, delay time.Duration) {
					select {
					case <-ctx.Done():
					case <-time.After(delay):
					}
				}
				ctx := context.Background()
				ctxTimeout, cancel := context.WithCancel(ctx)
				sleepContext(ctxTimeout, duration)
				asyncCb(right(nil))()
				return func() Any {
					return Sync{eff: func() Any { cancel(); return right(nil) }}
				}
			}
		}
		return Async{asyncFn: asyncFn}
	}

	exports["generalBracket"] = func(acquire Any) Any {
		return func(options Any) Any {
			bracketConditions := options.(Dict)
			return func(k Any) Any {
				withResource := k.(func(Any) Any)
				return Bracket{
					acquire:           acquire,
					bracketConditions: bracketConditions,
					withResource:      withResource,
				}
			}
		}
	}

	// ParAff ~> Aff
	exports["_sequential"] = func(oldAff Any) Any {
		return oldAff
	}

}
