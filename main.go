package main

import (
	"fmt"
	"sort"
)

func main() {
	nums := []int{2, 2, 1}
	nums1 := []int{4, 1, 2, 1, 2}
	result := singleNumber(nums)
	result1 := singleNumber(nums1)
	fmt.Println(" 只出现一次的数字:", result)
	fmt.Println(" 只出现一次的数字", result1)
	fmt.Println("-----------------")
	x := 121
	result3 := isPalindrome(x)
	fmt.Println("回文数判断:", result3)
	x1 := -121
	result4 := isPalindrome(x1)
	fmt.Println("回文数判断:", result4)
	fmt.Println("-----------------")
	s := "([)]"
	result5 := isValid(s)
	fmt.Println("有效括号判断：", result5)
	s1 := "()[]{}"
	result6 := isValid(s1)
	fmt.Println("有效括号判断：", result6)
	fmt.Println("-----------------")
	strs := []string{"flower", "flow", "flight"}
	result7 := longestCommonPrefix(strs)
	fmt.Println("最长公共前缀：", result7)
	strs1 := []string{"dog", "racecar", "car"}
	result8 := longestCommonPrefix(strs1)
	fmt.Println("最长公共前缀：", result8)
	digits := []int{1, 2, 3}
	res1 := plusOne(digits)
	fmt.Println("-----------------")
	fmt.Println("加1：", res1)
	digits2 := []int{4, 3, 2, 1}
	res2 := plusOne(digits2)
	fmt.Println("加1：", res2)
	digits3 := []int{9}
	res3 := plusOne(digits3)
	fmt.Println("加1：", res3)
	nums4 := []int{1, 1, 2}
	res4 := removeDuplicates(nums4)
	fmt.Println("删除重复项：", res4)
	test1 := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	res5 := merge(test1)
	fmt.Println("合并区间：", res5)
	nums5 := []int{2, 7, 11, 15}
	target5 := 9
	res6 := twoSum(nums5, target5)
	fmt.Println("两数之和：", res6)
}

// 只出现一次的数字
func singleNumber(nums []int) int {
	result := 0
	for _, num := range nums {
		result ^= num
	}
	return result
}

// 回文数
func isPalindrome(x int) bool {
	if x < 0 || (x%10 == 0 && x != 0) {
		return false
	}
	reversedHalf := 0
	for x > reversedHalf {
		reversedHalf = reversedHalf*10 + x%10
		x = x / 10
	}
	return x == reversedHalf || x == reversedHalf/10
}

// 有效的括号
func isValid(s string) bool {
	stack := []rune{}
	pairs := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}

	for _, char := range s {
		switch char {
		case '(', '{', '[':
			stack = append(stack, char) // 左括号入栈
		case ')', '}', ']':
			if len(stack) == 0 || stack[len(stack)-1] != pairs[char] {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0 // 栈为空表示所有括号都匹配
}

// 最长公共前缀
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for _, str := range strs[1:] {
		i := 0
		for i < len(prefix) && i < len(str) && prefix[i] == str[i] {
			i++
		}
		prefix = prefix[:i]
		if len(prefix) == 0 {
			return ""
		}
	}
	return prefix
}

// 加一
func plusOne(digits []int) []int {
	n := len(digits)
	for i := n - 1; i >= 0; i-- {
		if digits[i] < 9 {
			digits[i]++
			return digits
		}
		digits[i] = 0
	}
	return append([]int{1}, digits...)
}

// 删除有序数组中的重复项
func removeDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	slow := 0
	//便利数组，找出新的唯一的元素
	for i := 0; i < len(nums); i++ {
		if nums[i] != nums[slow] {
			slow++
			nums[slow] = nums[i]
		}
	}
	return slow + 1
}

// 合并区间
func merge(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return [][]int{}
	}
	//区间按照从小到达排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	result := [][]int{intervals[0]}
	for i := 1; i < len(intervals); i++ {
		currentInterval := intervals[i]
		currentStart := currentInterval[0] // 当前区间的起始值
		currentEnd := currentInterval[1]   // 当前区间的结束值

		lastResultIndex := len(result) - 1 // 最后一个区间的索引
		lastInterval := result[lastResultIndex]
		lastEnd := lastInterval[1]
		if currentStart <= lastEnd {
			newEnd := lastEnd
			if currentEnd > newEnd {
				newEnd = currentEnd
			}
			result[lastResultIndex][1] = newEnd
		} else {
			result = append(result, currentInterval)
		}

	}
	return result
}

// 两数之和
func twoSum(nums []int, target int) []int {
	numMap := make(map[int]int)
	for i, num := range nums {
		complement := target - num
		if j, exist := numMap[complement]; exist {
			return []int{j, i}
		}
		numMap[num] = i
	}
	return nil
}
