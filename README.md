## Synchronizing Goroutines

```go
func task(name string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(name, "completed")
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go task("Task 1", &wg)
	go task("Task 2", &wg)

	wg.Wait()
	fmt.Println("All tasks completed")
}
```

## Communicating Between Goroutines

```go
func sendData(ch chan<- int, data int) {
	fmt.Println("Sending", data)
	ch <- data // Block until received
	fmt.Println("Finished sending", data)
}

func receiveData(ch <-chan int) {
	val := <-ch // Block until value is sent
	fmt.Println("Received", val)
}

func main() {
	ch := make(chan int)

	go sendData(ch, 10)
	go receiveData(ch)

	time.Sleep(time.Second)
	fmt.Println("Done")
}
```

## Avoiding Race Conditions

```go
type safeCounter struct {
	mu    sync.Mutex
	count int
}

func (sc *safeCounter) increment() {
	sc.mu.Lock()
	sc.count++
	sc.mu.Unlock()
}

func main() {
	sc := &safeCounter{}
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			sc.increment()
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			sc.increment()
		}
	}()

	wg.Wait()
	fmt.Println("Final count:", sc.count)
}
```

```sh
go run -race main.go

go test -race
```

<!-- ## Fan-Out/Fan-In -->

## Worker Pools

```go
func worker(id int, tasks <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range tasks {
		// Do some work, e.g., multiply by 2
		results <- t * 2
	}
}

func main() {
	tasks := make(chan int, 10)
	results := make(chan int, 10)
	var wg sync.WaitGroup

	// Create 3 workers
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, tasks, results, &wg)
	}

	// Send tasks
	for i := 1; i <= 5; i++ {
		tasks <- i
	}
	close(tasks)

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for r := range results {
		fmt.Println("Result:", r)
	}
}
```

> [!NOTE] Worker Pools
> **_Create a fixed number of workers (goroutines) that read tasks from a channel._**  
> **_This approach prevents spawning a huge number of goroutines if tasks spike._**

## Web Server

```go
// factorial computes factorial of n in a naive way.
func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

func main() {
	var wg sync.WaitGroup

	http.HandleFunc("/factorial", func(w http.ResponseWriter, r *http.Request) {
		nStr := r.URL.Query().Get("n")
		n, err := strconv.Atoi(nStr)
		if err != nil {
			http.Error(w, "Invalid number", http.StatusBadRequest)
			return
		}

		// Increment WaitGroup for each request
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			result := factorial(num)
			// Write the response (this part is safe if you only write from one goroutine at a time)
			fmt.Fprintf(w, "Factorial(%d) = %d\n", num, result)
		}(n)
	})

	log.Println("Server starting at :8080")
	// Start the server (blocking call)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

> [**Goroutines in Go: A Practical Guide to Concurrency**](https://getstream.io/blog/goroutines-go-concurrency-guide/)

## Race Condition

```go
type post struct {
	views int
}

func (p *post) inc(wg *sync.WaitGroup) {
	defer wg.Done()

	p.views += 1
}

func main() {
	var wg sync.WaitGroup

	p := post{views: 0}

	for range 100 {
		wg.Add(1)
		go p.inc(&wg)
	}

	wg.Wait()
	fmt.Println(p.views)
}
```

```go
type post struct {
	views int
	mu    sync.Mutex
}

func (p *post) inc(wg *sync.WaitGroup) {

	defer func() {
		wg.Done()
		p.mu.Unlock()
	}()

	p.mu.Lock()
	time.Sleep(time.Second * 2)
	p.views += 1
}

func main() {
	var wg sync.WaitGroup

	p := post{views: 0}

	for range 2 {
		wg.Add(1)
		go p.inc(&wg)
	}

	wg.Wait()
	fmt.Println(p.views)
}
```

> [!NOTE] Mutual Exclusion
> **_Enable access synchronization to a resource thorough memory. It enusres that only one go routine can access the shared resource at a time._**  
> **_Critical section refers to a portion of your code that modifies or accesses a shared resource since only one goroutine can be inside a critical section protected by mutexes at a time it's important to keep the section as efficient as possible and unnecessary operations within the critical section can become a bottleneck hindering the overall performance of you concrete progra, this is hwere the read and write mutexes comes into play._**  
> **_Unlike traditional mutuxes which enforces exclusive access the read write mutuxes allow for a more granular approach to synchronization by using read write mutuxes effectively you can optimize access to shared resource allowing concurrent reads while ensuring exclusive rights this can greatly improve performance and scalability of you go application specially in scenarios with frequent read operations._**

> [!TIP] RWMutex
> **_RWMutex is very similar to the Mutex struct._**  
> **_RWMutex provides a little more flexibility as compared to the Mutex._**  
> **_RWMutex provides better performances for frequent reads compared to a regualer Mutex._**  
> **_Read locks are non-blocking, aloowing multiple goroutines to read the shared resource concurrenlty._**

> [**Go sync.Mutex: Normal and Starvation Mode**](https://victoriametrics.com/blog/go-sync-mutex/)
