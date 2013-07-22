// game369 project main.go
package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	//"unicode/utf8"
)

type (
	TBoard [9][9]bool
	TCo    struct {
		x, y int
	}
)

var (
	fg termbox.Attribute
	bg termbox.Attribute
)

var (
	board TBoard
	cur   TCo
)

const (
	cDefault termbox.Attribute = iota
	cBlack
	cRed
	cGreen
	cYellow
	cBlue
	cMagenta
	cCyan
	cWhite
)

const (
	aBold termbox.Attribute = 1 << (iota + 4)
	aUnderline
	aReverse
)

const (
	boardTop  = (24 - (9*2 + 2)) / 2
	boardLeft = (80 - (9*4 + 3)) / 2
	crossMark = '╳'
)

func printStr(x, y int, data string) {
	//k := 0
	var k int
	for _, char := range data {
		termbox.SetCell(x+k, y, char, fg, bg)
		if (char < 1<<7) || (char > 9471) && (char < 9600) {
			k++
		} else {
			k += 2
		}
	}
}

func textcolor(newColor termbox.Attribute) {
	fg = newColor
}

func textbackground(newColor termbox.Attribute) {
	bg = newColor
}

func (b *TBoard) init() (err error) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			b[i][j] = false
		}
	}
	err = nil
	return
}

func (b *TBoard) drawBox() {
	printStr(boardLeft, boardTop, "    1   2   3   4   5   6   7   8   9")
	printStr(boardLeft, boardTop+1, "  ┌───┬───┬───┬───┬───┬───┬───┬───┬───┐")
	for i := 0; i < 9; i++ {
		printStr(boardLeft, boardTop+i*2+2, string('1'+rune(i))+" │   │   │   │   │   │   │   │   │   │")
		if i != 8 {
			printStr(boardLeft, boardTop+i*2+3, "  ├───┼───┼───┼───┼───┼───┼───┼───┼───┤")
		} else {
			printStr(boardLeft, boardTop+i*2+3, "  └───┴───┴───┴───┴───┴───┴───┴───┴───┘")
		}
	}
	return
}

func (co TCo) pointTo() (ex TCo) {
	ex.x = boardLeft + 4 + co.x*4
	ex.y = boardTop + 2 + co.y*2
	return
}

func (b TBoard) drawCell(co TCo) {
	var ex TCo
	ex = co.pointTo()
	textcolor(cWhite | aBold)
	printStr(ex.x, ex.y, string(crossMark))
	textcolor(cDefault)
}

func (b TBoard) highlightCell(co TCo) {
	var (
		mark string
		ex   TCo
	)
	ex = co.pointTo()
	textcolor(cWhite | aBold)
	textbackground(cGreen)
	if b[co.x][co.y] {
		mark = string(crossMark)
	} else {
		mark = " "
	}
	printStr(ex.x, ex.y, mark)
	textcolor(cDefault)
	textbackground(cDefault)
}

func readkey() termbox.Key {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			return ev.Key
		}
	}
}

func (a TCo) add(b TCo) (c TCo) {
	return TCo{a.x + b.x, a.y + b.y}
}

func (a TCo) isInRange(b TCo, c TCo) bool {
	return (a.x >= b.x) && (a.x <= c.x) && (a.y >= b.y) && (a.y <= c.y)
}

func (b TBoard) gen() {
	var (
		mark rune
		ex   TCo
	)
	err := termbox.Clear(cDefault, cDefault)
	if err != nil {
		panic(err)
	}
	textcolor(cCyan | aBold)
	textbackground(cDefault)
	b.drawBox()
	textcolor(cDefault)
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if b[i][j] {
				mark = rune(crossMark)
				textcolor(cWhite | aBold)
			} else {
				mark = ' '
				textcolor(cDefault)
			}
			if (i == cur.x) && (j == cur.y) {
				textbackground(cGreen)
			} else {
				textbackground(cDefault)
			}
			ex = TCo{i, j}.pointTo()
			termbox.SetCell(ex.x, ex.y, mark, fg, bg)
		}
	}
	termbox.Flush()
	return
}

func (b *TBoard) playMove(co TCo) {
	b[co.x][co.y] = true
}

func main() {
	/*fmt.Println("🂡")
	fmt.Println(len("🂡"))
	fmt.Println(utf8.RuneCountInString("🂡"))
	return*/
	board.init()
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	fg = termbox.ColorDefault
	bg = termbox.ColorDefault
	textcolor(cCyan | aBold)
	board.drawBox()
	board.drawCell(TCo{0, 0})
	termbox.Flush()
	board.gen()
readinput:
	for {
		tempKey := readkey()
		tempCur := TCo{-1, -1}
		switch tempKey {
		case termbox.KeyArrowDown:
			tempCur = cur.add(TCo{0, 1})
		case termbox.KeyArrowUp:
			tempCur = cur.add(TCo{0, -1})
		case termbox.KeyArrowLeft:
			tempCur = cur.add(TCo{-1, 0})
		case termbox.KeyArrowRight:
			tempCur = cur.add(TCo{1, 0})
		case termbox.KeySpace:
			(&board).playMove(cur)
		case termbox.KeyEsc:
			break readinput
		}
		if tempCur.isInRange(TCo{0, 0}, TCo{8, 8}) {
			cur = tempCur
		}
		board.gen()
	}
	fmt.Println("Hello World!")
}
