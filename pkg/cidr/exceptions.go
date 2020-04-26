package cidr

func (n *Network) has31Exception() bool {
	ones, bits := n.Mask.Size()
	return bits == 32 && ones == 31
}
func (n *Network) has32Exception() bool {
	ones, bits := n.Mask.Size()
	return bits == 32 && ones == 32
}
