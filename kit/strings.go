package kit

func ReverseString(s string) string {
	forwardS := []rune(s)

	var reversed []rune

	for i := len(forwardS) - 1; i > -1; i -= 1 {
		reversed = append(reversed, forwardS[i])
	}
	reversedS := string(reversed)
	return reversedS

}
