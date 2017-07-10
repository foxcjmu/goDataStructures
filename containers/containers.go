// containers.go -- Implementation of the containers package
// author: C. Fox
// version: 1/2016
//
// containers provides a standard set of basic containers and collections, including
// two kinds of stacks, maps, and sets, four kinds of queus, and three kinds of lists.
// All containers store values of type interface{}.
// The container types are arranged into a hierachy pictured below.
//
//   Container -- import container
//     Dispenser -- non-iterable containers (not realized in code)
//       Stack -- import containers/stack
//         ArrayStack
//         LinkedStack
//       Queue -- import containers/queue
//         ArrayQueue
//         LinkedQueue
//         ArrayRandomizer
//         LinkedRandomizer
//     Collection -- iterable containers: import container
//       List -- import containers/list
//         ArrayList
//         LinkedList
//         SinglyLinkedList
//       Set -- import containers/set
//         HashSet
//         TreeSet
//       Map import containers/map
//         HashMap
//         TreeMap
//     Iterator -- import containers
//

package containers

// Container is the root type in the containers hierarchy.
// Every Container includes these operations.
type Container interface {
	Size() int   // return the number of items in the container
	Empty() bool // true iff the container is empty
	Clear()      // make a container empty
}

// Collection is the ancestor type of all traversible containers in the hierarchy.
// Every Collection includes these operations.
type Collection interface {
	Container                    // include Size, Clear, and Empty
	Contains(e interface{}) bool // return true iff element e is in the collection
	NewIterator() Iterator       // return a new external Iterator entity
	Apply(f func(interface{}))   // internally iterate and apply f to every element
}

// Iterator is the interface for all Collection external iterators.
type Iterator interface {
	Reset()                    // prepare for another iteration
	Done() bool                // return true iff this iterator is finished
	Next() (interface{}, bool) // return the next element and ok indication
}

// Values stored in sets or maps must be Equalers
type Equaler interface {
	Equal(x interface{}) bool // true iff x is identical to the receiver
}

// Values stored in HashMaps or HashSets must be Hashers
type Hasher interface {
	Equaler
	Hash(s int) int // hash the receiver into 0..(s-1)
}

// Values stored in TreeMaps or TreeSets must be Comparers
type Comparer interface {
	Equaler
	Less(x interface{}) bool // true iff the receiver is less than x
}
