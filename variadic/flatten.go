package variadic

import (
    "reflect"
)

// Flatten takes a variable number of arguments and returns an array
// containing, in sequence, each argument or - if an argument is an array or
// a slice - that argument's contents.
//
// e.g. Flatten(1, []int{2, 3, 4}, 5) => 1, 2, 3, 4, 5
func Flatten(args ... interface{}) []interface{} {
    result := make([]interface{}, 0)
    
    for _, arg := range args {
        v := reflect.ValueOf(arg)
        vk := v.Kind()
        
        if (vk == reflect.Array) || (vk == reflect.Slice) {
            for i := 0; i < v.Len(); i++ {
                result = append(result, v.Index(i).Interface())
            }
        } else {
            result = append(result, arg)
        }
    }
    
    return result
}

// FlattenRecursive takes a variable number of arguments and returns an array
// containing, in sequence, each argument or - if an argument is an array or
// a slice - that argument's contents, flattened recursively.
//
// e.g. FlattenRecursive([]interface{}{1, 2, []int{3, 4}}, 5) => 1, 2, 3, 4, 5
func FlattenRecursive(args ... interface{}) []interface{} {
    result := make([]interface{}, 0)
    
    for _, arg := range args {
        v := reflect.ValueOf(arg)
        vk := v.Kind()
        
        if (vk == reflect.Array) || (vk == reflect.Slice) {
            for i := 0; i < v.Len(); i++ {
                result = append(result, FlattenRecursive(v.Index(i).Interface())...)
            }
        } else {
            result = append(result, arg)
        }
    }
    
    return result
}

// FlattenExcludingNils takes a variable number of arguments and returns an
// array containing, in sequence, each non-nil argument or - if an argument is
// an array or a slice - that argument's non-nil contents.
//
// e.g. FlattenExcludingNils(1, nil, []int{2, nil, 3, 4}, 5) => 1, 2, 3, 4, 5
func FlattenExcludingNils(args ... interface{}) []interface{} {
    result := make([]interface{}, 0)
    
    for _, arg := range args {
        if arg == nil { continue }
        
        v := reflect.ValueOf(arg)
        vk := v.Kind()
        
        if (vk == reflect.Array) || (vk == reflect.Slice) {
            for i := 0; i < v.Len(); i++ {
                element := v.Index(i)
                if element.IsNil() { continue }
                result = append(result, element.Interface())
            }
        } else {
            result = append(result, arg)
        }
    }
    
    return result
}
