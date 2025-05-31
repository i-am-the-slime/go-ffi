package purescript_aff

import (
	"fmt"
	"time"

	. "github.com/purescript-native/go-runtime"
)

// Pure a
type Pure struct {
	value Any
}

// Throw Error
type Throw struct {
	err error
}

// Catch (Aff a) (Error -> Aff a)
type Catch struct {
	aff        Any
	errorToAff func(Any) Any
}

// Sync (Effect a)
type Sync struct {
	eff func() Any
}

type AsyncCallback = Any // func(Any) func()
type Canceler = Any      // func(Any) Any // Error -> Aff Unit
type AsyncFn = Any       // func(Any /* really AsyncCallback */) Any /* really func() Canceler */

// Async ((Either Error a -> Effect Unit) -> Effect (Canceler))
type Async struct {
	asyncFn AsyncFn
}

// forall b. Bind (Aff b) (b -> Aff a)
type Bind struct {
	affOfB Any
	bToAff func(Any) Any
}

// forall b. Bracket (Aff b) (BracketConditions b) (b -> Aff a)
type Bracket struct {
	acquire           Any
	bracketConditions Dict
	withResource      func(Any) Any
}

// forall b. Fork Boolean (Aff b) ?(Fiber b -> a)
type Fork struct {
	questionableBool bool //Uhh
	affOfB           Any
	fiberBToA        func(Any) Any
}

// Sequential (ParAff a)
type Sequential struct {
	parAff Any
}

type Return struct{}

type Resume struct {
	b Cons
}

type Release struct {
	bracketConditions Dict
	result            Any
}

type Finalizer struct {
	finalizer Any
}

type Finalized struct {
	step Any
	fail Any
}

// forall b. Map (b -> a) (ParAff b)
type ParMap struct {
	bToA      func(Any) Any
	parAffOfB Any
}

// forall b. Apply (ParAff (b -> a)) (ParAff b)
type ParApply struct {
	parAffOfBToA Any
	parAffOfB    Any
}

// Alt (ParAff a) (ParAff a)
type ParAlt struct {
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
	interrupt Any
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
	return Pure{value: Dict{}}
}

func runSync(left func(Any) Any, right func(Any) Any, eff func() Any) Any {
	return right(eff())
}

func runAsync(
	left func(Any) Any, //  left :: Error -> Either Error a
	asyncFn Any, //  asyncFn :: (Either Error a -> Effect Unit) -> Effect Canceler
	cb Any, //  cb       :: Either Error a -> Effect Unit
) Any {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[runAsync] Caught panic:", err)
			msg := fmt.Sprintf("Error in runAsync: %v", err)
			// Optionally call cb(left(err)) here:
			if cbFunc, ok := cb.(func(Any) func() Any); ok {
				cbFunc(left(err))()
			}
			panic(msg)
		}
	}()

	fmt.Printf("[runAsync] asyncFn=%T, cb=%T\n", asyncFn, cb)

	// 1) Partial application #1
	fmt.Println("[runAsync] -- Attempting Apply(asyncFn, cb)...")
	part := debugApply(asyncFn, cb)
	if part == nil {
		panic("[runAsync] Apply(asyncFn, cb) returned nil!")
	}
	fmt.Printf("[runAsync]    => part result: %T %+v\n", part, part)

	// 2) Now we expect part is an Effect Canceler => a zero-arg function
	fmt.Println("[runAsync] -- Type-asserting part as func() Any")
	effFunc, ok := part.(func() Any)
	if !ok {
		panic("[runAsync] part is not func() Any (Effect Canceler).")
	}

	// 3) Call it
	fmt.Println("[runAsync] -- Calling the zero-arg function.")
	result := effFunc()
	fmt.Printf("[runAsync]    => effFunc() result: %T %+v\n", result, result)

	return result
}

