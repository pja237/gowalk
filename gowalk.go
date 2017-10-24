package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

func walk(cur string, sync chan int) {
	var i int = 0
	sync_ch:=make(chan int)
	//fmt.Println("In ", cur)
	file,err:=os.Open(cur)
	if err!=nil {
		fmt.Println("ERROR: ", err)
		sync<-1
		return
		//panic(err)
	}
	filesall,err:=file.Readdir(-1)
	if err!=nil {
		fmt.Println("ERROR: ", err)
		sync<-1
		return
		//panic(err)
	}
	for _,v:=range filesall {
		fmt.Println(filepath.Join(cur, v.Name()), v.IsDir(), v.ModTime())
		if v.IsDir()==true {
			//fmt.Println("ENTERing ", filepath.Join(cur, v.Name()))
			i++
			//go walk(filepath.Join(cur, v.Name()), sync_ch)
			go walk(cur+"/"+v.Name(), sync_ch)
		}
	}
	//fmt.Println(cur, " SPAWNED ", i)
	if i==0 {
		//fmt.Println(cur, " DONE")
		sync<-1
		return
	} else {
		for i>0 {
			//fmt.Println(cur, " WAITING: ", i)
			//fmt.Println(cur, "GOT: ", <-sync_ch)
			<-sync_ch
			i--
		}
	}
	sync<-1
}


func main() {
	procs,_:=strconv.Atoi(os.Args[1])
	runtime.GOMAXPROCS(procs)
	//fmt.Println("GOMAXPROCS=",runtime.GOMAXPROCS(procs))
	sync:=make(chan int)
	go walk(os.Args[2], sync)
	<-sync
}
