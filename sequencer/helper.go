package sequencer

// Helper functions for deep copying pointers
func copyIntPtr(ptr *int) *int {
	if ptr == nil {
		return nil
	}
	c := *ptr
	return &c
}

func copyUint8Ptr(ptr *uint8) *uint8 {
	if ptr == nil {
		return nil
	}
	c := *ptr
	return &c
}

func copyUint8SlicePtr(ptr *[]uint8) *[]uint8 {
	if ptr == nil {
		return nil
	}
	c := make([]uint8, len(*ptr))
	copy(c, *ptr)
	return &c
}
