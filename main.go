package main

import (
	"rhionin.com/Rhionin/SanderServer/progress"
	"fmt"
)

func main() {
	wips := progress.CheckProgress();
	fmt.Println(wips)
}