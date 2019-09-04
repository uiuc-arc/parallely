package main

import "fmt"
import "parallely"



func func_0() {
  defer parallely.Wg.Done()
  var x int;
var y int;
x = 0;
y = 0;
__flag_1 := false;
 x = parallely.RandchoiceFlag(float32(0.9), y, 3, &__flag_1);
__flag_2 := false;
 y = parallely.RandchoiceFlag(float32(0.99), x, 5, &__flag_2);

 if __flag_2 {
 __flag_2 = false;
 y = parallely.RandchoiceFlag(float32(0.99), x, 5, &__flag_2);

 }
 __flag_1 = __flag_1 || __flag_2;


 if __flag_1 {
 __flag_1 = false;
 x = parallely.RandchoiceFlag(float32(0.9), y, 3, &__flag_1);
__flag_2 := false;
 y = parallely.RandchoiceFlag(float32(0.99), x, 5, &__flag_2);

 if __flag_2 {
 __flag_2 = false;
 y = parallely.RandchoiceFlag(float32(0.99), x, 5, &__flag_2);

 }
 __flag_1 = __flag_1 || __flag_2;


 }
 

  fmt.Println("Ending thread : ", 0);
}

func main() {
	fmt.Println("Starting main thread");
	
	parallely.InitChannels(1);

  go func_0();


	fmt.Println("Main thread waiting for others to finish");  
	parallely.Wg.Wait()
  fmt.Println("Done!");
}
