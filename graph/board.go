package graph

import (
    "strconv"
)

type Board struct {
    s    []uint8
    size uint8
}

func (b *Board) Play(color, x, y uint8) (score uint8, err string) {
    if b.isTesuki(color, x, y) {
        return 0, "Tesuki"
    }
    b.place(color, x, y)
    score = b.searchKills(color, x, y)
    return score, ""
}

func (b *Board) Remove(x, y uint8) {
    b.place(0, x, y)
}

func (b *Board) Create(size uint16) {
    b.s = make([]uint8, size*size)
    b.size = uint8(size)
}

func (b *Board) place(c, x, y uint8) {
    if y >= b.size || x >= b.size || c > 2 {
        panic("Invalid board.place operation") //TODO don't panic
    }
    b.s[x+y*b.size] = c
}

func (b *Board) isTesuki(c, x, y uint8) bool {
    var legalCheck Board
    legalCheck.Create(uint16(b.size))
    copy(legalCheck.s, b.s)
    legalCheck.place(c, x, y)
    legalCheck.searchKills(c, x, y)
    xs := []uint8{x}
    ys := []uint8{y}
    _ = legalCheck.GetGroup(c, 0, xs, ys)
    return !legalCheck.hasFreedom(0, xs, ys)
}

func (b *Board) IsEqual(c *Board) bool {
    for i, v := range b.s {
        if v != c.s[i] {
            return false
        }
    }
    return true
}

func GetAdjecent(x, y uint8) (xa, xb, ya, yb uint8) {
    xa = x - 1
    xb = x + 1
    ya = y - 1
    yb = y + 1
    return xa, xb, ya, yb
}

func (b *Board) hasFreedom(i uint8, xs, ys []uint8) bool {
    xa, xb, ya, yb := GetAdjecent(xs[i], ys[i])

    if xa < b.size && b.GetColor(xa, ys[i]) == 0 {
        return true
    }
    if xb < b.size && b.GetColor(xb, ys[i]) == 0 {
        return true
    }
    if ya < b.size && b.GetColor(xs[i], ya) == 0 {
        return true
    }
    if yb < b.size && b.GetColor(xs[i], yb) == 0 {
        return true
    }
    i += 1
    if len(xs) == int(i) {
        return false
    }
    return b.hasFreedom(i, xs, ys)
}

func (b *Board) GetEmpty() (empty []uint8) {
    for i, v := range b.s {
        if v == 0 {
            y := i / int(b.size)
            x := i % int(b.size)
            empty = append(empty, uint8(x))
            empty = append(empty, uint8(y))
        }
    }
    return empty
}

func (b *Board) GetGroup(c, i uint8, xs, ys []uint8) bool {
    xa, xb, ya, yb := GetAdjecent(xs[i], ys[i])

    if xa < b.size && b.GetColor(xa, ys[i]) == c && !IsPresentinGroup(xa, ys[i], xs, ys) {
        xs = append(xs, xa)
        ys = append(ys, ys[i])
    }
    if xb < b.size && b.GetColor(xb, ys[i]) == c && !IsPresentinGroup(xb, ys[i], xs, ys) {
        xs = append(xs, xb)
        ys = append(ys, ys[i])
    }
    if ya < b.size && b.GetColor(xs[i], ya) == c && !IsPresentinGroup(xs[i], ya, xs, ys) {
        xs = append(xs, xs[i])
        ys = append(ys, ya)
    }
    if yb < b.size && b.GetColor(xs[i], yb) == c && !IsPresentinGroup(xs[i], yb, xs, ys) {
        xs = append(xs, xs[i])
        ys = append(ys, yb)
    }
    i += 1
    if len(xs) == int(i) {
        return true
    }
    return b.GetGroup(c, i, xs, ys)
}

func IsPresentinGroup(x, y uint8, xs, ys []uint8) bool {
    for i, v := range xs {
        if v == x && ys[i] == y {
            return true
        }
    }
    return false
}

func (b *Board) GetColor(x, y uint8) uint8 {
    if y >= b.size || x >= b.size {
        panic("Invalid board.GetColor operation") //TODO don't panic
    }
    return b.s[x+y*b.size]
}

func (b *Board) searchKills(color, x, y uint8) uint8 {
    var opponent uint8
    var score uint8 = 0
    if color == 1 {
        opponent = 2
    } else {
        opponent = 1
    }
    xa, xb, ya, yb := GetAdjecent(x, y)
    if xa < b.size && b.GetColor(xa, y) == opponent {
        xs := []uint8{xa}
        ys := []uint8{y}
        _ = b.GetGroup(opponent, 0, xs, ys)
        if !b.hasFreedom(0, xs, ys) {
            score += b.kill(xs, ys)
        }
    }
    if xb < b.size && b.GetColor(xb, y) == opponent {
        xs := []uint8{xb}
        ys := []uint8{y}
        _ = b.GetGroup(opponent, 0, xs, ys)
        if !b.hasFreedom(0, xs, ys) {
            score += b.kill(xs, ys)
        }
    }
    if ya < b.size && b.GetColor(x, ya) == opponent {
        xs := []uint8{x}
        ys := []uint8{ya}
        _ = b.GetGroup(opponent, 0, xs, ys)
        if !b.hasFreedom(0, xs, ys) {
            score += b.kill(xs, ys)
        }
    }
    if yb < b.size && b.GetColor(x, yb) == opponent {
        xs := []uint8{x}
        ys := []uint8{yb}
        _ = b.GetGroup(opponent, 0, xs, ys)
        if !b.hasFreedom(0, xs, ys) {
            score += b.kill(xs, ys)
        }
    }
    return score
}

func (b *Board) kill(xs, ys []uint8) uint8 {
    for i, v := range xs {
        b.place(0, v, ys[i])
    }
    return uint8(len(xs))
}

func (b *Board) isEnclosed(c, x, y uint8) bool {
    xa, xb, ya, yb := GetAdjecent(x, y)
    found2Colors := false
    if xa < b.size && b.GetColor(xa, y) != c {
        found2Colors = true
    }
    if xb < b.size && b.GetColor(xb, y) != c {
        found2Colors = true
    }
    if ya < b.size && b.GetColor(x, ya) != c {
        found2Colors = true
    }
    if yb < b.size && b.GetColor(x, yb) != c {
        found2Colors = true
    }
    return !found2Colors
}

func (b *Board) tostr() string {
    var board string
    var i uint8 = 1
    for _, c := range b.s {
        board = board + strconv.Itoa(int(c)) + " "
        if i == b.size {
            board = board + "\n"
            i = 1
        } else {
            i += 1
        }
    }
    return board
}
