package sequencer

// Helper functions for deep copying pointers
func copyIntPtr(ptr *int) *int {
	if ptr == nil {
		return nil
	}
	copy := *ptr
	return &copy
}

func copyUint8Ptr(ptr *uint8) *uint8 {
	if ptr == nil {
		return nil
	}
	copy := *ptr
	return &copy
}

func copyUint8SlicePtr(ptr *[]uint8) *[]uint8 {
	if ptr == nil {
		return nil
	}
	copy := make([]uint8, len(*ptr))
	for i, v := range *ptr {
		copy[i] = v
	}
	return &copy
}
