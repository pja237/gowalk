package main

import (
	"fmt"
	"os"
	"path/filepath"
    "sync"
    "runtime"
    "strconv"
)

// sync attempt with waitgroups
var wg sync.WaitGroup

func walk(cur string) {
	var i int = 0
	//fmt.Println("In ", cur)
    defer wg.Done()

	file,err:=os.Open(cur)
	if err!=nil {
		fmt.Println("ERROR: ", err)
		return
		//panic(err)
	}
	filesall,err:=file.Readdir(-1)
	if err!=nil {
		fmt.Println("ERROR: ", err)
		return
		//panic(err)
	}
	for _,v:=range filesall {
		fmt.Println(filepath.Join(cur, v.Name()), v.IsDir(), v.ModTime())
		if v.IsDir()==true {
			//fmt.Println("ENTERing ", filepath.Join(cur, v.Name()))
			i++
			//go walk(filepath.Join(cur, v.Name()), sync_ch)
            wg.Add(1)
			go walk(filepath.Join(cur, v.Name()))
		}
	}
}


func main() {
    // this is sent through ENV{'GOMAXPROCS'}
	procs,_:=strconv.Atoi(os.Args[1])
	runtime.GOMAXPROCS(procs)
	//fmt.Println("GOMAXPROCS=",runtime.GOMAXPROCS(procs))

    wg.Add(1)
    go walk(os.Args[2])
    wg.Wait()

}
