package toolkit

type Random struct{
	m_b, m_a uint64
}

func (r Random)NextInt(rangeInt uint64)uint64{
	// Use rejection to find a random value in a range that is a multiple of "range"
	n := (0xffffffffffffffff % rangeInt) + 1;
	var x uint64;
	for	; x+n < n; {
		x = r.Next();
	}
	// Use modulus to return the final value
	return x % rangeInt;
}
func (r Random)Next()uint64{
	// Use rejection to find a random value in a range that is a multiple of "range"
	r.m_a = 0x141F2B69 * (r.m_a & 0x3ffffffff) + (r.m_a >> 32);
	r.m_b = 0xC2785A6B * (r.m_b & 0x3ffffffff) + (r.m_b >> 32);
	return r.m_a ^ r.m_b;
}

func (r Random)SetSeed(seed uint64){
	r.m_b = 0xCA535ACA9535ACB2 + seed;
	r.m_a = 0x6CCF6660A66C35E7 + (seed << 24);
}

