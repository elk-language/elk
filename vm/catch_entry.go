package vm

type CatchEntry struct {
	From        int  // index of the first instruction that can be handled by this catch
	To          int  // index of the last instruction that can be handled by this catch
	JumpAddress int  // index of the byte that the VM should jump to
	Finally     bool // whether this entry is for a finally clause
}

// Number of bytes this catch covers
func (c *CatchEntry) ByteRange() int {
	return c.To - c.From
}

func NewCatchEntry(from, to, jumpAddress int, finally bool) *CatchEntry {
	return &CatchEntry{
		From:        from,
		To:          to,
		JumpAddress: jumpAddress,
		Finally:     finally,
	}
}
