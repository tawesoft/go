package operator

import (
    "fmt"
    "math"
)

const maxUint = ^uint(0)
const minUint = 0
const maxInt = int(maxUint >> 1)
const minInt = -maxInt - 1

const maxInt8  = math.MaxInt8
const maxInt16 = math.MaxInt16
const maxInt32 = math.MaxInt32
const maxInt64 = math.MaxInt64

const minInt8  = math.MinInt8
const minInt16 = math.MinInt16
const minInt32 = math.MinInt32
const minInt64 = math.MinInt64

const maxUint8  = math.MaxUint8
const maxUint16 = math.MaxUint16
const maxUint32 = math.MaxUint32
const maxUint64 = math.MaxUint64

const minUint8  = 0
const minUint16 = 0
const minUint32 = 0
const minUint64 = 0

const maxFloat32 = math.MaxFloat32
const maxFloat64 = math.MaxFloat64

const minFloat32 = -math.MaxFloat32
const minFloat64 = -math.MaxFloat64

// ErrorOverflow is the return type used for integer overflow in a checked function
var ErrorOverflow  = fmt.Errorf("Overflow")

// ErrorOverflow is the return type used for an undefined operation in a checked function such as shifting a negative
// integer left or divide by zero
var ErrorUndefined = fmt.Errorf("Undefined operation")

// ErrorNaN is the return type used for an operation on a float that is Not a Number in a checked function
var ErrorNaN       = fmt.Errorf("Not a number (NaN)")
