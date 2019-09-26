func flipBit(x float64) float64 {
  bits := math.Float64bits(x)
  bitnum := rand.Intn(52)
  mask := uint64(1) << uint(bitnum)
  bits ^= mask
  return math.Float64frombits(bits)
}
