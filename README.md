# go-worker-pool
Thank you for all the good libraries and Articles:
1. [Tunny](https://github.com/Jeffail/tunny)
2. [Wpool](https://github.com/gotohr/wpool)
3. [Handling 1 millon requests per minute with golang](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/)
4. [Visually understanding worker pool](https://medium.com/coinmonks/visually-understanding-worker-pool-48a83b7fc1f5)

## How to use this?
#### Install package
```bash
$ go get github.com/jayhuang75/go-worker-pool
```

#### In your application main.go, import the package
```go
import (
    "github.com/jayhuang75/go-worker-pool"
)
```

#### Example how to use the worker pool
```go
// Person Struct
type Person struct {
	Name string
	Age  int
}

// ResourceProcessor func
func ResourceProcessor(resource interface{}) error {
	// fmt.Printf("worker: started, working for %s\n", resource)
	fmt.Printf(">>>>>>>>>>>>>> %s \n", resource.(Person).Name+" ok")
	return nil
}

// ResultProcessor func
func ResultProcessor(result worker.Result) error {
	fmt.Printf("Result processor got error: %s\n", result.Err)
	fmt.Printf("Result processor got result: %d\n", result.Job)
	return nil
}

func main() {

	p1 := Person{"apple ", 3}
	p2 := Person{"orange", 8}
	p3 := Person{"pear", 35}
	p4 := Person{"pizza", 3}
	p5 := Person{"cafe", 8}

	persons := []Person{p1, p2, p3, p4, p5}

	numCPUs := runtime.NumCPU()

	// convert the Struct to the interface
	resources := make([]interface{}, len(persons))
	for i, s := range persons {
		resources[i] = s
	}

	pool := worker.NewPool(numCPUs)
	pool.Start(resources, ResourceProcessor, ResultProcessor)

}
```
