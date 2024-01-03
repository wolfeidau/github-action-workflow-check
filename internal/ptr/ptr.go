package ptr

func ToString(v *string) string {
	if v == nil {
		return ""
	}

	return *v
}
