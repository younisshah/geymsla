## Geymsla  
### A [distributed] key-value store based on BoltDB.


Usage

```go
    if err != nil {
    		panic(err)
    	}
    	defer geymsla.Close()
    	type Person struct {
    		Name string
    	}
    	person := Person{Name: "Younis"}
    	if err := geymsla.Set("me", person); err != nil {
    		panic(err)
    	}
    	value, err := geymsla.Get("me")
    	if err != nil {
    		panic(err)
    	}
    	fmt.Println(value)
```

__TODO__

1) Add Raft or SWIM protocol support.
2) Use Hashicorp's Raft implementation or Serf.
3) Write a comprehensive doc