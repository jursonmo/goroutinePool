package main

import(
	"fmt"
	"goropool"
)

func main(){

	p1 := Person{"mo", 29}
	p2 := Person{"chen", 26}
	
	tt := make([]goropool.Task, 4)
	tt = []goropool.Task{
		goropool.NewTask(ShowName, p1, p2),		
		goropool.NewTask(ShowName, p2),
		goropool.NewTask(ShowAge, p2),	
		goropool.NewTask(ShowAll, p1, p2),	
	}

	mygp := goropool.NewGoRoPool(2, 3)
	mygp.Run()
	fmt.Printf("len task %d\n", len(tt))
	for i := 0; i < len(tt); i++ {
		mygp.AddTask(tt[i])
	}
	fmt.Println(mygp.GetFree())
	mygp.WaitJobDone()
}

type Person struct {
	Name string
	Age int
}

func ShowName(args ...interface{}) interface{} {
	fmt.Printf("ShowName args len=%d\n", len(args))
	for _, arg := range args {
		if _, ok := arg.(Person); !ok {
			fmt.Printf(" not person\n")
			return -1
		}	
		fmt.Printf(" name =%s\n", arg.(Person).Name)		
	}
	return nil
}
func ShowAge(args ...interface{}) interface{} {
	fmt.Printf("ShowAge args len=%d\n", len(args))
	for _, arg := range args {
		if _, ok := arg.(Person); !ok {
			fmt.Printf(" not person\n")
			return -1
		}
		fmt.Printf(" age =%d\n", arg.(Person).Age)		
	}
	return nil
}

func ShowAll(args ...interface{}) interface{} {
	fmt.Printf("ShowAll args len=%d\n", len(args))
	for _, arg := range args {
		if _, ok := arg.(Person); !ok {
			fmt.Printf(" not person\n")
			return -1
		}
		fmt.Printf("name=%s, age =%d\n", arg.(Person).Name, arg.(Person).Age)		
	}
	return nil
}