func debugApply(fn Any, arg Any) Any {
	// Print the shapes
	fmt.Printf("[debugApply] calling Apply(...) with fn=%T, arg=%T\n", fn, arg)

	res := Apply(fn, arg)
	switch x := res.(type) {
	case nil:
		fmt.Println("[debugApply] result => nil!")
	case func() Any:
		fmt.Printf("[debugApply] result => zero-arg function func() Any\n")
	case func(Any) Any:
		fmt.Printf("[debugApply] result => one-arg function func(Any) Any\n")
	default:
		fmt.Printf("[debugApply] result => some other type: %T\n", x)
	}
	return res
}

// function sequential(util, supervisor, par) {
//     return new Aff(ASYNC, function (cb) {
//       return function () {
//         return runPar(util, supervisor, par, cb);
//       };
//     });
//   }

func sequential(util Any, supervisor Any, par Any) Async {
	return Async{asyncFn: func(cb Any /*AsyncCallback*/) Any {
		return func() Canceler {
			return runPar(util, supervisor, par, cb)
		}
	}}
}

func runPar(util Any, supervisor Any, par Any, cb Any) Canceler {
	// TODO
	panic("Implement parallel")
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
		// fmt.Println("Draining...", size)
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
		// fmt.Println("Enqueuing...")
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

func Fiber(util_ Any, supervisor Any, aff Any) Any {
	var util Dict = util_.(Dict)
	var isLeft func(Any) Any = func(x Any) Any {
		return util["isLeft"].(func(Any) Any)(x)
	}
	var left func(Any) Any = util["left"].(func(Any) Any)
	var right func(Any) Any = util["right"].(func(Any) Any)
	var fromRight func(Any) Any = util["fromRight"].(func(Any) Any)
	var fromLeft func(Any) Any = util["fromLeft"].(func(Any) Any)
	runTick := 0

	step := aff       // Successful step
	var fail Any      // Failure step
	var interrupt Any // Asynchronouse interrupt

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
	joinId := 0
	joins := map[int]OnComplete{}
	var rethrow = true

	status := SUSPENDED

	var run func(int) Any
	run = func(localRunTick int) Any {
		// fmt.Println("New round", "localRunTick", localRunTick, "bracketCount", bracketCount)
		// fmt.Println("Step", step)
		var tmp Any
		var result Any
		var attempt Any
		for {
			tmp = nil
			result = nil
			attempt = nil
			switch status {
			case STEP_BIND:
				// fmt.Println("STEP_BIND", b.head)
				status = CONTINUE // next step
				headFn := b.head.(func(Any) Any)
				// fmt.Printf("STEP_BIND headFn: %T\n", headFn)
				// Makeshift try catch block
				func() {
					// defer func() {
					// 	if err := recover(); err != nil {
					// 		// early return on error
					// 		status = RETURN
					// 		fail = left(err)
					// 		step = nil
					// 	}
					// }()
					newStep := headFn(step)
					// fmt.printf("STEP_BIND result: %T, value: %v\n", newStep, newStep)
					step = newStep
					if b.tail == nil {
						b.head = nil
					} else {
						b.head = b.tail.head
						b.tail = b.tail.tail
					}
				}()
			case STEP_RESULT:
				if step == nil {
					status = RETURN
				} else if isLeft(step).(bool) {
					status = RETURN
					fail = step
					step = nil
				} else if b.head == nil {
					status = RETURN
				} else {
					status = STEP_BIND
					step = fromRight(step)
				}

			case CONTINUE:
				// fmt.Printf("CONTINUE step: %T, value: %+v\n", step, step)

				switch currentStep := step.(type) {

				case Bind:
					// fmt.Println("\tBind")
					if b.head != nil {
						b.tail = &Cons{head: b.head, tail: b.tail}
					}
					b.head = currentStep.bToAff
					status = CONTINUE
					step = currentStep.affOfB

				case Pure:
					// fmt.Println("\tPure")
					if b.head == nil {
						// fmt.Println("\t> Head nil")
						// we're done
						status = RETURN
						step = right(currentStep.value)
					} else {
						// fmt.println("\t> Head exists", currentStep.value)
						// this happens after a bind
						status = STEP_BIND
						step = currentStep.value
					}

				case Sync:
					// fmt.Println("\tSync", currentStep)
					status = STEP_RESULT
					step = runSync(left, right, currentStep.eff)

				case Async:
					fmt.Println("\tAsync")
					fmt.Printf("\tAsyncFn type: %T\n", currentStep.asyncFn)
					status = PENDING
					step = runAsync(left, currentStep.asyncFn, func(theResult Any) func() {
						return func() {
							if runTick != localRunTick {
								return
							}
							runTick++
							scheduler.enqueue(func() {
								if runTick != localRunTick+1 {
									return
								}
								status = STEP_RESULT
								step = theResult
								run(runTick)
							})
						}
					})
					return nil

				case Throw:
					// fmt.Println("\tThrow")
					status = RETURN
					fail = left(currentStep.err)
					step = nil

				case Catch:
					// fmt.Println("\tCatch", currentStep.aff)
					if b.head == nil {
						// fmt.println("\t\tHead is nil")
						attempts = &InterruptCons{interrupt: interrupt, head: step, tail: attempts}
					} else {
						// fmt.println("\t\tHead is not nil", b.head)
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
					// fmt.Println("\tBracket")
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
					// fmt.Println("\tFork")
					status = STEP_RESULT
					tmp = Fiber(util, supervisor, currentStep.affOfB)
					if supervisor != nil {
						supervisor.(Dict)["register"].(func(Any))(tmp)
					}
					if currentStep.questionableBool {
						tmp.(Dict)["run"].(func())()
					}
					step = util["right"].(func(Any) Any)(tmp)
				case Sequential:
					// fmt.Println("\tSequential")
					status = CONTINUE
					step = sequential(util, supervisor, currentStep.parAff)
				case func() Any:
					// fmt.println("Step is a function, executing it")
					step = step.(func() Any)() // Execute the function and get the result
					// fmt.println("Result of function execution:", step)
					status = STEP_RESULT // Or set appropriate state
					// fmt.printf("Unhandled step after execution: %T, value: %+v\n", step, step)
					panic("Unhandled step in CONTINUE")
				default:
					// fmt.println("Unhandled step in CONTINUE:", step)

				}

			case RETURN:
				// fmt.Printf("RETURN: step=%+v, fail=%+v, interrupt=%+v\n", step, fail, interrupt)
				b.head = nil
				b.tail = nil
				if attempts == nil || attempts.head == nil {
					// fmt.printf("\tNo attempts: setting status COMPLETED\n")
					status = COMPLETED
					if interrupt != nil {
						step = interrupt
					} else if fail != nil {
						step = fail
					}
				} else {
					// fmt.printf("\tHave attempts: head=%+v\n", attempts.head)
					tmp = attempts.interrupt
					attempt = attempts.head
					attempts = attempts.tail

					switch currentAttempt := attempt.(type) {
					case Catch:
						// fmt.println("\tReturn Catch")
						if (interrupt != nil) && interrupt != tmp && bracketCount == 0 {
							// fmt.println("\t\tGonna RETURN")
							status = RETURN
						} else if fail != nil {
							// fmt.println("\t\tGonna CONTINUE")

							status = CONTINUE
							step = currentAttempt.errorToAff(fromLeft(fail))
							fail = nil
						}
						// We cannot resume from an unmasked interrupt or exception.
					case Resume:
						// fmt.println("\tResume")
						if (interrupt != nil) && interrupt != tmp && bracketCount == 0 || fail != nil {
							status = RETURN
						} else {
							b.head = currentAttempt.b.head
							b.tail = currentAttempt.b.tail
							status = STEP_BIND
							step = fromRight(step)
							// fmt.println("\t\tNext step", step)
						}

					case Bracket:
						// fmt.println("\tBracket")
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
						// fmt.println("\tRelease")
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
						if (interrupt != nil) && interrupt != tmp && bracketCount == 0 {
							step = currentAttempt.bracketConditions["killed"].(func(Any) Any)(fromLeft(interrupt)).(func(Any) Any)(currentAttempt.result)
						} else if fail != nil {
							step = currentAttempt.bracketConditions["failed"].(func(Any) Any)(fromLeft(fail)).(func(Any) Any)(currentAttempt.result)
						} else {
							step = currentAttempt.bracketConditions["completed"].(func(Any) Any)(fromRight(step)).(func(Any) Any)(currentAttempt.result)
						}
						fail = nil
						bracketCount++
					case Finalizer:
						// fmt.println("\tFinalizer")
						bracketCount++
						attempts = &InterruptCons{
							head:      Finalized{step: step, fail: fail},
							tail:      attempts,
							interrupt: interrupt,
						}
						status = CONTINUE
						step = currentAttempt.finalizer
					case Finalized:
						// fmt.println("\tFinalized")
						bracketCount--
						status = RETURN
						step = currentAttempt.step
						fail = currentAttempt.fail
					default:
						break
					}
				}
			case COMPLETED:
				// fmt.println("COMPLETED", joins)
				for _, join := range joins {
					rethrow = rethrow && join.rethrow
					join.handler(step)()
				}
				joins = nil
				if (interrupt != nil) && fail != nil {
					go panic(fromLeft(fail))
				} else if isLeft(step).(bool) && rethrow {
					go panic(fromLeft(step))
				}
				return nil

			case SUSPENDED:
				// fmt.println("SUSPENDED")
				status = CONTINUE
			case PENDING:
				// fmt.Println("PENDING")
				return nil

			default:
				break
			}
		}

	}
	runFn := func() Any {
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

	// function onComplete(join) {
	//      return function () {
	//        if (status === COMPLETED) {
	//          rethrow = rethrow && join.rethrow;
	//          join.handler(step)();
	//          return function () {};
	//        }

	//        var jid    = joinId++;
	//        joins      = joins || {};
	//        joins[jid] = join;

	//        return function() {
	//          if (joins !== null) {
	//            delete joins[jid];
	//          }
	//        };
	//      };
	//    }

	onComplete := func(join OnComplete) func() Any {
		return func() Any {
			if status == COMPLETED {
				rethrow = rethrow && join.rethrow
				join.handler(step)()
				return func() Any { return nil }
			}
			joinId += 1
			jid := joinId
			if joins == nil {
				joins = map[int]OnComplete{}
			}
			joins[jid] = join

			return func() Any {
				if joins != nil {
					delete(joins, jid)
				}
				return nil
			}
		}
	}

	join := func(cb func(Any) func() Any) func() Any {

		return func() Any {
			canceler := onComplete(OnComplete{rethrow: false, handler: cb})()
			if status == SUSPENDED {
				run(runTick)
			}
			return canceler
		}
	}
	kill := func(error Any, cbAny Any) Any {
		// fmt.Println("kill:")
		cb := cbAny.(AsyncCallback)
		return func() Any {
			if status == COMPLETED {
				// If the fiber is already completed, notify the callback
				Run(Apply(cb, right(nil)))
				return func() Any { return nil }
			}

			// Register a completion handler that will notify the callback once done
			canceler := onComplete(OnComplete{
				rethrow: false,
				handler: func(result Any) func() Any {
					return func() Any {
						Run(Apply(cb, right(nil)))
						return nil
					}
				},
			})()

			// Handle based on the current status of the fiber
			switch status {
			case SUSPENDED:
				// fmt.Println("SUSPENDED")
				// Interrupt the fiber and complete it
				interrupt = left(error)
				status = COMPLETED
				step = interrupt
				run(runTick)

			case PENDING:
				// fmt.Println("PENDING")
				// If the fiber is pending, mark it as interrupted
				if interrupt == nil {
					interrupt = left(error)
				}
				// If no bracket is protecting this fiber, add a finalizer and return
				if bracketCount == 0 {
					if status == PENDING {
						// Nur wenn step wirklich eine Funktion ist
						if stepFn, ok := step.(func(Any) Any); ok {
							attempts = &InterruptCons{
								head:      Finalizer{finalizer: stepFn(error)},
								tail:      attempts,
								interrupt: interrupt,
							}
						}
					}
					status = RETURN
					step = nil
					fail = nil
					run(runTick + 1)
				}

			default:
				// fmt.Println("Other status")
				// For other statuses, mark as interrupted and transition to RETURN
				if interrupt == nil {
					interrupt = left(error)
				}
				if bracketCount == 0 {
					status = RETURN
					step = nil
					fail = nil
				}
			}

			return canceler
		}
	}
	return Dict{
		"run":         runFn,
		"kill":        kill,
		"join":        join,
		"onComplete":  onComplete,
		"isSuspended": func() Any { return status == SUSPENDED },
	}
}

func init() {
	exports := Foreign("Effect.Aff")

	// ∷ ∀	 a b. Aff a → (a → Aff b) → Aff b
	exports["_pure"] = func(value Any) Any {
		// fmt.Println("func: _pure")
		return Pure{value: value}
	}
	// ∷ ∀ a. Error → Aff a
	exports["_throwError"] = func(e_ Any) Any {
		// fmt.Println("func: _throwError")
		e := e_.(error)
		return Throw{err: e}
	}
	// ∷ ∀ a. Aff a → (Error → Aff a) → Aff a
	exports["_catchError"] = func(aff Any) Any {
		// fmt.Println("func: _catchError")
		return func(k Any) Any {
			return Catch{aff: aff, errorToAff: k.(func(Any) Any)}
		}
	}
	// ∷ ∀ a b. (a → b) → Aff a → Aff b
	exports["_map"] = func(f_ Any) Any {
		f := f_.(func(Any) Any)
		return func(aff Any) Any {
			switch a := aff.(type) {
			case Pure:
				return Pure{value: f(a.value)}
			default:
				return Bind{
					affOfB: aff,
					bToAff: func(b Any) Any { return Pure{value: f(b)} }}
			}
		}
	}

	// ∷ ∀ a b. Aff a → (a → Aff b) → Aff b
	exports["_bind"] = func(aff Any) Any {
		return func(f_ Any) Any {
			f := f_.(func(Any) Any)
			bToAff := func(b Any) Any {
				return f(b)
			}
			return Bind{affOfB: aff, bToAff: bToAff}
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
		// fmt.Println("func: _makeFiber")
		return func() Any {
			return Fiber(util, nil, aff)
		}
	}

	// ∷ ∀ a. ((Either Error a -> Effect Unit) -> Effect Canceler) -> Aff a
	// exports["makeAff"] = func(asyncFn AsyncFn) Async { return Async{asyncFn: asyncFn} }

	exports["makeAff"] = func(fn Any) Any {
		return Async{
			asyncFn: func(cb Any) Any {
				return func() Any {
					// The user function is one arg, returning a zero-arg effect:
					effVal := Apply(fn, cb) // partial application fn(cb)
					if effVal == nil {
						panic("makeAff: Apply(fn, cb) was nil (mismatch)")
					}

					// That result is a zero-arg function: run it
					canc := Run(effVal)
					return canc
				}
			},
		}
	}

	// ∀ a. Effect a → Aff a
	exports["_liftEffect"] = func(effect_ Any) Any {
		// fmt.Println("func: _liftEffect")
		effect := effect_.(func() Any)
		return Sync{eff: effect}
	}

	// ∀ a. Fn.Fn2 (Unit → Either a Unit) Number (Aff Unit)
	exports["_delay"] = func(right_ Any, millis_ Any) Any {
		right := right_.(func(Any) Any)
		millis := int(millis_.(float64))

		return Async{asyncFn: func(cb Any) Any {
			timer := time.NewTimer(time.Duration(millis) * time.Millisecond)

			// Äquivalent zu setDelay(ms, cb(right()))
			go func() {
				<-timer.C
				Run(Apply(cb, right(nil)))
			}()

			// Äquivalent zu return function() { return Aff.Sync... }
			return func() Canceler {
				return func(error Any) Any {
					return Sync{eff: func() Any {
						stopped := timer.Stop()
						return right(stopped) // Geben wir zurück ob Timer gestoppt wurde, ähnlich wie clearTimeout
					}}
				}
			}
		}}
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
