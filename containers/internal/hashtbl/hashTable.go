// hashTable.go: Implementation of hash tables for use in sets and maps.
// This implementation uses chaining and is not dynamic (that is, once the
// table size is set, it does not change).
//
// author: C. Fox
// version: 8/2012

// hashtbl provides an implementation of a hash table with the key and value
// aggregated (a hash table) for use in containers.
package hashtbl

import (
	"containers"
	"math"
)

const DefaultTableSize = 991 // how big the make the hash table by default

// These hash tables use chaining, so the hash table is an array of list heads
// whose nodes are tableNodes.
type tableNode struct {
	key   containers.Hasher // key used to locate key-value pair
	value interface{}       // value that goes with the key
	next  *tableNode        // link to the next node
}

// newTableNode creates a new hash table linked list node and
// initializes it to its arguments.
func newTableNode(key containers.Hasher, value interface{}, link *tableNode) *tableNode {
	result := new(tableNode)
	result.key, result.value, result.next = key, value, link
	return result
}

// HashTable is the data structure for a hash table instance.
type HashTable struct {
	tableSize int          // how many slots in the table
	count     int          // how many values are stored in the table
	table     []*tableNode // the hash table itself
}

// Create and return a new empty hash table with an optionally specified
// tableSize. The hash table size will be the first prime number >= tableSize
// if tableSize is specified and is at least 3; otherwise it will be
// DefaultTableSize.
func NewHashTable(tableSize ...int) *HashTable {
	result := new(HashTable)
	result.tableSize = DefaultTableSize
	if 0 < len(tableSize) && 2 < tableSize[0] {
		result.tableSize = nextPrime(tableSize[0])
	}
	result.table = make([]*tableNode, result.tableSize)
	return result
}

// Empty returns true iff this hash table is empty.
func (t *HashTable) Empty() bool { return t.count == 0 }

// TableSize returns the number of slots in the hash table.
func (t *HashTable) TableSize() int { return t.tableSize }

// Size returns the number of values in the hash table.
func (t *HashTable) Size() int { return t.count }

// Clear makes the hash table empty.
func (t *HashTable) Clear() {
	if t.tableSize < 3 {
		t.tableSize = DefaultTableSize
	}
	t.table = make([]*tableNode, t.tableSize)
	t.count = 0
}

// Get retrieves a value from a from a table given its key.
// Precondition: key is in the table.
// Precondition violation: return nil, false.
// Normal return: return valuev, true.
func (t *HashTable) Get(key containers.Hasher) (interface{}, bool) {
	if t.tableSize < 3 {
		t.Clear()
	}
	node := t.table[key.Hash(t.tableSize)]
	for node != nil {
		if node.key.Equal(key) {
			return node.value, true
		}
		node = node.next
	}
	return nil, false
}

// Insert puts v into the table, or replaces v if its is already there.
func (t *HashTable) Insert(key containers.Hasher, value interface{}) {
	if t.tableSize < 3 {
		t.Clear()
	}
	index := key.Hash(t.tableSize)
	node := t.table[index]
	for node != nil {
		if node.key.Equal(key) {
			node.value = value
			return
		}
		node = node.next
	}
	t.table[index] = newTableNode(key, value, t.table[index])
	t.count++
}

// Delete removes v from the table, or does nothing if it is not there.
func (t *HashTable) Delete(key containers.Hasher) {
	if t.tableSize < 3 {
		t.Clear()
	}
	index := key.Hash(t.tableSize)
	node := t.table[index]
	if node == nil {
		return
	}
	if key.Equal(node.key) {
		t.table[index] = node.next
		t.count--
		return
	}
	lastNode := node
	node = node.next
	for node != nil {
		if node.key.Equal(key) {
			lastNode.next = node.next
			t.count--
			return
		}
		lastNode = node
		node = node.next
	}
}

/////////////////////////////////////////////////////////////////////////////
// hashTableIterator keeps track of where we are in a table during iteration.
type hashTableIterator struct {
	table []*tableNode // reference to the table traversed
	index int          // table index for the next value
	node  *tableNode   // pointer to the node for the next value
}

// Reset prepares for a new iteration.
func (iter *hashTableIterator) Reset() {
	for iter.index = 0; iter.index < len(iter.table); iter.index++ {
		iter.node = iter.table[iter.index]
		if iter.node != nil {
			break
		}
	}
}

// Done returns true iff iteration is complete.
func (iter *hashTableIterator) Done() bool {
	return iter.node == nil
}

// Next returns the next value in the iteration.
// Precondition: there is a next value.
// Precondition violation: return nil and false.
// Normal return: return the next value and true.
func (iter *hashTableIterator) Next() (interface{}, bool) {
	if iter.node == nil {
		return nil, false
	}
	result := iter.node.value
	iter.node = iter.node.next
	if iter.node == nil {
		iter.index++
		for ; iter.index < len(iter.table); iter.index++ {
			iter.node = iter.table[iter.index]
			if iter.node != nil {
				break
			}
		}
	}
	return result, true
}

// NewIterator creates and returns a new external value iterator.
func (t *HashTable) NewIterator() containers.Iterator {
	result := new(hashTableIterator)
	result.table = t.table
	result.Reset()
	return result
}

/////////////////////////////////////////////////////////////////////////////
// hashTableKeyIterator keeps track of where we are in a table during iteration.
type hashTableKeyIterator struct {
	table []*tableNode // reference to the table traversed
	index int          // table index for the next value
	node  *tableNode   // pointer to the node for the next value
}

// Reset prepares for a new iteration.
func (iter *hashTableKeyIterator) Reset() {
	for iter.index = 0; iter.index < len(iter.table); iter.index++ {
		iter.node = iter.table[iter.index]
		if iter.node != nil {
			break
		}
	}
}

// Done returns true iff iteration is complete.
func (iter *hashTableKeyIterator) Done() bool {
	return iter.node == nil
}

// Next returns the next key in the iteration.
// Precondition: there is a next value.
// Precondition violation: return nil and false.
// Normal return: return the next value and true.
func (iter *hashTableKeyIterator) Next() (interface{}, bool) {
	if iter.node == nil {
		return nil, false
	}
	result := iter.node.key
	iter.node = iter.node.next
	if iter.node == nil {
		iter.index++
		for ; iter.index < len(iter.table); iter.index++ {
			iter.node = iter.table[iter.index]
			if iter.node != nil {
				break
			}
		}
	}
	return result, true
}

// NewKeyIterator creates and returns a new external key iterator.
func (t *HashTable) NewKeyIterator() containers.Iterator {
	result := new(hashTableKeyIterator)
	result.table = t.table
	result.Reset()
	return result
}

/////////////////////////////////////////////////////////////
// Helper functions /////////////////////////////////////////

// nextPrime finds the next prime greater than or equal to n.
func nextPrime(n int) int {
	if n <= 2 {
		return 2
	}
	if n%2 == 0 {
		n++
	}
	var result int
	for result = n; !isPrime(result); {
		result += 2
	}
	return result
}

// isPrime returns true iff n is prime.
func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	maxFactor := int(math.Sqrt(float64(n)))
	for k := 3; k <= maxFactor; k += 2 {
		if n%k == 0 {
			return false
		}
	}
	return true
}
