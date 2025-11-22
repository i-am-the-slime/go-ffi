package purescript_aff

import (
	"fmt"
	"sync"
	"time"

	. "github.com/purescript-native/go-runtime"
)

// Re-export types for clarity
type Effect = EffFn // Effect is func() Any

// Internal queue for goroutines to schedule effects on main thread
var effectQueue = make(chan EffFn, 100)

// DrainEffectQueue processes any pending effects (call from main thread/loop)
func DrainEffectQueue() {
	for {
		select {
		case eff := <-effectQueue:
			Run(eff)
		default:
			return
		}
	}
}

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
	result    Any // Memoized result
}

func (p ParMap) getResult() Any {
	return p.result
}

// forall b. Apply (ParAff (b -> a)) (ParAff b)
type ParApply struct {
	parAffOfBToA Any
	parAffOfB    Any
	result       Any // Memoized result
}

func (p ParApply) getResult() Any {
	return p.result
}

// Alt (ParAff a) (ParAff a)
type ParAlt struct {
	option1 Any
	option2 Any
	result  Any // Memoized result
}

func (p ParAlt) getResult() Any {
	return p.result
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
	if cb == nil {
		panic("runAsync: callback is nil")
	}
	if asyncFn == nil {
		panic("runAsync: asyncFn is nil")
	}

	// Apply asyncFn to cb, returning an Effect Canceler
	part := Apply(asyncFn, cb)
	if part == nil {
		panic("[runAsync] Apply(asyncFn, cb) returned nil!")
	}

	// Run the Effect to get the Canceler
	result := Run(part)
	return result
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

// Forked represents a leaf node in the parallel tree
type Forked struct {
	fid    int
	resume *Cons // Stack to resume (head, tail)
	result Any   // Memoized result (EMPTY if not yet resolved)
}

var EMPTY = struct{}{}

func runPar(util Any, supervisor Any, par Any, cb Any) Canceler {
	utilDict := util.(Dict)
	isLeft := utilDict["isLeft"].(func(Any) Any)
	fromRight := utilDict["fromRight"].(func(Any) Any)
	left := utilDict["left"].(func(Any) Any)
	right := utilDict["right"].(func(Any) Any)

	// Mutable state - protected by mutex for thread safety
	var mu sync.Mutex
	fiberId := 0
	fibers := make(map[int]Any)
	killId := 0
	kills := make(map[int]map[int]Any)
	
	// Early exit error for Alt cancellation
	earlyExit := fmt.Errorf("[ParAff] Early exit")
	
	var interrupt Any = nil
	var root Any = EMPTY

	// kill cancels a subtree
	var kill func(error Any, par Any, cb func(Any) func() Any) map[int]Any
	kill = func(killError Any, par Any, cb func(Any) func() Any) map[int]Any {
		step := par
		var head Any = nil
		var tail *Cons = nil
		count := 0
		killsMap := make(map[int]Any)
		
		for {
			switch currentStep := step.(type) {
			case *Forked:
				if currentStep.result == EMPTY {
					mu.Lock()
					if fiber, ok := fibers[currentStep.fid]; ok {
						idx := count
						count++
						fiberDict := fiber.(Dict)
						killFn := fiberDict["kill"].(func(Any, Any) Any)
						// Call killFn to get the effect, then call the effect to get the canceler
						killEffect := killFn(killError, func(result Any) func() Any {
							return func() Any {
								mu.Lock()
								count--
								if count == 0 {
									mu.Unlock()
									cb(result)()
								} else {
									mu.Unlock()
								}
								return nil
							}
						})
						// Execute the effect to get the canceler
						if effectFn, ok := killEffect.(func() Any); ok {
							killsMap[idx] = effectFn()
						}
					}
					mu.Unlock()
				}
				
				// Terminal case
				if head == nil {
					goto done
				}
				
				// Go down the right side of the tree
				switch h := head.(type) {
				case *ParApply:
					step = h.parAffOfB
				case *ParAlt:
					step = h.option2
				default:
					goto done
				}
				
				// Move to next head from stack
				if tail == nil {
					head = nil
				} else {
					head = tail.head
					tail = tail.tail
				}
				
			case *ParMap:
				step = currentStep.parAffOfB
				
			case *ParApply:
				if head != nil {
					tail = &Cons{head: head, tail: tail}
				}
				head = step
				step = currentStep.parAffOfBToA
				
			case *ParAlt:
				if head != nil {
					tail = &Cons{head: head, tail: tail}
				}
				head = step
				step = currentStep.option1
			
			default:
				goto done
			}
		}
		
	done:
		if count == 0 {
			cb(right(nil))()
		} else {
			// Run all cancellation effects and store the resulting cancelers
			for i := 0; i < count; i++ {
				if cancelFn, ok := killsMap[i].(func() Any); ok {
					killsMap[i] = cancelFn()
				}
			}
		}
		
		return killsMap
	}

	// Helper to get result from Any type (reads result field which may be accessed by other goroutines)
	// Note: results are write-once, so we can read without mutex after checking != EMPTY
	getResult := func(node Any) Any {
		switch n := node.(type) {
		case *ParMap:
			return n.result
		case *ParApply:
			return n.result
		case *ParAlt:
			return n.result
		case *Forked:
			return n.result
		default:
			return EMPTY
		}
	}

	// join bubbles results back up the tree
	var join func(result Any, head Any, tail *Cons)
	join = func(result Any, head Any, tail *Cons) {
		var fail Any
		var step Any
		
		if isLeft(result).(bool) {
			fail = result
			step = nil
		} else {
			step = result
			fail = nil
		}
		
		for {
			mu.Lock()
			if interrupt != nil {
				mu.Unlock()
				return
			}
			mu.Unlock()
			
			// Reached root
			if head == nil {
				if fail != nil {
					cb.(func(Any) func() Any)(fail)()
				} else {
					cb.(func(Any) func() Any)(step)()
				}
				return
			}
			
			// Check if already computed (with double-checked locking)
			if getResult(head) != EMPTY {
				return
			}
			
			switch h := head.(type) {
			case *ParMap:
				// Lock to prevent race when multiple fibers complete simultaneously
				mu.Lock()
				if h.result != EMPTY {
					// Already computed by another fiber
					mu.Unlock()
					return
				}
				if fail == nil {
					mapped := right(h.bToA(fromRight(step)))
					h.result = mapped
					step = mapped
				} else {
					h.result = fail
				}
				mu.Unlock()
				
			case *ParApply:
				mu.Lock()
				if h.result != EMPTY {
					// Already computed by another fiber
					mu.Unlock()
					return
				}
				lhs := getResult(h.parAffOfBToA)
				rhs := getResult(h.parAffOfB)
				mu.Unlock()
				
				if fail != nil {
					// Set result under lock
					mu.Lock()
					if h.result != EMPTY {
						mu.Unlock()
						return
					}
					h.result = fail
					kid := killId
					killId++
					var toKill Any
					if fail == lhs {
						toKill = h.parAffOfB
					} else {
						toKill = h.parAffOfBToA
					}
					mu.Unlock()
					
					// Kill the other side (without holding mutex)
					done := false
					currentKid := kid
					innerKills := kill(earlyExit, toKill, func(Any) func() Any {
						return func() Any {
							mu.Lock()
							delete(kills, currentKid)
							mu.Unlock()
							if !done {
								done = true
								return nil
							}
							if tail == nil {
								join(fail, nil, nil)
							} else {
								join(fail, tail.head, tail.tail)
							}
							return nil
						}
					})
					
					mu.Lock()
					kills[currentKid] = innerKills
					mu.Unlock()
					
					if done {
						// Kill completed synchronously, callback will handle join
						return
					}
					// Kill is pending, we return and wait for callback
					return
				} else if lhs == EMPTY || rhs == EMPTY {
					// Can't proceed yet
					return
				} else {
					// Apply the function
					mu.Lock()
					if h.result != EMPTY {
						mu.Unlock()
						return
					}
					fn := fromRight(lhs).(func(Any) Any)
					arg := fromRight(rhs)
					applied := right(fn(arg))
					h.result = applied
					step = applied
					mu.Unlock()
				}
				
			case *ParAlt:
				mu.Lock()
				if h.result != EMPTY {
					// Already computed by another fiber
					mu.Unlock()
					return
				}
				lhs := getResult(h.option1)
				rhs := getResult(h.option2)
				mu.Unlock()
				
				// Wait for at least one side or both errors
				if lhs == EMPTY && (rhs == EMPTY || !isLeft(rhs).(bool)) {
					return
				}
				if rhs == EMPTY && (lhs == EMPTY || !isLeft(lhs).(bool)) {
					return
				}
				
				// Both errors - use first
				if lhs != EMPTY && isLeft(lhs).(bool) && rhs != EMPTY && isLeft(rhs).(bool) {
					mu.Lock()
					if h.result != EMPTY {
						mu.Unlock()
						return
					}
					if step == lhs {
						fail = rhs
					} else {
						fail = lhs
					}
					step = nil
					h.result = fail
					mu.Unlock()
				} else {
					// One succeeded - use it and kill the other
					mu.Lock()
					if h.result != EMPTY {
						mu.Unlock()
						return
					}
					h.result = step
					kid := killId
					killId++
					var toKill Any
					if step == lhs {
						toKill = h.option2
					} else {
						toKill = h.option1
					}
					mu.Unlock()
					
					// Kill the other side (without holding mutex)
					done := false
					currentKid := kid
					innerKills := kill(earlyExit, toKill, func(Any) func() Any {
						return func() Any {
							mu.Lock()
							delete(kills, currentKid)
							mu.Unlock()
							if !done {
								done = true
								return nil
							}
							if tail == nil {
								join(step, nil, nil)
							} else {
								join(step, tail.head, tail.tail)
							}
							return nil
						}
					})
					
					mu.Lock()
					kills[currentKid] = innerKills
					mu.Unlock()
					
					if done {
						// Kill completed synchronously, callback will handle join
						return
					}
					// Kill is pending, we return and wait for callback
					return
				}
			}
			
			// Move up the tree
			if tail == nil {
				head = nil
			} else {
				head = tail.head
				tail = tail.tail
			}
		}
	}

	// resolve creates completion handler for a fiber
	resolve := func(forked *Forked) func(Any) func() Any {
		return func(result Any) func() Any {
			return func() Any {
				mu.Lock()
				delete(fibers, forked.fid)
				forked.result = result
				resumeHead := forked.resume.head
				var resumeTail *Cons = nil
				if forked.resume.tail != nil {
					resumeTail = forked.resume.tail
				}
				mu.Unlock()
				
				join(result, resumeHead, resumeTail)
				return nil
			}
		}
	}

	// run walks the parallel tree and forks fibers
	run := func() {
		const (
			RUN_CONTINUE = 1
			RUN_RETURN   = 2
		)
		
		status := RUN_CONTINUE
		step := par
		var head Any = nil
		var tail *Cons = nil
		
		for {
			switch status {
			case RUN_CONTINUE:
				switch currentStep := step.(type) {
				case ParMap:
					if head != nil {
						tail = &Cons{head: head, tail: tail}
					}
					head = &ParMap{bToA: currentStep.bToA, parAffOfB: EMPTY, result: EMPTY}
					step = currentStep.parAffOfB
					
				case ParApply:
					if head != nil {
						tail = &Cons{head: head, tail: tail}
					}
					head = &ParApply{parAffOfBToA: EMPTY, parAffOfB: currentStep.parAffOfB, result: EMPTY}
					step = currentStep.parAffOfBToA
					
				case ParAlt:
					if head != nil {
						tail = &Cons{head: head, tail: tail}
					}
					head = &ParAlt{option1: EMPTY, option2: currentStep.option2, result: EMPTY}
					step = currentStep.option1
					
				default:
					// Leaf node - create a fiber
					mu.Lock()
					fid := fiberId
					fiberId++
					mu.Unlock()
					
					forked := &Forked{
						fid:    fid,
						resume: &Cons{head: head, tail: tail},
						result: EMPTY,
					}
					
					fiber := Fiber(util, supervisor, step)
					fiberDict := fiber.(Dict)
					onCompleteFn := fiberDict["onComplete"].(func(OnComplete) func() Any)
					onCompleteFn(OnComplete{
						rethrow: false,
						handler: resolve(forked),
					})()
					
					mu.Lock()
					fibers[fid] = fiber
					mu.Unlock()
					
					if supervisor != nil {
						supervisorDict := supervisor.(Dict)
						registerFn := supervisorDict["register"].(func(Any))
						registerFn(fiber)
					}
					
					status = RUN_RETURN
					step = forked
				}
				
			case RUN_RETURN:
				// Terminal case
				if head == nil {
					goto done
				}
				
				// Fill in the left side
				switch h := head.(type) {
				case *ParMap:
					if h.parAffOfB == EMPTY {
						h.parAffOfB = step
						status = RUN_CONTINUE
						step = h.parAffOfB
					} else {
						step = head
						if tail == nil {
							head = nil
						} else {
							head = tail.head
							tail = tail.tail
						}
					}
					
				case *ParApply:
					if h.parAffOfBToA == EMPTY {
						h.parAffOfBToA = step
						status = RUN_CONTINUE
						step = h.parAffOfB
					} else {
						h.parAffOfB = step
						step = head
						if tail == nil {
							head = nil
						} else {
							head = tail.head
							tail = tail.tail
						}
					}
					
				case *ParAlt:
					if h.option1 == EMPTY {
						h.option1 = step
						status = RUN_CONTINUE
						step = h.option2
					} else {
						h.option2 = step
						step = head
						if tail == nil {
							head = nil
						} else {
							head = tail.head
							tail = tail.tail
						}
					}
				}
			}
		}
		
	done:
		root = step
		
		// Start all fibers
		mu.Lock()
		fibersCopy := make(map[int]Any)
		for k, v := range fibers {
			fibersCopy[k] = v
		}
		mu.Unlock()
		
	for _, fiber := range fibersCopy {
		fiberDict := fiber.(Dict)
		runFn := fiberDict["run"].(func() interface{})
		Run(runFn)
	}
	}

	// cancel kills the entire tree
	cancel := func(cancelError Any, cb func(Any) func() Any) Any {
		mu.Lock()
		interrupt = left(cancelError)
		
		// Cancel all pending kills
		for _, innerKills := range kills {
			for _, killFn := range innerKills {
				if fn, ok := killFn.(func() Any); ok {
					fn()
				}
			}
		}
		kills = make(map[int]map[int]Any)
		mu.Unlock()
		
		newKills := kill(cancelError, root, cb)
		
		return func(killError Any) Any {
			return Async{asyncFn: func(killCb Any) Any {
				return func() Any {
					for _, killFn := range newKills {
						if fn, ok := killFn.(func() Any); ok {
							fn()
						}
					}
					return nonCanceler
				}
			}}
		}
	}

	run()

	// Return canceler
	return func(killError Any) Any {
		return Async{asyncFn: func(killCb Any) Any {
			return func() Any {
				return cancel(killError, func(result Any) func() Any {
					return func() Any {
						if kcb, ok := killCb.(func(Any) func() Any); ok {
							kcb(result)()
						}
						return nil
					}
				})
			}
		}}
	}
}

func SupervisorNew(util Any) Any {
	utilDict := util.(Dict)
	isLeft := utilDict["isLeft"].(func(Any) Any)
	fromLeft := utilDict["fromLeft"].(func(Any) Any)
	
	var mu sync.Mutex
	fibers := make(map[int]Any)
	fiberId := 0
	count := 0
	
	register := func(fiber Any) {
		mu.Lock()
		defer mu.Unlock()
		
		fid := fiberId
		fiberId++
		
		fiberDict := fiber.(Dict)
		onCompleteFn := fiberDict["onComplete"].(func(OnComplete) func() Any)
		onCompleteFn(OnComplete{
			rethrow: true,
			handler: func(result Any) func() Any {
				return func() Any {
					mu.Lock()
					defer mu.Unlock()
					count--
					delete(fibers, fid)
					return nil
				}
			},
		})()
		
		fibers[fid] = fiber
		count++
	}
	
	isEmpty := func() Any {
		mu.Lock()
		defer mu.Unlock()
		return count == 0
	}
	
	killAll := func(killError Any, cb func()) func() Any {
		return func() Any {
			mu.Lock()
			if count == 0 {
				mu.Unlock()
				cb()
				return nil
			}
			
			killCount := 0
			kills := make(map[int]Any)
			fibersCopy := make(map[int]Any)
			for k, v := range fibers {
				fibersCopy[k] = v
			}
			
			// Clear state
			fibers = make(map[int]Any)
			fiberId = 0
			count = 0
			mu.Unlock()
			
			// Kill each fiber
			for fid, fiber := range fibersCopy {
				fiberDict := fiber.(Dict)
				killFn := fiberDict["kill"].(func(Any, Any) Any)
				
				currentFid := fid
				// killFn returns func() Any, which when called returns the canceler
				killEffect := killFn(killError, func(result Any) func() Any {
					return func() Any {
						mu.Lock()
						delete(kills, currentFid)
						killCount--
						if isLeft(result).(bool) && fromLeft(result) != nil {
							err := fromLeft(result)
							go func() {
								time.Sleep(0)
								panic(err)
							}()
						}
						shouldCallback := killCount == 0
						mu.Unlock()
						
						if shouldCallback {
							cb()
						}
						return nil
					}
				})
				// Call the effect to get the canceler
				if killEffectFn, ok := killEffect.(func() Any); ok {
					kills[currentFid] = killEffectFn()
				}
				// Increment killCount AFTER setting up the callback
				// This ensures the callback can safely decrement it
				mu.Lock()
				killCount++
				mu.Unlock()
			}
			
			// Return canceler for the killAll operation
			return func(error Any) Any {
				return Sync{eff: func() Any {
					mu.Lock()
					defer mu.Unlock()
					for _, killFn := range kills {
						if fn, ok := killFn.(func() Any); ok {
							fn()
						}
					}
					return nil
				}}
			}
		}
	}
	
	return Dict{
		"register": register,
		"isEmpty":  isEmpty,
		"killAll":  killAll,
	}
}


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
					step = runAsync(left, currentStep.asyncFn, func(theResult Any) func() Any {
						return func() Any {
					if runTick != localRunTick {
						return nil
					}
					runTick++
					// Create Effect with explicit type and enqueue it
					var eff EffFn = func() Any {
						if runTick != localRunTick+1 {
							return nil
						}
						status = STEP_RESULT
						step = theResult
						run(runTick)
						return nil
					}
					Run(eff)
					return nil
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
				status = STEP_RESULT
				tmp = Fiber(util, supervisor, currentStep.affOfB)
			if supervisor != nil {
				supervisor.(Dict)["register"].(func(Any))(tmp)
			}
			if currentStep.questionableBool {
				Run(tmp.(Dict)["run"].(func() interface{}))
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
				panic(fromLeft(fail))
			} else if isLeft(step).(bool) && rethrow {
				panic(fromLeft(step))
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
			// Just run directly - no need for scheduler in Go
			run(runTick)
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

	join := func(cb Any) Any {
		return func() Any {
			cbFn := cb.(func(Any) func() Any)
			canceler := onComplete(OnComplete{rethrow: false, handler: cbFn})()
			if status == SUSPENDED {
				run(runTick)
			}
			return canceler
		}
	}
	kill := func(error Any, cbAny Any) Any {
		// fmt.Println("kill:")
		if cbAny == nil {
			panic("Fiber.kill: callback is nil")
		}
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
			if cb == nil {
				panic("_delay: callback is nil")
			}
			timer := time.NewTimer(time.Duration(millis) * time.Millisecond)
			
			// Construct the result and Effect on main thread BEFORE goroutine starts
			result := right(nil)
			effectToRun := Apply(cb, result)

		// Goroutine waits, then queues the pre-constructed effect for main thread
		go func() {
			<-timer.C
			// Queue the effect for main thread to run
			effectQueue <- effectToRun.(EffFn)
		}()

			// Return canceler
			return func() Canceler {
				return func(error Any) Any {
					return Sync{eff: func() Any {
						stopped := timer.Stop()
						return right(stopped)
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

	// ∷ ∀ a. Boolean → Aff a → Aff (Fiber a)
	exports["_fork"] = func(immediate Any) Any {
		return func(aff Any) Any {
			immediateBool := immediate.(bool)
			return Fork{
				questionableBool: immediateBool,
				affOfB:           aff,
			}
		}
	}

	// ParAff map
	exports["_parAffMap"] = func(f Any) Any {
		return func(parAff Any) Any {
			return ParMap{
				bToA:      f.(func(Any) Any),
				parAffOfB: parAff,
				result:    EMPTY,
			}
		}
	}

	// ParAff apply
	exports["_parAffApply"] = func(parAff1 Any) Any {
		return func(parAff2 Any) Any {
			return ParApply{
				parAffOfBToA: parAff1,
				parAffOfB:    parAff2,
				result:       EMPTY,
			}
		}
	}

	// ParAff alt
	exports["_parAffAlt"] = func(parAff1 Any) Any {
		return func(parAff2 Any) Any {
			return ParAlt{
				option1: parAff1,
				option2: parAff2,
				result:  EMPTY,
			}
		}
	}

	// ParAff ~> Aff
	exports["_sequential"] = func(par Any) Any {
		return Sequential{parAff: par}
	}

	// ∷ ∀ a. Fn.Fn2 FFIUtil (Aff a) (Effect { fiber :: Fiber a, supervisor :: Supervisor })
	exports["_makeSupervisedFiber"] = func(util Any, aff Any) Any {
		return func() Any {
			supervisor := SupervisorNew(util)
			fiber := Fiber(util, supervisor, aff)
			return Dict{
				"fiber":      fiber,
				"supervisor": supervisor,
			}
		}
	}

	// ∷ ∀ a. Fn.Fn3 Error Supervisor (Effect Unit) (Effect (Canceler))
	exports["_killAll"] = func(error Any, supervisor Any, cb Any) Any {
		supervisorDict := supervisor.(Dict)
		killAllFn := supervisorDict["killAll"].(func(Any, func()) func() Any)
		cbFunc := cb.(func() Any)
		return killAllFn(error, func() {
			cbFunc()
		})
	}

	exports["nonCanceler"] = nonCanceler
	
	// Internal function to process queued effects from goroutines
	exports["drainEffectQueueImpl"] = func() Any {
		return func() Any {
			DrainEffectQueue()
			return nil
		}
	}

}
