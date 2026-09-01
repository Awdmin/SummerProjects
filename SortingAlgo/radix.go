package main

import(
	"fmt"
	"time"
	"math/rand"
)

func radixPass(base int, nums []uint32, layer int) []uint32 {
	crates := make([][]uint32, base)

	for i := range nums{
		var idx uint32
		if layer == 1 {
			idx = nums[i] & 0xFFFF
		} else {
			idx = (nums[i] >> 16) & 0xFFFF
		}
		crates[idx] = append(crates[idx], nums[i])
	}
	
	temp := make([]uint32, len(nums))
	idx := 0

	for i := range crates {
		for j := range crates[i] {
			temp[idx] = crates[i][j]
			idx++
		}
	}

	return temp

}



func main() {
	n_nums := 100000000
	base := 65536
	nums := make([]uint32, n_nums)

	for i := range n_nums{
		x := rand.Uint32()
		nums[i] = x
	}

	//fmt.Println(nums)

	start := time.Now()

	temp := radixPass(base, nums, 1)	
	sol := radixPass(base, temp, 2)

	elapsed := time.Since(start)
	fmt.Println("Sorting time: ", elapsed)

	for i := 0; i < len(sol)-1; i++ {
		if sol[i] > sol[i+1] {
			fmt.Println("Not sorted!!!")
			break
		}
	}


	fmt.Println("Done!")
}

