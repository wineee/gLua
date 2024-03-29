package state

// 返回栈顶索引
func (self *luaState) GetTop() int {
	return self.stack.top
}

// 把索引转换为绝对索引
func (self *luaState) AbsIndex(idx int) int {
	return self.stack.absIndex(idx)
}


// 看是否还可以推入n个值而不会导致溢出
func (self *luaState) CheckStack(n int) bool {
	self.stack.check(n)
	return true // never fails
}

// 从栈顶弹出n个值
func (self *luaState) Pop(n int) {
	for i := 0; i < n; i++ {
		self.stack.pop()
	}
}

// 值从一个位置复制到另一个位置
func (self *luaState) Copy(fromIdx, toIdx int) {
	val := self.stack.get(fromIdx)
	self.stack.set(toIdx, val)
}

// 把指定索引处的值推入栈
func (self *luaState) PushValue(idx int) {
	val := self.stack.get(idx)
	self.stack.push(val)
}

// 将栈顶值弹出，然后写入指定位置
func (self *luaState) Replace(idx int) {
	val := self.stack.pop()
	self.stack.set(idx, val)
}

// 将指定值移到栈顶
func (self *luaState) Insert(idx int) {
	self.Rotate(idx, 1)
}

// 删除指定索引处的值，然后将该值上面的值全部下移一个位置
func (self *luaState) Remove(idx int) {
	self.Rotate(idx, -1)
	self.Pop(1)
}

// 将[idx, top]索引区间内的值朝栈顶方向旋转n个位置。如果n是负数，那么实际效果就是朝栈底方向旋转
// 使用三次翻转的方法
func (self *luaState) Rotate(idx, n int) {
	t := self.stack.top - 1           /* end of stack segment being rotated */
	p := self.stack.absIndex(idx) - 1 /* start of segment */
	var m int                         /* end of prefix */
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	self.stack.reverse(p, m)   /* reverse the prefix with length 'n' */
	self.stack.reverse(m+1, t) /* reverse the suffix */
	self.stack.reverse(p, t)   /* reverse the entire segment */
}

// 将栈顶索引设置为指定值。如果指定值小于当前栈顶索引，效果则相当于弹出操作
func (self *luaState) SetTop(idx int) {
	newTop := self.stack.absIndex(idx)
	if newTop < 0 {
		panic("stack underflow!")
	}

	n := self.stack.top - newTop
	if n > 0 {
		for i := 0; i < n; i++ {
			self.stack.pop()
		}
	} else {
		for i := 0; i > n; i-- {
			self.stack.push(nil)
		}
	}
}
