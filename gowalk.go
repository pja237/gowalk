package main

import (
	"fmt"
	"os"
	"path/filepath"
	//	"runtime"
	"strconv"
	"sync"
)

// sync attempt with waitgroups
var wg sync.WaitGroup
var block chan int

func walk(cur string) {
	var i int = 0
	//fmt.Println("In ", cur)
	defer wg.Done()
	block <- 1

	file, err := os.Open(cur)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
		//panic(err)
	}
	filesall, err := file.Readdir(-1)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
		//panic(err)
	}
	for _, v := range filesall {
		fmt.Println(filepath.Join(cur, v.Name()), v.IsDir(), v.ModTime())
		if v.IsDir() == true {
			//fmt.Println("ENTERing ", filepath.Join(cur, v.Name()))
			i++
			//go walk(filepath.Join(cur, v.Name()), sync_ch)
			wg.Add(1)
			go walk(filepath.Join(cur, v.Name()))
		}
	}
	<-block
}

func main() {
	// this is sent through ENV{'GOMAXPROCS'}
	//procs, _ := strconv.Atoi(os.Args[1])
	//runtime.GOMAXPROCS(procs)
	//fmt.Println("GOMAXPROCS=",runtime.GOMAXPROCS(procs))
	ngophers, _ := strconv.Atoi(os.Args[1])
	block = make(chan int, ngophers)

	wg.Add(1)
	go walk(os.Args[2])
	wg.Wait()

}
