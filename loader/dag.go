package loader

import (
    "fmt"
    "strings"
)

// dag is directed acyclic graph of dependencies, made up of disconnected
// subgraphs.
type dag struct {
    // nodes is an array of Tasks.
    // The array index is the index in the DAG.
    nodes      []*Task

    // requires means that this Task cannot complete until these Tasks complete.
    // The array index is the index in the DAG.
    requires   [][]int

    // requiredBy means that this Task blocks these Tasks from completing.
    // The array index is the index in the DAG.
    requiredBy [][]int

    // resultRequires means that this Task requires the result of these Tasks
    // (which may have completed). NOTE: This array is ordered!
    // The array index is the index in the DAG.
    resultRequires [][]int

    // resultRequiredBy means that this task still needs to provide a result
    // to these Tasks.
    // The array index is the index in the DAG.
    resultRequiredBy [][]int

    // every leaf node until it is removed from pending
    pending []int

    // every result collected, which may be set to nil when no longer needed.
    // The array index is the index in the DAG.
    results []interface{}
}

// scope provides scoping for task names across subsequent sibling tasks and
// child subtasks, additionally scoped by a single invocation of Loader Add
type scope map[string]int

// dup creates a copy of scope. In this way the function callstack can
// implement a stack of scope states.
func (s scope) dup() scope {
    n := make(map[string]int)
    for k, v := range s {
        n[k] = v
    }
    return n
}

// scope stringer for debugging
func (s scope) String() string {
    result := make([]string, 0)
    for k, v := range s {
        result = append(result, fmt.Sprintf("%q=>%d", k, v))
    }
    return strings.Join(result, ", ")
}

func (dag *dag) isLeaf(idx int) bool {
    return len(dag.requires[idx]) == 0
}

// addEdge such that "a requires b" and "a requires the result of b"
func (dag *dag) addEdge(a int, b int) {
    dag.requires[a]         = append(dag.requires[a],         b)
    dag.requiredBy[b]       = append(dag.requiredBy[b],       a)
    dag.resultRequires[a]   = append(dag.resultRequires[a],   b)
    dag.resultRequiredBy[b] = append(dag.resultRequiredBy[b], a)
}

// add: see the loader Add method
func (dag *dag) add(tasks []Task) error {
    return dag.addTasks(tasks, -1, make(scope))
}

// addTasks: see the loader Add method
func (dag *dag) addTasks(tasks []Task, parent int, scope scope) error {
    subscope := scope.dup()

    dag.results = append(dag.results, make([]interface{}, len(tasks))...)

    for i, task := range tasks {

        dag.nodes            = append(dag.nodes,            &tasks[i])
        dag.requires         = append(dag.requires,         make([]int, 0))
        dag.requiredBy       = append(dag.requiredBy,       make([]int, 0))
        dag.resultRequires = append(dag.resultRequires,   make([]int, 0))
        dag.resultRequiredBy = append(dag.resultRequiredBy, make([]int, 0))

        idx := len(dag.nodes) - 1

        if parent >= 0 {
            // parent always < idx
            dag.addEdge(parent, idx)
        }

        if task.Name != "" {
            subscope[task.Name] = idx
        }

        for _, named := range task.RequiresNamed {
            namedDep, exists := subscope[named]
            if !exists {
                return fmt.Errorf("error adding task: named requirement %q not in scope", named)
            }
            dag.addEdge(idx, namedDep)
        }

        if tasks[i].RequiresDirect != nil {
            err := dag.addTasks(tasks[i].RequiresDirect, idx, subscope)
            if err != nil { return err }
        }

        if dag.isLeaf(idx) {
            dag.pending = append(dag.pending, idx)
        }
    }

    return nil
}

func (dag *dag) inputs(idx int) []interface{} {
    var inputs []interface{}
    for _, requirement := range dag.resultRequires[idx] {
        inputs = append(inputs, dag.results[requirement])
    }
    return inputs
}

func (dag *dag) registerResult(idx int, result interface{}) {
    dag.results[idx] = result

    // remove the edges between the task and its parents
    parents := dag.requiredBy[idx]
    for _, parent := range parents {
        dag.requires[parent] = intArrayFindAndDeleteElement(dag.requires[parent], idx)
    }
    dag.requiredBy[idx] = dag.requiredBy[idx][0:0] // empty

    // if any parent is now a leaf, it is added to pending
    for _, parent := range parents {
        if dag.isLeaf(parent) {
            dag.pending = append(dag.pending, parent)
        }
    }
}

/*
func edgesString(xs [][]int) string {
    result := make([]string, 0)
    edgesStringF(xs, func(s string) {
        result = append(result, s)
    })
    return strings.Join(result, "")
}

func edgesStringF(xs [][]int, result func(string)) {
    result("{\n")
    for i, x := range xs {
        result(fmt.Sprintf("    %d => {", i))

        for _, y := range x {
            result(fmt.Sprintf("%d, ", y))
        }

        result("}\n")
    }
    result("}")
}
*/
