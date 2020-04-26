package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"fmt"
	"log"
)

const maxNoteInBlock = 10

type Mempool struct {
	NotesMap map[string]*Note
}

func (m *Mempool) CalNotesHash(notes []*Note) [20]byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(notes)
	if err != nil {
		log.Panic(err)
	}
	return sha1.Sum(result.Bytes())
}
func (m *Mempool) AddNote(note *Note) {
	//be careful of memory lea
	fmt.Println("[*]Add a New Note to Mempool")
	hash := note.HashID[:]
	m.NotesMap[hash] = note
}

func (m *Mempool) GetBlockNotes() []*Note {
	//fmt.Println("[~]Before A new Block Notesmap count is ", len(m.NotesMap))
	count := 0
	var res []*Note
	for _, v := range m.NotesMap {
		res = append(res, &Note{
			HashID:    v.HashID,
			Timestamp: v.Timestamp,
		})
		count += 1
		if count == maxNoteInBlock {
			break
		}
	}
	//fmt.Println("[~]Return A new Block Notesmap count is ", len(m.NotesMap))
	return res
}

func (m *Mempool) HasNote(hash string) bool {
	_, ok := m.NotesMap[hash]
	return ok
}

func (m *Mempool) DeleteNote(hash string) {
	delete(m.NotesMap, hash)
}

func (m *Mempool) PrintMempool() {
	fmt.Println("--------Notes Mempool is ", len(m.NotesMap)," -------")
	for _, v := range m.NotesMap {
		fmt.Println("+++++")
		fmt.Println("note HashID : ", v.HashID)
		fmt.Println("note timestamp : ", v.Timestamp)
		fmt.Println("+++++")
	}
}
