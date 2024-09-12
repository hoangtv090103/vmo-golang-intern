## Solution 1: Using Channel

- `c <- value`: Sends data through a channel. Only one goroutine can receive the data at a time, preventing race conditions.
- `go func() { for val := range ch { counter += val } }()`: Goroutine that listens on the channel and updates the counter, ensuring only one goroutine modifies it at a time.
- **Closing channels**: When a channel is closed, no further data can be sent, ensuring that goroutines do not mistakenly write to it.

## Solution 2: Using Mutex

- Use `sync.Mutex` to lock a shared variable during updates and unlock it after.
- `mutex.Lock()`: Acquires the lock, ensuring that only one goroutine can access the shared resource (critical section).
- `mutex.Unlock()`: Releases the lock after the critical section is executed, allowing other waiting goroutines to proceed.
- **Blocking behavior**: If a goroutine tries to lock while another holds the lock, it will block until the lock is released.

## Solution 3: Read/Write Locks

- `rwMutex.RLock()`: Allows multiple goroutines to concurrently read the data. If no write lock is active, any number of readers can proceed.
- `rwMutex.Lock()`: Acquires a write lock, preventing any readers or writers from accessing the data until the write operation is completed.
- `rwMutex.RUnlock()` and `rwMutex.Unlock()`: Release the read and write locks, respectively, allowing other goroutines to proceed.

## Solution 4: Using Atomic Operations

- `atomic.AddInt32(&counter, 23)`: Atomically increments the value of counter by 23. This ensures that the operation happens as a single, indivisible step.
- `atomic.LoadInt32(&counter)`: Reads the value of counter atomically, ensuring that no other goroutine modifies it during the read.
- **Atomic nature**: The operation is guaranteed to be completed without interference from other goroutines, preventing race conditions.

## Solution 5: Immutability

- Each goroutine has its own `localCounter`, which is immutable (each goroutine does not share this variable).
- After the local calculations, the result is sent to a channel. The main goroutine then sums the results to get the final counter.
- By avoiding shared mutable state, there is no risk of race conditions.

## Final Results

- **Final Counter: 38571**
  - Solving RC using Channel
- **Final Counter: 46000**
  - Solving RC using Mutex
- **Final Counter: 46000**
  - Solving RC using Read/Write Lock
- **Final Counter: 46000**
  - Solving RC using Atomic Operation
- **Final Counter: 46000**
  - Solving RC using Immutability
