package loader

// intArrayDeleteElement removes an item at index idx from an integer array
// without preserving order by swapping with the last item, then shrinking
func intArrayDeleteElement(xs []int, idx int) []int {
    xs[idx] = xs[len(xs) - 1]
    return xs[:len(xs) - 1]
}

// intArrayFindAndDeleteElement finds and removes an item from an integer array
// without preserving order by swapping with the last item, then shrinking
func intArrayFindAndDeleteElement(xs []int, element int) []int {
    for i, x := range xs {
        if x == element {
            return intArrayDeleteElement(xs, i)
        }
    }
    panic("element does not exist in array as expected")
}
