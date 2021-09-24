// https://play.golang.org/p/Qg_uv_inCek
// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func main() {
	s := []string{"James", "Wagner", "Christene", "Mike"}
	fmt.Println(contains(s, "James")) // true
	fmt.Println(contains(s, "Jack")) // false
}
