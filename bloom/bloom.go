package bloom

const (
	// SHBF is SingleHash BloomFilter
	SHBF = iota
	// LHBF is LessHashing BloomFilter
	LHBF
)

// BloomFilter ...
type BloomFilter struct {
	width  int
	kNum   int
	ftype  uint8
	Filter []uint16
}

// JHash hashing function
func JHash(str []byte, length uint) uint {
	var hash uint = 1315423911

	for i := 0; uint(i) < length; i++ {
		hash ^= ((hash << 5) + uint(str[i]) + (hash >> 2))
	}
	return hash
}

// OCaml hashing function
func OCaml(str []byte, length uint) uint {
	var hash uint

	for i := 0; uint(i) < length; i++ {
		hash = hash*19 + uint(str[i])
	}
	return hash
}

// NewBF ...
func NewBF(w, k int, t uint8) *BloomFilter {
	return &BloomFilter{width: w, kNum: k, ftype: t, Filter: make([]uint16, (w/16)+1)}
}

// Insert ...
func (bf *BloomFilter) Insert(str []byte, length uint) {
	switch bf.ftype {
	case 0:
		hash := JHash(str, length)

		for i := 0; i < bf.kNum; i++ {
			v1 := hash >> 16
			v2 := hash << uint(i)
			hVal := (v1 ^ v2) % uint(bf.width)

			p1 := hVal / 16
			p2 := hVal % 16
			tmp := 1 << (16 - p2)
			if p2 == 0 {
				if p1 <= 0 {
					p1 = uint(len(bf.Filter))
				}
				bf.Filter[p1-1] = bf.Filter[p1-1] | 1
			} else {
				bf.Filter[p1] = bf.Filter[p1] | uint16(tmp)
				v1 = v1 << 1
			}
		}
	case 1:
		v1 := JHash(str, length) % uint(bf.width)
		v2 := OCaml(str, length) % uint(bf.width)

		for i := 0; i < bf.kNum; i++ {
			hVal := (v1 + uint(i)*v2) % uint(bf.width)
			p1 := hVal / 16
			p2 := hVal % 16
			tmp := 1 << (16 - p2)

			if p2 == 0 {
				if p1 <= 0 {
					p1 = uint(len(bf.Filter))
				}
				bf.Filter[p1-1] = bf.Filter[p1-1] | 1
			} else {
				bf.Filter[p1] = bf.Filter[p1] | uint16(tmp)
			}
		}
	}
}

// Query ...
func (bf *BloomFilter) Query(str []byte, length uint) int {
	ans := 1

	switch bf.ftype {
	case 0:
		hash := JHash(str, length)
		var tmp uint16

		for i := 0; i < bf.kNum; i++ {
			v1 := hash >> 16
			v2 := hash << uint(i)
			hVal := (v1 ^ v2) % uint(bf.width)

			p1 := hVal / 16
			p2 := hVal % 16

			if p2 == 0 {
				if p1 <= 0 {
					p1 = uint(len(bf.Filter))
				}
				tmp = bf.Filter[p1-1] & 1
			} else {
				tmp = (bf.Filter[p1] >> (16 - p2)) & 1
			}

			if tmp == 0 {
				ans = 0
				break
			}
			v1 = v1 << 1
		}
		return ans
	case 1:
		v1 := JHash(str, length) % uint(bf.width)
		v2 := OCaml(str, length) % uint(bf.width)
		var tmp uint16

		for i := 0; i < bf.kNum; i++ {
			hVal := (v1 + uint(i)*v2) % uint(bf.width)
			p1 := hVal / 16
			p2 := hVal % 16

			if p2 == 0 {
				if p1 <= 0 {
					p1 = uint(len(bf.Filter))
				}
				tmp = bf.Filter[p1-1] & 1
			} else {
				tmp = (bf.Filter[p1] >> (16 - p2)) & 1
			}

			if tmp == 0 {
				ans = 0
				break
			}
		}
		return ans
	}
	return ans
}
