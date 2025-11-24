package main

import (
	"fmt"
	"sort"
	"strconv"
)

func main() {
	nums := []int{1, 2, 2}
	fmt.Println("出现一次的数字:", singleNumber(nums))
	fmt.Println("是否回文数:", isPalindrome(12))
	fmt.Println("是否有效的括号:", isValid("([]{})"))
	strs := []string{"aa", "aab", "ac"}
	fmt.Println("最大公共前缀:", longestCommonPrefix(strs))
	fmt.Println("加1后的数组:", plusOne([]int{9, 9, 9}))
	num1 := []int{0, 1, 1, 2, 2, 3, 3, 3, 3, 3}
	fmt.Println("去重后的数组长度:", removeDuplicates(num1))
	intervals := [][]int{{1, 2}, {2, 6}, {8, 10}, {8, 18}}
	fmt.Println("合并后的区间:", merge(intervals))
	fmt.Println(twoSum([]int{2, 6, 15, 16}, 17))

}

func singleNumber(nums []int) int {
	numberAndCnt := make(map[int]int)
	for _, num := range nums {
		numberAndCnt[num] = numberAndCnt[num] + 1
	}
	for k, v := range numberAndCnt {
		if v == 1 {
			return k
		}
	}
	return -1
}

func isPalindrome(x int) bool {
	if x < 0 || x == 10 {
		return false
	}
	if x < 10 {
		return true
	}
	xStr := strconv.Itoa(x)
	i, j := 0, len(xStr)-1
	for i < j {
		if xStr[i] != xStr[j] {
			return false
		}
		i++
		j--
	}
	return true
}

// 有效的括号
func isValid(s string) bool {
	if s == "" {
		return false
	}
	if len(s)%2 != 0 {
		return false
	}
	symbolMap := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}
	stack := make([]rune, 0, len(s))
	for _, char := range s {
		switch char {
		case '(', '[', '{':
			stack = append(stack, char)
		case ')', ']', '}':
			if len(stack) == 0 || stack[len(stack)-1] != symbolMap[char] {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}

func longestCommonPrefix(strs []string) string {
	if len(strs) <= 1 {
		return ""
	}
	// 假设第一个字符串是公共前缀
	prefix := strs[0]
	// 公共前缀的最大长度
	maxLen := 0
	minLen := 0
	strs = strs[1:]
	for _, str := range strs {
		for i := 0; i < len(str) && i < len(prefix); i++ {
			if str[i] != prefix[i] {
				maxLen = i
				break
			}
		}
		if minLen == 0 {
			minLen = maxLen
		} else if maxLen < minLen {
			minLen = maxLen
		}
	}
	// 公共前缀的最大长度就是公共前缀的长度
	return prefix[:minLen]
}

func plusOne(digits []int) []int {
	// 从后往前遍历数组
	for i := len(digits) - 1; i >= 0; i-- {
		// 如果当前位不是9，直接加1返回
		if digits[i] != 9 {
			digits[i]++
			return digits
		}
		// 如果当前位是9，设为0，继续遍历前一位
		digits[i] = 0
	}
	// 如果所有位都是9，需要在数组最前面插入1
	digits = append([]int{1}, digits...)
	return digits
}

func removeDuplicates(nums []int) int {
	if len(nums) <= 1 {
		return len(nums)
	}
	newLen := 0
	for i, j := 0, 1; i < len(nums) && j < len(nums); {
		if nums[i] != nums[j] {
			i++
			nums[i] = nums[j]
		}
		j++
		if j == len(nums) {
			newLen = i + 1
		}
	}
	return newLen
}

func merge(intervals [][]int) [][]int {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	targetArr := [][]int{intervals[0]}
	for i := 1; i < len(intervals); i++ {
		if intervals[i][0] <= targetArr[len(targetArr)-1][1] &&
			intervals[i][1] > targetArr[len(targetArr)-1][1] {
			targetArr[len(targetArr)-1][1] = intervals[i][1]
		} else {
			targetArr = append(targetArr, intervals[i])
		}
	}
	return targetArr
}

func twoSum(nums []int, target int) []int {
	// 1 2 3 4
	mapSum := make(map[int]int)
	for _, num := range nums {
		if _, ok := mapSum[target-num]; ok {
			return []int{mapSum[target-num], num}
		}
		mapSum[num] = num
	}
	return nil
}
