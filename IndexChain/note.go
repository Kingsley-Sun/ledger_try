package main

type Note struct {
	HashID    string
	Timestamp string
}

func IsSameNotes(n1 []*Note,n2 []*Note) bool {
	if len(n1) != len(n2) {
		//fmt.Println("[~]Notes length don't match")
		return false
	}
	length := len(n1)
	for i := 0; i<length;i++ {
		note1 := *(n1[i])
		note2 := *(n2[i])
		if note1.Timestamp != note2.Timestamp {
			return false
		}
		if note1.HashID != note2.HashID {
			return false
		}
	}
	return true
}