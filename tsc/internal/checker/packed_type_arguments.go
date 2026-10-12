package checker

// Cache families use separate maps. Packed inputs are injective; a collision with
// a hashed fallback has the same 128-bit collision model as the original cache.
func packedTypeArgumentKey(tag uint64, types []*Type, flag bool) CacheHashKey {
	var ids [3]uint64
	for i, t := range types {
		ids[i] = uint64(t.id)
	}
	hi := tag<<48 | uint64(len(types))<<32 | ids[2]
	if flag {
		hi |= 1 << 34
	}
	return CacheHashKey{Hi: hi, Lo: ids[0] | ids[1]<<32}
}
